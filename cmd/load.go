/*
Copyright © 2025 Alve Reduan <hey@alvereduan.com>
*/
package cmd

import (
	"github.com/reduan2660/swapenv/internal/cmd_loader"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Loads all the environment files",
	RunE: func(cmd *cobra.Command, args []string) error {
		envName := viper.GetString("env")
		cmd_loader.Load(envName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(loadCmd)
	loadCmd.Flags().String("env", "*", "Specific environment to load")
}
