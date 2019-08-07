package nntp

import (
	"crypto/tls"
	"io"
	"net/textproto"

	"github.com/pkg/errors"
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
		return nil, errors.Wrapf(err, "Failed to connect to %s", address)
	}

	_, _, err = conn.ReadCodeLine(20)
	if err != nil {
		conn.Close()
		return nil, errors.Wrap(err, "Could not read 20 while establishing connection")
	}

	return conn, nil
}

// DialTLS will establish a TLS connection to an NNTP server.
func DialTLS(address string) (*Conn, error) {
	conn := new(Conn)
	tlsConn, err := tls.Dial("tcp", address, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to establish TLS connection")
	}
	conn.Conn = textproto.NewConn(tlsConn)

	_, _, err = conn.ReadCodeLine(20)
	if err != nil {
		conn.Close()
		return nil, errors.Wrap(err, "Could not read 20")
	}

	return conn, nil
}

// Authenticate will authenticate with the server using the given username and password
func (conn *Conn) Authenticate(user string, pass string) error {
	id, err := conn.Cmd("AUTHINFO USER %s", user)
	if err != nil {
		return errors.Wrap(err, "Failed to send username")
	}
	conn.StartResponse(id)
	code, _, err := conn.ReadCodeLine(381)
	conn.EndResponse(id)

	switch code {
	case 481, 482, 502:
		return errors.Wrap(err, "Authenticate out of sequence or command not available")
	case 281:
		// Authenticated without password.
		return nil
	case 381:
		// Send Password
		break
	default:
		return errors.Wrap(err, "Failed reading 381 while authenticating")
	}
	id, err = conn.Cmd("AUTHINFO PASS %s", pass)
	if err != nil {
		return errors.Wrap(err, "Failed to send password")
	}
	conn.StartResponse(id)
	code, _, err = conn.ReadCodeLine(281)
	conn.EndResponse(id)
	return errors.Wrap(err, "Failed reading 281 while authenticating")
}

// ReadMessage will return a Reader onto the body of a message
func (conn *Conn) ReadMessage(messageID string) (io.Reader, error) {
	id, err := conn.Cmd("BODY <%s>", messageID)
	conn.StartResponse(id)
	defer conn.EndResponse(id)
	if err != nil {
		return nil, errors.Wrap(err, "BODY command failed")
	}

	_, _, err = conn.ReadCodeLine(222)
	if err != nil {
		return nil, errors.Wrap(err, "Could not read 222")
	}
	return conn.DotReader(), nil
}
