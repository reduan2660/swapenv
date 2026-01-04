package cmd

import (
	"github.com/reduan2660/swapenv/internal/cmd_receive"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var receiveCmd = &cobra.Command{
	Use:   "receive",
	Short: "Receive environment from another device",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		serverURL := viper.GetString("server")
		return cmd_receive.Receive(serverURL)
	},
}

func init() {
	rootCmd.AddCommand(receiveCmd)
	receiveCmd.Flags().String("server", "https://swapenv.sh", "swapenv server URL")
}

func GetReceiveCmd() *cobra.Command {
	return receiveCmd
}
