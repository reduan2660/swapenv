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

func Load(env string, replace bool) error {

	projectName, localOwner, localDirectory, homeDirectory, projectPath, err := GetBasicInfo(GetBasicInfoOptions{ReadOnly: false})
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

	fmt.Printf("replace %v\n", replace)

	if !replace {
		for envName := range envs {
			existingEnvValues, err := filehandler.ReadProjectEnv(projectPath, envName)
			if err == nil {
				envs[envName] = MergeEnv(envs[envName], existingEnvValues, false)
			}
		}
	}

	newProject := MarshalProject(projectName, localOwner, localDirectory, envs)

	projectJson, err := newProject.MarshalJSON()
	if err != nil {
		return fmt.Errorf("error generating json: %w", err)
	}

	if err := filehandler.WriteProject(homeDirectory, projectPath, projectJson); err != nil {
		return err
	}

	if err := filehandler.DeleteEnvFiles(files); err != nil {
		return err
	}

	fmt.Print("loaded environment:")
	for envName := range envs {
		fmt.Printf(" %s", envName)
	}

	return nil
}

type GetBasicInfoOptions struct {
	ReadOnly bool
}

func GetBasicInfo(opts GetBasicInfoOptions) (projectName, localOwner, localDirectory, homeDirectory, projectPath string, err error) {

	localDirectory, err = GetCurrentDirectory()
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("error getting local directory: %w", err)
	}

	existingProject, err := filehandler.FindProjectByLocalPath(localDirectory)
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("error checking project map: %w", err)
	}

	if existingProject != nil {
		projectName = existingProject.ProjectName
	} else if !opts.ReadOnly {

		projectName, err = GetProjectName()
		if err != nil {
			return "", "", "", "", "", fmt.Errorf("error getting project name: %w", err)
		}

		dirs, err := filehandler.ReadProjectDirs()
		if err != nil {
			return "", "", "", "", "", fmt.Errorf("error reading project map: %w", err)
		}

		nameExists := false
		for _, dir := range dirs {
			if dir.ProjectName == projectName {
				nameExists = true
				break
			}
		}

		if nameExists {
			parentDir := filepath.Base(filepath.Dir(localDirectory))
			projectName = parentDir + "/" + projectName
		}

		newProjectDir := types.ProjectDir{
			ProjectName: projectName,
			CurrentEnv:  "",
			LocalPath:   localDirectory,
			RemotePath:  "",
		}
		if err := filehandler.UpsertProjectDir(newProjectDir); err != nil {
			return "", "", "", "", "", fmt.Errorf("error adding project to map: %w", err)
		}

	}

	localOwner, err = GetLocalOwner()
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("error getting local user: %w", err)
	}
	if projectName != "" {
		homeDirectory, err = GetHomeDirectory(projectName)
		if err != nil {
			return "", "", "", "", "", fmt.Errorf("error getting home directory: %w", err)
		}

		projectPath = filepath.Join(homeDirectory, "project.json")
	}
	return
}

func GetProjectName() (string, error) {

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
