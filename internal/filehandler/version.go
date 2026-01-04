package filehandler

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strconv"

	"github.com/spf13/viper"
)

func MigrateProjectIfNeeded(projectName string) error {
	dir, err := FindProjectByName(projectName)
	if err != nil || dir == nil {
		return err
	}

	if dir.LatestVersion > 0 {
		return nil
	}

	homeDir, err := GetHomeDirectory(projectName)
	if err != nil {
		return err
	}

	oldPath := filepath.Join(homeDir, "project.json")
	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		return nil
	}

	newPath := filepath.Join(homeDir, "v1.json")
	if err := os.Rename(oldPath, newPath); err != nil {
		return err
	}

	return UpdateProjectVersion(projectName, 1, 1)
}

func UpdateProjectVersion(projectName string, currentVersion, latestVersion int) error {
	dirs, err := ReadProjectDirs()
	if err != nil {
		return err
	}

	for i, dir := range dirs {
		if dir.ProjectName == projectName {
			dirs[i].CurrentVersion = currentVersion
			dirs[i].LatestVersion = latestVersion
			if dirs[i].VersionNames == nil {
				dirs[i].VersionNames = make(map[string]string)
			}
			return WriteProjectDirs(dirs)
		}
	}

	return fmt.Errorf("project not found: %s", projectName)
}

func GetVersionFilePath(projectName string, version int) (string, error) {
	homeDir, err := GetHomeDirectory(projectName)
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, fmt.Sprintf("v%d.json", version)), nil

}

func ListVersions(projectName string) ([]int, error) {
	homeDir, err := GetHomeDirectory(projectName)
	if err != nil {
		return nil, err
	}

	files, err := filepath.Glob(filepath.Join(homeDir, "v*.json"))
	if err != nil {
		return nil, err
	}

	versions := make([]int, 0, len(files))
	for _, f := range files {
		name := filepath.Base(f)
		var v int
		if _, err := fmt.Sscanf(name, "v%d.json", &v); err == nil {
			versions = append(versions, v)
		}
	}

	sort.Ints(versions)
	return versions, nil
}

func ResolveVersion(projectName, versionStr string) (int, error) {
	dir, err := FindProjectByName(projectName)
	if err != nil {
		return 0, err
	}
	if dir == nil {
		return 0, fmt.Errorf("project not found: %s", projectName)
	}

	// Empty string = current version
	if versionStr == "" {
		return dir.CurrentVersion, nil
	}

	// "latest" = latest version
	if versionStr == "latest" {
		if dir.LatestVersion > 0 {
			return dir.LatestVersion, nil
		}
		// Fallback to current version
		if dir.CurrentVersion > 0 {
			return dir.CurrentVersion, nil
		}
		return 0, fmt.Errorf("no version found for project: %s", projectName)
	}

	// Try as number first
	if v, err := strconv.Atoi(versionStr); err == nil {
		versions, err := ListVersions(projectName)
		if err != nil {
			return 0, err
		}
		if slices.Contains(versions, v) {
			return v, nil
		}
		return 0, fmt.Errorf("version %d not found", v)
	}

	// Try as name
	for v, name := range dir.VersionNames {
		if name == versionStr {
			ver, _ := strconv.Atoi(v)
			return ver, nil
		}
	}

	return 0, fmt.Errorf("version '%s' not found", versionStr)
}

func BumpVersion(projectName string) (int, error) {
	if err := MigrateProjectIfNeeded(projectName); err != nil {
		return 0, err
	}

	dirs, err := ReadProjectDirs()
	if err != nil {
		return 0, err
	}

	for i, dir := range dirs {
		if dir.ProjectName == projectName {
			if dirs[i].VersionNames == nil {
				dirs[i].VersionNames = make(map[string]string)
			}

			dirs[i].LatestVersion++
			dirs[i].CurrentVersion = dirs[i].LatestVersion

			if err := WriteProjectDirs(dirs); err != nil {
				return 0, err
			}
			return dirs[i].LatestVersion, nil
		}
	}

	return 0, fmt.Errorf("project not found: %s", projectName)
}

func PruneVersions(projectName string) error {
	dir, err := FindProjectByName(projectName)
	if err != nil || dir == nil {
		return err
	}

	versions, err := ListVersions(projectName)
	if err != nil {
		return err
	}

	maxVersions := viper.GetInt("max_versions")
	if len(versions) <= maxVersions {
		return nil
	}

	toDelete := versions[:len(versions)-maxVersions]

	for _, v := range toDelete {
		// Skip named (protected) versions
		if _, isNamed := dir.VersionNames[strconv.Itoa(v)]; isNamed {
			continue
		}
		// Skip current version
		if v == dir.CurrentVersion {
			continue
		}

		path, _ := GetVersionFilePath(projectName, v)
		os.Remove(path)
	}

	return nil
}

func SetCurrentVersion(projectName string, version int) error {
	dirs, err := ReadProjectDirs()
	if err != nil {
		return err
	}

	for i, dir := range dirs {
		if dir.ProjectName == projectName {
			dirs[i].CurrentVersion = version
			return WriteProjectDirs(dirs)
		}
	}

	return fmt.Errorf("project not found: %s", projectName)
}

func RenameVersion(projectName string, version int, name string) error {
	dirs, err := ReadProjectDirs()
	if err != nil {
		return err
	}

	for i, dir := range dirs {
		if dir.ProjectName == projectName {
			if dirs[i].VersionNames == nil {
				dirs[i].VersionNames = make(map[string]string)
			}
			dirs[i].VersionNames[strconv.Itoa(version)] = name
			return WriteProjectDirs(dirs)
		}
	}

	return fmt.Errorf("project not found: %s", projectName)
}
