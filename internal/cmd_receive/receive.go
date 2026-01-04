package cmd_receive

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

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
		return fmt.Errorf("failed to load credentials: %w. we've logged you out, try loggin in again", err)
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

	_, payload, err := conn.ReadMessage()
	if err != nil {
		return err
	}

	encrypted, err := base64.StdEncoding.DecodeString(string(payload))
	if err != nil {
		return fmt.Errorf("invalid payload encoding: %w", err)
	}

	decrypted, err := crypto.Decrypt(encrypted, privKey)
	if err != nil {
		return fmt.Errorf("decryption failed: %w", err)
	}

	var envValues []types.EnvValue
	if err := json.Unmarshal(decrypted, &envValues); err != nil {
		return fmt.Errorf("invalid env data: %w", err)
	}

	outputPath, err := saveReceived(envValues)
	if err != nil {
		return fmt.Errorf("failed to save: %w", err)
	}

	fmt.Printf("Environment saved to %s\n", outputPath)
	return nil

}

// saveReceived stores env data - todo: project mapping
func saveReceived(envValues []types.EnvValue) (string, error) {
	baseDir, err := filehandler.GetBaseDir()
	if err != nil {
		return "", err
	}

	receivedDir := filepath.Join(baseDir, "received")
	if err := os.MkdirAll(receivedDir, 0755); err != nil {
		return "", err
	}

	filename := fmt.Sprintf("%d.json", time.Now().Unix())
	outputPath := filepath.Join(receivedDir, filename)

	data, err := json.MarshalIndent(envValues, "", "  ")
	if err != nil {
		return "", err
	}

	if err := os.WriteFile(outputPath, data, 0600); err != nil {
		return "", err
	}

	return outputPath, nil
}
