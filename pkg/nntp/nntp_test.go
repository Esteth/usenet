package nntp

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/textproto"
	"regexp"
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
	defer func() {
		if r := recover(); r != nil {
			t.Fatal(r)
		}
	}()
	go func() {
		for {
			conn, err := server.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				// Reply 201 on connection established.
				// 201 is server ready, posting not allowed.
				writer := textproto.NewWriter(bufio.NewWriter(c))
				writer.PrintfLine("201 ")

				reader := textproto.NewReader(bufio.NewReader(c))
				cmd, err := reader.ReadLine()
				if err != nil && err != io.EOF {
					panic(fmt.Errorf("failed to read command: %w", err))
				}

				if !regexp.MustCompile("BODY <(.+?)>").MatchString(cmd) {
					panic(fmt.Errorf("failed to read BODY command: %w", err))
				}

				// 222 is command expected, response follows.
				writer.PrintfLine("222 ")
				messageWriter := writer.DotWriter()
				io.WriteString(messageWriter, "Sample message content.")
				messageWriter.Close()

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
