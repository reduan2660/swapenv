/*
Copyright © 2025 Alve Reduan <hey@alvereduan.com>
*/
package cmd

import (
	"github.com/reduan2660/swapenv/internal/cmd_spit"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var spitCmd = &cobra.Command{
	Use:   "spit",
	Short: "exports environment(s) back to .env files",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		envName := viper.GetString("env")
		return cmd_spit.Spit(envName)
	},
}

func init() {
	rootCmd.AddCommand(spitCmd)
	spitCmd.Flags().String("env", "*", "to spit specific env")
}

func GetSpitCmd() *cobra.Command {
	return spitCmd
}
