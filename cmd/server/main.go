package main

import (
	"flag"
	"log"
	"net"

	"github.com/sandeep-sarkar/tic-tac-toe/internal/server"
)

const (
	defaultAddress = ":8080"
	boardSize      = 3
)

func main() {
	address := flag.String("address", defaultAddress, "server listen address")
	flag.Parse()

	listener, err := net.Listen("tcp", *address)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()

	gameServer := server.NewServer(boardSize)
	if err := gameServer.Serve(listener); err != nil {
		log.Fatalf("Server stopped: %v", err)
	}
}
