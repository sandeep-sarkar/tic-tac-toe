package main

import (
	"log"
	"net"

	"github.com/sandeep-sarkar/tic-tac-toe/internal/server"
)

const (
	serverAddress = ":8080"
	boardSize     = 3
)

func main() {
	listener, err := net.Listen("tcp", serverAddress)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()

	gameServer := server.NewServer(boardSize)
	if err := gameServer.Serve(listener); err != nil {
		log.Fatalf("Server stopped: %v", err)
	}
}
