// testing implementation
package netcat

import (
	"log"
	"net"
	"testing"
	"time"
)

// Global var srv for Server struct
var srv Server

func init() {
	// Start the new Server
	srv, err := NewServer("tcp", ":8080")
	if err != nil {
		log.Println("error starting TCP server")
		return
	}

	// Run the server in goroutine to stop blocking
	go func() {
		srv.Run()
	}()
}

// TestNetServerRunnning simply check that the
// test server is up and can accept connections
func TestNetServerRunning(t *testing.T) {
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		t.Error("could not connect to server: ", err)
	}
	defer conn.Close()
}

// TestNetcat runs the main Netcat routine against
// dummy server defined in this test case
func TestNetcat(t *testing.T) {
	// server represents textual properties of our dummy server
	server := target{
		"tcp",
		"localhost",
		8080,
	}
	timeout := 5000 * time.Millisecond
	err := Netcat(&server, timeout, false)
	if err != nil {
		t.Error("could not netcat to server: ", err)
	}
}
