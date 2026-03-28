package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/dotbrains/glimpse/internal/comments"
	"github.com/dotbrains/glimpse/internal/git"
	"github.com/dotbrains/glimpse/internal/instance"
	"github.com/dotbrains/glimpse/internal/server"
	"github.com/spf13/cobra"
)

func newTreeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tree",
		Short: "Browse project files with syntax highlighting and comments",
		Long:  "Opens a file tree browser in your browser. Browse your repo, read files with syntax highlighting, and leave comments on any file or line.",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

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

			files, err := gc.ListFiles(ctx, "HEAD")
			if err != nil {
				return fmt.Errorf("listing files: %w", err)
			}

			port := flagPort
			if port == 0 {
				port = instance.NextPort(defaultBasePort)
			}

			treeData := server.TreeData{
				Repo:    repoName,
				Files:   files,
				RepoDir: repoDir,
			}

			store := comments.NewStore("tree-" + strconv.Itoa(port))
			srv := server.NewTreeServer(treeData, port, store, repoDir)

			info := instance.Info{
				PID:     os.Getpid(),
				Port:    port,
				RepoDir: repoDir,
				Base:    "tree",
				Started: time.Now(),
			}
			_ = instance.Register(info)

			if !flagQuiet {
				cmd.Printf("→ %d files in %s\n", len(files), repoName)
				cmd.Printf("→ Serving at %s/tree.html\n", srv.Addr())
			}

			if !flagNoOpen {
				openBrowser(srv.Addr() + "/tree.html")
			}

			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
			go func() {
				<-sigCh
				instance.Unregister(port)
				os.Exit(0)
			}()

			return srv.ListenAndServe()
		},
	}
	return cmd
}
