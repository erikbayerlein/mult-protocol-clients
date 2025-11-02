package auth

import (
	"fmt"
	"regexp"

	"github.com/erikbayerlein/mult-protocol-clients/internal/tcp"
)

func Auth(alunoID int, host string, port int) (string, error) {
	authRequest := fmt.Sprintf("AUTH|aluno_id=%d|FIM", alunoID)
	authResponse, err := tcp.Request(authRequest, host, port)
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile(`token=([^|]*)`)
	match := re.FindStringSubmatch(authResponse)
	if len(match) < 2 {
		return "", fmt.Errorf("token not found in AUTH response")
	}
	return match[1], nil
}

func RequireLogin() (TokenRecord, error) {
	rec, err := LoadToken()
	if err != nil {
		return TokenRecord{}, fmt.Errorf("no active session. Please 'login <aluno_id>' first")
	}
	return rec, nil
}

func LogoutRemote(token string, host string, port int) error {
	logoutRequest := fmt.Sprintf("LOGOUT|token=%s|FIM", token)

	_, err := tcp.Request(logoutRequest, host, port)
	if err != nil {
		return err
	}

	return nil
}
