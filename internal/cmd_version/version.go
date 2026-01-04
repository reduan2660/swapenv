package cmd_version

import (
	"fmt"
	"strconv"

	"github.com/reduan2660/swapenv/internal/cmd_loader"
	"github.com/reduan2660/swapenv/internal/filehandler"
)

func Show() error {
	projectName, _, localDirectory, _, _, err := cmd_loader.GetBasicInfo(cmd_loader.GetBasicInfoOptions{ReadOnly: true})
	if err != nil {
		return err
	}

	if projectName == "" {
		fmt.Println("no project under current directory, use swapenv load to initiate.")
		return nil
	}

	project, err := filehandler.FindProjectByLocalPath(localDirectory)
	if err != nil {
		return err
	}

	name := ""
	if n, ok := project.VersionNames[strconv.Itoa(project.CurrentVersion)]; ok {
		name = fmt.Sprintf(" (%s)", n)
	}
	fmt.Printf("current: v%d%s\n", project.CurrentVersion, name)
	fmt.Printf("latest:  v%d\n", project.LatestVersion)
	return nil
}

func Set(versionStr string) error {
	projectName, _, _, _, _, err := cmd_loader.GetBasicInfo(cmd_loader.GetBasicInfoOptions{ReadOnly: true})
	if err != nil {
		return err
	}

	if projectName == "" {
		fmt.Println("no project under current directory, use swapenv load to initiate.")
		return nil
	}

	version, err := filehandler.ResolveVersion(projectName, versionStr)
	if err != nil {
		return err
	}

	if err := filehandler.SetCurrentVersion(projectName, version); err != nil {
		return err
	}

	fmt.Printf("switched to v%d\n", version)
	return nil
}

func List() error {
	projectName, _, localDirectory, _, _, err := cmd_loader.GetBasicInfo(cmd_loader.GetBasicInfoOptions{ReadOnly: true})
	if err != nil {
		return err
	}

	if projectName == "" {
		fmt.Println("no project under current directory, use swapenv load to initiate.")
		return nil
	}

	project, err := filehandler.FindProjectByLocalPath(localDirectory)
	if err != nil {
		return err
	}

	versions, err := filehandler.ListVersions(projectName)
	if err != nil {
		return err
	}

	for _, v := range versions {
		marker := "  "
		if v == project.CurrentVersion {
			marker = "* "
		}

		name := ""
		if n, ok := project.VersionNames[strconv.Itoa(v)]; ok {
			name = fmt.Sprintf(" (%s)", n)
		}

		latest := ""
		if v == project.LatestVersion {
			latest = " [latest]"
		}

		fmt.Printf("%s%d%s%s\n", marker, v, name, latest)
	}

	return nil
}

func Rename(versionStr, name string) error {
	projectName, _, _, _, _, err := cmd_loader.GetBasicInfo(cmd_loader.GetBasicInfoOptions{ReadOnly: true})
	if err != nil {
		return err
	}

	if projectName == "" {
		fmt.Println("no project under current directory, use swapenv load to initiate.")
		return nil
	}

	if name == "latest" {
		return fmt.Errorf("cannot use 'latest' as version name")
	}

	version, err := filehandler.ResolveVersion(projectName, versionStr)
	if err != nil {
		return err
	}

	if err := filehandler.RenameVersion(projectName, version, name); err != nil {
		return err
	}

	fmt.Printf("v%d named '%s'\n", version, name)
	return nil
}

func Rollback(steps int) error {
	projectName, _, localDirectory, _, _, err := cmd_loader.GetBasicInfo(cmd_loader.GetBasicInfoOptions{ReadOnly: true})
	if err != nil {
		return err
	}

	if projectName == "" {
		fmt.Println("no project under current directory, use swapenv load to initiate.")
		return nil
	}

	project, err := filehandler.FindProjectByLocalPath(localDirectory)
	if err != nil {
		return err
	}

	versions, err := filehandler.ListVersions(projectName)
	if err != nil {
		return err
	}

	currentIdx := -1
	for i, v := range versions {
		if v == project.CurrentVersion {
			currentIdx = i
			break
		}
	}

	if currentIdx == -1 {
		return fmt.Errorf("current version not found")
	}

	newIdx := currentIdx - steps
	if newIdx < 0 {
		newIdx = 0
	}

	newVersion := versions[newIdx]
	if err := filehandler.SetCurrentVersion(projectName, newVersion); err != nil {
		return err
	}

	fmt.Printf("rolled back to v%d\n", newVersion)
	return nil
}
