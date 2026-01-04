package test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/reduan2660/swapenv/cmd"
	"github.com/reduan2660/swapenv/internal/filehandler"
	"github.com/reduan2660/swapenv/internal/types"
	"github.com/spf13/cobra"
)

// TestMapProjectNotFound tests map fails when project doesn't exist
func TestMapProjectNotFound(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	// Try to map a nonexistent project
	err := mapProject("nonexistent-project")
	if err == nil {
		t.Error("map should fail when project doesn't exist")
	}
}

// TestMapReceivedProject tests mapping a received project to current directory
// Note: This tests that after mapping, the directory is findable by local path
func TestMapReceivedProject(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	projectName := "received-project"

	// Create a project without local path (simulating received project)
	newProject := types.ProjectDir{
		ProjectName:    projectName,
		LocalPath:      "", // No local path initially
		CurrentEnv:     "",
		CurrentVersion: 1,
		LatestVersion:  1,
		VersionNames:   make(map[string]string),
	}
	if err := filehandler.UpsertProjectDir(newProject); err != nil {
		t.Fatal(err)
	}

	// Create version file
	homeDir, err := filehandler.GetHomeDirectory(projectName)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(homeDir, 0755); err != nil {
		t.Fatal(err)
	}

	now := time.Now().UTC().Unix()
	projectData := types.Project{
		Id:             uuid.New().String(),
		Name:           projectName,
		Owner:          "received",
		LocalDirectory: "",
		CreatedAt:      now,
		ModifiedAt:     now,
		Envs: map[string][]types.EnvValue{
			"dev": {{Key: "ENV_1", Val: "value"}},
		},
	}
	projectJSON, _ := projectData.MarshalJSON()
	versionPath := filepath.Join(homeDir, "v1.json")
	if err := os.WriteFile(versionPath, projectJSON, 0644); err != nil {
		t.Fatal(err)
	}

	// Verify project exists before mapping
	projectBefore, err := filehandler.FindProjectByName(projectName)
	if err != nil {
		t.Fatalf("failed to find project before mapping: %v", err)
	}
	if projectBefore == nil {
		t.Fatal("project should exist before mapping")
	}

	// Map the project to current directory
	if err := mapProject(projectName); err != nil {
		t.Fatalf("map failed: %v", err)
	}

	// Verify project can now be found by local path
	// (UpsertProjectDir creates new entry when LocalPath changes)
	projectByPath, err := filehandler.FindProjectByLocalPath(testProjectDir)
	if err != nil {
		t.Fatal(err)
	}
	if projectByPath == nil {
		t.Fatal("project should be findable by local path after mapping")
	}
	if projectByPath.ProjectName != projectName {
		t.Errorf("expected project name %s, got %s", projectName, projectByPath.ProjectName)
	}
}

// TestMapToAnotherProject tests remapping directory from one project to another
func TestMapToAnotherProject(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	// Create first project mapped to current dir
	project1 := types.ProjectDir{
		ProjectName:    "project1",
		LocalPath:      testProjectDir,
		CurrentEnv:     "",
		CurrentVersion: 1,
		LatestVersion:  1,
		VersionNames:   make(map[string]string),
	}
	if err := filehandler.UpsertProjectDir(project1); err != nil {
		t.Fatal(err)
	}
	setupVersionFile(t, "project1")

	// Verify project1 is mapped to current directory
	found, err := filehandler.FindProjectByLocalPath(testProjectDir)
	if err != nil {
		t.Fatal(err)
	}
	if found == nil {
		t.Fatal("project1 should be mapped to testProjectDir")
	}
	if found.ProjectName != "project1" {
		t.Errorf("expected project1, got %s", found.ProjectName)
	}

	// Now map the same directory to a different project
	project2 := types.ProjectDir{
		ProjectName:    "project2",
		LocalPath:      testProjectDir, // Same path, different project
		CurrentEnv:     "",
		CurrentVersion: 1,
		LatestVersion:  1,
		VersionNames:   make(map[string]string),
	}
	if err := filehandler.UpsertProjectDir(project2); err != nil {
		t.Fatal(err)
	}
	setupVersionFile(t, "project2")

	// Verify directory now maps to project2 (UpsertProjectDir updates by LocalPath)
	foundAfter, err := filehandler.FindProjectByLocalPath(testProjectDir)
	if err != nil {
		t.Fatal(err)
	}
	if foundAfter == nil {
		t.Fatal("should find project by local path")
	}
	if foundAfter.ProjectName != "project2" {
		t.Errorf("expected project2 after remap, got %s", foundAfter.ProjectName)
	}
}

// TestMapThenTo tests the flow: receive -> map -> to
func TestMapThenTo(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	projectName := "received-for-to"

	// Create a received project
	newProject := types.ProjectDir{
		ProjectName:    projectName,
		LocalPath:      "",
		CurrentEnv:     "",
		CurrentVersion: 1,
		LatestVersion:  1,
		VersionNames:   make(map[string]string),
	}
	if err := filehandler.UpsertProjectDir(newProject); err != nil {
		t.Fatal(err)
	}

	// Create version file with envs
	homeDir, err := filehandler.GetHomeDirectory(projectName)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(homeDir, 0755); err != nil {
		t.Fatal(err)
	}

	now := time.Now().UTC().Unix()
	projectData := types.Project{
		Id:             uuid.New().String(),
		Name:           projectName,
		Owner:          "received",
		LocalDirectory: "",
		CreatedAt:      now,
		ModifiedAt:     now,
		Envs: map[string][]types.EnvValue{
			"dev": {
				{Key: "APP_ENV", Val: "development"},
				{Key: "DEBUG", Val: "true"},
			},
		},
	}
	projectJSON, _ := projectData.MarshalJSON()
	versionPath := filepath.Join(homeDir, "v1.json")
	if err := os.WriteFile(versionPath, projectJSON, 0644); err != nil {
		t.Fatal(err)
	}

	// Map project to current directory
	if err := mapProject(projectName); err != nil {
		t.Fatalf("map failed: %v", err)
	}

	// Now "to dev" should work
	toCmd := getToCmd()
	toCmd.Flags().Set("replace", "true")
	if err := toCmd.RunE(toCmd, []string{"dev"}); err != nil {
		t.Fatalf("to dev failed after map: %v", err)
	}

	// Verify .env was created
	content, err := os.ReadFile(".env")
	if err != nil {
		t.Fatal(err)
	}
	contentStr := string(content)

	if !contains(contentStr, "APP_ENV=development") {
		t.Error("APP_ENV=development should be in .env")
	}
	if !contains(contentStr, "DEBUG=true") {
		t.Error("DEBUG=true should be in .env")
	}
}

// TestFindProjectByLocalPath tests finding project by directory
func TestFindProjectByLocalPath(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	// Create a mapped project
	project := types.ProjectDir{
		ProjectName:    "mapped-project",
		LocalPath:      testProjectDir,
		CurrentEnv:     "",
		CurrentVersion: 1,
		LatestVersion:  1,
		VersionNames:   make(map[string]string),
	}
	if err := filehandler.UpsertProjectDir(project); err != nil {
		t.Fatal(err)
	}

	// Find by local path
	found, err := filehandler.FindProjectByLocalPath(testProjectDir)
	if err != nil {
		t.Fatal(err)
	}
	if found == nil {
		t.Fatal("should find project by local path")
	}
	if found.ProjectName != "mapped-project" {
		t.Errorf("expected project name 'mapped-project', got '%s'", found.ProjectName)
	}
}

// TestFindProjectByLocalPathNotFound tests no project at path
func TestFindProjectByLocalPathNotFound(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	// No project exists
	found, err := filehandler.FindProjectByLocalPath("/nonexistent/path")
	if err != nil {
		t.Fatal(err)
	}
	if found != nil {
		t.Error("should not find project for nonexistent path")
	}
}

// Helper to simulate map command without stdin prompts
func mapProject(projectName string) error {
	project, err := filehandler.FindProjectByName(projectName)
	if err != nil {
		return err
	}
	if project == nil {
		return os.ErrNotExist
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	project.LocalPath = cwd
	return filehandler.UpsertProjectDir(*project)
}

// Helper to create version file for a project
func setupVersionFile(t *testing.T, projectName string) {
	t.Helper()

	homeDir, err := filehandler.GetHomeDirectory(projectName)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(homeDir, 0755); err != nil {
		t.Fatal(err)
	}

	now := time.Now().UTC().Unix()
	projectData := types.Project{
		Id:             uuid.New().String(),
		Name:           projectName,
		Owner:          "test",
		LocalDirectory: "",
		CreatedAt:      now,
		ModifiedAt:     now,
		Envs: map[string][]types.EnvValue{
			"dev": {{Key: "ENV_1", Val: "value"}},
		},
	}
	projectJSON, _ := projectData.MarshalJSON()
	versionPath := filepath.Join(homeDir, "v1.json")
	if err := os.WriteFile(versionPath, projectJSON, 0644); err != nil {
		t.Fatal(err)
	}
}

// Helper to get to command (avoiding import cycle with cmd package)
func getToCmd() *cobra.Command {
	return cmd.GetToCmd()
}
