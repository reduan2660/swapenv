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
