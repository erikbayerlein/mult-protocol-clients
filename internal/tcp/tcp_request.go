package tcp

import (
	"bytes"
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

	buff := make([]byte, 64*1024)
	_, err = conn.Read(buff)
	if err != nil && err.Error() != "EOF" {
		return "", fmt.Errorf("read error: %v", err)
	}

	trimmedData := bytes.TrimRight(buff, "\x00")

	response := strings.TrimSpace(string(trimmedData))
	return response, nil
}
