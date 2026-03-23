package instance

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestRegisterAndList(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_DATA_HOME", tmp)

	info := Info{
		PID:     os.Getpid(), // current process, so isAlive will be true
		Port:    5391,
		RepoDir: "/tmp/test-repo",
		Base:    "main",
		Compare: "feature",
		Started: time.Now(),
	}
	if err := Register(info); err != nil {
		t.Fatalf("Register error: %v", err)
	}

	instances := List()
	if len(instances) != 1 {
		t.Fatalf("expected 1 instance, got %d", len(instances))
	}
	if instances[0].Port != 5391 {
		t.Errorf("port = %d, want 5391", instances[0].Port)
	}
	if instances[0].RepoDir != "/tmp/test-repo" {
		t.Errorf("repoDir = %q", instances[0].RepoDir)
	}
}

func TestUnregister(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_DATA_HOME", tmp)

	info := Info{PID: os.Getpid(), Port: 5555, RepoDir: "/tmp/r", Started: time.Now()}
	Register(info)

	Unregister(5555)

	instances := List()
	if len(instances) != 0 {
		t.Errorf("expected 0 instances after unregister, got %d", len(instances))
	}
}

func TestFindByRepo(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_DATA_HOME", tmp)

	Register(Info{PID: os.Getpid(), Port: 5391, RepoDir: "/a", Started: time.Now()})
	Register(Info{PID: os.Getpid(), Port: 5392, RepoDir: "/b", Started: time.Now()})

	found := FindByRepo("/b")
	if found == nil {
		t.Fatal("expected to find instance for /b")
	}
	if found.Port != 5392 {
		t.Errorf("port = %d, want 5392", found.Port)
	}

	notFound := FindByRepo("/nonexistent")
	if notFound != nil {
		t.Error("expected nil for nonexistent repo")
	}
}

func TestNextPort(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_DATA_HOME", tmp)

	// No instances: should return base port.
	port := NextPort(5391)
	if port != 5391 {
		t.Errorf("expected 5391, got %d", port)
	}

	// Register one on 5391.
	Register(Info{PID: os.Getpid(), Port: 5391, RepoDir: "/x", Started: time.Now()})

	port = NextPort(5391)
	if port != 5392 {
		t.Errorf("expected 5392, got %d", port)
	}
}

func TestListCleansStaleInstances(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("XDG_DATA_HOME", tmp)

	// Register with a PID that almost certainly doesn't exist.
	info := Info{PID: 999999999, Port: 5391, RepoDir: "/stale", Started: time.Now()}
	Register(info)

	instances := List()
	// Should be cleaned up since PID 999999999 isn't alive.
	if len(instances) != 0 {
		t.Errorf("expected stale instance to be cleaned, got %d", len(instances))
	}
}

func TestFormatTable_Empty(t *testing.T) {
	out := FormatTable(nil)
	if out != "No running instances." {
		t.Errorf("expected 'No running instances.', got %q", out)
	}
}

func TestFormatTable_WithInstances(t *testing.T) {
	instances := []Info{
		{PID: 123, Port: 5391, RepoDir: "/home/user/project", Base: "", Compare: ""},
		{PID: 456, Port: 5392, RepoDir: "/home/user/api", Base: "main", Compare: "feature"},
	}
	out := FormatTable(instances)
	if !strings.Contains(out, "5391") {
		t.Error("expected port 5391 in output")
	}
	if !strings.Contains(out, "working tree") {
		t.Error("expected 'working tree' for empty refs")
	}
	if !strings.Contains(out, "main..feature") {
		t.Error("expected 'main..feature' for refs")
	}
	if !strings.Contains(out, "PORT") {
		t.Error("expected header row")
	}
}

func TestFormatTable_BaseOnlyRef(t *testing.T) {
	instances := []Info{
		{PID: 1, Port: 5391, RepoDir: "/r", Base: "HEAD~3", Compare: ""},
	}
	out := FormatTable(instances)
	if !strings.Contains(out, "HEAD~3") {
		t.Error("expected 'HEAD~3' for base-only ref")
	}
}

func TestTruncate(t *testing.T) {
	if truncate("short", 40) != "short" {
		t.Error("short string should not be truncated")
	}
	long := "/very/long/path/that/exceeds/the/maximum/allowed/length/for/display"
	result := truncate(long, 40)
	if len(result) != 40 {
		t.Errorf("expected length 40, got %d", len(result))
	}
	if !strings.HasPrefix(result, "...") {
		t.Error("truncated string should start with ...")
	}
}

func TestDataDir_XDG(t *testing.T) {
	t.Setenv("XDG_DATA_HOME", "/custom/data")
	dir := DataDir()
	if dir != "/custom/data/glimpse" {
		t.Errorf("DataDir = %q, want /custom/data/glimpse", dir)
	}
}

func TestDataDir_Default(t *testing.T) {
	t.Setenv("XDG_DATA_HOME", "")
	dir := DataDir()
	if !strings.Contains(dir, ".local/share/glimpse") {
		t.Errorf("DataDir = %q, expected to contain .local/share/glimpse", dir)
	}
}
