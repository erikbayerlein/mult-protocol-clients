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
	token, err := auth.Auth(studentId, sc.Host, sc.Port)
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
		resp, err := sc.doOperation("echo", token, map[string]interface{}{"mensagem": strings.Join(args, " ")})
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
		resp, err := sc.doOperation("soma", token, map[string]interface{}{"numeros": ints})
		fmt.Println("→", resp)
		return err

	case "timestamp":
		resp, err := sc.doOperation("timestamp", token, map[string]interface{}{})
		fmt.Println("→", resp)
		return err

	case "status":
		resp, err := sc.doOperation("status", token, map[string]interface{}{"detalhado": true})
		fmt.Println("→", resp)
		return err

	case "history":
		limit := 10
		if len(args) >= 1 {
			if v, convErr := strconv.Atoi(args[0]); convErr == nil {
				limit = v
			}
		}
		resp, err := sc.doOperation("historico", token, map[string]interface{}{"limite": limit})
		fmt.Println("→", resp)
		return err

	default:
		return fmt.Errorf("unknown operation: %s", op)
	}
}

func (sc *StringClient) Logout(token string) error {
	return auth.LogoutRemote(token, sc.Host, sc.Port)
}

func (sc *StringClient) doOperation(op, token string, params map[string]interface{}) (string, error) {
	args := []string{"OP", "operacao=" + op, "token=" + token}
	for key, value := range params {
		switch v := value.(type) {
		case []int:
			nums := make([]string, len(v))
			for i, n := range v {
				nums[i] = fmt.Sprintf("%d", n)
			}
			args = append(args, fmt.Sprintf("%s=[%s]", key, strings.Join(nums, ",")))
		default:
			args = append(args, fmt.Sprintf("%s=%v", key, v))
		}
	}
	args = append(args, "FIM")

	message := strings.Join(args, "|")
	fmt.Println(message)
	return tcp.Request(message, sc.Host, sc.Port)
}
