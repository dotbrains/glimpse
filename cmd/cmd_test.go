package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestExecute_Version(t *testing.T) {
	os.Args = []string{"glimpse", "--version"}
	err := Execute("0.0.1-test")
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}
}

func TestNewRootCmd(t *testing.T) {
	root := newRootCmd("0.1.0")
	if root.Use != "glimpse [ref] [ref]" {
		t.Errorf("Use = %q", root.Use)
	}

	// Verify subcommands.
	cmds := make(map[string]bool)
	for _, c := range root.Commands() {
		cmds[c.Name()] = true
	}
	for _, want := range []string{"config", "list", "review", "resolve", "tree", "resolve-tree"} {
		if !cmds[want] {
			t.Errorf("missing subcommand %q", want)
		}
	}
}

func TestNewRootCmd_Version(t *testing.T) {
	root := newRootCmd("1.2.3")
	if root.Version != "1.2.3" {
		t.Errorf("expected version 1.2.3, got %q", root.Version)
	}
}

func TestExecute_Help(t *testing.T) {
	root := newRootCmd("dev")
	root.SetArgs([]string{"--help"})
	var out bytes.Buffer
	root.SetOut(&out)

	err := root.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := out.String()
	if !strings.Contains(output, "glimpse") {
		t.Error("expected project name in help output")
	}
	if !strings.Contains(output, "config") {
		t.Error("expected 'config' subcommand in help")
	}
	if !strings.Contains(output, "review") {
		t.Error("expected 'review' subcommand in help")
	}
	if !strings.Contains(output, "resolve") {
		t.Error("expected 'resolve' subcommand in help")
	}
}

func TestListCmd_Help(t *testing.T) {
	root := newRootCmd("dev")
	root.SetArgs([]string{"list", "--help"})
	var out bytes.Buffer
	root.SetOut(&out)

	err := root.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "running glimpse instances") {
		t.Error("expected list help text")
	}
}

func TestReviewCmd_Help(t *testing.T) {
	root := newRootCmd("dev")
	root.SetArgs([]string{"review", "--help"})
	var out bytes.Buffer
	root.SetOut(&out)

	err := root.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	output := out.String()
	if !strings.Contains(output, "AI agent") {
		t.Error("expected review help text")
	}
	if !strings.Contains(output, "--focus") {
		t.Error("expected --focus flag in help")
	}
}

func TestResolveCmd_Help(t *testing.T) {
	root := newRootCmd("dev")
	root.SetArgs([]string{"resolve", "--help"})
	var out bytes.Buffer
	root.SetOut(&out)

	err := root.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "unresolved") {
		t.Error("expected resolve help text")
	}
}

func TestTreeCmd_Help(t *testing.T) {
	root := newRootCmd("dev")
	root.SetArgs([]string{"tree", "--help"})
	var out bytes.Buffer
	root.SetOut(&out)

	err := root.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "file tree browser") {
		t.Error("expected tree help text")
	}
}

func TestResolveTreeCmd_Help(t *testing.T) {
	root := newRootCmd("dev")
	root.SetArgs([]string{"resolve-tree", "--help"})
	var out bytes.Buffer
	root.SetOut(&out)

	err := root.Execute()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "tree browser") {
		t.Error("expected resolve-tree help text")
	}
}

func TestRootCmd_HasUnifiedFlag(t *testing.T) {
	root := newRootCmd("dev")
	f := root.Flags().Lookup("unified")
	if f == nil {
		t.Error("missing --unified flag")
	}
}

func TestRootCmd_HasDarkFlag(t *testing.T) {
	root := newRootCmd("dev")
	f := root.Flags().Lookup("dark")
	if f == nil {
		t.Error("missing --dark flag")
	}
	if f.DefValue != "true" {
		t.Errorf("--dark default = %q, want true", f.DefValue)
	}
}

func TestRunConfigInit(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	root := newRootCmd("test")
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"config", "init"})

	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}

	// Config file should exist.
	configPath := filepath.Join(tmp, ".config", "glimpse", "config.yaml")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("config file not created")
	}

	if !strings.Contains(buf.String(), "Wrote default config") {
		t.Error("expected success message")
	}
}

func TestRunConfigInit_AlreadyExists(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	// Pre-create config.
	configDir := filepath.Join(tmp, ".config", "glimpse")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(configDir, "config.yaml"), []byte("existing"), 0o644); err != nil {
		t.Fatal(err)
	}

	root := newRootCmd("test")
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"config", "init"})

	err := root.Execute()
	if err == nil {
		t.Fatal("expected error when config exists")
	}
}

func TestRunConfigInit_Force(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)

	// Pre-create config.
	configDir := filepath.Join(tmp, ".config", "glimpse")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(configDir, "config.yaml"), []byte("existing"), 0o644); err != nil {
		t.Fatal(err)
	}

	root := newRootCmd("test")
	buf := &bytes.Buffer{}
	root.SetOut(buf)
	root.SetArgs([]string{"config", "init", "--force"})

	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}

	if !strings.Contains(buf.String(), "Wrote default config") {
		t.Error("expected success message with --force")
	}
}
