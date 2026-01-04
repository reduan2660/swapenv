package cmd_ls

import (
	"fmt"

	"github.com/reduan2660/swapenv/internal/cmd_loader"
	"github.com/reduan2660/swapenv/internal/filehandler"
)

func ListEnv(showVersions bool) error {
	projectName, _, localDirectory, _, projectPath, err := cmd_loader.GetBasicInfo(cmd_loader.GetBasicInfoOptions{ReadOnly: true})
	if err != nil {
		return err
	}

	if projectName == "" {
		fmt.Printf("no project under current directory, use swapenv load to initiate.")
		return nil
	}

	envNames, err := filehandler.ListProjectEnv(projectPath)
	if err != nil {
		return fmt.Errorf("error reading project file: %w", err)
	}

	fmt.Printf("available environments:")
	for _, en := range envNames {
		fmt.Printf(" %s", en)
	}

	if showVersions {
		project, err := filehandler.FindProjectByLocalPath(localDirectory)
		if err != nil {
			return err
		}

		versions, err := filehandler.ListVersions(projectName)
		if err != nil {
			return err
		}

		fmt.Printf("\nversions:\n")
		for _, v := range versions {
			marker := "  "
			if v == project.CurrentVersion {
				marker = "* "
			}

			name := ""
			if n, ok := project.VersionNames[fmt.Sprintf("%d", v)]; ok {
				name = fmt.Sprintf(" (%s)", n)
			}

			latest := ""
			if v == project.LatestVersion {
				latest = " [latest]"
			}

			fmt.Printf("%s%d%s%s\n", marker, v, name, latest)
		}
	}

	return nil
}
