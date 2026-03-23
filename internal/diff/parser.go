package diff

import (
	"fmt"
	"strconv"
	"strings"
)

// LineType represents a diff line kind.
type LineType string

const (
	LineContext  LineType = "context"
	LineAdded   LineType = "added"
	LineRemoved LineType = "removed"
)

// Line is a single line in a diff hunk.
type Line struct {
	Type    LineType `json:"type"`
	Content string   `json:"content"`
	OldNum  int      `json:"oldNum,omitempty"` // line number in old file (0 = N/A)
	NewNum  int      `json:"newNum,omitempty"` // line number in new file (0 = N/A)
}

// Hunk is a section of a diff.
type Hunk struct {
	OldStart int    `json:"oldStart"`
	OldCount int    `json:"oldCount"`
	NewStart int    `json:"newStart"`
	NewCount int    `json:"newCount"`
	Header   string `json:"header"`
	Lines    []Line `json:"lines"`
}

// FileDiff represents the diff for a single file.
type FileDiff struct {
	OldName  string `json:"oldName"`
	NewName  string `json:"newName"`
	Status   string `json:"status"` // "modified", "added", "deleted", "renamed"
	Language string `json:"language"`
	Hunks    []Hunk `json:"hunks"`
}

// Stats returns (additions, deletions) for this file.
func (f *FileDiff) Stats() (int, int) {
	var add, del int
	for _, h := range f.Hunks {
		for _, l := range h.Lines {
			switch l.Type {
			case LineAdded:
				add++
			case LineRemoved:
				del++
			}
		}
	}
	return add, del
}

// Parse parses unified diff output into structured FileDiff slices.
func Parse(raw string) []FileDiff {
	var files []FileDiff
	var current *FileDiff

	lines := strings.Split(raw, "\n")
	i := 0

	for i < len(lines) {
		line := lines[i]

		// New file diff header.
		if strings.HasPrefix(line, "diff --git ") {
			if current != nil {
				files = append(files, *current)
			}
			current = &FileDiff{}
			current.OldName, current.NewName = parseGitDiffHeader(line)
			current.Status = "modified"
			current.Language = guessLanguage(current.NewName)
			i++
			continue
		}

		if current == nil {
			i++
			continue
		}

		// File status indicators.
		if strings.HasPrefix(line, "new file mode") {
			current.Status = "added"
			i++
			continue
		}
		if strings.HasPrefix(line, "deleted file mode") {
			current.Status = "deleted"
			i++
			continue
		}
		if strings.HasPrefix(line, "similarity index") || strings.HasPrefix(line, "rename from") || strings.HasPrefix(line, "rename to") {
			current.Status = "renamed"
			i++
			continue
		}

		// --- and +++ headers.
		if strings.HasPrefix(line, "--- ") {
			if strings.HasPrefix(line, "--- a/") {
				current.OldName = line[6:]
			}
			i++
			continue
		}
		if strings.HasPrefix(line, "+++ ") {
			if strings.HasPrefix(line, "+++ b/") {
				current.NewName = line[6:]
				current.Language = guessLanguage(current.NewName)
			}
			i++
			continue
		}

		// Hunk header.
		if strings.HasPrefix(line, "@@") {
			hunk := parseHunkHeader(line)
			oldLine := hunk.OldStart
			newLine := hunk.NewStart
			i++

			for i < len(lines) {
				hl := lines[i]
				if strings.HasPrefix(hl, "diff --git ") || strings.HasPrefix(hl, "@@") {
					break
				}

				if strings.HasPrefix(hl, "+") {
					hunk.Lines = append(hunk.Lines, Line{
						Type:    LineAdded,
						Content: hl[1:],
						NewNum:  newLine,
					})
					newLine++
				} else if strings.HasPrefix(hl, "-") {
					hunk.Lines = append(hunk.Lines, Line{
						Type:    LineRemoved,
						Content: hl[1:],
						OldNum:  oldLine,
					})
					oldLine++
				} else if strings.HasPrefix(hl, " ") {
					hunk.Lines = append(hunk.Lines, Line{
						Type:    LineContext,
						Content: hl[1:],
						OldNum:  oldLine,
						NewNum:  newLine,
					})
					oldLine++
					newLine++
				} else if strings.HasPrefix(hl, "\\") {
					// "\ No newline at end of file" — skip.
					i++
					continue
				} else {
					// Empty context line (just a newline).
					hunk.Lines = append(hunk.Lines, Line{
						Type:    LineContext,
						Content: "",
						OldNum:  oldLine,
						NewNum:  newLine,
					})
					oldLine++
					newLine++
				}
				i++
			}

			current.Hunks = append(current.Hunks, hunk)
			continue
		}

		i++
	}

	if current != nil {
		files = append(files, *current)
	}
	return files
}

// parseGitDiffHeader extracts file paths from "diff --git a/foo b/bar".
func parseGitDiffHeader(line string) (old, new string) {
	// "diff --git a/path b/path"
	line = strings.TrimPrefix(line, "diff --git ")
	parts := strings.SplitN(line, " b/", 2)
	if len(parts) == 2 {
		old = strings.TrimPrefix(parts[0], "a/")
		new = parts[1]
	}
	return
}

// parseHunkHeader parses "@@ -old,count +new,count @@ optional context".
func parseHunkHeader(line string) Hunk {
	h := Hunk{Header: line}

	// Extract the range part between @@ markers.
	parts := strings.SplitN(line, "@@", 3)
	if len(parts) < 2 {
		return h
	}
	ranges := strings.TrimSpace(parts[1])

	for _, r := range strings.Fields(ranges) {
		if strings.HasPrefix(r, "-") {
			start, count := parseRange(r[1:])
			h.OldStart = start
			h.OldCount = count
		} else if strings.HasPrefix(r, "+") {
			start, count := parseRange(r[1:])
			h.NewStart = start
			h.NewCount = count
		}
	}

	if len(parts) == 3 && strings.TrimSpace(parts[2]) != "" {
		h.Header = strings.TrimSpace(parts[2])
	}

	return h
}

func parseRange(s string) (start, count int) {
	parts := strings.SplitN(s, ",", 2)
	start, _ = strconv.Atoi(parts[0])
	if len(parts) == 2 {
		count, _ = strconv.Atoi(parts[1])
	} else {
		count = 1
	}
	return
}

// guessLanguage returns a rough language identifier from a file extension.
func guessLanguage(name string) string {
	ext := strings.ToLower(name)
	if idx := strings.LastIndex(ext, "."); idx >= 0 {
		ext = ext[idx:]
	}
	m := map[string]string{
		".go":    "go",
		".js":    "javascript",
		".jsx":   "jsx",
		".ts":    "typescript",
		".tsx":   "tsx",
		".py":    "python",
		".rb":    "ruby",
		".rs":    "rust",
		".java":  "java",
		".c":     "c",
		".h":     "c",
		".cpp":   "cpp",
		".cs":    "csharp",
		".swift": "swift",
		".kt":    "kotlin",
		".sh":    "bash",
		".bash":  "bash",
		".fish":  "fish",
		".zsh":   "bash",
		".yaml":  "yaml",
		".yml":   "yaml",
		".json":  "json",
		".toml":  "toml",
		".xml":   "xml",
		".html":  "html",
		".css":   "css",
		".scss":  "scss",
		".sql":   "sql",
		".md":    "markdown",
		".proto": "protobuf",
		".tf":    "hcl",
	}
	if lang, ok := m[ext]; ok {
		return lang
	}
	// Check full filenames.
	base := name
	if idx := strings.LastIndex(base, "/"); idx >= 0 {
		base = base[idx+1:]
	}
	switch base {
	case "Makefile", "Dockerfile", "Jenkinsfile":
		return strings.ToLower(base)
	}
	return ""
}

// Summary returns a human-readable one-line summary of the parsed diff.
func Summary(files []FileDiff) string {
	var totalAdd, totalDel int
	for _, f := range files {
		a, d := f.Stats()
		totalAdd += a
		totalDel += d
	}
	return fmt.Sprintf("%d files changed, %d insertions(+), %d deletions(-)", len(files), totalAdd, totalDel)
}
