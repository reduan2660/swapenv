package cmd_setter

import (
	"fmt"
	"os"
	"sort"

	"github.com/reduan2660/swapenv/internal/cmd_loader"
	"github.com/reduan2660/swapenv/internal/filehandler"
)

func Set(env string, replace bool, skipCommon bool, versionStr string, nowrap bool) error {

	projectName, _, _, _, projectPath, err := cmd_loader.GetBasicInfo(cmd_loader.GetBasicInfoOptions{ReadOnly: false})
	if err != nil {
		return err
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

	incomingEnvValues, err := filehandler.ReadProjectEnv(projectPath, env)
	if err != nil {
		return err
	}

	sort.Slice(incomingEnvValues, func(i, j int) bool {
		return incomingEnvValues[i].Order < incomingEnvValues[j].Order
	})

	if env == "common" {
		// TODO: define behaviour
		return fmt.Errorf("setting to common isnt allowed")
	}

	if !skipCommon {
		commonEnvValues, err := filehandler.ReadProjectEnv(projectPath, "common")

		if err == nil {
			sort.Slice(commonEnvValues, func(i, j int) bool {
				return commonEnvValues[i].Order < commonEnvValues[j].Order
			})

			// merge incoming (dev) with common: dev order first, dev values win for conflicts
			incomingEnvValues = cmd_loader.MergeEnv(commonEnvValues, incomingEnvValues, cmd_loader.MergeEnvConfig{
				ConflictPriority: "current",
			})
		}
	}

	envFilePath := ".env" // TODO: consider parent

	if _, err := os.Stat(envFilePath); os.IsNotExist(err) {
		if err := os.WriteFile(envFilePath, []byte{}, 0644); err != nil {
			return err
		}
	}

	curEnvFile, err := os.ReadFile(envFilePath)
	if err != nil {
		return err
	}

	curEnvValues, err := cmd_loader.ParseEnv(curEnvFile)
	if err != nil {
		return fmt.Errorf("error parsing .env: %w", err)
	}

	mergedEnv := cmd_loader.MergeEnv(incomingEnvValues, curEnvValues, cmd_loader.MergeEnvConfig{
		Replace:          replace,
		ConflictPriority: "incoming",
	})
	if err := filehandler.WriteEnv(mergedEnv, envFilePath, !nowrap); err != nil {
		return fmt.Errorf("error writing .env: %w", err)
	}

	if err := filehandler.UpdateCurrentEnv(projectName, env); err != nil {
		return fmt.Errorf("error updating current env: %w", err)
	}

	fmt.Printf("Swapped environment to: %v\n", env)
	return nil
}
