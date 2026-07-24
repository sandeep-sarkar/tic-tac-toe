package game

import "errors"

type Board struct {
	boardSize int
	gameBoard [][]byte
	turn      byte
	winner    byte
	status    string
	turnCount int
}

func NewBoard(boardSize int) (*Board, error) {
	if boardSize < 1 {
		return nil, errors.New("board size must be greater than 0")
	}
	gameBoard := make([][]byte, boardSize)
	for i := 0; i < boardSize; i++ {
		gameBoard[i] = make([]byte, boardSize)
	}

	return &Board{
		boardSize: boardSize,
		gameBoard: gameBoard,
		turn:      'X',
		status:    StatusReady,
	}, nil
}

func (b *Board) StartGame() {
	if b.status == StatusReady {
		b.status = StatusInProgress
	}
}

func (b *Board) GetStatus() string {
	return b.status
}

func (b *Board) GetWinner() string {
	if b.winner == 0 {
		return ""
	}
	return string(b.winner)
}

func (b *Board) GetTurn() string {
	return string(b.turn)
}

func (b *Board) switchTurn() {
	if b.turn == 'X' {
		b.turn = 'O'
	} else {
		b.turn = 'X'
	}
}

func (b *Board) PlayMove(pos int) error {
	if b.status != StatusInProgress {
		return errors.New("game not in progress")
	}

	if err := b.putMark(pos); err != nil {
		return err
	}

	if b.CheckWinner() {
		b.winner = b.turn
		b.status = StatusWon
		return nil
	}

	if b.turnCount >= (b.boardSize * b.boardSize) {
		b.status = StatusDraw
		return nil
	}

	b.switchTurn()

	return nil
}

// positionToCoordinates converts a board position to row & column coordinates.
func (b *Board) positionToCoordinates(pos int) (int, int, error) {
	maxBoardPos := b.boardSize * b.boardSize
	if pos < 1 || pos > maxBoardPos {
		return 0, 0, errors.New("position out of range")
	}
	pos = pos - 1

	row := pos / b.boardSize
	col := pos % b.boardSize
	return row, col, nil
}

// validatePos validates whether the position can be used.
func (b *Board) validatePos(pos int) error {
	maxBoardPos := b.boardSize * b.boardSize
	row, col, err := b.positionToCoordinates(pos)
	if err != nil {
		return err
	}

	if b.gameBoard[row][col] != 0 {
		return errors.New("position already marked")
	}

	if b.turnCount >= maxBoardPos {
		return errors.New("maximum moves already played")
	}
	return nil
}

// putMark places the current player's mark on the board.
func (b *Board) putMark(pos int) error {
	if err := b.validatePos(pos); err != nil {
		return err
	}

	row, col, err := b.positionToCoordinates(pos)
	if err != nil {
		return err
	}

	b.gameBoard[row][col] = b.turn

	b.turnCount++
	return nil
}

// CheckWinner checks whether current player has won.
func (b *Board) CheckWinner() bool {

	// Check each row
	won := true
	for row := 0; row < b.boardSize; row++ {
		won = true
		for col := 0; col < b.boardSize; col++ {
			if b.gameBoard[row][col] != b.turn {
				won = false
				break
			}
		}
		if won {
			return true
		}
	}

	// Check each col
	for col := 0; col < b.boardSize; col++ {
		won = true
		for row := 0; row < b.boardSize; row++ {
			if b.gameBoard[row][col] != b.turn {
				won = false
				break
			}
		}

		if won {
			return won
		}
	}

	// Check main diagonal
	won = true
	for i := 0; i < b.boardSize; i++ {
		if b.gameBoard[i][i] != b.turn {
			won = false
			break
		}
	}

	if won {
		return won
	}

	// Check opposite diagonal
	won = true

	for row, col := 0, b.boardSize-1; col >= 0 && row < b.boardSize; row, col = row+1, col-1 {
		if b.gameBoard[row][col] != b.turn {
			won = false
			break
		}
	}

	return won
}
