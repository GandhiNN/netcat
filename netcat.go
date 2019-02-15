// Package netcat :
// netcat.go is a read-only, TCP-only netcat client
// Inspired by https://notes.shichao.io/gopl/ch8/
// and https://devtheweb.io/blog/2018/07/testing-net-connections-part-one/
package netcat

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type target struct {
	IPProto   string
	IPaddress string
	port      int
}

func (t *target) getSocketAddr() (socket string) {
	portStr := strconv.Itoa(t.port)
	return t.IPaddress + ":" + portStr
}

func main() {

	// Parse CLI flags
	protoPtr := flag.String("proto", "tcp", "IP Protocol/Layer 4 connection type")
	ipaddrPtr := flag.String("address", "0.0.0.0", "IPv4 Address")
	portPtr := flag.Int("port", 22, "IP Address port")
	waitPtr := flag.Bool("wait", false, "If set to wait, it will wait for a server-side program to send bytes to the receiving end")
	flag.Parse()

	// Build the target struct
	targetNode := target{
		IPProto:   *protoPtr,
		IPaddress: *ipaddrPtr,
		port:      *portPtr,
	}

	// Setup CTRL-C handler
	setupInterruptHandler()

	// Establish connection
	timeout := 5000 * time.Millisecond
	Netcat(&targetNode, timeout, *waitPtr)
}

// Netcat initiate the netcat
func Netcat(t *target, timeout time.Duration, wait bool) (err error) {
	socket := t.getSocketAddr()
	log.Printf("Starting connection to %s port %d", t.IPaddress, t.port)
	go spinner(100 * time.Millisecond) // visual indicator that the program is running
	conn, err := net.DialTimeout(strings.ToLower(t.IPProto), socket, timeout)
	if err != nil {
		return err
	}
	defer conn.Close()
	log.Printf("Connection opened")
	if wait {
		StdOutCopy(os.Stdout, conn)
	}
	return nil
}

// StdOutCopy is a helper function to copy standard output to our conn object
func StdOutCopy(dst io.Writer, src io.Reader) {
	bWritten, err := io.Copy(dst, src)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("wrote %d", bWritten)
}

// spinner is just a cosmetic function to provide
// the user with a visual indication that program is still running
func spinner(delay time.Duration) {
	for {
		for _, r := range `-\|/-` {
			fmt.Printf("\r%c", r)
			time.Sleep(delay)
		}
	}
}

// setupInterruptHandler catch CTRL-C from the terminal
func setupInterruptHandler() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\rCTRL-C is pressed in Terminal")
		log.Printf("Exiting Program\n")
		os.Exit(0)
	}()
}
