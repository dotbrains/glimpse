package server

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"

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
	Data DiffData
	Port int
}

// NewServer creates a server ready to serve.
func NewServer(data DiffData, port int) *Server {
	return &Server{Data: data, Port: port}
}

// ListenAndServe starts the HTTP server.
func (s *Server) ListenAndServe() error {
	mux := http.NewServeMux()

	// API endpoint.
	mux.HandleFunc("/api/diff", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode(s.Data)
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
