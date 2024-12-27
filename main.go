package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mdlayher/vsock"
)

func main() {
	// Load environment variables
	enclaveCID := os.Getenv("CID")
	if enclaveCID == "" {
		enclaveCID = "16" // Default CID for the enclave
	}

	port := 5000 // Port for communication with the enclave

	// Parse CID as a uint32
	var cid uint32
	_, err := fmt.Sscanf(enclaveCID, "%d", &cid)
	if err != nil {
		log.Fatalf("Invalid CID: %v", err)
	}

	// Periodically ping the enclave every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	fmt.Printf("Starting periodic pings to Nitro Enclave (CID: %d, Port: %d)...\n", cid, port)

	for range ticker.C {
		pingEnclave(cid, uint32(port))
	}
}

// pingEnclave connects to the Nitro Enclave, sends a ping, and reads the response
func pingEnclave(cid uint32, port uint32) {
	config := &vsock.Config{}
	// Connect to the Nitro Enclave via vsock
	conn, err := vsock.Dial(cid, port, config)
	if err != nil {
		log.Printf("Failed to connect to enclave: %v", err)
		return
	}
	defer conn.Close()

	// Send a ping message to the enclave
	message := "ping"
	fmt.Printf("Sending message: %s\n", message)
	_, err = conn.Write([]byte(message + "\n"))
	if err != nil {
		log.Printf("Failed to send data to enclave: %v", err)
		return
	}

	// Read the response from the enclave
	reader := bufio.NewReader(conn)
	response, err := reader.ReadString('\n')
	if err != nil {
		log.Printf("Failed to read response from enclave: %v", err)
		return
	}

	fmt.Printf("Response from enclave: %s\n", response)
}
