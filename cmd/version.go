package cmd

import (
	"fmt"
	"strconv"

	"github.com/reduan2660/swapenv/internal/cmd_version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version [version]",
	Short: "Show or switch version",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd_version.Show()
		}
		return cmd_version.Set(args[0])
	},
}

var versionLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List all versions",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd_version.List()
	},
}

var versionRenameCmd = &cobra.Command{
	Use:   "rename <version> <name>",
	Short: "Name a version (protects from auto-deletion)",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd_version.Rename(args[0], args[1])
	},
}

var versionRollbackCmd = &cobra.Command{
	Use:   "rollback [steps]",
	Short: "Go back N versions (default 1)",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		steps := 1
		if len(args) > 0 {
			var err error
			steps, err = strconv.Atoi(args[0])
			if err != nil || steps < 1 {
				return fmt.Errorf("steps must be a positive integer")
			}
		}
		return cmd_version.Rollback(steps)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.AddCommand(versionLsCmd)
	versionCmd.AddCommand(versionRenameCmd)
}
