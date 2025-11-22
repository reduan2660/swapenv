package cmd_ls

import (
	"fmt"

	"github.com/reduan2660/swapenv/internal/cmd_loader"
	"github.com/reduan2660/swapenv/internal/filehandler"
)

func ListEnv() error {
	projectName, _, _, _, projectPath, err := cmd_loader.GetBasicInfo(cmd_loader.GetBasicInfoOptions{ReadOnly: true})
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

	return nil
}
