package tcp

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
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

func RequestBytes(payload []byte, host string, port int) ([]byte, error) {
	address := net.JoinHostPort(host, strconv.Itoa(port))

	if sharedConn == nil || sharedConn.conn == nil {
		c, err := net.Dial("tcp", address)
		if err != nil {
			return nil, fmt.Errorf("connection error: %v", err)
		}
		sharedConn = &Connection{conn: c}
	}

	sharedConn.mu.Lock()
	defer sharedConn.mu.Unlock()

	w := bufio.NewWriter(sharedConn.conn)
	r := bufio.NewReader(sharedConn.conn)

	var hdr [4]byte
	binary.BigEndian.PutUint32(hdr[:], uint32(len(payload)))
	if _, err := w.Write(hdr[:]); err != nil {
		return nil, fmt.Errorf("send header error: %v", err)
	}
	if _, err := w.Write(payload); err != nil {
		return nil, fmt.Errorf("send payload error: %v", err)
	}
	if err := w.Flush(); err != nil {
		return nil, fmt.Errorf("flush error: %v", err)
	}

	if _, err := io.ReadFull(r, hdr[:]); err != nil {
		if !strings.Contains(err.Error(), "EOF") {
			return nil, fmt.Errorf("read header error: %v", err)
		}
		return nil, err
	}

	n := binary.BigEndian.Uint32(hdr[:])
	if n == 0 {
		return []byte{}, nil
	}

	resp := make([]byte, n)
	if _, err := io.ReadFull(r, resp); err != nil {
		if !strings.Contains(err.Error(), "EOF") {
			return nil, fmt.Errorf("read body error: %v", err)
		}
		return nil, err
	}

	return resp, nil
}

func Close() error {
	if sharedConn == nil || sharedConn.conn == nil {
		return nil
	}
	err := sharedConn.conn.Close()
	sharedConn.conn = nil
	return err
}
