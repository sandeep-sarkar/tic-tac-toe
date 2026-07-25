package main

import (
	"fmt"
	"net"
	"testing"
)

func TestReadBoardSize(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()
	defer clientConn.Close()

	client := NewClient(clientConn)

	go func() {
		fmt.Fprintln(serverConn, "boardSize: 3")
	}()

	if err := client.readBoardSize(); err != nil {
		t.Fatalf("readBoardSize() returned unexpected error: %v", err)
	}

	if client.boardSize != 3 {
		t.Errorf("boardSize = %d; want 3", client.boardSize)
	}
}

func TestReadBoardSizeRejectsInvalidMessage(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()
	defer clientConn.Close()

	client := NewClient(clientConn)

	go func() {
		fmt.Fprintln(serverConn, "size: 3")
	}()

	err := client.readBoardSize()
	if err == nil {
		t.Fatal("readBoardSize() expected an error")
	}

	if err.Error() != "Invalid board size message" {
		t.Errorf(
			"error = %q; want %q",
			err.Error(),
			"Invalid board size message",
		)
	}
}

func TestReadBoardSizeRejectsNonPositiveSize(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()
	defer clientConn.Close()

	client := NewClient(clientConn)

	go func() {
		fmt.Fprintln(serverConn, "boardSize: 0")
	}()

	err := client.readBoardSize()
	if err == nil {
		t.Fatal("readBoardSize() expected an error")
	}

	if err.Error() != "Invalid board size" {
		t.Errorf(
			"error = %q; want %q",
			err.Error(),
			"Invalid board size",
		)
	}
}

func TestReadBoardSizeHandlesDisconnect(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	client := NewClient(clientConn)

	serverConn.Close()
	defer clientConn.Close()

	err := client.readBoardSize()
	if err == nil {
		t.Fatal("readBoardSize() expected an error")
	}

	if err.Error() != "Server disconnected before sending board size" {
		t.Errorf(
			"error = %q; want disconnect error",
			err.Error(),
		)
	}
}
