package cmd_setter

import (
	"fmt"
	"os"
	"sort"

	"github.com/reduan2660/swapenv/internal/cmd_loader"
	"github.com/reduan2660/swapenv/internal/filehandler"
	"github.com/reduan2660/swapenv/internal/types"
)

func Set(env string, replace bool) error {

	projectName, _, _, _, projectPath, err := cmd_loader.GetBasicInfo()
	if err != nil {
		return err
	}

	// TODO: try to validate information like if localDirectory matches
	// fmt.Printf("%v %v %v %v", projectName, localOwner, localDirectory, homeDirectory)

	incomingEnvValues, err := filehandler.ReadProjectEnv(projectPath, env)
	if err != nil {
		return err
	}

	sort.Slice(incomingEnvValues, func(i, j int) bool {
		return incomingEnvValues[i].Order < incomingEnvValues[j].Order
	})

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

	mergedEnv := mergeEnv(incomingEnvValues, curEnvValues, replace)
	if err := filehandler.WriteEnv(mergedEnv, envFilePath); err != nil {
		return fmt.Errorf("error writing .env: %w", err)
	}

	if err := filehandler.UpdateCurrentEnv(projectName, env); err != nil {
		return fmt.Errorf("error updating current env: %w", err)
	}

	fmt.Printf("Swapped environment to: %v\n", env)
	return nil
}

func mergeEnv(incoming, current []types.EnvValue, replace bool) []types.EnvValue {

	if replace {
		return incoming
	}

	incomingMap := make(map[string]types.EnvValue)
	for _, ev := range incoming {
		incomingMap[ev.Key] = ev
	}

	marked := make(map[string]bool)
	merged := make([]types.EnvValue, 0)

	for _, ev := range current {
		if incomingVal, exists := incomingMap[ev.Key]; exists { // overwrite if exists
			incomingVal.Order = ev.Order

			merged = append(merged, incomingVal)
			marked[ev.Key] = true
		} else {
			merged = append(merged, ev)
		}
	}

	for idx, ev := range incoming {
		if !marked[ev.Key] {
			ev.Order = len(current) + idx
			merged = append(merged, ev)
		}
	}

	return merged
}
