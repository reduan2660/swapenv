package cmd_map

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/reduan2660/swapenv/internal/filehandler"
)

func Map(projectName string) error {
	project, err := filehandler.FindProjectByName(projectName)
	if err != nil {
		return err
	}
	if project == nil {
		return fmt.Errorf("project '%s' not found", projectName)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// Check if directory already mapped to another project
	existing, err := filehandler.FindProjectByLocalPath(cwd)
	if err != nil {
		return err
	}

	if existing != nil && existing.ProjectName != projectName {
		fmt.Printf("Directory already mapped to '%s'. Overwrite? [y/N]: ", existing.ProjectName)
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		if strings.ToLower(strings.TrimSpace(input)) != "y" {
			fmt.Println("Aborted.")
			return nil
		}
		// Clear old project's localPath
		existing.LocalPath = ""
		if err := filehandler.UpsertProjectDir(*existing); err != nil {
			return err
		}
	}

	// Check if this project already has a localPath
	if project.LocalPath != "" && project.LocalPath != cwd {
		fmt.Printf("Project already mapped to '%s'. Overwrite? [y/N]: ", project.LocalPath)
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		if strings.ToLower(strings.TrimSpace(input)) != "y" {
			fmt.Println("Aborted.")
			return nil
		}
	}

	project.LocalPath = cwd
	if err := filehandler.UpsertProjectDir(*project); err != nil {
		return err
	}

	fmt.Printf("Mapped '%s' → %s\n", projectName, cwd)
	return nil
}
