package diff

import (
	"strings"
	"testing"
)

const sampleDiff = `diff --git a/main.go b/main.go
index abc1234..def5678 100644
--- a/main.go
+++ b/main.go
@@ -10,6 +10,7 @@ func main() {
 	fmt.Println("hello")
 	fmt.Println("world")
+	fmt.Println("new line")
 	os.Exit(0)
 }
diff --git a/README.md b/README.md
new file mode 100644
--- /dev/null
+++ b/README.md
@@ -0,0 +1,3 @@
+# My Project
+
+A description.
`

func TestParse_BasicDiff(t *testing.T) {
	files := Parse(sampleDiff)
	if len(files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(files))
	}

	// First file: main.go (modified).
	f := files[0]
	if f.NewName != "main.go" {
		t.Errorf("file[0].NewName = %q, want main.go", f.NewName)
	}
	if f.Status != "modified" {
		t.Errorf("file[0].Status = %q, want modified", f.Status)
	}
	if f.Language != "go" {
		t.Errorf("file[0].Language = %q, want go", f.Language)
	}
	if len(f.Hunks) != 1 {
		t.Fatalf("file[0] expected 1 hunk, got %d", len(f.Hunks))
	}

	h := f.Hunks[0]
	if h.OldStart != 10 || h.NewStart != 10 {
		t.Errorf("hunk starts: old=%d new=%d", h.OldStart, h.NewStart)
	}

	// Count added lines.
	add, del := f.Stats()
	if add != 1 || del != 0 {
		t.Errorf("stats: +%d -%d, want +1 -0", add, del)
	}

	// Second file: README.md (added).
	f2 := files[1]
	if f2.Status != "added" {
		t.Errorf("file[1].Status = %q, want added", f2.Status)
	}
	if f2.Language != "markdown" {
		t.Errorf("file[1].Language = %q, want markdown", f2.Language)
	}
	add2, _ := f2.Stats()
	if add2 != 3 {
		t.Errorf("README add count = %d, want 3", add2)
	}
}

func TestParse_EmptyDiff(t *testing.T) {
	files := Parse("")
	if len(files) != 0 {
		t.Errorf("expected 0 files for empty diff, got %d", len(files))
	}
}

func TestParse_DeletedFile(t *testing.T) {
	raw := `diff --git a/old.txt b/old.txt
deleted file mode 100644
--- a/old.txt
+++ /dev/null
@@ -1,2 +0,0 @@
-line one
-line two
`
	files := Parse(raw)
	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}
	if files[0].Status != "deleted" {
		t.Errorf("status = %q, want deleted", files[0].Status)
	}
	_, del := files[0].Stats()
	if del != 2 {
		t.Errorf("deletions = %d, want 2", del)
	}
}

func TestParse_RenamedFile(t *testing.T) {
	raw := `diff --git a/old.go b/new.go
similarity index 95%
rename from old.go
rename to new.go
`
	files := Parse(raw)
	if len(files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(files))
	}
	if files[0].Status != "renamed" {
		t.Errorf("status = %q, want renamed", files[0].Status)
	}
}

func TestSummary(t *testing.T) {
	files := Parse(sampleDiff)
	s := Summary(files)
	if !strings.Contains(s, "2 files changed") {
		t.Errorf("summary = %q, expected '2 files changed'", s)
	}
	if !strings.Contains(s, "4 insertions(+)") {
		t.Errorf("summary = %q, expected '4 insertions(+)'", s)
	}
}

func TestGuessLanguage(t *testing.T) {
	tests := map[string]string{
		"main.go":       "go",
		"app.tsx":       "tsx",
		"style.css":     "css",
		"Makefile":      "makefile",
		"Dockerfile":    "dockerfile",
		"config.yaml":   "yaml",
		"unknown.xyz":   "",
		"src/lib/a.rs":  "rust",
		"scripts/x.sh":  "bash",
	}
	for name, want := range tests {
		got := guessLanguage(name)
		if got != want {
			t.Errorf("guessLanguage(%q) = %q, want %q", name, got, want)
		}
	}
}

func TestParseHunkHeader(t *testing.T) {
	h := parseHunkHeader("@@ -10,6 +10,7 @@ func main() {")
	if h.OldStart != 10 || h.OldCount != 6 {
		t.Errorf("old: %d,%d want 10,6", h.OldStart, h.OldCount)
	}
	if h.NewStart != 10 || h.NewCount != 7 {
		t.Errorf("new: %d,%d want 10,7", h.NewStart, h.NewCount)
	}
}
