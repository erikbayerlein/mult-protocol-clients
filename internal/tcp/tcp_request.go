package tcp

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

func Request(message string, host string, port int) (string, error) {
	address := net.JoinHostPort(host, strconv.Itoa(port))

	conn, err := net.Dial("tcp", address)
	if err != nil {
		return "", fmt.Errorf("connection error: %v", err)
	}
	defer conn.Close()

	_, err = fmt.Fprintf(conn, "%s\n", message)
	if err != nil {
		return "", fmt.Errorf("send error: %v", err)
	}

	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil && err.Error() != "EOF" {
		return "", fmt.Errorf("read error: %v", err)
	}

	response = strings.TrimSpace(response)
	return response, nil
}
