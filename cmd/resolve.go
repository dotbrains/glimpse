package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/dotbrains/glimpse/internal/comments"
	"github.com/dotbrains/glimpse/internal/git"
	"github.com/dotbrains/glimpse/internal/instance"
	"github.com/spf13/cobra"
)

func newResolveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resolve [comment-id]",
		Short: "Output open comments for your agent to resolve",
		Long:  "Reads all open (unresolved) comments from the running glimpse instance and outputs them as structured text. Your AI agent can parse this output and make the requested code changes.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			gc := git.NewClient(cwd)
			repoDir, _ := gc.TopLevel(cmd.Context())

			inst := instance.FindByRepo(repoDir)
			if inst == nil {
				return fmt.Errorf("no running glimpse instance for %s — run `glimpse` first", cwd)
			}

			sessionKey := strconv.Itoa(inst.Port)
			store := comments.NewStore(sessionKey)

			// If a specific ID is given, output just that comment.
			if len(args) == 1 {
				c, ok := store.Get(args[0])
				if !ok {
					return fmt.Errorf("comment %q not found", args[0])
				}
				fmt.Fprint(cmd.OutOrStdout(), comments.FormatForAgent([]comments.Comment{*c}))
				return nil
			}

			// Output all open comments.
			open := store.List(true)
			fmt.Fprint(cmd.OutOrStdout(), comments.FormatForAgent(open))
			return nil
		},
	}
	return cmd
}
