/*
Copyright © 2025 Alve Reduan <hey@alvereduan.com>

*/
package cmd

import (
	"os"
	"fmt"
	"errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)


var (
	cfgFile string
)



var rootCmd = &cobra.Command{
	Use:   "swapenv",
	Short: "Switch and sync your environment convinently",

	PersistentPreRunE: func(cmd *cobra.Command, arge []string) error {
		return initializeConfig(cmd)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/swapenv/default.yaml)")
}

func initializeConfig(cmd *cobra.Command) error {

  // TODO: re-think : do we want to parse flags from envs that will already manage env?
	// viper.SetEnvPrefix("SWAPENV") 
	// viper.SetEnvKeyReplacer(strings.NewReplacer(".", "*", "-", "*"))
	// viper.AutomaticEnv()
	
	// Load The Config File
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(".")
		viper.AddConfigPath(home + "/.config/swapenv")
		viper.SetConfigName("default")
		viper.SetConfigType("yaml")
	}
	
	// Read the Config file
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError

		if !errors.As(err, &configFileNotFoundError) {
			return err
		}
	}

	// Bind cobra flags to viper
	err := viper.BindPFlags(cmd.Flags())
	if err != nil {
		return err
	}

	fmt.Println("Configuration initialized. Using config file:", viper.ConfigFileUsed())
	return nil
}
