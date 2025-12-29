package filehandler

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/reduan2660/swapenv/internal/types"
)

func ReadProjectEnv(projectPath, envName string) ([]types.EnvValue, error) {
	data, err := os.ReadFile(projectPath)
	if err != nil {
		return nil, err
	}

	var outer map[string]json.RawMessage
	if err := json.Unmarshal(data, &outer); err != nil {
		return nil, err
	}

	if len(outer) != 1 {
		return nil, fmt.Errorf("invalid project JSON: expected 1 project, got %d", len(outer))
	}

	var innerData json.RawMessage
	for _, data := range outer {
		innerData = data
		break
	}

	var inner map[string]json.RawMessage
	if err := json.Unmarshal(innerData, &inner); err != nil {
		return nil, err
	}

	envData, exists := inner[envName]
	if !exists {
		return nil, fmt.Errorf("environment '%s' not found in project", envName)
	}

	var envValues []types.EnvValue
	if err := json.Unmarshal(envData, &envValues); err != nil {
		return nil, fmt.Errorf("failed to parse environment '%s': %w", envName, err)
	}

	return envValues, nil
}

func ListProjectEnv(projectPath string) ([]string, error) {
	data, err := os.ReadFile(projectPath)
	if err != nil {
		return nil, err
	}

	var outer map[string]json.RawMessage
	if err := json.Unmarshal(data, &outer); err != nil {
		return nil, err
	}

	if len(outer) != 1 {
		return nil, fmt.Errorf("invalid project JSON: expected 1 project, got %d", len(outer))
	}

	var innerData json.RawMessage
	for _, data := range outer {
		innerData = data
		break
	}

	var inner map[string]json.RawMessage
	if err := json.Unmarshal(innerData, &inner); err != nil {
		return nil, err
	}

	excludeFields := map[string]bool{
		"id": true, "owner": true, "localDirectory": true,
		"createdAt": true, "modifiedAt": true,
	}

	envNames := make([]string, 0)
	for key := range inner {
		if !excludeFields[key] {
			envNames = append(envNames, key)
		}
	}

	return envNames, nil
}

func WriteProject(directory, filePath string, file_content []byte) error {
	if err := os.MkdirAll(directory, 0755); err != nil { // TODO - consider 0700 - rething permissions
		return err
	}

	return os.WriteFile(filePath, file_content, 0644)
}

func WriteEnv(envValues []types.EnvValue, filepath string) error {
	var builder strings.Builder

	for _, ev := range envValues {
		for i := 0; i < ev.Spacing; i++ {
			builder.WriteString("\n")
		}

		builder.WriteString(fmt.Sprintf("%s=%s\n", ev.Key, ev.Val))
	}

	content := builder.String()
	if len(content) > 0 && content[len(content)-1] == '\n' {
		content = content[:len(content)-1]
	}

	return os.WriteFile(filepath, []byte(content), 0644)

}

func DeleteEnvFiles(files []string) error {

	for _, file := range files {
		if err := os.Remove(file); err != nil {
			return fmt.Errorf("error deleting %s : %w", file, err)
		}
	}

	return nil
}
