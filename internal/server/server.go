package server

import (
	"fmt"
	"log"
	"net"

	"github.com/sandeep-sarkar/tic-tac-toe/internal/game"
)

type Server struct {
	waiting   chan net.Conn
	boardSize int
}

func NewServer(boardSize int) *Server {
	waiting := make(chan net.Conn)
	return &Server{
		waiting:   waiting,
		boardSize: boardSize,
	}
}

func (s *Server) Serve(listener net.Listener) error {
	go s.matchPlayers()

	log.Printf("Server listening on %s", listener.Addr())

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		log.Printf("Client connected: %s", conn.RemoteAddr())

		if err = s.sendBoardSize(conn); err != nil {
			log.Printf("Failed to send board size to %s: %s\n", conn.RemoteAddr(), err.Error())
			conn.Close()
			continue
		}

		s.waiting <- conn
	}
}

func (s *Server) matchPlayers() {
	for {
		playerX := <-s.waiting
		playerO := <-s.waiting

		go s.runSession(playerX, playerO)
	}
}

func (s *Server) runSession(playerX, playerO net.Conn) {
	log.Printf("Starting game: %s vs %s", playerX.RemoteAddr(), playerO.RemoteAddr())
	defer playerX.Close()
	defer playerO.Close()

	board, err := game.NewBoard(s.boardSize)
	if err != nil {
		log.Printf("Failed to create board %s", err.Error())
		return
	}

	board.StartGame()

	if err := s.notifyPlayers(playerX, playerO); err != nil {
		log.Printf("Failed to notify players: %s", err.Error())
		return
	}

	//run the game loop
	if err := s.runGameLoop(board, playerX, playerO); err != nil {
		log.Printf("Game ended: %s", err.Error())
	}
}

func (s *Server) notifyPlayers(playerX, playerO net.Conn) error {
	if err := s.sendMessage(playerX, "You are X\n"); err != nil {
		return err
	}

	if err := s.sendMessage(playerO, "You are O\n"); err != nil {
		return err
	}

	if err := s.sendMessage(playerX, "Your turn X\n"); err != nil {
		return err
	}

	return s.sendMessage(playerO, "Waiting for player X\n")

}

func (s *Server) sendMessage(conn net.Conn, message string) error {
	_, err := conn.Write([]byte(message))
	return err
}

func (s *Server) sendBoardSize(conn net.Conn) error {
	message := fmt.Sprintf("boardSize: %d\n", s.boardSize)
	_, err := conn.Write([]byte(message))
	return err
}
