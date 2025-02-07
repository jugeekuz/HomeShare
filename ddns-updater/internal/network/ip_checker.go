package network

import (
	"errors"
	"io"
	"net/http"
	"os"
)

func GetPublicIp() (string, error) {
	pingUrl := os.Getenv("PING_URL")

	resp, err := http.Get(pingUrl)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to fetch IP address")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	ip := string(body)
	return ip, nil
}