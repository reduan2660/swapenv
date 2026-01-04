package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/reduan2660/swapenv/cmd"
)

func TestLoad(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	devEnvContent := `ENV_1=dev


ENV_2=dev
# ENV_3=dev`
	createEnvFile(t, ".dev.env", devEnvContent)

	prodEnvContent := `ENV_1=prod

ENV_4=prod`
	createEnvFile(t, ".prod.env", prodEnvContent)

	loadCmd := cmd.GetLoadCmd()
	loadCmd.Flags().Set("env", "*")
	loadCmd.Flags().Set("replace", "false")

	if err := loadCmd.RunE(loadCmd, []string{}); err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if _, err := os.Stat(".dev.env"); !os.IsNotExist(err) {
		t.Error(".dev.env should be deleted after load")
	}
	if _, err := os.Stat(".prod.env"); !os.IsNotExist(err) {
		t.Error(".prod.env should be deleted after load")
	}

	// With versioning, files are stored as v1.json instead of project.json
	projectPath := filepath.Join(testHomeDir, "test-project", "v1.json")
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		t.Error("v1.json should exist after load")
	}
}

func TestLoadReplace(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	devEnvContent := `ENV_1=dev


ENV_2=dev
# ENV_3=dev`
	createEnvFile(t, ".dev.env", devEnvContent)

	prodEnvContent := `ENV_1=prod

ENV_4=prod`
	createEnvFile(t, ".prod.env", prodEnvContent)

	loadCmd := cmd.GetLoadCmd()
	loadCmd.Flags().Set("env", "*")
	loadCmd.Flags().Set("replace", "false")
	if err := loadCmd.RunE(loadCmd, []string{}); err != nil {
		t.Fatal(err)
	}

	toCmd := cmd.GetToCmd()
	toCmd.Flags().Set("replace", "false")
	if err := toCmd.RunE(toCmd, []string{"dev"}); err != nil {
		t.Fatal(err)
	}

	newDevEnvContent := `ENV_2=dev_updated

ENV_5=dev_updated`
	createEnvFile(t, ".dev.env", newDevEnvContent)

	loadCmd.Flags().Set("env", "*")
	loadCmd.Flags().Set("replace", "true")

	if err := loadCmd.RunE(loadCmd, []string{}); err != nil {
		t.Fatalf("load with replace failed: %v", err)
	}

	toCmd.Flags().Set("replace", "true")
	if err := toCmd.RunE(toCmd, []string{"dev"}); err != nil {
		t.Fatalf("to dev after replace failed: %v", err)
	}

	content, err := os.ReadFile(".env")
	if err != nil {
		t.Fatal(err)
	}

	contentStr := string(content)
	if contains(contentStr, "ENV_1=dev") {
		t.Error("ENV_1 should not exist after replace")
	}
	if !contains(contentStr, "ENV_2=dev_updated") {
		t.Error("ENV_2=dev_updated should be in .env")
	}
	if !contains(contentStr, "ENV_5=dev_updated") {
		t.Error("ENV_5=dev_updated should be in .env")
	}
}

func TestLoadCreatesNewVersion(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	devEnvContent := `ENV_1=dev`
	createEnvFile(t, ".dev.env", devEnvContent)

	loadCmd := cmd.GetLoadCmd()
	loadCmd.Flags().Set("env", "*")
	loadCmd.Flags().Set("replace", "false")

	if err := loadCmd.RunE(loadCmd, []string{}); err != nil {
		t.Fatalf("first load failed: %v", err)
	}

	// v1 should exist
	v1Path := filepath.Join(testHomeDir, "test-project", "v1.json")
	if _, err := os.Stat(v1Path); os.IsNotExist(err) {
		t.Error("v1.json should exist after first load")
	}

	// Load again - should create v2
	createEnvFile(t, ".dev.env", `ENV_1=dev_v2`)

	if err := loadCmd.RunE(loadCmd, []string{}); err != nil {
		t.Fatalf("second load failed: %v", err)
	}

	v2Path := filepath.Join(testHomeDir, "test-project", "v2.json")
	if _, err := os.Stat(v2Path); os.IsNotExist(err) {
		t.Error("v2.json should exist after second load")
	}
}

func TestLoadMergesWithPreviousVersion(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	// First load with two variables
	devEnvContent := `ENV_1=dev
ENV_2=dev`
	createEnvFile(t, ".dev.env", devEnvContent)

	loadCmd := cmd.GetLoadCmd()
	loadCmd.Flags().Set("env", "*")
	loadCmd.Flags().Set("replace", "false")

	if err := loadCmd.RunE(loadCmd, []string{}); err != nil {
		t.Fatalf("first load failed: %v", err)
	}

	// Second load with only one variable (should merge with v1)
	createEnvFile(t, ".dev.env", `ENV_1=dev_updated`)

	if err := loadCmd.RunE(loadCmd, []string{}); err != nil {
		t.Fatalf("second load failed: %v", err)
	}

	// Switch to dev and verify both variables exist
	toCmd := cmd.GetToCmd()
	toCmd.Flags().Set("replace", "true")
	if err := toCmd.RunE(toCmd, []string{"dev"}); err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(".env")
	if err != nil {
		t.Fatal(err)
	}

	contentStr := string(content)
	// ENV_1 should be updated
	if !contains(contentStr, "ENV_1=dev_updated") {
		t.Error("ENV_1=dev_updated should be in .env")
	}
	// ENV_2 should be preserved from v1
	if !contains(contentStr, "ENV_2=dev") {
		t.Error("ENV_2=dev should be merged from previous version")
	}
}
