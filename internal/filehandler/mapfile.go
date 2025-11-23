package filehandler

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/reduan2660/swapenv/internal/types"
	"github.com/spf13/viper"
)

func GetBaseDir() (string, error) {
	customHome := viper.GetString("home_directory")
	if customHome != "" {
		if err := os.MkdirAll(customHome, 0755); err != nil {
			return "", fmt.Errorf("failed to create custom home directory: %w", err)
		}
		return customHome, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	baseDir := filepath.Join(homeDir, ".swapenv")
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return "", err
	}

	return baseDir, nil

}

func GetHomeDirectory(projectName string) (string, error) {

	homeDir, err := GetBaseDir()
	if err != nil {
		return "", err
	}

	projectPath := filepath.Join(homeDir, projectName)
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		if err := os.MkdirAll(projectPath, 0755); err != nil {
			return "", err
		}
	}

	return projectPath, nil
}

func GetMapFilePath() (string, error) {
	homeDir, err := GetBaseDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, "map.json"), nil
}

func ReadProjectDirs() ([]types.ProjectDir, error) {
	mapPath, err := GetMapFilePath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(mapPath); os.IsNotExist(err) {
		return []types.ProjectDir{}, nil
	}

	data, err := os.ReadFile(mapPath)
	if err != nil {
		return nil, err
	}

	var dirs []types.ProjectDir
	if err := json.Unmarshal(data, &dirs); err != nil {
		return nil, err
	}

	return dirs, nil
}

func FindProjectByLocalPath(localPath string) (*types.ProjectDir, error) {
	dirs, err := ReadProjectDirs()
	if err != nil {
		return nil, err
	}

	for _, dir := range dirs {
		if dir.LocalPath == localPath {
			return &dir, nil
		}
	}

	return nil, nil
}

func ReadActiveEnv(localPath string) (string, error) {
	projectDir, err := FindProjectByLocalPath(localPath)
	if err != nil {
		return "", err
	}

	return projectDir.CurrentEnv, nil

}

func WriteProjectDirs(dirs []types.ProjectDir) error {
	mapPath, err := GetMapFilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(dirs, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(mapPath, data, 0644)
}

func UpsertProjectDir(projectDir types.ProjectDir) error {
	dirs, err := ReadProjectDirs()
	if err != nil {
		return err
	}

	found := false
	for i, dir := range dirs {
		if dir.LocalPath == projectDir.LocalPath {
			dirs[i] = projectDir
			found = true
			break
		}
	}

	if !found {
		dirs = append(dirs, projectDir)
	}

	return WriteProjectDirs(dirs)
}

func UpdateCurrentEnv(projectName, envName string) error {
	dirs, err := ReadProjectDirs()
	if err != nil {
		return err
	}

	found := false
	for i, dir := range dirs {
		if dir.ProjectName == projectName {
			dirs[i].CurrentEnv = envName
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("project not found in map: %s", projectName)
	}

	return WriteProjectDirs(dirs)
}
