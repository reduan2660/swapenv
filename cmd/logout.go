package cmd

import (
	"github.com/reduan2660/swapenv/internal/cmd_logout"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out from swapenv server",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		return cmd_logout.Logout()
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)
}

func GetLogoutCmd() *cobra.Command {
	return logoutCmd
}
