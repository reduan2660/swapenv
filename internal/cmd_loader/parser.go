package cmd_loader

import (
	"bufio"
	"bytes"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/reduan2660/swapenv/internal/types"
)

func ParseEnv(content []byte) ([]types.EnvValue, error) {

	envValues := make([]types.EnvValue, 0)

	scanner := bufio.NewScanner(bytes.NewReader(content))
	order := 1
	spacing := 0

	for scanner.Scan() {

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			spacing++
		}

		if strings.HasPrefix(line, "#") { // TODO: consider digits, and ?
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		if len(key) == 0 {
			continue
		}

		val = strings.Trim(val, `"'`)

		envValues = append(envValues, types.EnvValue{
			Key:     key,
			Val:     val,
			Order:   order,
			Spacing: spacing,
		})

		order++
		spacing = 0

	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return envValues, nil
}

type MergeEnvConfig struct {
	Replace          bool   // if true, just return incoming (ignore current)
	ConflictPriority string // "incoming" or "current" - which value wins for matching keys
}

// MergeEnv merges incoming and current environment values.
//
// Order: current's keys come first (in current's order), then incoming-only keys are appended.
// Conflict resolution: when a key exists in both, ConflictPriority determines which value wins:
//   - "incoming": use incoming's value
//   - "current": use current's value
//
// If Replace=true, just return incoming (ignore current entirely).
func MergeEnv(incoming, current []types.EnvValue, config MergeEnvConfig) []types.EnvValue {

	if config.Replace {
		return incoming
	}

	incomingMap := make(map[string]types.EnvValue)
	for _, ev := range incoming {
		incomingMap[ev.Key] = ev
	}

	marked := make(map[string]bool)
	merged := make([]types.EnvValue, 0)

	for _, ev := range current {
		if incomingVal, exists := incomingMap[ev.Key]; exists {
			marked[ev.Key] = true

			if config.ConflictPriority == "incoming" {
				incomingVal.Order = ev.Order
				merged = append(merged, incomingVal)
			} else {
				// "current" or default: current's value wins
				merged = append(merged, ev)
			}
		} else {
			merged = append(merged, ev)
		}
	}

	for idx, ev := range incoming {
		if !marked[ev.Key] {
			ev.Order = len(current) + idx
			merged = append(merged, ev)
		}
	}

	return merged
}

func MarshalProject(projectName, owner, localDirectory string, envs map[string][]types.EnvValue) types.Project {

	now := time.Now().UTC().Unix()

	return types.Project{
		Id:             uuid.New().String(),
		Name:           projectName,
		Owner:          owner,
		LocalDirectory: localDirectory,
		CreatedAt:      now,
		ModifiedAt:     now,
		Envs:           envs,
	}
}
