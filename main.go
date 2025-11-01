package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"syscall"

	"github.com/erikbayerlein/mult-protocol-clients/internal/auth"
	"github.com/erikbayerlein/mult-protocol-clients/internal/tcp"
)

const usageText = `
Commands:
  help                              Show this help
  clear                             Clear terminal screen
  login <aluno_id>                  Authenticate user and save token
  whoami                            Show current logged user
  logout                            Logout and clear token
  string <operation> [args...]      Run operation with string client
  exit / quit                       Exit program

Operations (string client):
  echo <text>            Echo server - Ex.: "Hello World!"
  sum  <n1,n2,...>       Sum a list of numbers - Ex.: "340,558"
  timestamp              Info about the server's time
  status                 Info about the server's status
  history [limit]        History of operations (optional limit, default 10)
`

func doOperation(op, token string, params map[string]interface{}) (string, error) {
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
	return tcp.Request(message)
}

func clearScreen() {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func gracefulShutdown() {
	rec, err := auth.LoadToken()
	if err == nil && rec.Token != "" {
		auth.LogoutRemote(rec.Token)
		_ = auth.ClearToken()
		fmt.Println("Logged out")
	}
}

// Command handlers

func runStringClient(op string, args []string) error {
	rec, err := auth.RequireLogin()
	if err != nil {
		return err
	}
	token := rec.Token

	switch op {
	case "echo":
		if len(args) < 1 {
			return fmt.Errorf("echo requires a message")
		}
		resp, err := doOperation("echo", token, map[string]interface{}{"mensagem": strings.Join(args, " ")})
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
		resp, err := doOperation("soma", token, map[string]interface{}{"numeros": ints})
		fmt.Println("→", resp)
		return err

	case "timestamp":
		resp, err := doOperation("timestamp", token, map[string]interface{}{})
		fmt.Println("→", resp)
		return err

	case "status":
		resp, err := doOperation("status", token, map[string]interface{}{"detalhado": true})
		fmt.Println("→", resp)
		return err

	case "history":
		limit := 10
		if len(args) >= 1 {
			if v, convErr := strconv.Atoi(args[0]); convErr == nil {
				limit = v
			}
		}
		resp, err := doOperation("historico", token, map[string]interface{}{"limite": limit})
		fmt.Println("→", resp)
		return err

	default:
		return fmt.Errorf("unknown operation: %s", op)
	}
}

// ===== Main interactive loop =====
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go func() {
		<-ctx.Done()
		fmt.Println()
		gracefulShutdown()
		os.Exit(0)
	}()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("GoClient Interactive Shell")
	fmt.Println("(type 'help' for commands, 'exit' to quit)")
	fmt.Print(usageText)

	for {
		fmt.Print("\n> ")
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if line == "exit" || line == "quit" {
			gracefulShutdown()
			break
		}

		parts := strings.Fields(line)
		cmd := parts[0]
		args := parts[1:]

		switch cmd {
		case "help":
			fmt.Print(usageText)

		case "login":
			if len(args) < 1 {
				fmt.Println("Usage: login <student_id>")
				continue
			}
			student_id, _ := strconv.Atoi(args[0])
			token, err := auth.Auth(student_id)
			if err != nil {
				fmt.Println("Login failed:", err)
				continue
			}
			auth.SaveToken(auth.TokenRecord{StudentId: student_id, Token: token})
			fmt.Println("Logged in with aluno_id:", student_id)

		case "whoami":
			rec, err := auth.LoadToken()
			if err != nil {
				fmt.Println("Not logged in.")
			} else {
				fmt.Printf("Logged in as aluno_id=%d\n", rec.StudentId)
			}

		case "clear":
			clearScreen()

		case "logout":
			rec, err := auth.LoadToken()
			if err == nil && rec.Token != "" {
				auth.LogoutRemote(rec.Token)
			}
			auth.ClearToken()
			fmt.Println("Logged out.")

		case "string":
			if len(args) < 1 {
				fmt.Println("Usage: string <operation> [args...]")
				continue
			}
			op := args[0]
			rest := args[1:]
			if err := runStringClient(op, rest); err != nil {
				fmt.Println("Error:", err)
			}

		default:
			fmt.Println("Unknown command:", cmd)
		}
	}
}
