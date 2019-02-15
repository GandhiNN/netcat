// Package netcat :
// server.go serves our purpose
// of building dummy target server
package netcat

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strings"
)

// Server interface defines the minimum contract
// our TCP and UDP server implementations must satisfy
type Server interface {
	Run() error
	Close() error
}

// TCPServer holds the structure of
// our TCP implementation
type TCPServer struct {
	addr   string
	server net.Listener
}

// Run starts the TCP Server
func (t *TCPServer) Run() (err error) {
	t.server, err = net.Listen("tcp", t.addr)
	if err != nil {
		return err
	}
	defer t.Close()

	for {
		conn, err := t.server.Accept()
		if err != nil {
			err = errors.New("could not accept connection")
			break
		}
		if conn == nil {
			err = errors.New("could not create connection")
			break
		}
		return t.handleConnections()
	}
	return
}

// handleConnections(plural) is used to accept connections on
// the TCPServer and handle each of them in separate
// goroutines
func (t *TCPServer) handleConnections() (err error) {
	for {
		conn, err := t.server.Accept()
		if err != nil || conn == nil {
			err = errors.New("could not accept connection")
			break
		}
		go t.handleConnection(conn)
	}
	return
}

// handleConnection(singular) deals with the business logic
// of each connection and their requests
func (t *TCPServer) handleConnection(conn net.Conn) {
	defer conn.Close()

	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	for {
		req, err := rw.ReadString('\n')
		if err != nil {
			rw.WriteString("failed to read input")
			rw.Flush()
			return
		}
		rw.WriteString(fmt.Sprintf("Request received: %s", req))
		rw.Flush()
	}
}

// Close shuts down the TCP Server
func (t *TCPServer) Close() (err error) {
	return t.server.Close()
}

// NewServer is a convenience function to create
// a new Server using given protocol and address
func NewServer(protocol, addr string) (Server, error) {
	switch strings.ToLower(protocol) {
	case "tcp":
		return &TCPServer{
			addr: addr,
		}, nil
	case "udp":
	}
	return nil, errors.New("Invalid protocol given")
}
