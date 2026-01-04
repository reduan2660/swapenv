package cmd_share

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/reduan2660/swapenv/internal/api"
	"github.com/reduan2660/swapenv/internal/cmd_loader"
	"github.com/reduan2660/swapenv/internal/cmd_login"
	"github.com/reduan2660/swapenv/internal/cmd_logout"
	"github.com/reduan2660/swapenv/internal/crypto"
	"github.com/reduan2660/swapenv/internal/filehandler"
	"github.com/reduan2660/swapenv/internal/types"
)

type wsMessage struct {
	Type    string `json:"type"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func Share(serverURL, projectName, envName, versionStr string) error {
	if !api.IsLoggedIn() {
		fmt.Println("Not logged in. Logging in...")
		if err := cmd_login.Login(serverURL); err != nil {
			return err
		}
	}

	creds, err := api.LoadCredentials()
	if err != nil {
		cmd_logout.Logout()
		return fmt.Errorf("failed to load credentials: %w. we've logged you out, try loggin in again", err)
	}

	if projectName == "" {
		name, _, _, _, _, err := cmd_loader.GetBasicInfo(cmd_loader.GetBasicInfoOptions{ReadOnly: true})
		if err != nil {
			return err
		}
		if name == "" {
			return fmt.Errorf("no project in current directory, use --project to specify")
		}
		projectName = name
	}

	project, err := filehandler.FindProjectByName(projectName)
	if err != nil {
		return fmt.Errorf("error finding project: %w", err)
	}
	if project == nil {
		return fmt.Errorf("project '%s' not found", projectName)
	}

	version, err := filehandler.ResolveVersion(projectName, versionStr)
	if err != nil {
		return err
	}

	projectPath, err := filehandler.GetVersionFilePath(projectName, version)
	if err != nil {
		return err
	}

	envNames, err := filehandler.ListProjectEnv(projectPath)
	if err != nil {
		return fmt.Errorf("failed to list environments: %w", err)
	}

	if len(envNames) == 0 {
		return fmt.Errorf("no environments found in project")
	}

	if envName != "" {
		if !slices.Contains(envNames, envName) {
			return fmt.Errorf("environment '%s' not found, available: %v", envName, envNames)
		}
		envNames = []string{envName}
	}

	envMap := make(map[string][]types.EnvValue)
	for _, name := range envNames {
		envValues, err := filehandler.ReadProjectEnv(projectPath, name)
		if err != nil {
			return fmt.Errorf("failed to read env '%s': %w", name, err)
		}
		envMap[name] = envValues
	}

	envData, err := json.Marshal(envMap)
	if err != nil {
		return err
	}

	conn, err := api.ConnectWS(serverURL, "/share", creds.Token)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	defer conn.Close()

	var msg wsMessage
	if err := conn.ReadJSON(&msg); err != nil {
		return err
	}

	if msg.Type != "waiting" {
		return fmt.Errorf("unexpected message: %s", msg.Type)
	}

	fmt.Printf("Session stream: %s\n", msg.Code)
	fmt.Printf("Sharing: %s (v%d) - envs: %v\n", projectName, version, envNames)
	fmt.Println("Waiting for receiver...")

	for {
		if err := conn.ReadJSON(&msg); err != nil {
			return err
		}

		switch msg.Type {
		case "ready":
			_, pubKeyBytes, err := conn.ReadMessage()
			if err != nil {
				return err
			}

			pubKey, err := crypto.ParsePublicKey(pubKeyBytes)
			if err != nil {
				return fmt.Errorf("invalid public key: %w", err)
			}

			encrypted, err := crypto.Encrypt(envData, pubKey)
			if err != nil {
				return fmt.Errorf("encryption failed: %w", err)
			}

			nameBytes := []byte(projectName)
			payload := append([]byte{byte(len(nameBytes))}, nameBytes...)
			payload = append(payload, encrypted...)

			encoded := base64.StdEncoding.EncodeToString(payload)
			if err := conn.WriteMessage(1, []byte(encoded)); err != nil {
				return err
			}

			fmt.Println("Environment shared successfully!")

		case "waiting":
			fmt.Printf("Session code: %s\n", msg.Code)
			fmt.Println("Waiting for receiver...")

		case "error":
			return fmt.Errorf("server error: %s", msg.Message)
		}
	}
}
