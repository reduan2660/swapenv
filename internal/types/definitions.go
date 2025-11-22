package types

import (
	"encoding/json"
	"fmt"
)

type ProjectDir struct {
	ProjectName string `json:"projectName"`
	CurrentEnv  string `json:"currentEnv"`
	LocalPath   string `json:"localPath"`
	RemotePath  string `json:"remotePath"`
}

type Project struct {
	Id             string                `json:"id"`
	Name           string                `json:"name"`
	Owner          string                `json:"owner"`
	LocalDirectory string                `json:"localDirectory"`
	CreatedAt      int64                 `json:"createdAt"`
	ModifiedAt     int64                 `json:"modifiedAt"`
	Envs           map[string][]EnvValue `json:"-"`
}

type EnvValue struct {
	Key     string `json:"key"`
	Val     string `json:"val"`
	Order   int    `json:"order"`
	Spacing int    `json:"spacing"`
}

func (e EnvValue) String() string {
	return fmt.Sprintf("%d - %s=%s", e.Order, e.Key, e.Val)
}

func (p *Project) MarshalJSON() ([]byte, error) {
	m := map[string]any{
		"id":             p.Id,
		"owner":          p.Owner,
		"localDirectory": p.LocalDirectory,
		"createdAt":      p.CreatedAt,
		"modifiedAt":     p.ModifiedAt,
	}

	for envName, envValues := range p.Envs {
		m[envName] = envValues
	}

	return json.Marshal(map[string]any{p.Name: m})
}
