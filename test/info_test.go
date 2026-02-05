package test

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"testing"

	"github.com/reduan2660/swapenv/cmd"
	"github.com/reduan2660/swapenv/internal/cmd_info"
)

func captureOutput(f func() error) (string, error) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String(), err
}

func TestInfoNoProject(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	infoCmd := cmd.GetInfoCmd()
	infoCmd.Flags().Set("format", "json")

	output, err := captureOutput(func() error {
		return infoCmd.RunE(infoCmd, []string{})
	})
	if err != nil {
		t.Fatalf("info failed: %v", err)
	}

	var info cmd_info.ProjectInfo
	if err := json.Unmarshal([]byte(output), &info); err != nil {
		t.Fatalf("failed to parse JSON output: %v", err)
	}

	if info.Project != "" {
		t.Error("project should be empty when no project exists")
	}
	if info.Env != "" {
		t.Error("env should be empty when no project exists")
	}
}

func TestInfoWithProject(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	// Load a project
	createEnvFile(t, ".dev.env", `ENV_1=dev`)
	createEnvFile(t, ".prod.env", `ENV_1=prod`)
	loadCmd := cmd.GetLoadCmd()
	loadCmd.Flags().Set("env", "*")
	loadCmd.Flags().Set("replace", "false")
	if err := loadCmd.RunE(loadCmd, []string{}); err != nil {
		t.Fatal(err)
	}

	infoCmd := cmd.GetInfoCmd()
	infoCmd.Flags().Set("format", "json")

	output, err := captureOutput(func() error {
		return infoCmd.RunE(infoCmd, []string{})
	})
	if err != nil {
		t.Fatalf("info failed: %v", err)
	}

	var info cmd_info.ProjectInfo
	if err := json.Unmarshal([]byte(output), &info); err != nil {
		t.Fatalf("failed to parse JSON output: %v", err)
	}

	if info.Project != "test-project" {
		t.Errorf("expected project 'test-project', got '%s'", info.Project)
	}
	if info.LatestVersion != 1 {
		t.Errorf("expected latest_version 1, got %d", info.LatestVersion)
	}
	if len(info.Envs) != 2 {
		t.Errorf("expected 2 envs, got %d", len(info.Envs))
	}
}

func TestInfoWithEnv(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	// Load a project
	createEnvFile(t, ".dev.env", `ENV_1=dev`)
	loadCmd := cmd.GetLoadCmd()
	loadCmd.Flags().Set("env", "*")
	loadCmd.Flags().Set("replace", "false")
	if err := loadCmd.RunE(loadCmd, []string{}); err != nil {
		t.Fatal(err)
	}

	// Switch to dev env
	toCmd := cmd.GetToCmd()
	toCmd.Flags().Set("replace", "true")
	if err := toCmd.RunE(toCmd, []string{"dev"}); err != nil {
		t.Fatal(err)
	}

	infoCmd := cmd.GetInfoCmd()
	infoCmd.Flags().Set("format", "json")

	output, err := captureOutput(func() error {
		return infoCmd.RunE(infoCmd, []string{})
	})
	if err != nil {
		t.Fatalf("info failed: %v", err)
	}

	var info cmd_info.ProjectInfo
	if err := json.Unmarshal([]byte(output), &info); err != nil {
		t.Fatalf("failed to parse JSON output: %v", err)
	}

	if info.Project != "test-project" {
		t.Errorf("expected project 'test-project', got '%s'", info.Project)
	}
	if info.Env != "dev" {
		t.Errorf("expected env 'dev', got '%s'", info.Env)
	}
}

func TestInfoPlainFormat(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	// Load a project
	createEnvFile(t, ".dev.env", `ENV_1=dev`)
	loadCmd := cmd.GetLoadCmd()
	loadCmd.Flags().Set("env", "*")
	loadCmd.Flags().Set("replace", "false")
	if err := loadCmd.RunE(loadCmd, []string{}); err != nil {
		t.Fatal(err)
	}

	// Switch to dev env
	toCmd := cmd.GetToCmd()
	toCmd.Flags().Set("replace", "true")
	if err := toCmd.RunE(toCmd, []string{"dev"}); err != nil {
		t.Fatal(err)
	}

	infoCmd := cmd.GetInfoCmd()
	infoCmd.Flags().Set("format", "plain")

	output, err := captureOutput(func() error {
		return infoCmd.RunE(infoCmd, []string{})
	})
	if err != nil {
		t.Fatalf("info failed: %v", err)
	}

	expected := "test-project:dev\n"
	if output != expected {
		t.Errorf("expected '%s', got '%s'", expected, output)
	}
}

func TestInfoPlainFormatNoEnv(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	// Load a project but don't switch to an env
	createEnvFile(t, ".dev.env", `ENV_1=dev`)
	loadCmd := cmd.GetLoadCmd()
	loadCmd.Flags().Set("env", "*")
	loadCmd.Flags().Set("replace", "false")
	if err := loadCmd.RunE(loadCmd, []string{}); err != nil {
		t.Fatal(err)
	}

	infoCmd := cmd.GetInfoCmd()
	infoCmd.Flags().Set("format", "plain")

	output, err := captureOutput(func() error {
		return infoCmd.RunE(infoCmd, []string{})
	})
	if err != nil {
		t.Fatalf("info failed: %v", err)
	}

	expected := "test-project:none\n"
	if output != expected {
		t.Errorf("expected '%s', got '%s'", expected, output)
	}
}

func TestInfoPlainFormatNoProject(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	infoCmd := cmd.GetInfoCmd()
	infoCmd.Flags().Set("format", "plain")

	output, err := captureOutput(func() error {
		return infoCmd.RunE(infoCmd, []string{})
	})
	if err != nil {
		t.Fatalf("info failed: %v", err)
	}

	// Plain format with no project should output nothing
	if output != "" {
		t.Errorf("expected empty output, got '%s'", output)
	}
}

func TestInfoWithMultipleVersions(t *testing.T) {
	cleanup := setupTestEnv(t)
	defer cleanup()

	// Create v1
	createEnvFile(t, ".dev.env", `ENV_1=v1`)
	loadCmd := cmd.GetLoadCmd()
	loadCmd.Flags().Set("env", "*")
	loadCmd.Flags().Set("replace", "false")
	// Discard load output
	captureOutput(func() error {
		return loadCmd.RunE(loadCmd, []string{})
	})

	// Create v2
	createEnvFile(t, ".dev.env", `ENV_1=v2`)
	// Discard load output
	captureOutput(func() error {
		return loadCmd.RunE(loadCmd, []string{})
	})

	infoCmd := cmd.GetInfoCmd()
	infoCmd.Flags().Set("format", "json") // Ensure JSON format

	output, err := captureOutput(func() error {
		return infoCmd.RunE(infoCmd, []string{})
	})
	if err != nil {
		t.Fatalf("info failed: %v", err)
	}

	var info cmd_info.ProjectInfo
	if err := json.Unmarshal([]byte(output), &info); err != nil {
		t.Fatalf("failed to parse JSON output: %v", err)
	}

	if info.Version != 2 {
		t.Errorf("expected version 2, got %d", info.Version)
	}
	if info.LatestVersion != 2 {
		t.Errorf("expected latest_version 2, got %d", info.LatestVersion)
	}
}
