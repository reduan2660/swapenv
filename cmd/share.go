package cmd

import (
	"github.com/reduan2660/swapenv/internal/cmd_share"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var shareCmd = &cobra.Command{
	Use:   "share",
	Short: "Share environment",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		serverURL := viper.GetString("server")
		envName := args[0]
		return cmd_share.Share(serverURL, envName)
	},
}

func init() {
	rootCmd.AddCommand(shareCmd)
	shareCmd.Flags().String("server", "https://swapenv.sh", "swapenv server URL")
}

func GetShareCmd() *cobra.Command {
	return shareCmd
}
