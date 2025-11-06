package tcp

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
)

type Connection struct {
	conn net.Conn
	mu   sync.Mutex
}

var sharedConn *Connection

func Request(message string, host string, port int) (string, error) {
	address := net.JoinHostPort(host, strconv.Itoa(port))

	if sharedConn == nil || sharedConn.conn == nil {
		c, err := net.Dial("tcp", address)
		if err != nil {
			return "", fmt.Errorf("connection error: %v", err)
		}
		sharedConn = &Connection{conn: c}
	}

	sharedConn.mu.Lock()
	defer sharedConn.mu.Unlock()

	_, err := fmt.Fprintf(sharedConn.conn, "%s\n", message)
	if err != nil {
		return "", fmt.Errorf("send error: %v", err)
	}

	buff := make([]byte, 64*1024)
	n, err := sharedConn.conn.Read(buff)
	if err != nil && !strings.Contains(err.Error(), "EOF") {
		return "", fmt.Errorf("read error: %v", err)
	}

	trimmedData := bytes.TrimRight(buff[:n], "\x00")
	response := strings.TrimSpace(string(trimmedData))

	return response, nil
}

func Close() error {
	if sharedConn == nil || sharedConn.conn == nil {
		return nil
	}
	err := sharedConn.conn.Close()
	sharedConn.conn = nil
	return err
}
