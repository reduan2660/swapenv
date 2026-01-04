/*
Copyright © 2025 Alve Reduan <hey@alvereduan.com>
*/
package cmd

import (
	"github.com/reduan2660/swapenv/internal/cmd_ls"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var lsCmd = &cobra.Command{
	Use:   "ls",
	Short: "list available environments for current project",
	RunE: func(cmd *cobra.Command, args []string) error {

		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		showVersions := viper.GetBool("version")
		return cmd_ls.ListEnv(showVersions)
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
	lsCmd.Flags().BoolP("version", "v", false, "show version information")
}

func GetLsCmd() *cobra.Command {
	return lsCmd
}
