package cmd

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/dotbrains/glimpse/internal/comments"
	"github.com/dotbrains/glimpse/internal/diff"
	"github.com/dotbrains/glimpse/internal/git"
	"github.com/dotbrains/glimpse/internal/instance"
	"github.com/dotbrains/glimpse/internal/review"
	"github.com/spf13/cobra"
)

func newReviewCmd() *cobra.Command {
	var flagFocus string

	cmd := &cobra.Command{
		Use:   "review [refs]",
		Short: "Run AI code review and post comments to the viewer",
		Long:  "Sends the diff to an AI agent (claude by default), parses structured review output, and writes inline comments to the running glimpse instance.",
		Args:  cobra.MaximumNArgs(2),
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

			// Find running instance for this repo.
			inst := instance.FindByRepo(repoDir)
			if inst == nil {
				return fmt.Errorf("no running glimpse instance for %s — run `glimpse` first", repoDir)
			}

			// Resolve refs.
			base, compare := git.ResolveRefs(args, flagBase, flagCompare)

			// Get diff.
			rawDiff, err := gc.Diff(ctx, base, compare)
			if err != nil {
				return fmt.Errorf("git diff failed: %w", err)
			}
			if base == "" && compare == "" {
				staged, _ := gc.DiffStaged(ctx)
				if staged != "" {
					rawDiff = staged + "\n" + rawDiff
				}
			}

			files := diff.Parse(rawDiff)
			if len(files) == 0 {
				cmd.Println("No changes to review.")
				return nil
			}

			cmd.Printf("→ Reviewing %s...\n", diff.Summary(files))

			// Run AI review.
			reviewer := review.NewReviewer()
			reviewComments, err := reviewer.Run(ctx, rawDiff, flagFocus)
			if err != nil {
				return err
			}

			if len(reviewComments) == 0 {
				cmd.Println("✓ No issues found.")
				return nil
			}

			// Write comments to store.
			sessionKey := strconv.Itoa(inst.Port)
			store := comments.NewStore(sessionKey)
			if err := store.AddBatch(reviewComments); err != nil {
				return fmt.Errorf("saving comments: %w", err)
			}

			// Tally by severity.
			counts := map[string]int{}
			for _, c := range reviewComments {
				counts[c.Severity]++
			}

			cmd.Printf("✓ %d comments posted to %s\n", len(reviewComments), fmt.Sprintf("http://localhost:%d", inst.Port))
			for _, sev := range []string{"must-fix", "suggestion", "nit", "question"} {
				if n := counts[sev]; n > 0 {
					cmd.Printf("  %d %s\n", n, sev)
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&flagFocus, "focus", "", "focus review on specific areas (security, performance, testing)")
	return cmd
}
