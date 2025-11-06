package auth

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/erikbayerlein/mult-protocol-clients/internal/tcp"
)

type authJSONResp struct {
	Token string `json:"token"`
}

var legacyTokenRe = regexp.MustCompile(`token=([^|]*)`)

func Auth(req string, host string, port int) (string, error) {
	authResponse, err := tcp.Request(req, host, port)
	if err != nil {
		return "", err
	}

	fmt.Printf("Received response: %s\n", authResponse)

	var j authJSONResp
	if err := json.Unmarshal([]byte(authResponse), &j); err == nil && j.Token != "" {
		return j.Token, nil
	}

	if m := legacyTokenRe.FindStringSubmatch(authResponse); len(m) == 2 {
		return m[1], nil
	}

	return "", fmt.Errorf("token not found in response")
}

func LogoutRemote(req string, host string, port int) error {
	_, err := tcp.Request(req, host, port)
	if err != nil {
		return err
	}

	return nil
}

func RequireLogin() (TokenRecord, error) {
	rec, err := LoadToken()
	if err != nil {
		return TokenRecord{}, fmt.Errorf("No active session. Please 'login <aluno_id>' first")
	}
	return rec, nil
}
