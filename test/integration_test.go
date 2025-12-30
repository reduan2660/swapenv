package test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/reduan2660/swapenv/cmd"
	"github.com/spf13/viper"
)

var (
	testProjectDir string
	testHomeDir    string
)

func setupTestEnv(t *testing.T) func() {
	t.Helper()

	tempDir := t.TempDir()
	testHomeDir = filepath.Join(tempDir, ".swapenv-test")

	configPath := filepath.Join(tempDir, "test-config.yaml")
	configContent := "home_directory: " + testHomeDir + "\n"
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	testProjectDir = filepath.Join(tempDir, "test-project")
	if err := os.MkdirAll(testProjectDir, 0755); err != nil {
		t.Fatal(err)
	}

	origDir, _ := os.Getwd()
	if err := os.Chdir(testProjectDir); err != nil {
		t.Fatal(err)
	}

	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		t.Fatal(err)
	}

	cleanup := func() {
		os.Chdir(origDir)
		viper.Reset()
	}

	return cleanup
}

func createEnvFile(t *testing.T, filename, content string) {
	t.Helper()
	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

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

	projectPath := filepath.Join(testHomeDir, "test-project", "project.json")
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		t.Error("project.json should exist after load")
	}
}

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

func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

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
