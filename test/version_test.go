package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/reduan2660/swapenv/cmd"
	"github.com/reduan2660/swapenv/internal/filehandler"
)

func TestVersionBump(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	// Load creates v1
	createEnvFile(t, ".dev.env", `ENV_1=dev`)
	loadCmd := cmd.GetLoadCmd()
	loadCmd.Flags().Set("env", "*")
	loadCmd.Flags().Set("replace", "false")
	if err := loadCmd.RunE(loadCmd, []string{}); err != nil {
		t.Fatal(err)
	}

	// Check v1 exists
	v1Path := filepath.Join(testHomeDir, "test-project", "v1.json")
	if _, err := os.Stat(v1Path); os.IsNotExist(err) {
		t.Error("v1.json should exist after first load")
	}

	// Load again creates v2
	createEnvFile(t, ".dev.env", `ENV_1=dev_v2`)
	if err := loadCmd.RunE(loadCmd, []string{}); err != nil {
		t.Fatal(err)
	}

	v2Path := filepath.Join(testHomeDir, "test-project", "v2.json")
	if _, err := os.Stat(v2Path); os.IsNotExist(err) {
		t.Error("v2.json should exist after second load")
	}

	// Load again creates v3
	createEnvFile(t, ".dev.env", `ENV_1=dev_v3`)
	if err := loadCmd.RunE(loadCmd, []string{}); err != nil {
		t.Fatal(err)
	}

	v3Path := filepath.Join(testHomeDir, "test-project", "v3.json")
	if _, err := os.Stat(v3Path); os.IsNotExist(err) {
		t.Error("v3.json should exist after third load")
	}
}

func TestVersionList(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	// Create multiple versions
	createEnvFile(t, ".dev.env", `ENV_1=v1`)
	loadCmd := cmd.GetLoadCmd()
	loadCmd.Flags().Set("env", "*")
	loadCmd.Flags().Set("replace", "false")
	if err := loadCmd.RunE(loadCmd, []string{}); err != nil {
		t.Fatal(err)
	}

	createEnvFile(t, ".dev.env", `ENV_1=v2`)
	if err := loadCmd.RunE(loadCmd, []string{}); err != nil {
		t.Fatal(err)
	}

	createEnvFile(t, ".dev.env", `ENV_1=v3`)
	if err := loadCmd.RunE(loadCmd, []string{}); err != nil {
		t.Fatal(err)
	}

	// List versions using filehandler
	versions, err := filehandler.ListVersions("test-project")
	if err != nil {
		t.Fatalf("ListVersions failed: %v", err)
	}

	if len(versions) != 3 {
		t.Errorf("expected 3 versions, got %d", len(versions))
	}

	if versions[0] != 1 || versions[1] != 2 || versions[2] != 3 {
		t.Errorf("expected versions [1, 2, 3], got %v", versions)
	}

	// Test version ls command
	versionLsCmd := cmd.GetVersionLsCmd()
	if err := versionLsCmd.RunE(versionLsCmd, []string{}); err != nil {
		t.Fatalf("version ls failed: %v", err)
	}
}

func TestVersionSet(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	// Create multiple versions
	createEnvFile(t, ".dev.env", `ENV_1=v1`)
	loadCmd := cmd.GetLoadCmd()
	loadCmd.Flags().Set("env", "*")
	loadCmd.Flags().Set("replace", "false")
	if err := loadCmd.RunE(loadCmd, []string{}); err != nil {
		t.Fatal(err)
	}

	createEnvFile(t, ".dev.env", `ENV_1=v2`)
	if err := loadCmd.RunE(loadCmd, []string{}); err != nil {
		t.Fatal(err)
	}

	// Current version should be 2 (latest)
	project, err := filehandler.FindProjectByLocalPath(testProjectDir)
	if err != nil {
		t.Fatal(err)
	}
	if project.CurrentVersion != 2 {
		t.Errorf("expected current version 2, got %d", project.CurrentVersion)
	}

	// Switch to version 1
	versionCmd := cmd.GetVersionCmd()
	if err := versionCmd.RunE(versionCmd, []string{"1"}); err != nil {
		t.Fatalf("version 1 failed: %v", err)
	}

	// Verify current version changed
	project, err = filehandler.FindProjectByLocalPath(testProjectDir)
	if err != nil {
		t.Fatal(err)
	}
	if project.CurrentVersion != 1 {
		t.Errorf("expected current version 1, got %d", project.CurrentVersion)
	}
}

func TestVersionSetLatest(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	// Create multiple versions
	createEnvFile(t, ".dev.env", `ENV_1=v1`)
	loadCmd := cmd.GetLoadCmd()
	loadCmd.Flags().Set("env", "*")
	loadCmd.Flags().Set("replace", "false")
	if err := loadCmd.RunE(loadCmd, []string{}); err != nil {
		t.Fatal(err)
	}

	createEnvFile(t, ".dev.env", `ENV_1=v2`)
	if err := loadCmd.RunE(loadCmd, []string{}); err != nil {
		t.Fatal(err)
	}

	createEnvFile(t, ".dev.env", `ENV_1=v3`)
	if err := loadCmd.RunE(loadCmd, []string{}); err != nil {
		t.Fatal(err)
	}

	// Switch to v1
	if err := filehandler.SetCurrentVersion("test-project", 1); err != nil {
		t.Fatal(err)
	}

	// Switch to "latest"
	versionCmd := cmd.GetVersionCmd()
	if err := versionCmd.RunE(versionCmd, []string{"latest"}); err != nil {
		t.Fatalf("version latest failed: %v", err)
	}

	// Verify current version is now 3 (latest)
	project, err := filehandler.FindProjectByLocalPath(testProjectDir)
	if err != nil {
		t.Fatal(err)
	}
	if project.CurrentVersion != 3 {
		t.Errorf("expected current version 3, got %d", project.CurrentVersion)
	}
}

func TestVersionRename(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	// Create a version
	createEnvFile(t, ".dev.env", `ENV_1=v1`)
	loadCmd := cmd.GetLoadCmd()
	loadCmd.Flags().Set("env", "*")
	loadCmd.Flags().Set("replace", "false")
	if err := loadCmd.RunE(loadCmd, []string{}); err != nil {
		t.Fatal(err)
	}

	// Rename version 1
	renameCmd := cmd.GetVersionRenameCmd()
	if err := renameCmd.RunE(renameCmd, []string{"1", "initial"}); err != nil {
		t.Fatalf("version rename failed: %v", err)
	}

	// Verify name was set
	project, err := filehandler.FindProjectByLocalPath(testProjectDir)
	if err != nil {
		t.Fatal(err)
	}
	if name, ok := project.VersionNames["1"]; !ok || name != "initial" {
		t.Errorf("expected version 1 to be named 'initial', got %v", project.VersionNames)
	}

	// Test resolving by name
	version, err := filehandler.ResolveVersion("test-project", "initial")
	if err != nil {
		t.Fatalf("ResolveVersion by name failed: %v", err)
	}
	if version != 1 {
		t.Errorf("expected version 1 when resolving 'initial', got %d", version)
	}
}

func TestVersionRenameCannotUseLatest(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	// Create a version
	createEnvFile(t, ".dev.env", `ENV_1=v1`)
	loadCmd := cmd.GetLoadCmd()
	loadCmd.Flags().Set("env", "*")
	loadCmd.Flags().Set("replace", "false")
	if err := loadCmd.RunE(loadCmd, []string{}); err != nil {
		t.Fatal(err)
	}

	// Try to name version "latest" - should fail
	renameCmd := cmd.GetVersionRenameCmd()
	err := renameCmd.RunE(renameCmd, []string{"1", "latest"})
	if err == nil {
		t.Fatal("renaming to 'latest' should return an error")
	}
}

func TestVersionRollback(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	// Create multiple versions
	createEnvFile(t, ".dev.env", `ENV_1=v1`)
	loadCmd := cmd.GetLoadCmd()
	loadCmd.Flags().Set("env", "*")
	loadCmd.Flags().Set("replace", "false")
	if err := loadCmd.RunE(loadCmd, []string{}); err != nil {
		t.Fatal(err)
	}

	createEnvFile(t, ".dev.env", `ENV_1=v2`)
	if err := loadCmd.RunE(loadCmd, []string{}); err != nil {
		t.Fatal(err)
	}

	createEnvFile(t, ".dev.env", `ENV_1=v3`)
	if err := loadCmd.RunE(loadCmd, []string{}); err != nil {
		t.Fatal(err)
	}

	// Current version should be 3
	project, err := filehandler.FindProjectByLocalPath(testProjectDir)
	if err != nil {
		t.Fatal(err)
	}
	if project.CurrentVersion != 3 {
		t.Errorf("expected current version 3, got %d", project.CurrentVersion)
	}

	// Rollback 1 step
	rollbackCmd := cmd.GetVersionRollbackCmd()
	if err := rollbackCmd.RunE(rollbackCmd, []string{}); err != nil {
		t.Fatalf("rollback failed: %v", err)
	}

	project, err = filehandler.FindProjectByLocalPath(testProjectDir)
	if err != nil {
		t.Fatal(err)
	}
	if project.CurrentVersion != 2 {
		t.Errorf("expected current version 2 after rollback, got %d", project.CurrentVersion)
	}

	// Rollback 2 steps (should hit v1, not go below)
	if err := rollbackCmd.RunE(rollbackCmd, []string{"2"}); err != nil {
		t.Fatalf("rollback 2 failed: %v", err)
	}

	project, err = filehandler.FindProjectByLocalPath(testProjectDir)
	if err != nil {
		t.Fatal(err)
	}
	if project.CurrentVersion != 1 {
		t.Errorf("expected current version 1 after rollback 2, got %d", project.CurrentVersion)
	}
}

func TestVersionPruning(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	// Create 7 versions (max_versions is 5)
	loadCmd := cmd.GetLoadCmd()
	loadCmd.Flags().Set("env", "*")
	loadCmd.Flags().Set("replace", "false")

	for i := 1; i <= 7; i++ {
		createEnvFile(t, ".dev.env", `ENV_1=dev`)
		if err := loadCmd.RunE(loadCmd, []string{}); err != nil {
			t.Fatalf("load %d failed: %v", i, err)
		}
	}

	// Should have 5 versions (v3-v7), v1 and v2 should be pruned
	versions, err := filehandler.ListVersions("test-project")
	if err != nil {
		t.Fatal(err)
	}

	if len(versions) != 5 {
		t.Errorf("expected 5 versions after pruning, got %d: %v", len(versions), versions)
	}

	// v1 and v2 should not exist
	v1Path := filepath.Join(testHomeDir, "test-project", "v1.json")
	if _, err := os.Stat(v1Path); !os.IsNotExist(err) {
		t.Error("v1.json should be pruned")
	}

	v2Path := filepath.Join(testHomeDir, "test-project", "v2.json")
	if _, err := os.Stat(v2Path); !os.IsNotExist(err) {
		t.Error("v2.json should be pruned")
	}
}

func TestVersionPruningProtectsNamedVersions(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	loadCmd := cmd.GetLoadCmd()
	loadCmd.Flags().Set("env", "*")
	loadCmd.Flags().Set("replace", "false")

	// Create v1
	createEnvFile(t, ".dev.env", `ENV_1=dev`)
	if err := loadCmd.RunE(loadCmd, []string{}); err != nil {
		t.Fatal(err)
	}

	// Name v1 to protect it
	if err := filehandler.RenameVersion("test-project", 1, "protected"); err != nil {
		t.Fatal(err)
	}

	// Create 6 more versions (v2-v7)
	for i := 2; i <= 7; i++ {
		createEnvFile(t, ".dev.env", `ENV_1=dev`)
		if err := loadCmd.RunE(loadCmd, []string{}); err != nil {
			t.Fatalf("load %d failed: %v", i, err)
		}
	}

	// v1 should still exist (protected by name)
	v1Path := filepath.Join(testHomeDir, "test-project", "v1.json")
	if _, err := os.Stat(v1Path); os.IsNotExist(err) {
		t.Error("v1.json should be protected from pruning due to name")
	}

	// Can resolve by name
	version, err := filehandler.ResolveVersion("test-project", "protected")
	if err != nil {
		t.Fatalf("ResolveVersion by name failed: %v", err)
	}
	if version != 1 {
		t.Errorf("expected version 1 for 'protected', got %d", version)
	}
}

func TestVersionShow(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	// Create a version
	createEnvFile(t, ".dev.env", `ENV_1=dev`)
	loadCmd := cmd.GetLoadCmd()
	loadCmd.Flags().Set("env", "*")
	loadCmd.Flags().Set("replace", "false")
	if err := loadCmd.RunE(loadCmd, []string{}); err != nil {
		t.Fatal(err)
	}

	// Show version (no args)
	versionCmd := cmd.GetVersionCmd()
	if err := versionCmd.RunE(versionCmd, []string{}); err != nil {
		t.Fatalf("version show failed: %v", err)
	}
}

func TestVersionResolveInvalidVersion(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	// Create a version
	createEnvFile(t, ".dev.env", `ENV_1=dev`)
	loadCmd := cmd.GetLoadCmd()
	loadCmd.Flags().Set("env", "*")
	loadCmd.Flags().Set("replace", "false")
	if err := loadCmd.RunE(loadCmd, []string{}); err != nil {
		t.Fatal(err)
	}

	// Try to resolve non-existent version number
	_, err := filehandler.ResolveVersion("test-project", "999")
	if err == nil {
		t.Error("resolving non-existent version should return error")
	}

	// Try to resolve non-existent version name
	_, err = filehandler.ResolveVersion("test-project", "nonexistent")
	if err == nil {
		t.Error("resolving non-existent version name should return error")
	}
}

func TestLsWithVersionFlag(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	// Create multiple versions
	createEnvFile(t, ".dev.env", `ENV_1=v1`)
	loadCmd := cmd.GetLoadCmd()
	loadCmd.Flags().Set("env", "*")
	loadCmd.Flags().Set("replace", "false")
	if err := loadCmd.RunE(loadCmd, []string{}); err != nil {
		t.Fatal(err)
	}

	createEnvFile(t, ".dev.env", `ENV_1=v2`)
	if err := loadCmd.RunE(loadCmd, []string{}); err != nil {
		t.Fatal(err)
	}

	// ls with -v flag should show versions
	lsCmd := cmd.GetLsCmd()
	lsCmd.Flags().Set("version", "true")
	if err := lsCmd.RunE(lsCmd, []string{}); err != nil {
		t.Fatalf("ls -v failed: %v", err)
	}
}

func TestSetVersionByName(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	// Create versions
	createEnvFile(t, ".dev.env", `ENV_1=v1`)
	loadCmd := cmd.GetLoadCmd()
	loadCmd.Flags().Set("env", "*")
	loadCmd.Flags().Set("replace", "false")
	if err := loadCmd.RunE(loadCmd, []string{}); err != nil {
		t.Fatal(err)
	}

	createEnvFile(t, ".dev.env", `ENV_1=v2`)
	if err := loadCmd.RunE(loadCmd, []string{}); err != nil {
		t.Fatal(err)
	}

	// Name v1
	if err := filehandler.RenameVersion("test-project", 1, "release-1"); err != nil {
		t.Fatal(err)
	}

	// Switch to version by name
	versionCmd := cmd.GetVersionCmd()
	if err := versionCmd.RunE(versionCmd, []string{"release-1"}); err != nil {
		t.Fatalf("version release-1 failed: %v", err)
	}

	// Verify current version is 1
	project, err := filehandler.FindProjectByLocalPath(testProjectDir)
	if err != nil {
		t.Fatal(err)
	}
	if project.CurrentVersion != 1 {
		t.Errorf("expected current version 1, got %d", project.CurrentVersion)
	}
}

func TestToWithNamedVersion(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	// Create v1
	createEnvFile(t, ".dev.env", `ENV_1=v1_value`)
	loadCmd := cmd.GetLoadCmd()
	loadCmd.Flags().Set("env", "*")
	loadCmd.Flags().Set("replace", "false")
	if err := loadCmd.RunE(loadCmd, []string{}); err != nil {
		t.Fatal(err)
	}

	// Name v1
	if err := filehandler.RenameVersion("test-project", 1, "stable"); err != nil {
		t.Fatal(err)
	}

	// Create v2
	createEnvFile(t, ".dev.env", `ENV_1=v2_value`)
	if err := loadCmd.RunE(loadCmd, []string{}); err != nil {
		t.Fatal(err)
	}

	// Switch to dev using named version "stable"
	toCmd := cmd.GetToCmd()
	toCmd.Flags().Set("replace", "true")
	toCmd.Flags().Set("version", "stable")
	if err := toCmd.RunE(toCmd, []string{"dev"}); err != nil {
		t.Fatalf("to dev --version stable failed: %v", err)
	}

	content, err := os.ReadFile(".env")
	if err != nil {
		t.Fatal(err)
	}

	contentStr := string(content)
	if !contains(contentStr, "ENV_1=v1_value") {
		t.Error("ENV_1=v1_value should be in .env when using --version stable")
	}
}
