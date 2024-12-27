package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/mdlayher/vsock"
)

const (
	BUFF_SIZE = 4096
	PORT      = 5000
	CID       = 16
)

// readAll reads from the connection until EOF or error
func readAll(conn *vsock.Conn) (string, error) {
	var result strings.Builder
	buffer := make([]byte, BUFF_SIZE)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", fmt.Errorf("read error: %v", err)
		}
		result.Write(buffer[:n])
		if n < BUFF_SIZE {
			break
		}
	}

	return result.String(), nil
}

// server starts a vsock server
func server(port uint32) error {
	// Create vsock listener
	listener, err := vsock.Listen(port, &vsock.Config{})
	if err != nil {
		return fmt.Errorf("failed to create listener: %v", err)
	}
	defer listener.Close()

	log.Printf("Started server on port %d", port)

	for {
		// Accept connection
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Accept failed: %v", err)
			continue
		}

		// Handle each connection in a goroutine
		go func(c *vsock.Conn) {
			defer c.Close()

			addr := c.RemoteAddr().String()
			log.Printf("New connection to addr: %s", addr)

			// Read message
			msg, err := readAll(c)
			if err != nil {
				log.Printf("Failed to read message: %v", err)
				return
			}

			log.Printf("Received: %s", msg)

			// Echo message back
			_, err = c.Write([]byte(msg))
			if err != nil {
				log.Printf("Failed to send response: %v", err)
				return
			}
		}(conn.(*vsock.Conn))
	}
}

// client connects to a vsock server and sends messages
func client(cid, port uint32) error {
	// Connect to the server
	conn, err := vsock.Dial(cid, port, &vsock.Config{})
	if err != nil {
		return fmt.Errorf("connection failed: %v", err)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		// Read input from user
		fmt.Print("Enter message (or 'quit' to exit): ")
		message, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read input: %v", err)
		}

		// Trim whitespace and check for quit command
		message = strings.TrimSpace(message)
		if message == "quit" {
			return nil
		}

		// Send message
		_, err = conn.Write([]byte(message))
		if err != nil {
			return fmt.Errorf("failed to send message: %v", err)
		}

		// Read response
		response, err := readAll(conn)
		if err != nil {
			return fmt.Errorf("failed to read response: %v", err)
		}

		fmt.Printf("Response from server: %s\n", response)
	}
}

func main() {
	// Parse command line arguments
	isServer := flag.Bool("server", false, "Run as server")
	port := flag.Uint("port", PORT, "Port to listen on/connect to")
	cid := flag.Uint("cid", CID, "CID to connect to (client only)")
	flag.Parse()

	var err error
	if *isServer {
		fmt.Println("Starting server...")
		err = server(uint32(*port))
	} else {
		fmt.Printf("Starting client, connecting to CID: %d, Port: %d...\n", *cid, *port)
		err = client(uint32(*cid), uint32(*port))
	}

	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
