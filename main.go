package main

import (
	"log"
	"time"

	"github.com/linuxkit/virtsock/pkg/vsock"
)

func main() {
	const (
		cid  = 3    // Replace with the Nitro Enclave CID
		port = 5000 // Port to connect to
	)

	for {
		conn, err := vsock.Dial(cid, port)
		if err != nil {
			log.Printf("Failed to connect to vsock server: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		log.Println("Connected to vsock server")

		message := "hello"
		_, err = conn.Write([]byte(message))
		if err != nil {
			log.Printf("Failed to send message: %v", err)
			conn.Close()
			time.Sleep(5 * time.Second)
			continue
		}

		log.Printf("Sent message: %s", message)

		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			log.Printf("Failed to read response: %v", err)
			conn.Close()
			time.Sleep(5 * time.Second)
			continue
		}

		log.Printf("Received response: %s", string(buffer[:n]))
		conn.Close()

		time.Sleep(5 * time.Second)
	}
}
