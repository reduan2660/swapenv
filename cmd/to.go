/*
Copyright © 2025 Alve Reduan <hey@alvereduan.com>
*/
package cmd

import (
	"github.com/reduan2660/swapenv/internal/cmd_setter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var toCmd = &cobra.Command{
	Use:   "to",
	Short: "Sets an environment for this project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := viper.BindPFlags(cmd.Flags()); err != nil {
			return err
		}
		envName := args[0]
		replace := viper.GetBool("replace")
		skipCommon := viper.GetBool("skip-common")
		return cmd_setter.Set(envName, replace, skipCommon)
	},
}

func init() {
	rootCmd.AddCommand(toCmd)
	toCmd.Flags().Bool("replace", false, "to replace the existing .env instead of overwriting")
	toCmd.Flags().Bool("skip-common", false, "dont append common env variables (if exists)")
}

func GetToCmd() *cobra.Command {
	return toCmd
}
