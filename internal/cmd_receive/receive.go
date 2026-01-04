package cmd_receive

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/reduan2660/swapenv/internal/api"
	"github.com/reduan2660/swapenv/internal/cmd_login"
	"github.com/reduan2660/swapenv/internal/cmd_logout"
	"github.com/reduan2660/swapenv/internal/crypto"
	"github.com/reduan2660/swapenv/internal/filehandler"
	"github.com/reduan2660/swapenv/internal/types"
)

type wsMessage struct {
	Type    string   `json:"type"`
	Code    string   `json:"code,omitempty"`
	Codes   []string `json:"codes,omitempty"`
	Message string   `json:"message,omitempty"`
}

func Receive(serverURL string) error {
	if !api.IsLoggedIn() {
		fmt.Println("Not logged in. Starting login flow...")
		if err := cmd_login.Login(serverURL); err != nil {
			return err
		}
	}

	creds, err := api.LoadCredentials()
	if err != nil {
		cmd_logout.Logout()
		return fmt.Errorf("failed to load credentials: %w. logged you out, try again", err)
	}

	conn, err := api.ConnectWS(serverURL, "/receive", creds.Token)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer conn.Close()

	var msg wsMessage
	if err := conn.ReadJSON(&msg); err != nil {
		return err
	}

	switch msg.Type {
	case "error":
		return fmt.Errorf("server error: %s", msg.Message)

	case "choose":
		fmt.Println("Multiple active streams:")
		for i, code := range msg.Codes {
			fmt.Printf("  %d. %s\n", i+1, code)
		}
		fmt.Print("Enter stream: ")

		var choice string
		fmt.Scanln(&choice)

		if err := conn.WriteJSON(map[string]string{"code": choice}); err != nil {
			return err
		}

		if err := conn.ReadJSON(&msg); err != nil {
			return err
		}

		if msg.Type == "error" {
			return fmt.Errorf("server error: %s", msg.Message)
		}

	case "connected":
		fmt.Printf("Connected to stream: %s\n", msg.Code)
	}

	privKey, err := crypto.GenerateKeyPair()
	if err != nil {
		return fmt.Errorf("failed to generate keypair: %w", err)
	}

	if err := conn.WriteMessage(1, privKey.PublicKey().Bytes()); err != nil {
		return err
	}

	fmt.Println("Waiting for encrypted data...")

	_, rawPayload, err := conn.ReadMessage()
	if err != nil {
		return err
	}

	decoded, err := base64.StdEncoding.DecodeString(string(rawPayload))
	if err != nil {
		return fmt.Errorf("invalid payload encoding: %w", err)
	}

	if len(decoded) < 2 {
		return fmt.Errorf("payload too short")
	}
	nameLen := int(decoded[0])
	if len(decoded) < 1+nameLen {
		return fmt.Errorf("invalid payload structure")
	}
	projectName := string(decoded[1 : 1+nameLen])
	encrypted := decoded[1+nameLen:]

	decrypted, err := crypto.Decrypt(encrypted, privKey)
	if err != nil {
		return fmt.Errorf("decryption failed: %w", err)
	}

	var envMap map[string][]types.EnvValue
	if err := json.Unmarshal(decrypted, &envMap); err != nil {
		return fmt.Errorf("invalid env data: %w", err)
	}

	version, err := saveReceived(projectName, envMap)
	if err != nil {
		return fmt.Errorf("failed to save: %w", err)
	}

	fmt.Printf("Received: %s (v%d)\n", projectName, version)
	for envName := range envMap {
		fmt.Printf("  - %s\n", envName)
	}
	return nil
}

func saveReceived(projectName string, envMap map[string][]types.EnvValue) (int, error) {
	project, err := filehandler.FindProjectByName(projectName)
	if err != nil {
		return 0, err
	}

	var version int

	if project == nil {
		newProject := types.ProjectDir{
			ProjectName:    projectName,
			LocalPath:      "",
			CurrentEnv:     "",
			CurrentVersion: 1,
			LatestVersion:  1,
			VersionNames:   make(map[string]string),
		}
		if err := filehandler.UpsertProjectDir(newProject); err != nil {
			return 0, err
		}
		version = 1
		fmt.Printf("New project: %s (use 'swapenv map %s' to assign directory)\n", projectName, projectName)
	} else {
		version, err = filehandler.BumpVersion(projectName)
		if err != nil {
			return 0, err
		}
	}

	versionPath, err := filehandler.GetVersionFilePath(projectName, version)
	if err != nil {
		return 0, err
	}

	homeDir, err := filehandler.GetHomeDirectory(projectName)
	if err != nil {
		return 0, err
	}

	localDir := ""
	if project != nil {
		localDir = project.LocalPath
	}

	now := time.Now().UTC().Unix()
	projectData := types.Project{
		Id:             uuid.New().String(),
		Name:           projectName,
		Owner:          "received",
		LocalDirectory: localDir,
		CreatedAt:      now,
		ModifiedAt:     now,
		Envs:           envMap,
	}

	projectJSON, err := projectData.MarshalJSON()
	if err != nil {
		return 0, err
	}

	if err := filehandler.WriteProject(homeDir, versionPath, projectJSON); err != nil {
		return 0, err
	}

	return version, nil
}
