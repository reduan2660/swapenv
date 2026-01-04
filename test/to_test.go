package test

import (
	"os"
	"testing"

	"github.com/reduan2660/swapenv/cmd"
)

func TestTo(t *testing.T) {
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
		t.Fatalf("to dev failed: %v", err)
	}

	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		t.Error(".env should exist after switching")
	}

	content, err := os.ReadFile(".env")
	if err != nil {
		t.Fatal(err)
	}

	contentStr := string(content)
	if !contains(contentStr, "ENV_1=dev") {
		t.Error("ENV_1=dev should be in .env")
	}
	if !contains(contentStr, "ENV_2=dev") {
		t.Error("ENV_2=dev should be in .env")
	}
}

func TestLs(t *testing.T) {
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

	lsCmd := cmd.GetLsCmd()
	if err := lsCmd.RunE(lsCmd, []string{}); err != nil {
		t.Fatalf("ls failed: %v", err)
	}
}

func TestRoot(t *testing.T) {
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

	rootCmd := cmd.GetRootCmd()
	if err := rootCmd.RunE(rootCmd, []string{}); err != nil {
		t.Fatalf("root failed: %v", err)
	}
}

func TestToWithSpecificVersion(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	// Create and load v1
	createEnvFile(t, ".dev.env", `ENV_1=v1_value`)

	loadCmd := cmd.GetLoadCmd()
	loadCmd.Flags().Set("env", "*")
	loadCmd.Flags().Set("replace", "false")
	if err := loadCmd.RunE(loadCmd, []string{}); err != nil {
		t.Fatal(err)
	}

	// Create and load v2
	createEnvFile(t, ".dev.env", `ENV_1=v2_value`)
	if err := loadCmd.RunE(loadCmd, []string{}); err != nil {
		t.Fatal(err)
	}

	// Switch to dev using v1
	toCmd := cmd.GetToCmd()
	toCmd.Flags().Set("replace", "true")
	toCmd.Flags().Set("version", "1")
	if err := toCmd.RunE(toCmd, []string{"dev"}); err != nil {
		t.Fatalf("to dev --version 1 failed: %v", err)
	}

	content, err := os.ReadFile(".env")
	if err != nil {
		t.Fatal(err)
	}

	contentStr := string(content)
	if !contains(contentStr, "ENV_1=v1_value") {
		t.Error("ENV_1=v1_value should be in .env when using --version 1")
	}

	// Now switch using latest
	toCmd.Flags().Set("version", "latest")
	if err := toCmd.RunE(toCmd, []string{"dev"}); err != nil {
		t.Fatalf("to dev --version latest failed: %v", err)
	}

	content, err = os.ReadFile(".env")
	if err != nil {
		t.Fatal(err)
	}

	contentStr = string(content)
	if !contains(contentStr, "ENV_1=v2_value") {
		t.Error("ENV_1=v2_value should be in .env when using --version latest")
	}
}
