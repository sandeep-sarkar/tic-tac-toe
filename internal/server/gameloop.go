package server

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/sandeep-sarkar/tic-tac-toe/internal/game"
)

func (s *Server) runGameLoop(board *game.Board, PlayerX, PlayerO net.Conn) error {

	var currentPlayer net.Conn
	var currentReader *bufio.Scanner

	readerX := bufio.NewScanner(PlayerX)
	readerO := bufio.NewScanner(PlayerO)

	if board.GetTurn() == "X" {
		currentPlayer = PlayerX
		currentReader = readerX
	} else {
		currentPlayer = PlayerO
		currentReader = readerO
	}

	if err := s.sendBoard(currentPlayer, board); err != nil {
		return err
	}

	for {
		move, err := s.readMove(currentReader)
		if err != nil {
			if errors.Is(err, errPlayerDisconnected) {
				return err
			}

			if err := s.sendMessage(currentPlayer, err.Error()+"\n"); err != nil {
				return err
			}
			continue
		}

		if err := board.PlayMove(move); err != nil {
			errMsg := fmt.Sprintf("Invalid move: %s\n", err.Error())
			if err := s.sendMessage(currentPlayer, errMsg); err != nil {
				return err
			}
			continue
		}
		if board.GetStatus() == game.StatusWon {
			message := fmt.Sprintf("Player %s wins\n", board.GetWinner())
			if err = s.sendBoardToBoth(PlayerX, PlayerO, board); err != nil {
				return err
			}

			return s.sendToBoth(PlayerX, PlayerO, message)
		}

		if board.GetStatus() == game.StatusDraw {

			if err = s.sendBoardToBoth(PlayerX, PlayerO, board); err != nil {
				return err
			}

			message := "Game ended in a draw\n"
			return s.sendToBoth(PlayerX, PlayerO, message)
		}

		turnMessage := ""
		if board.GetTurn() == "X" {
			currentPlayer = PlayerX
			currentReader = readerX
			turnMessage = "Your turn X\n"
		} else {
			currentPlayer = PlayerO
			currentReader = readerO
			turnMessage = "Your turn O\n"
		}

		if err = s.sendBoard(currentPlayer, board); err != nil {
			return err
		}

		if err = s.sendMessage(currentPlayer, turnMessage); err != nil {
			return err
		}
	}
}

// readMove reads and validates current player's move
func (s *Server) readMove(reader *bufio.Scanner) (int, error) {
	if !reader.Scan() {
		return 0, errPlayerDisconnected
	}

	input := strings.TrimSpace(reader.Text())
	move, err := strconv.Atoi(input)
	if err != nil {
		return 0, errors.New("input must be a number")
	}

	return move, nil
}

// sendBoard sends the current status of the board
func (s *Server) sendBoard(Player net.Conn, board *game.Board) error {
	message := strings.Join(board.GetBoard(), "\n") + "\n"
	return s.sendMessage(Player, message)
}

// sendBoardToBoth sends the current status of the board to both the players
func (s *Server) sendBoardToBoth(PlayerX, PlayerO net.Conn, board *game.Board) error {
	message := strings.Join(board.GetBoard(), "\n") + "\n"
	return s.sendToBoth(PlayerX, PlayerO, message)
}

// sendToBoth sends the same message to both the players
func (s *Server) sendToBoth(playerX, playerO net.Conn, message string) error {
	if err := s.sendMessage(playerX, message); err != nil {
		return err
	}
	return s.sendMessage(playerO, message)
}
