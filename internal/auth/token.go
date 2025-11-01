package auth

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type TokenRecord struct {
	StudentId int    `json:"student_id"`
	Token     string `json:"token"`
}

func TokenFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".goclient")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	return filepath.Join(dir, "token.json"), nil
}

func SaveToken(rec TokenRecord) error {
	path, err := TokenFilePath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

func LoadToken() (TokenRecord, error) {
	path, err := TokenFilePath()
	if err != nil {
		return TokenRecord{}, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return TokenRecord{}, err
	}
	var rec TokenRecord
	if err := json.Unmarshal(data, &rec); err != nil {
		return TokenRecord{}, err
	}
	if rec.Token == "" || rec.StudentId == 0 {
		return TokenRecord{}, errors.New("invalid token file")
	}
	return rec, nil
}

func ClearToken() error {
	path, err := TokenFilePath()
	if err != nil {
		return err
	}
	_ = os.Remove(path)
	return nil
}
