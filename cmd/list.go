package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/dotbrains/glimpse/internal/instance"
	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	var jsonOutput bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "Show all running glimpse instances",
		RunE: func(cmd *cobra.Command, args []string) error {
			instances := instance.List()

			if jsonOutput {
				data, err := json.MarshalIndent(instances, "", "  ")
				if err != nil {
					return err
				}
				fmt.Fprintln(cmd.OutOrStdout(), string(data))
				return nil
			}

			fmt.Fprint(cmd.OutOrStdout(), instance.FormatTable(instances))
			return nil
		},
	}

	cmd.Flags().BoolVar(&jsonOutput, "json", false, "machine-readable JSON output")
	return cmd
}
