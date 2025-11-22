package cmd_ls

import (
	"fmt"
	"os"

	"github.com/reduan2660/switchenv/internal/cmd_loader"
	"github.com/reduan2660/switchenv/internal/filehandler"
)

func ListEnv() error {
	projectName, _, _, _, projectPath, err := cmd_loader.GetBasicInfo()
	if err != nil {
		return err
	}

	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		fmt.Printf("%v project has not been initiated yet\n", projectName)
		return nil
	}

	envNames, err := filehandler.ListProjectEnv(projectPath)
	if err != nil {
		return fmt.Errorf("error reading project file: %w", err)
	}

	for _, en := range envNames {
		fmt.Println(en)
	}

	return nil
}
