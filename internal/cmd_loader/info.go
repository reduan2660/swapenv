package cmd_loader

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"

	"github.com/reduan2660/swapenv/internal/filehandler"
	"github.com/reduan2660/swapenv/internal/types"
)

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

		if err := filehandler.MigrateProjectIfNeeded(projectName); err != nil {
			return "", "", "", "", "", fmt.Errorf("error migrating project: %w", err)
		}
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
			ProjectName:    projectName,
			CurrentEnv:     "",
			LocalPath:      localDirectory,
			RemotePath:     "",
			CurrentVersion: 0,
			LatestVersion:  0,
			VersionNames:   make(map[string]string),
		}
		if err := filehandler.UpsertProjectDir(newProjectDir); err != nil {
			return "", "", "", "", "", fmt.Errorf("error adding project to map: %w", err)
		}

		existingProject, err = filehandler.FindProjectByLocalPath(localDirectory)
		if err != nil {
			return "", "", "", "", "", fmt.Errorf("error fetching new project: %w", err)
		}

	}

	localOwner, err = GetLocalOwner()
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("error getting local user: %w", err)
	}
	if projectName != "" {
		homeDirectory, err = filehandler.GetHomeDirectory(projectName)
		if err != nil {
			return "", "", "", "", "", fmt.Errorf("error getting home directory: %w", err)
		}

		existingProject, _ = filehandler.FindProjectByLocalPath(localDirectory) // Re-fetch after potential migration

		version := existingProject.CurrentVersion
		if version == 0 {
			version = 1 // fist load will create v1
		}

		// projectPath = filepath.Join(homeDirectory, "project.json")
		projectPath, err = filehandler.GetVersionFilePath(projectName, version)
		if err != nil {
			return "", "", "", "", "", fmt.Errorf("error gettting version path: %w", err)
		}
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
