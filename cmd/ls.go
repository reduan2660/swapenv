/*
Copyright © 2025 Alve Reduan <hey@alvereduan.com>
*/
package cmd

import (
	"github.com/reduan2660/swapenv/internal/cmd_ls"
	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "list available environments for current project",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd_ls.ListEnv()
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}

func GetLsCmd() *cobra.Command {
	return lsCmd
}
