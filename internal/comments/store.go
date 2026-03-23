package comments

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/dotbrains/glimpse/internal/instance"
)

// Severity levels for review comments.
const (
	SeverityMustFix    = "must-fix"
	SeveritySuggestion = "suggestion"
	SeverityNit        = "nit"
	SeverityQuestion   = "question"
)

// Comment represents a single inline comment on a diff.
type Comment struct {
	ID        string    `json:"id"`
	File      string    `json:"file"`
	Line      int       `json:"line"`
	Side      string    `json:"side"`      // "old" or "new"
	Body      string    `json:"body"`
	Severity  string    `json:"severity"`  // must-fix, suggestion, nit, question
	Author    string    `json:"author"`    // "user" or "ai"
	Resolved  bool      `json:"resolved"`
	CreatedAt time.Time `json:"createdAt"`
}

// Store manages comment persistence.
type Store struct {
	mu   sync.RWMutex
	path string
	data []Comment
}

// NewStore creates or loads a comment store for the given session key.
func NewStore(sessionKey string) *Store {
	dir := filepath.Join(instance.DataDir(), "comments")
	_ = os.MkdirAll(dir, 0o755)
	path := filepath.Join(dir, sessionKey+".json")

	s := &Store{path: path}
	s.load()
	return s
}

func (s *Store) load() {
	data, err := os.ReadFile(s.path)
	if err != nil {
		s.data = []Comment{}
		return
	}
	if json.Unmarshal(data, &s.data) != nil {
		s.data = []Comment{}
	}
}

func (s *Store) save() error {
	data, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}

// Add creates a new comment and persists it.
func (s *Store) Add(file string, line int, side, body, severity, author string) (*Comment, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	c := Comment{
		ID:        generateID(),
		File:      file,
		Line:      line,
		Side:      side,
		Body:      body,
		Severity:  normalizeSeverity(severity),
		Author:    author,
		Resolved:  false,
		CreatedAt: time.Now(),
	}
	s.data = append(s.data, c)
	return &c, s.save()
}

// AddBatch creates multiple comments at once (used by AI review).
func (s *Store) AddBatch(comments []Comment) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range comments {
		if comments[i].ID == "" {
			comments[i].ID = generateID()
		}
		if comments[i].CreatedAt.IsZero() {
			comments[i].CreatedAt = time.Now()
		}
		comments[i].Severity = normalizeSeverity(comments[i].Severity)
	}
	s.data = append(s.data, comments...)
	return s.save()
}

// List returns all comments, optionally filtered.
func (s *Store) List(onlyOpen bool) []Comment {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []Comment
	for _, c := range s.data {
		if onlyOpen && c.Resolved {
			continue
		}
		result = append(result, c)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].File != result[j].File {
			return result[i].File < result[j].File
		}
		return result[i].Line < result[j].Line
	})
	return result
}

// Resolve marks a comment as resolved.
func (s *Store) Resolve(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.data {
		if s.data[i].ID == id {
			s.data[i].Resolved = true
			return s.save()
		}
	}
	return fmt.Errorf("comment %q not found", id)
}

// Delete removes a comment.
func (s *Store) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.data {
		if s.data[i].ID == id {
			s.data = append(s.data[:i], s.data[i+1:]...)
			return s.save()
		}
	}
	return fmt.Errorf("comment %q not found", id)
}

// Get returns a single comment by ID.
func (s *Store) Get(id string) (*Comment, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, c := range s.data {
		if c.ID == id {
			return &c, true
		}
	}
	return nil, false
}

// FormatForAgent returns comments as structured text that an AI agent can parse.
func FormatForAgent(comments []Comment) string {
	if len(comments) == 0 {
		return "No open comments."
	}
	var out string
	for _, c := range comments {
		tag := ""
		if c.Severity != "" {
			tag = fmt.Sprintf("[%s] ", c.Severity)
		}
		out += fmt.Sprintf("%s:%d %s%s (id: %s)\n", c.File, c.Line, tag, c.Body, c.ID)
	}
	return out
}

func generateID() string {
	b := make([]byte, 6)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func normalizeSeverity(s string) string {
	switch s {
	case SeverityMustFix, SeveritySuggestion, SeverityNit, SeverityQuestion:
		return s
	case "critical", "error", "bug":
		return SeverityMustFix
	case "warning", "improvement":
		return SeveritySuggestion
	case "style", "nitpick":
		return SeverityNit
	default:
		if s == "" {
			return SeveritySuggestion
		}
		return s
	}
}
