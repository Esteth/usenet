package nntp

import (
	"bufio"
	"net"
	"fmt"
	"testing"
	"net/textproto"
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
				return
			}
			go func(c net.Conn) {
				writer := textproto.NewWriter(bufio.NewWriter(c))
				writer.PrintfLine("%d ", 201)
				c.Close()
			}(conn)
		}
	}()

	conn, err := Dial(fmt.Sprintf(":%d", port))
	if err != nil {
		t.Fatalf("failed to dial: %v", err)
	}

	_, err = conn.ReadMessage("testMessage")
	if err != nil {
		t.Fatalf("failed to read message: %v", err)
	}
}