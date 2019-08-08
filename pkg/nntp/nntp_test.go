package nntp

import (
	"net"
	"testing"
)

func TestConnect(t *testing.T) {
	server, err := net.Listen("tcp", ":9108")
	if err != nil {
		t.Fatalf("could not listen: %v", err)
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

	_, err = Dial(":9108")
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}
}