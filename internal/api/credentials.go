package api

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/reduan2660/swapenv/internal/filehandler"
)

type Credentials struct {
	Token     string `json:"token"`
	UserId    string `json:"user_id"`
	OrgId     string `json:"org_id"`
	ExpiresAt int64  `json:"expires_at"`
}

func GetCredentialsPath() (string, error) {
	baseDir, err := filehandler.GetBaseDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(baseDir, "credentials.json"), nil
}

func LoadCredentials() (*Credentials, error) {

	path, err := GetCredentialsPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var creds Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, err
	}

	return &creds, nil
}

func SaveCredentials(creds *Credentials) error {
	path, err := GetCredentialsPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

func DeleteCredentials() error {
	path, err := GetCredentialsPath()
	if err != nil {
		return err
	}
	return os.Remove(path)
}

func IsLoggedIn() bool {
	creds, err := LoadCredentials()
	if err != nil || creds.Token == "" {
		return false
	}

	return time.Now().Unix() < creds.ExpiresAt
}
