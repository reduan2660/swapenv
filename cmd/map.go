package cmd

import (
	"github.com/reduan2660/swapenv/internal/cmd_map"
	"github.com/spf13/cobra"
)

var mapCmd = &cobra.Command{
	Use:   "map <project-name>",
	Short: "Assign current directory to a project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd_map.Map(args[0])
	},
}

func init() {
	rootCmd.AddCommand(mapCmd)
}

func GetMapCmd() *cobra.Command {
	return mapCmd
}
