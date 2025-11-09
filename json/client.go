package json

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/erikbayerlein/mult-protocol-clients/internal/auth"
	"github.com/erikbayerlein/mult-protocol-clients/internal/tcp"
)

type JsonClient struct {
	Host string
	Port int
}

func (jc *JsonClient) Login(studentId int) error {
	authReq := Auth{
		Type:      "autenticar",
		StudentId: strconv.Itoa(studentId),
	}

	payload, err := json.Marshal(authReq)
	if err != nil {
		return fmt.Errorf("marshal auth request: %w", err)
	}

	token, err := auth.Auth(string(payload), jc.Host, jc.Port)
	if err != nil {
		return err
	}

	if err := auth.SaveToken(auth.TokenRecord{StudentId: studentId, Token: token}); err != nil {
		return fmt.Errorf("could not save token: %w", err)
	}
	return nil
}

func (jc *JsonClient) Logout(token string) error {
	logoutReq := Logout{
		Type:  "logout",
		Token: token,
	}

	payload, err := json.Marshal(logoutReq)
	if err != nil {
		return fmt.Errorf("marshal logout request: %w", err)
	}

	req := string(payload)
	fmt.Printf("Sending request: %s\n", req)
	return auth.LogoutRemote(req, jc.Host, jc.Port)
}

func (jc *JsonClient) Run(op string, args []string) error {
	rec, err := auth.RequireLogin()
	if err != nil {
		return err
	}
	token := rec.Token

	var params any
	switch op {
	case "echo":
		if len(args) < 1 {
			return fmt.Errorf("echo requires a message")
		}
		params = EchoParams{Message: strings.Join(args, " ")}

	case "sum":
		if len(args) < 1 {
			return fmt.Errorf("sum requires a comma-separated list")
		}
		parts := strings.Split(args[0], ",")
		nums := make([]int, 0, len(parts))
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}
			v, conv := strconv.Atoi(p)
			if conv != nil {
				return fmt.Errorf("invalid number %q", p)
			}
			nums = append(nums, v)
		}
		params = SumParams{Numeros: nums}

	case "timestamp":
		params = TimestampParams{}

	case "status":
		params = StatusParams{Detalhado: true}

	case "history":
		limit := 10
		if len(args) >= 1 {
			if v, conv := strconv.Atoi(args[0]); conv == nil && v > 0 {
				limit = v
			}
		}
		params = HistoryParams{Limite: limit}

	default:
		return fmt.Errorf("unknown operation: %s", op)
	}

	switch op {
	case "sum":
		op = "soma"
	case "history":
		op = "historico"
	}

	body := Operation{
		Type:      "operacao",
		Operation: op,
		Token:     token,
		Params:    params,
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("marshal operation request: %w", err)
	}

	req := string(payload)
	fmt.Printf("Sending request: %s\n", req)

	resp, err := tcp.Request(req, jc.Host, jc.Port)
	if err != nil {
		return err
	}

	fmt.Printf("Received: %s\n", resp)
	return nil
}
