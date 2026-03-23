package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/dotbrains/glimpse/internal/config"
	"github.com/dotbrains/glimpse/internal/diff"
	"github.com/dotbrains/glimpse/internal/git"
	"github.com/dotbrains/glimpse/internal/instance"
	"github.com/dotbrains/glimpse/internal/server"
	"github.com/spf13/cobra"
)

var (
	flagBase    string
	flagCompare string
	flagPort    int
	flagNoOpen  bool
	flagQuiet   bool
	flagNew     bool
)

const defaultBasePort = 5391

func newRootCmd(version string) *cobra.Command {
	root := &cobra.Command{
		Use:   "glimpse [ref] [ref]",
		Short: "GitHub-style git diff viewer CLI",
		Long:  "Browser-based, GitHub-style diff viewer for git changes. View uncommitted changes, branch comparisons, commit ranges, and more with syntax-highlighted split diffs.",
		Args:  cobra.MaximumNArgs(2),
		RunE:  runDiff,
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
		Version: version,
	}

	root.SetVersionTemplate(fmt.Sprintf("glimpse version %s\n", version))

	root.Flags().StringVar(&flagBase, "base", "", "base ref to compare from (e.g. main, HEAD~3)")
	root.Flags().StringVar(&flagCompare, "compare", "", "ref to compare against base")
	root.Flags().IntVar(&flagPort, "port", 0, "custom port (default: auto-assigned from 5391)")
	root.Flags().BoolVar(&flagNoOpen, "no-open", false, "don't open browser")
	root.Flags().BoolVar(&flagQuiet, "quiet", false, "minimal terminal output")
	root.Flags().BoolVar(&flagNew, "new", false, "stop existing instance and start fresh")

	root.AddCommand(newListCmd())
	root.AddCommand(newConfigCmd())

	return root
}

func runDiff(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	if !git.GitInstalled() {
		return fmt.Errorf("git is not installed or not on PATH")
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	gc := git.NewClient(cwd)
	if !gc.IsRepo(ctx) {
		return fmt.Errorf("not a git repository: %s", cwd)
	}

	repoDir, _ := gc.TopLevel(ctx)
	repoName := gc.RepoName(ctx)

	// Check for existing instance.
	if existing := instance.FindByRepo(repoDir); existing != nil && !flagNew {
		url := fmt.Sprintf("http://localhost:%d", existing.Port)
		if !flagQuiet {
			cmd.Printf("→ Existing instance at %s\n", url)
		}
		if !flagNoOpen {
			openBrowser(url)
		}
		return nil
	}

	// Kill existing if --new.
	if flagNew {
		instance.StopByRepo(repoDir)
	}

	// Resolve refs.
	base, compare := git.ResolveRefs(args, flagBase, flagCompare)

	// Get diff.
	rawDiff, err := gc.Diff(ctx, base, compare)
	if err != nil {
		return fmt.Errorf("git diff failed: %w", err)
	}

	// Also include staged changes if no refs specified.
	if base == "" && compare == "" {
		staged, _ := gc.DiffStaged(ctx)
		if staged != "" {
			rawDiff = staged + "\n" + rawDiff
		}
	}

	files := diff.Parse(rawDiff)

	// Determine display refs.
	displayBase := base
	displayCompare := compare
	if displayBase == "" && displayCompare == "" {
		displayBase = "working tree"
	}

	data := server.DiffData{
		Repo:    repoName,
		Base:    displayBase,
		Compare: displayCompare,
		Summary: diff.Summary(files),
		Files:   files,
	}

	// Pick port.
	port := flagPort
	if port == 0 {
		port = instance.NextPort(defaultBasePort)
	}

	srv := server.NewServer(data, port)

	// Register instance.
	info := instance.Info{
		PID:     os.Getpid(),
		Port:    port,
		RepoDir: repoDir,
		Base:    base,
		Compare: compare,
		Started: time.Now(),
	}
	if err := instance.Register(info); err != nil && !flagQuiet {
		cmd.PrintErrf("⚠ Could not register instance: %v\n", err)
	}

	if !flagQuiet {
		cmd.Printf("→ %s\n", diff.Summary(files))
		cmd.Printf("→ Serving at %s\n", srv.Addr())
	}

	if !flagNoOpen {
		openBrowser(srv.Addr())
	}

	// Handle shutdown.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		instance.Unregister(port)
		os.Exit(0)
	}()

	return srv.ListenAndServe()
}

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration",
	}

	var force bool

	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Create default config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath, err := config.ConfigPath()
			if err != nil {
				return err
			}

			if !force {
				if _, err := os.Stat(cfgPath); err == nil {
					return fmt.Errorf("config already exists at %s (use --force to overwrite)", cfgPath)
				}
			}

			if err := config.Save(config.DefaultConfig()); err != nil {
				return err
			}

			display := cfgPath
			if home, err := os.UserHomeDir(); err == nil {
				if rel, err := filepath.Rel(home, cfgPath); err == nil {
					display = "~/" + rel
				}
			}

			cmd.Printf("✓ Wrote default config to %s\nEdit the file to customize settings.\n", display)
			return nil
		},
	}
	initCmd.Flags().BoolVar(&force, "force", false, "overwrite existing config")

	cmd.AddCommand(initCmd)
	return cmd
}

func openBrowser(url string) {
	var c string
	switch runtime.GOOS {
	case "darwin":
		c = "open"
	case "linux":
		c = "xdg-open"
	default:
		return
	}
	_ = exec.Command(c, url).Start()
}

// Execute runs the root command.
func Execute(version string) error {
	return newRootCmd(version).Execute()
}
