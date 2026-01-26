package cmd_spit

import (
	"fmt"
	"slices"

	"github.com/reduan2660/swapenv/internal/cmd_loader"
	"github.com/reduan2660/swapenv/internal/filehandler"
)

func Spit(envPattern, versionStr string) error {
	projectName, _, _, _, projectPath, err := cmd_loader.GetBasicInfo(cmd_loader.GetBasicInfoOptions{ReadOnly: true})
	if err != nil {
		return err
	}

	if projectName == "" {
		fmt.Printf("no project under current directory, use swapenv load to initiate.")
		return nil
	}

	if versionStr != "" {
		version, err := filehandler.ResolveVersion(projectName, versionStr)
		if err != nil {
			return err
		}
		projectPath, err = filehandler.GetVersionFilePath(projectName, version)
		if err != nil {
			return err
		}
	}

	envNames, err := filehandler.ListProjectEnv(projectPath)
	if err != nil {
		return fmt.Errorf("error reading project file: %w", err)
	}

	targetEnvs := envNames
	if envPattern != "*" {
		if !slices.Contains(envNames, envPattern) {
			return fmt.Errorf("environment '%s' not found", envPattern)
		}
		targetEnvs = []string{envPattern}
	}

	for _, envName := range targetEnvs {
		envValues, err := filehandler.ReadProjectEnv(projectPath, envName)
		if err != nil {
			return fmt.Errorf("error reading %s: %w", envName, err)
		}

		outputFile := fmt.Sprintf(".%s.env", envName)
		if err := filehandler.WriteEnv(envValues, outputFile, false); err != nil {
			return fmt.Errorf("error writing %s: %w", outputFile, err)
		}

		fmt.Printf("spit %s to %s\n", envName, outputFile)
	}

	return nil
}
