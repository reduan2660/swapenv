/*
Copyright © 2025 Alve Reduan <hey@alvereduan.com>
*/
package cmd

import (
	"github.com/reduan2660/swapenv/internal/cmd_info"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "output project info as JSON (used for integrations)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		format := viper.GetString("format")
		envOnly := viper.GetBool("env-only")
		return cmd_info.Info(format, envOnly)
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
	infoCmd.Flags().String("format", "json", "output format (json|plain)")
	infoCmd.Flags().Bool("env-only", false, "output only the environment name (plain format)")
}

func GetInfoCmd() *cobra.Command {
	return infoCmd
}
