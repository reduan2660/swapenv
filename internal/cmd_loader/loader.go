package cmd_loader

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/reduan2660/swapenv/internal/filehandler"
	"github.com/reduan2660/swapenv/internal/types"
)

func Load(env string) error {

	projectName, localOwner, localDirectory, homeDirectory, projectPath, err := GetBasicInfo()
	if err != nil {
		return err
	}

	filename := "." + env + ".env"

	files, err := filepath.Glob(filename)
	if err != nil {
		return fmt.Errorf("error getting files: %w", err)
	}

	if len(files) == 0 {
		return nil
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

	newProject := MarshalProject(projectName, localOwner, localDirectory, envs)

	projectJson, err := newProject.MarshalJSON()
	if err != nil {
		return fmt.Errorf("error generating json: %w", err)
	}

	// Todo - load existing prject and merge
	// - 1 homeDirectory need to figure out .swapenv config -> read from config -> default to ~/.swapenv/projectname

	if err := filehandler.WriteProject(homeDirectory, projectPath, projectJson); err != nil {
		return err
	}

	if err := filehandler.DeleteEnvFiles(files); err != nil {
		return err
	}

	return nil
}

func GetBasicInfo() (projectName, localOwner, localDirectory, homeDirectory, projectPath string, err error) {

	projectName, err = GetProjectName()
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("error getting project name: %w", err)
	}

	localOwner, err = GetLocalOwner()
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("error getting local user: %w", err)
	}

	localDirectory, err = GetCurrentDirectory()
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("error getting local directory: %w", err)
	}

	homeDirectory, err = GetHomeDirectory(projectName)
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("error getting home directory: %w", err)
	}

	projectPath = filepath.Join(homeDirectory, "project.json")

	return
}

func GetProjectName() (string, error) {
	// TODO : need to consider --project-name, in case of failure need to ask user for input
	currentDirectory, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Base(currentDirectory), nil
}

func GetLocalOwner() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}

	return currentUser.Username, nil
}

func GetCurrentDirectory() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return cwd, nil
}

func GetHomeDirectory(projectName string) (string, error) {
	// TODO: need improvement, need to consider config

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	projectPath := filepath.Join(homeDir, ".swapenv", projectName)
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		if err := os.MkdirAll(projectPath, 0755); err != nil { // TODO: re-thing permission
			return "", err
		}
	}

	return projectPath, nil
}
