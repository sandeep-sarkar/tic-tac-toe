package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

type Client struct {
	conn        net.Conn
	serverInput *bufio.Scanner
	boardSize   int
}

func NewClient(conn net.Conn) *Client {
	return &Client{
		conn:        conn,
		serverInput: bufio.NewScanner(conn),
	}
}

func main() {
	serverAddress := flag.String(
		"server",
		"localhost:8080",
		"server address in the form host:port",
	)

	flag.Parse()

	conn, err := net.Dial("tcp", *serverAddress)

	if err != nil {
		log.Fatalf("Failed to connect to %s:%v", *serverAddress, err)
	}

	defer conn.Close()

	client := NewClient(conn)

	if err := client.readBoardSize(); err != nil {
		log.Fatalf("Failed to read board size: %s", err.Error())
	}

	fmt.Printf("Connected to server. Board Size: %d x %d\n", client.boardSize, client.boardSize)
	fmt.Println("Waiting for another player to connect...")

	go client.readServer()
	client.writeServer()
}

func (c *Client) readBoardSize() error {
	if !c.serverInput.Scan() {
		return errors.New("Server disconnected before sending board size")
	}

	if _, err := fmt.Sscanf(c.serverInput.Text(),
		"boardSize: %d", &c.boardSize); err != nil {
		return errors.New("Invalid board size message")
	}

	if c.boardSize <= 0 {
		return errors.New("Invalid board size")
	}
	return nil
}

func (c *Client) readServer() {
	for c.serverInput.Scan() {
		fmt.Println(c.serverInput.Text())
	}

	if err := c.serverInput.Err(); err != nil {
		log.Printf("Connection error: %s", err)
	}

	fmt.Println("Disconnected from server")
	os.Exit(0)
}

func (c *Client) writeServer() {
	input := bufio.NewScanner(os.Stdin)
	maxPosition := c.boardSize * c.boardSize

	for input.Scan() {
		value := strings.TrimSpace(input.Text())

		position, err := strconv.Atoi(value)
		if err != nil {
			fmt.Println("Please enter a number.")
			continue
		}

		if position < 1 || position > maxPosition {
			fmt.Printf(
				"Please enter a position between 1 and %d.\n",
				maxPosition,
			)
			continue
		}

		if _, err := fmt.Fprintln(c.conn, position); err != nil {
			log.Printf("Failed to send move: %v", err)
			return
		}
	}

	if err := input.Err(); err != nil {
		log.Printf("Failed to read input: %v", err)
	}
}
