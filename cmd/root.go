/*
Copyright © 2025 Alve Reduan <hey@alvereduan.com>
*/
package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/reduan2660/swapenv/internal/cmd_loader"
	"github.com/reduan2660/swapenv/internal/filehandler"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:   "swapenv",
	Short: "Switch and sync your environment conveniently",

	RunE: func(cmd *cobra.Command, args []string) error {
		return showProjectInfo()
	},

	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
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

	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError

		if !errors.As(err, &configFileNotFoundError) {
			return err
		}
	}

	err := viper.BindPFlags(cmd.Flags())
	if err != nil {
		return err
	}

	// fmt.Println("Configuration initialized. Using config file:", viper.ConfigFileUsed())
	return nil
}

func showProjectInfo() error {
	projectName, _, localDirectory, _, _, err := cmd_loader.GetBasicInfo(cmd_loader.GetBasicInfoOptions{ReadOnly: true})
	if err != nil {
		return fmt.Errorf("Error: %v\n", err)
	}

	if projectName == "" {
		fmt.Printf("no project under current directory, use swapenv load to initiate.")
		return nil
	}

	activeEnv, err := filehandler.ReadActiveEnv(localDirectory)
	if err != nil {
		return fmt.Errorf("error reading active environment: %w", err)
	}

	if activeEnv == "" {
		fmt.Printf("no active environment for %s - to list available environments run swapenv ls", projectName)
	} else {
		fmt.Printf("active environment: %s", activeEnv)
	}
	return nil
}
