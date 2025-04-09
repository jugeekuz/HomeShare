package network

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

type ApifyResponse struct {
	PublicIp 		string `json:"ip"`
}

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

	var apifyResponse ApifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&apifyResponse); err != nil {
		return "", err
	}

	return apifyResponse.PublicIp, nil
}