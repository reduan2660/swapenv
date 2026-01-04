package cmd_share

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/reduan2660/swapenv/internal/api"
	"github.com/reduan2660/swapenv/internal/cmd_loader"
	"github.com/reduan2660/swapenv/internal/cmd_login"
	"github.com/reduan2660/swapenv/internal/cmd_logout"
	"github.com/reduan2660/swapenv/internal/crypto"
	"github.com/reduan2660/swapenv/internal/filehandler"
)

type wsMessage struct {
	Type    string `json:"type"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

func Share(serverURL, envName string) error {
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
	_, _, _, _, projectPath, err := cmd_loader.GetBasicInfo(cmd_loader.GetBasicInfoOptions{ReadOnly: true})
	if err != nil {
		return err
	}

	envValues, err := filehandler.ReadProjectEnv(projectPath, envName)
	if err != nil {
		return fmt.Errorf("failed to read env '%s': %w", envName, err)
	}

	envData, err := json.Marshal(envValues)
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
	fmt.Println("Waiting for receiver...")

	for {
		if err := conn.ReadJSON(&msg); err != nil {
			return err
		}

		switch msg.Type {
		case "ready":
			// read receiver's public key
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

			// send encrypted payload
			payload := base64.StdEncoding.EncodeToString(encrypted)
			if err := conn.WriteMessage(1, []byte(payload)); err != nil {
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
