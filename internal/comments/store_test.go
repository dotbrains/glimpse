package comments

import (
	"os"
	"strings"
	"testing"
)

func TestStore_AddAndList(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_DATA_HOME", tmp)

	s := NewStore("test-session")

	c, err := s.Add("main.go", 42, "new", "nil deref here", "must-fix", "user")
	if err != nil {
		t.Fatalf("Add error: %v", err)
	}
	if c.ID == "" {
		t.Error("expected non-empty ID")
	}
	if c.File != "main.go" || c.Line != 42 {
		t.Errorf("got file=%q line=%d", c.File, c.Line)
	}
	if c.Severity != SeverityMustFix {
		t.Errorf("severity = %q, want must-fix", c.Severity)
	}

	all := s.List(false)
	if len(all) != 1 {
		t.Fatalf("expected 1 comment, got %d", len(all))
	}
}

func TestStore_Resolve(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_DATA_HOME", tmp)

	s := NewStore("test-resolve")
	c, _ := s.Add("a.go", 1, "new", "fix this", "", "user")

	if err := s.Resolve(c.ID); err != nil {
		t.Fatalf("Resolve error: %v", err)
	}

	// Open-only list should be empty.
	open := s.List(true)
	if len(open) != 0 {
		t.Errorf("expected 0 open comments, got %d", len(open))
	}

	// All list should have 1.
	all := s.List(false)
	if len(all) != 1 {
		t.Errorf("expected 1 total comment, got %d", len(all))
	}
	if !all[0].Resolved {
		t.Error("comment should be resolved")
	}
}

func TestStore_Delete(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_DATA_HOME", tmp)

	s := NewStore("test-delete")
	c, _ := s.Add("a.go", 1, "new", "remove me", "", "user")

	if err := s.Delete(c.ID); err != nil {
		t.Fatalf("Delete error: %v", err)
	}

	all := s.List(false)
	if len(all) != 0 {
		t.Errorf("expected 0 comments after delete, got %d", len(all))
	}
}

func TestStore_DeleteNotFound(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_DATA_HOME", tmp)

	s := NewStore("test-notfound")
	err := s.Delete("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent ID")
	}
}

func TestStore_ResolveNotFound(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_DATA_HOME", tmp)

	s := NewStore("test-resolve-notfound")
	err := s.Resolve("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent ID")
	}
}

func TestStore_Get(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_DATA_HOME", tmp)

	s := NewStore("test-get")
	c, _ := s.Add("a.go", 10, "new", "found it", "nit", "ai")

	got, ok := s.Get(c.ID)
	if !ok {
		t.Fatal("expected to find comment")
	}
	if got.Body != "found it" {
		t.Errorf("body = %q", got.Body)
	}

	_, ok = s.Get("nonexistent")
	if ok {
		t.Error("expected not found for nonexistent ID")
	}
}

func TestStore_AddBatch(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_DATA_HOME", tmp)

	s := NewStore("test-batch")
	batch := []Comment{
		{File: "a.go", Line: 1, Body: "first", Severity: "critical", Author: "ai"},
		{File: "b.go", Line: 2, Body: "second", Severity: "nit", Author: "ai"},
		{File: "c.go", Line: 3, Body: "third", Severity: "", Author: "ai"},
	}
	if err := s.AddBatch(batch); err != nil {
		t.Fatalf("AddBatch error: %v", err)
	}

	all := s.List(false)
	if len(all) != 3 {
		t.Fatalf("expected 3 comments, got %d", len(all))
	}

	// "critical" should be normalized to "must-fix".
	for _, c := range all {
		if c.File == "a.go" && c.Severity != SeverityMustFix {
			t.Errorf("expected must-fix for critical, got %q", c.Severity)
		}
		// Empty severity should default to suggestion.
		if c.File == "c.go" && c.Severity != SeveritySuggestion {
			t.Errorf("expected suggestion for empty, got %q", c.Severity)
		}
		// IDs should be generated.
		if c.ID == "" {
			t.Error("expected non-empty ID")
		}
	}
}

func TestStore_ListSortedByFileAndLine(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_DATA_HOME", tmp)

	s := NewStore("test-sort")
	s.Add("z.go", 10, "new", "last file", "", "user")
	s.Add("a.go", 20, "new", "first file high line", "", "user")
	s.Add("a.go", 5, "new", "first file low line", "", "user")

	all := s.List(false)
	if len(all) != 3 {
		t.Fatalf("expected 3, got %d", len(all))
	}
	if all[0].File != "a.go" || all[0].Line != 5 {
		t.Errorf("first = %s:%d, want a.go:5", all[0].File, all[0].Line)
	}
	if all[1].File != "a.go" || all[1].Line != 20 {
		t.Errorf("second = %s:%d, want a.go:20", all[1].File, all[1].Line)
	}
	if all[2].File != "z.go" {
		t.Errorf("third = %s, want z.go", all[2].File)
	}
}

func TestStore_Persistence(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_DATA_HOME", tmp)

	// Write with one store instance.
	s1 := NewStore("test-persist")
	s1.Add("x.go", 1, "new", "persisted", "nit", "user")

	// Read with a new store instance.
	s2 := NewStore("test-persist")
	all := s2.List(false)
	if len(all) != 1 {
		t.Fatalf("expected 1 persisted comment, got %d", len(all))
	}
	if all[0].Body != "persisted" {
		t.Errorf("body = %q", all[0].Body)
	}
}

func TestFormatForAgent_Empty(t *testing.T) {
	out := FormatForAgent(nil)
	if out != "No open comments." {
		t.Errorf("expected 'No open comments.', got %q", out)
	}
}

func TestFormatForAgent_WithComments(t *testing.T) {
	cs := []Comment{
		{ID: "abc", File: "main.go", Line: 42, Severity: "must-fix", Body: "nil deref"},
		{ID: "def", File: "lib.go", Line: 10, Severity: "nit", Body: "rename var"},
	}
	out := FormatForAgent(cs)
	if !strings.Contains(out, "main.go:42 [must-fix] nil deref (id: abc)") {
		t.Errorf("missing first comment in output: %q", out)
	}
	if !strings.Contains(out, "lib.go:10 [nit] rename var (id: def)") {
		t.Errorf("missing second comment in output: %q", out)
	}
}

func TestNormalizeSeverity(t *testing.T) {
	tests := map[string]string{
		"must-fix":    "must-fix",
		"suggestion":  "suggestion",
		"nit":         "nit",
		"question":    "question",
		"critical":    "must-fix",
		"error":       "must-fix",
		"bug":         "must-fix",
		"warning":     "suggestion",
		"improvement": "suggestion",
		"style":       "nit",
		"nitpick":     "nit",
		"":            "suggestion",
		"unknown":     "unknown",
	}
	for input, want := range tests {
		got := normalizeSeverity(input)
		if got != want {
			t.Errorf("normalizeSeverity(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestStore_LoadCorruptFile(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_DATA_HOME", tmp)

	// Write corrupt JSON.
	dir := tmp + "/glimpse/comments"
	os.MkdirAll(dir, 0o755)
	os.WriteFile(dir+"/corrupt.json", []byte("{{not json"), 0o644)

	s := NewStore("corrupt")
	all := s.List(false)
	if len(all) != 0 {
		t.Errorf("expected 0 comments for corrupt file, got %d", len(all))
	}
}
