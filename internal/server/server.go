package server

import (
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strings"

	"github.com/dotbrains/glimpse/internal/comments"
	"github.com/dotbrains/glimpse/internal/diff"
	"github.com/dotbrains/glimpse/internal/gh"
)

// NewGHClient creates a GitHub CLI client (thin wrapper to avoid import cycle).
func NewGHClient() *gh.Client {
	return gh.NewClient()
}

//go:embed static/*
var staticFS embed.FS

// DiffData is the JSON payload served at /api/diff.
type DiffData struct {
	Repo     string          `json:"repo"`
	Base     string          `json:"base"`
	Compare  string          `json:"compare"`
	Summary  string          `json:"summary"`
	Files    []diff.FileDiff `json:"files"`
	ViewMode string          `json:"viewMode"` // "split" or "unified"
	Theme    string          `json:"theme"`    // "dark" or "light"
	PROwner  string          `json:"prOwner,omitempty"`
	PRRepo   string          `json:"prRepo,omitempty"`
	PRNumber string          `json:"prNumber,omitempty"`
}

// TreeData is the JSON payload served at /api/tree.
type TreeData struct {
	Repo    string   `json:"repo"`
	Files   []string `json:"files"`
	RepoDir string   `json:"-"`
}

// Server serves the diff viewer UI and API.
type Server struct {
	Data     DiffData
	Tree     *TreeData
	Port     int
	Comments *comments.Store
	GitDir   string // repo root for tree file serving
}

// NewServer creates a server ready to serve.
func NewServer(data DiffData, port int, store *comments.Store) *Server {
	return &Server{Data: data, Port: port, Comments: store}
}

// NewTreeServer creates a server for the file tree browser.
func NewTreeServer(tree TreeData, port int, store *comments.Store, gitDir string) *Server {
	return &Server{Tree: &tree, Port: port, Comments: store, GitDir: gitDir}
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

	// Tree API.
	mux.HandleFunc("/api/tree", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if s.Tree != nil {
			_ = json.NewEncoder(w).Encode(s.Tree)
		} else {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"files": []string{}})
		}
	})

	// File content API for tree browser.
	mux.HandleFunc("/api/file", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		path := r.URL.Query().Get("path")
		if path == "" || s.GitDir == "" {
			http.Error(w, "missing path", http.StatusBadRequest)
			return
		}
		// Read file from working tree.
		import_path := s.GitDir + "/" + path
		data, err := os.ReadFile(import_path)
		if err != nil {
			http.Error(w, "file not found", http.StatusNotFound)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"path": path, "content": string(data)})
	})

	// GitHub PR push/pull endpoints.
	mux.HandleFunc("/api/gh/push", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if s.Data.PROwner == "" {
			http.Error(w, "not a PR diff", http.StatusBadRequest)
			return
		}
		all := s.Comments.List(false)
		ghClient := NewGHClient()
		if err := ghClient.PostReviewComments(r.Context(), s.Data.PROwner, s.Data.PRRepo, s.Data.PRNumber, all); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]int{"pushed": len(all)})
	})

	mux.HandleFunc("/api/gh/pull", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		if s.Data.PROwner == "" {
			http.Error(w, "not a PR diff", http.StatusBadRequest)
			return
		}
		ghClient := NewGHClient()
		imported, err := ghClient.FetchPRComments(r.Context(), s.Data.PROwner, s.Data.PRRepo, s.Data.PRNumber)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(imported) > 0 {
			_ = s.Comments.AddBatch(imported)
		}
		_ = json.NewEncoder(w).Encode(map[string]int{"pulled": len(imported)})
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
