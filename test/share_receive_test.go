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
)

// TestShareProjectNotFound tests share fails when project doesn't exist
func TestShareProjectNotFound(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	shareCmd := cmd.GetShareCmd()
	shareCmd.Flags().Set("project", "nonexistent-project")

	err := shareCmd.RunE(shareCmd, []string{})
	if err == nil {
		t.Error("share should fail when project doesn't exist")
	}
}

// TestShareEnvNotFound tests share fails when specified env doesn't exist
func TestShareEnvNotFound(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	// Create a project with dev env
	createEnvFile(t, ".dev.env", `ENV_1=dev`)
	loadCmd := cmd.GetLoadCmd()
	loadCmd.Flags().Set("env", "*")
	loadCmd.Flags().Set("replace", "false")
	if err := loadCmd.RunE(loadCmd, []string{}); err != nil {
		t.Fatal(err)
	}

	shareCmd := cmd.GetShareCmd()
	shareCmd.Flags().Set("env", "nonexistent-env")

	err := shareCmd.RunE(shareCmd, []string{})
	if err == nil {
		t.Error("share should fail when specified env doesn't exist")
	}
}

// TestShareNoEnvsInProject tests share fails when project has no environments
func TestShareNoEnvsInProject(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	// Create project directory structure without any envs
	projectDir := types.ProjectDir{
		ProjectName:    "empty-project",
		LocalPath:      testProjectDir,
		CurrentEnv:     "",
		CurrentVersion: 1,
		LatestVersion:  1,
		VersionNames:   make(map[string]string),
	}
	if err := filehandler.UpsertProjectDir(projectDir); err != nil {
		t.Fatal(err)
	}

	// Create version file with empty envs
	homeDir, err := filehandler.GetHomeDirectory("empty-project")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(homeDir, 0755); err != nil {
		t.Fatal(err)
	}

	now := time.Now().UTC().Unix()
	project := types.Project{
		Id:             uuid.New().String(),
		Name:           "empty-project",
		Owner:          "test",
		LocalDirectory: testProjectDir,
		CreatedAt:      now,
		ModifiedAt:     now,
		Envs:           make(map[string][]types.EnvValue),
	}
	projectJSON, _ := project.MarshalJSON()
	versionPath := filepath.Join(homeDir, "v1.json")
	if err := os.WriteFile(versionPath, projectJSON, 0644); err != nil {
		t.Fatal(err)
	}

	shareCmd := cmd.GetShareCmd()
	shareCmd.Flags().Set("project", "empty-project")

	err = shareCmd.RunE(shareCmd, []string{})
	if err == nil {
		t.Error("share should fail when project has no environments")
	}
}

// TestShareVersionNotFound tests share fails when specified version doesn't exist
func TestShareVersionNotFound(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	// Create a project
	createEnvFile(t, ".dev.env", `ENV_1=dev`)
	loadCmd := cmd.GetLoadCmd()
	loadCmd.Flags().Set("env", "*")
	loadCmd.Flags().Set("replace", "false")
	if err := loadCmd.RunE(loadCmd, []string{}); err != nil {
		t.Fatal(err)
	}

	shareCmd := cmd.GetShareCmd()
	shareCmd.Flags().Set("version", "999")

	err := shareCmd.RunE(shareCmd, []string{})
	if err == nil {
		t.Error("share should fail when specified version doesn't exist")
	}
}

// TestShareNoProjectInCurrentDir tests share fails when no project flag and not in a project dir
func TestShareNoProjectInCurrentDir(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	// Don't create any project, just try to share
	shareCmd := cmd.GetShareCmd()

	err := shareCmd.RunE(shareCmd, []string{})
	if err == nil {
		t.Error("share should fail when not in a project directory and no --project flag")
	}
}

// TestSaveReceivedNewProject tests receiving creates new project correctly
func TestSaveReceivedNewProject(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	// Simulate receiving a new project
	envMap := map[string][]types.EnvValue{
		"dev": {
			{Key: "DB_HOST", Val: "localhost"},
			{Key: "DB_PORT", Val: "5432"},
		},
		"prod": {
			{Key: "DB_HOST", Val: "prod.db.com"},
			{Key: "DB_PORT", Val: "5432"},
		},
	}

	// Call the internal save function by simulating what receive does
	projectName := "received-project"

	// Create project (simulating saveReceived logic)
	newProject := types.ProjectDir{
		ProjectName:    projectName,
		LocalPath:      "", // No local path for received projects
		CurrentEnv:     "",
		CurrentVersion: 1,
		LatestVersion:  1,
		VersionNames:   make(map[string]string),
	}
	if err := filehandler.UpsertProjectDir(newProject); err != nil {
		t.Fatalf("failed to create project: %v", err)
	}

	versionPath, err := filehandler.GetVersionFilePath(projectName, 1)
	if err != nil {
		t.Fatal(err)
	}

	homeDir, err := filehandler.GetHomeDirectory(projectName)
	if err != nil {
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
		Envs:           envMap,
	}

	projectJSON, err := projectData.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	if err := filehandler.WriteProject(homeDir, versionPath, projectJSON); err != nil {
		t.Fatalf("failed to write project: %v", err)
	}

	// Verify project was created
	project, err := filehandler.FindProjectByName(projectName)
	if err != nil {
		t.Fatal(err)
	}
	if project == nil {
		t.Fatal("received project should exist")
	}
	if project.LocalPath != "" {
		t.Error("received project should have empty LocalPath")
	}
	if project.CurrentVersion != 1 {
		t.Errorf("expected version 1, got %d", project.CurrentVersion)
	}

	// Verify envs were saved
	envNames, err := filehandler.ListProjectEnv(versionPath)
	if err != nil {
		t.Fatal(err)
	}
	if len(envNames) != 2 {
		t.Errorf("expected 2 envs, got %d", len(envNames))
	}
}

// TestSaveReceivedExistingProject tests receiving bumps version on existing project
func TestSaveReceivedExistingProject(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	// Create initial project with v1
	createEnvFile(t, ".dev.env", `ENV_1=v1`)
	loadCmd := cmd.GetLoadCmd()
	loadCmd.Flags().Set("env", "*")
	loadCmd.Flags().Set("replace", "false")
	if err := loadCmd.RunE(loadCmd, []string{}); err != nil {
		t.Fatal(err)
	}

	projectName := "test-project"

	// Verify v1 exists
	project, err := filehandler.FindProjectByName(projectName)
	if err != nil {
		t.Fatal(err)
	}
	if project.CurrentVersion != 1 {
		t.Fatalf("expected initial version 1, got %d", project.CurrentVersion)
	}

	// Simulate receiving new version (bumping version)
	newVersion, err := filehandler.BumpVersion(projectName)
	if err != nil {
		t.Fatalf("failed to bump version: %v", err)
	}
	if newVersion != 2 {
		t.Errorf("expected bumped version 2, got %d", newVersion)
	}

	// Save received envs to new version
	envMap := map[string][]types.EnvValue{
		"dev": {
			{Key: "ENV_1", Val: "v2_received"},
		},
	}

	versionPath, err := filehandler.GetVersionFilePath(projectName, newVersion)
	if err != nil {
		t.Fatal(err)
	}

	homeDir, err := filehandler.GetHomeDirectory(projectName)
	if err != nil {
		t.Fatal(err)
	}

	now := time.Now().UTC().Unix()
	projectData := types.Project{
		Id:             uuid.New().String(),
		Name:           projectName,
		Owner:          "received",
		LocalDirectory: project.LocalPath,
		CreatedAt:      now,
		ModifiedAt:     now,
		Envs:           envMap,
	}

	projectJSON, _ := projectData.MarshalJSON()
	if err := filehandler.WriteProject(homeDir, versionPath, projectJSON); err != nil {
		t.Fatal(err)
	}

	// Verify version was bumped
	project, err = filehandler.FindProjectByName(projectName)
	if err != nil {
		t.Fatal(err)
	}
	if project.CurrentVersion != 2 {
		t.Errorf("expected current version 2, got %d", project.CurrentVersion)
	}
	if project.LatestVersion != 2 {
		t.Errorf("expected latest version 2, got %d", project.LatestVersion)
	}

	// Verify both versions exist
	versions, err := filehandler.ListVersions(projectName)
	if err != nil {
		t.Fatal(err)
	}
	if len(versions) != 2 {
		t.Errorf("expected 2 versions, got %d", len(versions))
	}
}
