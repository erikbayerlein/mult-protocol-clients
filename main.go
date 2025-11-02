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
	sc "github.com/erikbayerlein/mult-protocol-clients/strings"
)

const usageText = `
Commands:
  help                              Show this help
  clear                             Clear terminal screen
  login <client> <student_id>       Authenticate user to server and save
  whoami                            Show current logged user and client
  logout                            Logout and clear token
  string <operation> [args...]      Run operation with string client
  json <operation> [args...]      	Run operation with json client (TODO)
  protobuff <operation> [args...]   Run operation with protobuff client (TODO)
  exit / quit                       Exit program

Operations (client):
  echo <text>            Echo server - Ex.: "Hello World!"
  sum  <n1,n2,...>       Sum a list of numbers - Ex.: "340,558"
  timestamp              Info about the server's time
  status                 Info about the server's status
  history [limit]        History of operations (optional limit, default 10)
`

const (
	host           = "54.174.195.77"
	string_port    = 8080
	json_port      = 8081
	protobuff_port = 8082
)

var (
	currentClient = ""

	string_client = sc.StringClient{
		Host: host,
		Port: string_port,
	}
)

func clearScreen() {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	default:
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	}
}

func gracefulShutdown() {
	rec, err := auth.LoadToken()
	if err == nil && rec.Token != "" {

		switch currentClient {
		case "string":
			_ = string_client.Logout(rec.Token)

		case "json":

		case "protobuff":

		}

		_ = auth.ClearToken()
		fmt.Println("Logged out")
	}
}

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
	fmt.Println("Go Multiprotocol Clients")
	fmt.Println("(type 'help' for commands, 'exit' to quit)")
	fmt.Print(usageText)

	for {
		prompt := "> "
		if currentClient != "" {
			prompt = currentClient + " > "
		}
		fmt.Print("\n" + prompt)

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

		case "clear":
			clearScreen()

		case "login":
			if len(args) < 2 {
				fmt.Println("Usage: login <client> <student_id>")
				continue
			}
			clientArg := strings.ToLower(args[0])
			studentID, err := strconv.Atoi(args[1])
			if err != nil || studentID <= 0 {
				fmt.Printf("Invalid student_id: %q\n", args[1])
				continue
			}

			switch clientArg {
			case "string":
				if err := string_client.Login(studentID); err != nil {
					fmt.Println("Login failed:", err)
					continue
				}
				currentClient = clientArg
				fmt.Printf("Logged in on %s server as student_id=%d\n", currentClient, studentID)

			case "json":
				fmt.Println("JSON client not implemented yet. (TODO)")
			case "protobuff":
				fmt.Println("ProtoBuff client not implemented yet. (TODO)")
			default:
				fmt.Printf("Invalid client: %s\nUse: string | json | protobuff\n", clientArg)
			}

		case "whoami":
			rec, err := auth.LoadToken()
			if err != nil || rec.Token == "" {
				fmt.Println("Not logged in.")
			} else if currentClient == "" {
				fmt.Printf("Logged in as aluno_id=%d (no client selected in this session)\n", rec.StudentId)
			} else {
				fmt.Printf("Logged in as aluno_id=%d on client=%s\n", rec.StudentId, currentClient)
			}

		case "logout":
			rec, err := auth.LoadToken()
			if rec.Token == "" || err != nil {
				fmt.Println("You're not logged.")
				continue
			}
			switch currentClient {
			case "string":
				if err := string_client.Logout(rec.Token); err != nil {
					fmt.Println("Logout error:", err)
				}
			case "json":
				fmt.Println("Implement. TODO")
			case "protobuff":
				fmt.Println("Implement. TODO")
			}
			_ = auth.ClearToken()
			currentClient = ""
			fmt.Println("Logged out.")

		default:
			if currentClient != "" {
				op := cmd
				rest := args
				switch currentClient {
				case "string":
					if err := string_client.Run(op, rest); err != nil {
						fmt.Println("Error:", err)
					}
				case "json":
					fmt.Println("JSON client not implemented yet. (TODO)")
				case "protobuff":
					fmt.Println("ProtoBuff client not implemented yet. (TODO)")
				default:
					fmt.Println("Unknown command:", cmd)
				}
				continue
			}

			fmt.Println("Unknown command:", cmd)
		}
	}
}
