package cmd_login

import (
	"errors"
	"fmt"
	"time"

	"github.com/reduan2660/swapenv/internal/api"
)

func hyperlink(url string) string {
	return fmt.Sprintf("\x1b]8;;%s\x07%s\x1b]8;;\x07", url, url)
}

func Login(serverURL string) error {
	if api.IsLoggedIn() {
		fmt.Println("already logged in")
		return nil
	}

	client := api.NewClient(serverURL)
	deviceResp, err := client.RequestDeviceCode()
	if err != nil {
		return fmt.Errorf("failed to request device code: %w", err)
	}

	// fmt.Printf("\nOpen %s in your browser\n", deviceResp.VerificationURI)
	fmt.Printf("\nOpen %s in your browser\n", hyperlink(deviceResp.VerificationURI))
	fmt.Printf("Enter code: %s\n\n", deviceResp.UserCode)

	interval := time.Duration(deviceResp.Interval) * time.Second
	timeout := time.After(time.Duration(deviceResp.ExpiresIn) * time.Second)

	for {

		select {
		case <-timeout:
			return fmt.Errorf("login timed out")
		case <-time.After(interval):
			authResp, err := client.PollAuth(deviceResp.DeviceCode)

			if errors.Is(err, fmt.Errorf("slow_down")) {
				interval += 5 * time.Second
				continue
			}

			if err != nil {
				return fmt.Errorf("auth failed: %w", err)
			}

			if authResp == nil {
				continue // pending
			}

			creds := &api.Credentials{
				Token:     authResp.Token,
				UserId:    authResp.UserID,
				OrgId:     authResp.OrgID,
				ExpiresAt: authResp.ExpiresAt,
			}

			if err := api.SaveCredentials(creds); err != nil {
				return fmt.Errorf("failed to save credentials: %w", err)
			}

			fmt.Println("logged in successfully")
			return nil
		}
	}
}
