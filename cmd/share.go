package cmd

import (
	"github.com/reduan2660/swapenv/internal/cmd_share"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var shareCmd = &cobra.Command{
	Use:   "share",
	Short: "Share environment",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}

		serverURL := viper.GetString("server")
		projectName := viper.GetString("project")
		envName := viper.GetString("env")
		version := viper.GetString("version")

		return cmd_share.Share(serverURL, projectName, envName, version)
	},
}

func init() {
	rootCmd.AddCommand(shareCmd)
	shareCmd.Flags().String("project", "", "project to share (default: current directory)")
	shareCmd.Flags().String("env", "", "specific environment to share (default: all)")
	shareCmd.Flags().String("version", "latest", "version to share")
}

func GetShareCmd() *cobra.Command {
	return shareCmd
}
