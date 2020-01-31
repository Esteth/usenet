package nntp

import (
	"net"
	"fmt"
	"testing"
)

func TestConnect(t *testing.T) {
	server, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("could not listen: %v", err)
	}
	var port int
	_, err = fmt.Sscanf(server.Addr().String(), "127.0.0.1:%d", &port)
	if err != nil {
		t.Fatalf("could not parse port: %v", err)
	}
	defer server.Close()
	go func() {
		for {
			conn, err := server.Accept()
			if err != nil {
				t.Fatalf("failed to accept connection: %v", err)
			}
			go func(c net.Conn) {
				c.Write([]byte(" 20 \r\n"))
				c.Close()
			}(conn)
		}
	}()

	_, err = net.Dial("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
}