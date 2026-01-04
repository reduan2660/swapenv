package test

import (
	"os"
	"strings"
	"testing"

	"github.com/reduan2660/swapenv/cmd"
)

func TestToWithCommon(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	// Create common env with shared variables
	commonEnvContent := `SHARED_VAR=common
DATABASE_URL=common_db`
	createEnvFile(t, ".common.env", commonEnvContent)

	// Create dev env with its own variables and one override
	devEnvContent := `ENV_1=dev
DATABASE_URL=dev_db`
	createEnvFile(t, ".dev.env", devEnvContent)

	// Load both environments
	loadCmd := cmd.GetLoadCmd()
	loadCmd.Flags().Set("env", "*")
	loadCmd.Flags().Set("replace", "false")
	if err := loadCmd.RunE(loadCmd, []string{}); err != nil {
		t.Fatal(err)
	}

	// Switch to dev (should merge with common)
	toCmd := cmd.GetToCmd()
	toCmd.Flags().Set("replace", "false")
	toCmd.Flags().Set("skip-common", "false")

	if err := toCmd.RunE(toCmd, []string{"dev"}); err != nil {
		t.Fatalf("to dev failed: %v", err)
	}

	content, err := os.ReadFile(".env")
	if err != nil {
		t.Fatal(err)
	}

	contentStr := string(content)

	// Should have dev's ENV_1
	if !contains(contentStr, "ENV_1=dev") {
		t.Error("ENV_1=dev should be in .env")
	}

	// Should have common's SHARED_VAR
	if !contains(contentStr, "SHARED_VAR=common") {
		t.Error("SHARED_VAR=common should be in .env (from common)")
	}

	// When key exists in both, target env (dev) overrides common
	if !contains(contentStr, "DATABASE_URL=dev_db") {
		t.Error("DATABASE_URL=dev_db should be in .env (dev overrides common)")
	}

	// Verify ordering: dev's variables should come before common-only variables
	env1Pos := strings.Index(contentStr, "ENV_1=dev")
	dbPos := strings.Index(contentStr, "DATABASE_URL=dev_db")
	sharedPos := strings.Index(contentStr, "SHARED_VAR=common")

	if sharedPos < env1Pos {
		t.Error("SHARED_VAR (common-only) should appear after ENV_1 (dev)")
	}
	if sharedPos < dbPos {
		t.Error("SHARED_VAR (common-only) should appear after DATABASE_URL (dev)")
	}
}

func TestToSkipCommon(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	// Create common env with shared variables
	commonEnvContent := `SHARED_VAR=common
DATABASE_URL=common_db`
	createEnvFile(t, ".common.env", commonEnvContent)

	// Create dev env
	devEnvContent := `ENV_1=dev
DATABASE_URL=dev_db`
	createEnvFile(t, ".dev.env", devEnvContent)

	// Load both environments
	loadCmd := cmd.GetLoadCmd()
	loadCmd.Flags().Set("env", "*")
	loadCmd.Flags().Set("replace", "false")
	if err := loadCmd.RunE(loadCmd, []string{}); err != nil {
		t.Fatal(err)
	}

	// Switch to dev with skip-common flag
	toCmd := cmd.GetToCmd()
	toCmd.Flags().Set("replace", "true")
	toCmd.Flags().Set("skip-common", "true")

	if err := toCmd.RunE(toCmd, []string{"dev"}); err != nil {
		t.Fatalf("to dev failed: %v", err)
	}

	content, err := os.ReadFile(".env")
	if err != nil {
		t.Fatal(err)
	}

	contentStr := string(content)

	// Should have dev's ENV_1
	if !contains(contentStr, "ENV_1=dev") {
		t.Error("ENV_1=dev should be in .env")
	}

	// Should NOT have common's SHARED_VAR (skip-common was set)
	if contains(contentStr, "SHARED_VAR=common") {
		t.Error("SHARED_VAR should NOT be in .env when skip-common is set")
	}

	// Should have dev's DATABASE_URL
	if !contains(contentStr, "DATABASE_URL=dev_db") {
		t.Error("DATABASE_URL=dev_db should be in .env")
	}
}

func TestToCommonEnvNotAllowed(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	// Create common env
	commonEnvContent := `SHARED_VAR=common
DATABASE_URL=common_db`
	createEnvFile(t, ".common.env", commonEnvContent)

	// Load common environment
	loadCmd := cmd.GetLoadCmd()
	loadCmd.Flags().Set("env", "*")
	loadCmd.Flags().Set("replace", "false")
	if err := loadCmd.RunE(loadCmd, []string{}); err != nil {
		t.Fatal(err)
	}

	// Attempt to switch to common env - should fail
	toCmd := cmd.GetToCmd()
	toCmd.Flags().Set("replace", "true")
	toCmd.Flags().Set("skip-common", "false")

	err := toCmd.RunE(toCmd, []string{"common"})
	if err == nil {
		t.Fatal("switching to common should return an error")
	}
}
