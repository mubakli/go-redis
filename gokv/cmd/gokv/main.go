package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

const address = ":6379"

type Store struct {
	mu sync.RWMutex
	data map[string]string
}

func NewStore() *Store {
	return &Store{
		data: make(map[string]string),
	}
}

func (s *Store) Set(key string, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = value
}

func (s *Store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, exists := s.data[key]
	return value, exists
}
 
func main() {

	store := NewStore()

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
	 go handleConnection(conn, store)
	}
}
func handleConnection(conn net.Conn, store *Store) {
	defer conn.Close()

	log.Printf("client connected: %s\n", conn.RemoteAddr())

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		input := strings.TrimSpace(scanner.Text())

		if input == "" {
			continue
		}

		response := handleCommand(input, store)
		fmt.Fprintln(conn, response)
	}

	if err := scanner.Err(); err != nil {
		log.Println("connection error:", err)
	}

	log.Printf("client disconnected: %s\n", conn.RemoteAddr())
}

func handleCommand(input string, store *Store) string {
	parts := strings.Fields(input)

	if len(parts) == 0 {
		return "ERR empty command"
	}

	command := strings.ToUpper(parts[0])

	switch command {
	case "PING":
		return "PONG"

	case "SET":
		if len(parts) < 3 {
			return "ERR usage: SET key value"
		}

		key := parts[1]
		value := parts[2]

		store.Set(key, value)

		return "OK"

	case "GET":
		if len(parts) < 2 {
			return "ERR usage: GET key"
		}

		key := parts[1]

		value, exists := store.Get(key)
		if !exists {
			return "(nil)"
		}

		return value

	case "QUIT":
		return "OK"

	default:
		return "ERR unknown command: " + command
	}
}
