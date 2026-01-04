package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

var (
	testProjectDir string
	testHomeDir    string
)

func setupTestEnv(t *testing.T) func() {
	t.Helper()

	tempDir := t.TempDir()
	// Resolve symlinks to match os.Getwd() behavior (e.g., /var -> /private/var on macOS)
	tempDir, _ = filepath.EvalSymlinks(tempDir)
	testHomeDir = filepath.Join(tempDir, ".swapenv-test")

	configPath := filepath.Join(tempDir, "test-config.yaml")
	configContent := "home_directory: " + testHomeDir + "\nmax_versions: 5\n"
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
