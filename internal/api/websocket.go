package api

import (
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

// ConnectWS establishes authenticated websocket connection
func ConnectWS(baseURL, path, token string) (*websocket.Conn, error) {
	wsURL := strings.Replace(baseURL, "https://", "wss://", 1)
	wsURL = strings.Replace(wsURL, "http://", "ws://", 1)
	wsURL = wsURL + path

	header := http.Header{}
	header.Set("Authorization", "Bearer "+token)

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, header)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
