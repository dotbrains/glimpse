package cmd

import (
	"fmt"
	"os"

	"github.com/dotbrains/glimpse/internal/comments"
	"github.com/dotbrains/glimpse/internal/git"
	"github.com/dotbrains/glimpse/internal/instance"
	"github.com/spf13/cobra"
)

func newResolveTreeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resolve-tree [comment-id]",
		Short: "Output open comments from the tree browser for your agent to resolve",
		Long:  "Reads all open (unresolved) comments from the tree browser instance and outputs them as structured text for AI agents.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			gc := git.NewClient(cwd)
			repoDir, _ := gc.TopLevel(cmd.Context())

			// Find a tree instance for this repo.
			var treeInst *instance.Info
			for _, inst := range instance.List() {
				if inst.RepoDir == repoDir && inst.Base == "tree" {
					treeInst = &inst
					break
				}
			}
			if treeInst == nil {
				return fmt.Errorf("no running glimpse tree instance for %s — run `glimpse tree` first", cwd)
			}

			sessionKey := "tree-" + fmt.Sprintf("%d", treeInst.Port)
			store := comments.NewStore(sessionKey)

			if len(args) == 1 {
				c, ok := store.Get(args[0])
				if !ok {
					return fmt.Errorf("comment %q not found", args[0])
				}
				fmt.Fprint(cmd.OutOrStdout(), comments.FormatForAgent([]comments.Comment{*c}))
				return nil
			}

			open := store.List(true)
			fmt.Fprint(cmd.OutOrStdout(), comments.FormatForAgent(open))
			return nil
		},
	}
	return cmd
}
