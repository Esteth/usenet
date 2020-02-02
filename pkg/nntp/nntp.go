package nntp

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/textproto"
	"os"

	"github.com/esteth/usenet/pkg/yenc"
)

// Conn represents an NNTP connection
type Conn struct {
	*textproto.Conn
}

// Dial will establish a connection to an NNTP server.
func Dial(address string) (*Conn, error) {
	conn := new(Conn)
	var err error

	conn.Conn, err = textproto.Dial("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to %s: %w", address, err)
	}

	_, _, err = conn.ReadCodeLine(20)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("Could not read 20X while establishing connection: %w", err)
	}

	return conn, nil
}

// DialTLS will establish a TLS connection to an NNTP server.
func DialTLS(address string) (*Conn, error) {
	conn := new(Conn)
	tlsConn, err := tls.Dial("tcp", address, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to establish TLS connection: %w", err)
	}
	conn.Conn = textproto.NewConn(tlsConn)

	_, _, err = conn.ReadCodeLine(20)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("Could not read 20X while establshing TLS connection: %w", err)
	}

	return conn, nil
}

// Authenticate will authenticate with the server using the given username and password
func (conn *Conn) Authenticate(user string, pass string) error {
	id, err := conn.Cmd("AUTHINFO USER %s", user)
	if err != nil {
		return fmt.Errorf("Failed to send username: %w", err)
	}
	conn.StartResponse(id)
	code, _, err := conn.ReadCodeLine(381)
	conn.EndResponse(id)

	switch code {
	case 481, 482, 502:
		return fmt.Errorf("Authenticate out of sequence or command not available: %w", err)
	case 281:
		// Authenticated without password.
		return nil
	case 381:
		// Send Password
		break
	default:
		return fmt.Errorf("Failed reading 381 while authenticating: %w", err)
	}
	id, err = conn.Cmd("AUTHINFO PASS %s", pass)
	if err != nil {
		return fmt.Errorf("Failed to send password: %w", err)
	}
	conn.StartResponse(id)
	code, _, err = conn.ReadCodeLine(281)
	conn.EndResponse(id)
	if err != nil {
		return fmt.Errorf("Failed reading 281 while authenticating: %w", err)
	}

	return nil
}

// ReadMessage will return a Reader onto the body of a message
func (conn *Conn) ReadMessage(messageID string) (io.Reader, error) {
	id, err := conn.Cmd("BODY <%s>", messageID)
	conn.StartResponse(id)
	defer conn.EndResponse(id)
	if err != nil {
		return nil, fmt.Errorf("BODY command failed: %w", err)
	}

	_, _, err = conn.ReadCodeLine(222)
	if err != nil {
		return nil, fmt.Errorf("Could not read 222: %w", err)
	}
	return conn.DotReader(), nil
}

// ReadMessageToFile downloads and writes the appropriate segment of the file the message represents.
// Returns the number of bytes written to the file, or an error.
func (conn *Conn) ReadMessageToFile(messageID string) (int64, error) {
	reader, err := conn.ReadMessage(messageID)
	if err != nil {
		return 0, fmt.Errorf("Could not read message: %w", err)
	}
	yencReader, err := yenc.NewReader(reader)
	if err != nil {
		return 0, fmt.Errorf("Could not create reader: %w", err)
	}

	filename, err := yencReader.Filename()
	if err != nil {
		return 0, fmt.Errorf("Could not get filename: %w", err)
	}

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
	defer file.Close()
	if err != nil {
		return 0, fmt.Errorf("Could not open output file: %w", err)
	}

	offset, err := yencReader.Offset()
	if err != nil {
		return 0, fmt.Errorf("Could not read offset from file: %w", err)
	}
	file.Seek(offset, 0)

	bytesWritten, err := io.Copy(file, yencReader)
	if err != nil {
		return 0, fmt.Errorf("Could not copy data to file: %w", err)
	}

	return bytesWritten, nil
}
