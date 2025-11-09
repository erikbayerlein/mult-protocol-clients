package strings

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/erikbayerlein/mult-protocol-clients/internal/auth"
	"github.com/erikbayerlein/mult-protocol-clients/internal/tcp"
)

type StringClient struct {
	Host string
	Port int
}

func (sc *StringClient) Login(studentId int) error {
	authRequest := fmt.Sprintf("AUTH|aluno_id=%d|FIM", studentId)
	token, err := auth.Auth(authRequest, sc.Host, sc.Port)
	if err != nil {
		return err
	}
	if err := auth.SaveToken(auth.TokenRecord{StudentId: studentId, Token: token}); err != nil {
		fmt.Println("Could not save token:", err)
		return err
	}
	return nil
}

func (sc *StringClient) Run(op string, args []string) error {
	rec, err := auth.RequireLogin()
	if err != nil {
		return err
	}
	token := rec.Token
	fmt.Println(token)

	switch op {
	case "echo":
		if len(args) < 1 {
			return fmt.Errorf("echo requires a message")
		}
		resp, err := sc.DoOperation("echo", token, map[string]any{"mensagem": strings.Join(args, " ")})
		fmt.Println("→", resp)
		return err

	case "sum":
		if len(args) < 1 {
			return fmt.Errorf("sum requires a comma-separated list")
		}
		parts := strings.Split(args[0], ",")
		ints := make([]int, 0, len(parts))
		for _, p := range parts {
			n, _ := strconv.Atoi(strings.TrimSpace(p))
			ints = append(ints, n)
		}
		resp, err := sc.DoOperation("soma", token, map[string]any{"nums": ints})
		fmt.Println("→", resp)
		return err

	case "timestamp":
		resp, err := sc.DoOperation("timestamp", token, map[string]any{})
		fmt.Println("→", resp)
		return err

	case "status":
		resp, err := sc.DoOperation("status", token, map[string]any{"detalhado": true})
		fmt.Println("→", resp)
		return err

	case "history":
		limit := 10
		if len(args) >= 1 {
			if v, convErr := strconv.Atoi(args[0]); convErr == nil {
				limit = v
			}
		}
		resp, err := sc.DoOperation("historico", token, map[string]any{"limite": limit})
		fmt.Println("→", resp)
		return err

	default:
		return fmt.Errorf("unknown operation: %s", op)
	}
}

func (sc *StringClient) Logout(token string) error {
	req := fmt.Sprintf("LOGOUT|token=%s|FIM", token)
	return auth.LogoutRemote(req, sc.Host, sc.Port)
}

func (sc *StringClient) DoOperation(op, token string, params map[string]any) (string, error) {
	args := []string{"OP", "token=" + token, "operacao=" + op}
	for key, value := range params {
		switch v := value.(type) {
		case []int:
			nums := make([]string, len(v))
			for i, n := range v {
				nums[i] = fmt.Sprintf("%d", n)
			}
			args = append(args, fmt.Sprintf("%s=%s", key, strings.Join(nums, ",")))
		default:
			args = append(args, fmt.Sprintf("%s=%v", key, v))
		}
	}
	args = append(args, "FIM")

	message := strings.Join(args, "|")
	fmt.Println(message)
	return tcp.Request(message, sc.Host, sc.Port)
}
