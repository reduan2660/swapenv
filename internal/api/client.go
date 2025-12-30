package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

type DeviceCodeResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

type AuthResponse struct {
	Token     string `json:"token"`
	UserID    string `json:"user_id"`
	OrgID     string `json:"org_id"`
	ExpiresAt int64  `json:"expires_at"`
}

func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}
}

func (c *Client) RequestDeviceCode() (*DeviceCodeResponse, error) {

	resp, err := c.HTTPClient.Post(c.BaseURL+"/auth/device/", "application/json", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned %d", resp.StatusCode)
	}

	var result DeviceCodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) PollAuth(deviceCode string) (*AuthResponse, error) {
	body := map[string]string{"device_code": deviceCode}
	jsonBody, _ := json.Marshal(body)

	resp, err := c.HTTPClient.Post(c.BaseURL+"/auth/poll/", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusAccepted {

		var status struct {
			Status string `json:"status"`
		}

		json.NewDecoder(resp.Body).Decode(&status)

		if status.Status == "slow_down" {
			return nil, fmt.Errorf("slow_down")
		}

		return nil, nil

	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned %d", resp.StatusCode)
	}

	var result AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
