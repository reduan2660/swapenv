package test

import (
	"os"
	"testing"

	"github.com/reduan2660/swapenv/cmd"
)

func TestSpit(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	devEnvContent := `ENV_1=dev


ENV_2=dev`
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

	if _, err := os.Stat(".dev.env"); !os.IsNotExist(err) {
		t.Error(".dev.env should be deleted after load")
	}

	spitCmd := cmd.GetSpitCmd()
	spitCmd.Flags().Set("env", "*")
	if err := spitCmd.RunE(spitCmd, []string{}); err != nil {
		t.Fatalf("spit failed: %v", err)
	}

	if _, err := os.Stat(".dev.env"); os.IsNotExist(err) {
		t.Error(".dev.env should exist after spit")
	}
	if _, err := os.Stat(".prod.env"); os.IsNotExist(err) {
		t.Error(".prod.env should exist after spit")
	}

	devContent, err := os.ReadFile(".dev.env")
	if err != nil {
		t.Fatal(err)
	}
	devStr := string(devContent)
	if !contains(devStr, "ENV_1=dev") {
		t.Error("ENV_1=dev should be in .dev.env")
	}
	if !contains(devStr, "ENV_2=dev") {
		t.Error("ENV_2=dev should be in .dev.env")
	}

	prodContent, err := os.ReadFile(".prod.env")
	if err != nil {
		t.Fatal(err)
	}
	prodStr := string(prodContent)
	if !contains(prodStr, "ENV_1=prod") {
		t.Error("ENV_1=prod should be in .prod.env")
	}
	if !contains(prodStr, "ENV_4=prod") {
		t.Error("ENV_4=prod should be in .prod.env")
	}
}

func TestSpitWithSpecificVersion(t *testing.T) {
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

	// Spit from v1
	spitCmd := cmd.GetSpitCmd()
	spitCmd.Flags().Set("env", "*")
	spitCmd.Flags().Set("version", "1")
	if err := spitCmd.RunE(spitCmd, []string{}); err != nil {
		t.Fatalf("spit --version 1 failed: %v", err)
	}

	content, err := os.ReadFile(".dev.env")
	if err != nil {
		t.Fatal(err)
	}

	contentStr := string(content)
	if !contains(contentStr, "ENV_1=v1_value") {
		t.Error("ENV_1=v1_value should be in .dev.env when spitting from version 1")
	}

	// Spit from latest (v2)
	spitCmd.Flags().Set("version", "latest")
	if err := spitCmd.RunE(spitCmd, []string{}); err != nil {
		t.Fatalf("spit --version latest failed: %v", err)
	}

	content, err = os.ReadFile(".dev.env")
	if err != nil {
		t.Fatal(err)
	}

	contentStr = string(content)
	if !contains(contentStr, "ENV_1=v2_value") {
		t.Error("ENV_1=v2_value should be in .dev.env when spitting from latest version")
	}
}

func TestSpitSpecificEnv(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	devEnvContent := `ENV_1=dev`
	createEnvFile(t, ".dev.env", devEnvContent)

	prodEnvContent := `ENV_1=prod`
	createEnvFile(t, ".prod.env", prodEnvContent)

	loadCmd := cmd.GetLoadCmd()
	loadCmd.Flags().Set("env", "*")
	loadCmd.Flags().Set("replace", "false")
	if err := loadCmd.RunE(loadCmd, []string{}); err != nil {
		t.Fatal(err)
	}

	// Only spit dev
	spitCmd := cmd.GetSpitCmd()
	spitCmd.Flags().Set("env", "dev")
	if err := spitCmd.RunE(spitCmd, []string{}); err != nil {
		t.Fatalf("spit --env dev failed: %v", err)
	}

	if _, err := os.Stat(".dev.env"); os.IsNotExist(err) {
		t.Error(".dev.env should exist after spit --env dev")
	}
	if _, err := os.Stat(".prod.env"); !os.IsNotExist(err) {
		t.Error(".prod.env should NOT exist after spit --env dev")
	}
}
