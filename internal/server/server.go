package server

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"strings"

	"github.com/dotbrains/glimpse/internal/comments"
	"github.com/dotbrains/glimpse/internal/diff"
)

//go:embed static/*
var staticFS embed.FS

// DiffData is the JSON payload served at /api/diff.
type DiffData struct {
	Repo    string          `json:"repo"`
	Base    string          `json:"base"`
	Compare string          `json:"compare"`
	Summary string          `json:"summary"`
	Files   []diff.FileDiff `json:"files"`
}

// Server serves the diff viewer UI and API.
type Server struct {
	Data     DiffData
	Port     int
	Comments *comments.Store
}

// NewServer creates a server ready to serve.
func NewServer(data DiffData, port int, store *comments.Store) *Server {
	return &Server{Data: data, Port: port, Comments: store}
}

// ListenAndServe starts the HTTP server.
func (s *Server) ListenAndServe() error {
	mux := http.NewServeMux()

	// Diff data.
	mux.HandleFunc("/api/diff", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		_ = json.NewEncoder(w).Encode(s.Data)
	})

	// Comments API.
	mux.HandleFunc("/api/comments", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PATCH,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		switch r.Method {
		case http.MethodGet:
			_ = json.NewEncoder(w).Encode(s.Comments.List(false))

		case http.MethodPost:
			var req struct {
				File     string `json:"file"`
				Line     int    `json:"line"`
				Side     string `json:"side"`
				Body     string `json:"body"`
				Severity string `json:"severity"`
			}
			body, _ := io.ReadAll(r.Body)
			if json.Unmarshal(body, &req) != nil {
				http.Error(w, "invalid JSON", http.StatusBadRequest)
				return
			}
			c, err := s.Comments.Add(req.File, req.Line, req.Side, req.Body, req.Severity, "user")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(c)

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Single comment operations: /api/comments/{id}
	mux.HandleFunc("/api/comments/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,PATCH,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		id := strings.TrimPrefix(r.URL.Path, "/api/comments/")
		if id == "" {
			http.Error(w, "missing id", http.StatusBadRequest)
			return
		}

		switch r.Method {
		case http.MethodPatch:
			if err := s.Comments.Resolve(id); err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]string{"status": "resolved"})

		case http.MethodDelete:
			if err := s.Comments.Delete(id); err != nil {
				http.Error(w, err.Error(), http.StatusNotFound)
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Static files.
	staticSub, err := fs.Sub(staticFS, "static")
	if err != nil {
		return fmt.Errorf("embedding static files: %w", err)
	}
	mux.Handle("/", http.FileServer(http.FS(staticSub)))

	return http.ListenAndServe(fmt.Sprintf(":%d", s.Port), mux)
}

// Addr returns the full URL.
func (s *Server) Addr() string {
	return fmt.Sprintf("http://localhost:%d", s.Port)
}
