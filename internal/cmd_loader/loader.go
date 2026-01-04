package cmd_loader

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/reduan2660/swapenv/internal/filehandler"
	"github.com/reduan2660/swapenv/internal/types"
)

func Load(env string, replace bool) error {

	projectName, localOwner, localDirectory, homeDirectory, _, err := GetBasicInfo(GetBasicInfoOptions{ReadOnly: false})
	if err != nil {
		return err
	}

	filename := "." + env + ".env"

	files, err := filepath.Glob(filename)
	if err != nil {
		return fmt.Errorf("error getting files: %w", err)
	}

	if len(files) == 0 {
		fmt.Print("no environment to load")
		return nil
	}

	newVersion, err := filehandler.BumpVersion(projectName)
	if err != nil {
		return fmt.Errorf("error bumping version: %w", err)
	}

	versionPath, err := filehandler.GetVersionFilePath(projectName, newVersion)
	if err != nil {
		return fmt.Errorf("error getting version path: %w", err)
	}

	envs := map[string][]types.EnvValue{}

	for _, file := range files {

		envName := strings.ToLower(strings.Split(file, ".")[1])

		fileContent, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("error reading file: %w", err)
		}

		envValues, err := ParseEnv(fileContent)
		if err != nil {
			return err
		}
		envs[envName] = envValues
	}

	if !replace {

		prevVersion := newVersion - 1
		if prevVersion > 0 {
			prevPath, _ := filehandler.GetVersionFilePath(projectName, prevVersion)

			for envName := range envs {
				existingEnvValues, err := filehandler.ReadProjectEnv(prevPath, envName)
				if err == nil {
					envs[envName] = MergeEnv(envs[envName], existingEnvValues, MergeEnvConfig{
						ConflictPriority: "incoming",
					})
				}
			}

		}
	}

	newProject := MarshalProject(projectName, localOwner, localDirectory, envs)

	projectJson, err := newProject.MarshalJSON()
	if err != nil {
		return fmt.Errorf("error generating json: %w", err)
	}

	if err := filehandler.WriteProject(homeDirectory, versionPath, projectJson); err != nil {
		return err
	}

	if err := filehandler.DeleteEnvFiles(files); err != nil {
		return err
	}

	if err := filehandler.PruneVersions(projectName); err != nil {
		return fmt.Errorf("error pruning versions: %w", err)
	}

	fmt.Print("loaded environment:")
	for envName := range envs {
		fmt.Printf(" %s", envName)
	}

	return nil
}
