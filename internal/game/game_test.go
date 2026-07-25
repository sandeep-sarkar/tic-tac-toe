package game

import "testing"

func newStartedBoard(t *testing.T, size int) *Board {
	t.Helper()

	board, err := NewBoard(size)
	if err != nil {
		t.Fatalf("NewBoard(%d) returned unexpected error: %v", size, err)
	}

	board.StartGame()
	return board
}

func playMoves(t *testing.T, board *Board, moves ...int) {
	t.Helper()

	for _, move := range moves {
		if err := board.PlayMove(move); err != nil {
			t.Fatalf("PlayMove(%d) returned unexpected error: %v", move, err)
		}
	}
}

func TestNewBoard(t *testing.T) {
	board, err := NewBoard(3)
	if err != nil {
		t.Fatalf("NewBoard() returned unexpected error: %v", err)
	}

	if board.boardSize != 3 {
		t.Errorf("boardSize = %d; want 3", board.boardSize)
	}

	if len(board.gameBoard) != 3 {
		t.Fatalf("number of rows = %d; want 3", len(board.gameBoard))
	}

	for row := range board.gameBoard {
		if len(board.gameBoard[row]) != 3 {
			t.Errorf(
				"row %d length = %d; want 3",
				row,
				len(board.gameBoard[row]),
			)
		}
	}

	if board.GetStatus() != StatusReady {
		t.Errorf(
			"status = %q; want %q",
			board.GetStatus(),
			StatusReady,
		)
	}

	if board.GetTurn() != "X" {
		t.Errorf("turn = %q; want %q", board.GetTurn(), "X")
	}

	if board.GetWinner() != "" {
		t.Errorf("winner = %q; want empty string", board.GetWinner())
	}
}

func TestNewBoardRejectsInvalidSize(t *testing.T) {
	board, err := NewBoard(0)

	if err == nil {
		t.Fatal("NewBoard(0) expected an error")
	}

	if board != nil {
		t.Fatal("NewBoard(0) expected a nil board")
	}
}

func TestStartGame(t *testing.T) {
	board, err := NewBoard(3)
	if err != nil {
		t.Fatalf("NewBoard() returned unexpected error: %v", err)
	}

	board.StartGame()

	if board.GetStatus() != StatusInProgress {
		t.Errorf(
			"status = %q; want %q",
			board.GetStatus(),
			StatusInProgress,
		)
	}
}

func TestStartGameDoesNotRestartCompletedGame(t *testing.T) {
	board := newStartedBoard(t, 3)

	playMoves(t, board, 1, 4, 2, 5, 3)

	if board.GetStatus() != StatusWon {
		t.Fatalf("status = %q; want %q", board.GetStatus(), StatusWon)
	}

	board.StartGame()

	if board.GetStatus() != StatusWon {
		t.Errorf(
			"status after StartGame() = %q; want %q",
			board.GetStatus(),
			StatusWon,
		)
	}
}

func TestPlayMoveBeforeStartReturnsError(t *testing.T) {
	board, err := NewBoard(3)
	if err != nil {
		t.Fatalf("NewBoard() returned unexpected error: %v", err)
	}

	err = board.PlayMove(1)
	if err == nil {
		t.Fatal("PlayMove() before StartGame() expected an error")
	}

	if err.Error() != "game not in progress" {
		t.Errorf("error = %q; want %q", err.Error(), "game not in progress")
	}
}

func TestPlayMovePlacesMarkAndSwitchesTurn(t *testing.T) {
	board := newStartedBoard(t, 3)

	if err := board.PlayMove(1); err != nil {
		t.Fatalf("PlayMove(1) returned unexpected error: %v", err)
	}

	if board.gameBoard[0][0] != 'X' {
		t.Errorf("position 1 = %q; want X", board.gameBoard[0][0])
	}

	if board.GetTurn() != "O" {
		t.Errorf("turn = %q; want O", board.GetTurn())
	}

	if board.turnCount != 1 {
		t.Errorf("turnCount = %d; want 1", board.turnCount)
	}
}

func TestPlayMoveRejectsOutOfRangePosition(t *testing.T) {
	testPos := 11

	board := newStartedBoard(t, 3)

	err := board.PlayMove(testPos)
	if err == nil {
		t.Fatalf("PlayMove(%d) expected an error", testPos)
	}

	if err.Error() != "position out of range" {
		t.Errorf(
			"error = %q; want %q",
			err.Error(),
			"position out of range",
		)
	}
}

func TestPlayMoveRejectsOccupiedPosition(t *testing.T) {
	board := newStartedBoard(t, 3)

	if err := board.PlayMove(1); err != nil {
		t.Fatalf("first PlayMove(1) returned unexpected error: %v", err)
	}

	err := board.PlayMove(1)
	if err == nil {
		t.Fatal("second PlayMove(1) expected an error")
	}

	if err.Error() != "position already marked" {
		t.Errorf(
			"error = %q; want %q",
			err.Error(),
			"position already marked",
		)
	}

	if board.turnCount != 1 {
		t.Errorf("turnCount = %d; want 1", board.turnCount)
	}

	if board.GetTurn() != "O" {
		t.Errorf("turn = %q; want O", board.GetTurn())
	}
}

func TestMoveRejectedAfterDraw(t *testing.T) {
	board := newStartedBoard(t, 3)

	playMoves(t, board, 1, 2, 3, 5, 4, 6, 8, 7, 9)

	err := board.PlayMove(1)
	if err == nil {
		t.Fatal("PlayMove() after draw expected an error")
	}

	if err.Error() != "game not in progress" {
		t.Errorf(
			"error = %q; want %q",
			err.Error(),
			"game not in progress",
		)
	}
}

func TestFinalMoveWinIsNotMarkedAsDraw(t *testing.T) {
	board := newStartedBoard(t, 3)

	playMoves(t, board, 1, 2, 3, 4, 5, 7, 6, 8, 9)

	if board.GetStatus() != StatusWon {
		t.Errorf("status = %q; want %q", board.GetStatus(), StatusWon)
	}

	if board.GetWinner() != "X" {
		t.Errorf("winner = %q; want X", board.GetWinner())
	}
}

func TestRowWin(t *testing.T) {
	board := newStartedBoard(t, 3)

	playMoves(t, board, 1, 4, 2, 5, 3)

	if board.GetStatus() != StatusWon {
		t.Errorf("status = %q; want %q", board.GetStatus(), StatusWon)
	}

	if board.GetWinner() != "X" {
		t.Errorf("winner = %q; want X", board.GetWinner())
	}

	if board.GetTurn() != "X" {
		t.Errorf("turn = %q; want X", board.GetTurn())
	}
}

func TestColumnWin(t *testing.T) {
	board := newStartedBoard(t, 3)

	playMoves(t, board, 2, 1, 5, 3, 8)

	if board.GetStatus() != StatusWon {
		t.Errorf("status = %q; want %q", board.GetStatus(), StatusWon)
	}

	if board.GetWinner() != "X" {
		t.Errorf("winner = %q; want X", board.GetWinner())
	}
}

func TestMainDiagonalWin(t *testing.T) {
	board := newStartedBoard(t, 3)

	playMoves(t, board, 1, 2, 5, 3, 9)

	if board.GetStatus() != StatusWon {
		t.Errorf("status = %q; want %q", board.GetStatus(), StatusWon)
	}

	if board.GetWinner() != "X" {
		t.Errorf("winner = %q; want X", board.GetWinner())
	}
}

func TestOppositeDiagonalWin(t *testing.T) {
	board := newStartedBoard(t, 3)

	playMoves(t, board, 3, 1, 5, 2, 7)

	if board.GetStatus() != StatusWon {
		t.Errorf("status = %q; want %q", board.GetStatus(), StatusWon)
	}

	if board.GetWinner() != "X" {
		t.Errorf("winner = %q; want X", board.GetWinner())
	}
}

func TestOWins(t *testing.T) {
	board := newStartedBoard(t, 3)

	playMoves(t, board, 1, 4, 2, 5, 9, 6)

	if board.GetStatus() != StatusWon {
		t.Errorf("status = %q; want %q", board.GetStatus(), StatusWon)
	}

	if board.GetWinner() != "O" {
		t.Errorf("winner = %q; want O", board.GetWinner())
	}

	if board.GetTurn() != "O" {
		t.Errorf("turn = %q; want O", board.GetTurn())
	}
}

func TestDraw(t *testing.T) {
	board := newStartedBoard(t, 3)

	playMoves(t, board, 1, 2, 3, 5, 4, 6, 8, 7, 9)

	if board.GetStatus() != StatusDraw {
		t.Errorf("status = %q; want %q", board.GetStatus(), StatusDraw)
	}

	if board.GetWinner() != "" {
		t.Errorf("winner = %q; want empty string", board.GetWinner())
	}
}

func TestMoveRejectedAfterWin(t *testing.T) {
	board := newStartedBoard(t, 3)

	playMoves(t, board, 1, 4, 2, 5, 3)

	err := board.PlayMove(6)
	if err == nil {
		t.Fatal("PlayMove() after win expected an error")
	}

	if err.Error() != "game not in progress" {
		t.Errorf("error = %q; want %q", err.Error(), "game not in progress")
	}
}

func TestGetBoard(t *testing.T) {
	board := newStartedBoard(t, 3)

	playMoves(t, board, 1, 5, 9)

	got := board.GetBoard()
	want := []string{
		"X__",
		"_O_",
		"__X",
	}

	if len(got) != len(want) {
		t.Fatalf("GetBoard() returned %d rows; want %d", len(got), len(want))
	}

	for row := range want {
		if got[row] != want[row] {
			t.Errorf(
				"GetBoard()[%d] = %q; want %q",
				row,
				got[row],
				want[row],
			)
		}
	}
}

func TestGetBoardBeforeStart(t *testing.T) {
	board, err := NewBoard(3)
	if err != nil {
		t.Fatalf("NewBoard() returned unexpected error: %v", err)
	}

	got := board.GetBoard()
	want := []string{
		"___",
		"___",
		"___",
	}

	if len(got) != len(want) {
		t.Fatalf("GetBoard() returned %d rows; want %d", len(got), len(want))
	}

	for row := range want {
		if got[row] != want[row] {
			t.Errorf(
				"GetBoard()[%d] = %q; want %q",
				row,
				got[row],
				want[row],
			)
		}
	}
}
