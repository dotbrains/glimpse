package instance

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// Info describes a running glimpse instance.
type Info struct {
	PID     int       `json:"pid"`
	Port    int       `json:"port"`
	RepoDir string    `json:"repoDir"`
	Base    string    `json:"base,omitempty"`
	Compare string    `json:"compare,omitempty"`
	Started time.Time `json:"started"`
}

// DataDir returns the glimpse data directory.
func DataDir() string {
	if xdg := os.Getenv("XDG_DATA_HOME"); xdg != "" {
		return filepath.Join(xdg, "glimpse")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "share", "glimpse")
}

func instanceDir() string {
	return filepath.Join(DataDir(), "instances")
}

func pidFile(port int) string {
	return filepath.Join(instanceDir(), fmt.Sprintf("%d.json", port))
}

// Register writes a PID file for the current instance.
func Register(info Info) error {
	dir := instanceDir()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}
	return os.WriteFile(pidFile(info.Port), data, 0o644)
}

// Unregister removes the PID file for the given port.
func Unregister(port int) {
	os.Remove(pidFile(port))
}

// List returns all running instances (stale ones are cleaned up).
func List() []Info {
	dir := instanceDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	var result []Info
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		var info Info
		if json.Unmarshal(data, &info) != nil {
			continue
		}
		// Check if process is still alive.
		if !isAlive(info.PID) {
			os.Remove(filepath.Join(dir, e.Name()))
			continue
		}
		result = append(result, info)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Port < result[j].Port
	})
	return result
}

// FindByRepo returns an existing instance for the given repo directory, or nil.
func FindByRepo(repoDir string) *Info {
	for _, inst := range List() {
		if inst.RepoDir == repoDir {
			return &inst
		}
	}
	return nil
}

// NextPort returns the next available port starting from basePort.
func NextPort(basePort int) int {
	used := make(map[int]bool)
	for _, inst := range List() {
		used[inst.Port] = true
	}
	port := basePort
	for used[port] {
		port++
	}
	return port
}

// StopByRepo kills the instance for a given repo and removes its PID file.
func StopByRepo(repoDir string) bool {
	inst := FindByRepo(repoDir)
	if inst == nil {
		return false
	}
	p, err := os.FindProcess(inst.PID)
	if err == nil {
		_ = p.Signal(syscall.SIGTERM)
	}
	Unregister(inst.Port)
	return true
}

func isAlive(pid int) bool {
	p, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	err = p.Signal(syscall.Signal(0))
	return err == nil
}

// FormatTable returns a human-readable table of running instances.
func FormatTable(instances []Info) string {
	if len(instances) == 0 {
		return "No running instances."
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%-6s %-7s %-40s %s\n", "PORT", "PID", "REPO", "REFS"))
	sb.WriteString(strings.Repeat("─", 80) + "\n")
	for _, inst := range instances {
		refs := "working tree"
		if inst.Base != "" && inst.Compare != "" {
			refs = inst.Base + ".." + inst.Compare
		} else if inst.Base != "" {
			refs = inst.Base
		}
		sb.WriteString(fmt.Sprintf("%-6s %-7s %-40s %s\n",
			strconv.Itoa(inst.Port),
			strconv.Itoa(inst.PID),
			truncate(inst.RepoDir, 40),
			refs,
		))
	}
	return sb.String()
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return "..." + s[len(s)-max+3:]
}
