package cmd_info

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/reduan2660/swapenv/internal/filehandler"
)

type ProjectInfo struct {
	Project       string   `json:"project,omitempty"`
	Env           string   `json:"env,omitempty"`
	Envs          []string `json:"envs,omitempty"`
	Version       int      `json:"version,omitempty"`
	LatestVersion int      `json:"latest_version,omitempty"`
}

func Info(format string) error {
	info := ProjectInfo{}

	cwd, err := os.Getwd()
	if err != nil {
		return outputInfo(info, format)
	}

	project, err := filehandler.FindProjectByLocalPath(cwd)
	if err != nil || project == nil {
		return outputInfo(info, format)
	}

	info.Project = project.ProjectName
	info.Env = project.CurrentEnv
	info.Version = project.CurrentVersion
	info.LatestVersion = project.LatestVersion

	version := project.CurrentVersion
	if version == 0 {
		version = project.LatestVersion
	}
	if version > 0 {
		projectPath, err := filehandler.GetVersionFilePath(project.ProjectName, version)
		if err == nil {
			envs, err := filehandler.ListProjectEnv(projectPath)
			if err == nil {
				info.Envs = envs
			}
		}
	}

	return outputInfo(info, format)
}

func outputInfo(info ProjectInfo, format string) error {
	if format == "plain" {
		if info.Project == "" {
			return nil
		}
		env := info.Env
		if env == "" {
			env = "none"
		}
		fmt.Printf("%s:%s\n", info.Project, env)
		return nil
	}

	data, err := json.Marshal(info)
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
