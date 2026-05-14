package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

const address = ":6379"

func main() {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", address, err)
	}
	defer listener.Close()

	log.Printf("gokv server is listening on %s", address)

	for {
	conn, err := listener.Accept()
	if err != nil {
		log.Printf("Failed to accept connection: %v", err)
		continue
	}	
	 go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	log.Printf("Client connected: %s", conn.RemoteAddr())

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() { 
		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		command := strings.ToUpper(input)
		switch command {
		case "PING":
			fmt.Fprintln(conn, "PONG")
		case "QUIT":
			fmt.Fprintln(conn, "Goodbye!")
			return
		default:
			fmt.Fprintf(conn, "Unknown command: %s\n", input)
		}
		if err := scanner.Err(); err != nil {
			log.Printf("Error reading from client: %v", err)
			return
		}
	}
}
