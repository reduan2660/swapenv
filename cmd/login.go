package cmd

import (
	"github.com/reduan2660/swapenv/internal/cmd_login"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with swapenv server",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		serverURL := viper.GetString("server")
		return cmd_login.Login(serverURL)
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	loginCmd.Flags().String("server", "https://swapenv.sh", "swapenv server URL")
}

func GetLoginCmd() *cobra.Command {
	return loginCmd
}
