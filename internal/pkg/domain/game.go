package domain

import (
	"time"

	"github.com/google/uuid"
)

const (
	PLAYER1_CHAR = 'x'
	PLAYER2_CHAR = 'o'
)

const (
	S_WAITING      status_type = "waiting"
	S_PLAYER1_TURN status_type = "player1_turn"
	S_PLAYER2_TURN status_type = "player2_turn"
	S_PLAYER1_WIN  status_type = "player1_win"
	S_PLAYER2_WIN  status_type = "player2_win"
	S_DRAW         status_type = "draw"
)

type status_type string

type PlayerChar int

type Game struct {
	id         uuid.UUID
	board      Board
	turn       PlayerChar
	single     int
	status     string
	player1_id uuid.UUID
	player2_id uuid.UUID
	player1    PlayerChar
	player2    PlayerChar
	created_at time.Time
}

func NewGame() *Game {
	game := &Game{
		id:         uuid.New(),
		board:      [3][3]int{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
		turn:       PLAYER1_CHAR,
		single:     1,
		status:     "waiting",
		player1_id: uuid.Nil,
		player2_id: uuid.Nil,
		player1:    PLAYER1_CHAR,
		player2:    PLAYER2_CHAR,
		created_at: time.Now(),
	}
	return game
}

func (t *Game) GetCreatedAt() time.Time {
	return t.created_at
}

func (t *Game) GetSingle() int {
	return t.single
}

func (t *Game) GetPlayer(player int) uuid.UUID {
	if player == 1 {
		return t.player1_id
	} else if player == 2 {
		return t.player2_id
	}
	return uuid.UUID{}
}

func (t *Game) GetId() uuid.UUID {
	return t.id
}

func (t *Game) GetBoard() Board {
	return t.board
}

func (t Game) GetStatus() string {
	return t.status
}

func (t *Game) GetActiveTurn() PlayerChar {
	return t.turn
}

func (t *Game) GetPlayerChar() PlayerChar {
	return t.player1
}

func (t *Game) GetOpponentChar() PlayerChar {
	return t.player2
}

func (t *Game) SetCreatedAt(createdAt int64) {
	t.created_at = time.Unix(createdAt, 0)
}

func (t *Game) SetSingle(s int) {
	t.single = s
}

func (t *Game) SetPlayer(player int, id uuid.UUID) {
	if player == 1 {
		t.player1_id = id
	} else if player == 2 {
		t.player2_id = id
	}
}

func (t *Game) SetId(id uuid.UUID) {
	t.id = id
}

func (t *Game) SetBoard(b Board) {
	t.board = b
}

func (t *Game) SetStatus(s string) {
	t.status = s
}

func (t *Game) SetActiveTurn(turn PlayerChar) {
	t.turn = turn
}

func (t *Game) SetPlayerChar(ch PlayerChar) {
	t.player1 = ch
}

func (t *Game) SetOpponentChar(ch PlayerChar) {
	t.player2 = ch
}

func (t *Game) Win(turn int) bool {
	for i := range 3 {
		if t.board[i][0] == turn &&
			t.board[i][1] == turn &&
			t.board[i][2] == turn {
			return true
		}
	}
	for j := range 3 {
		if t.board[0][j] == turn &&
			t.board[1][j] == turn &&
			t.board[2][j] == turn {
			return true
		}
	}
	win := (t.board[0][0] == turn &&
		t.board[1][1] == turn &&
		t.board[2][2] == turn) ||
		(t.board[0][2] == turn &&
			t.board[1][1] == turn &&
			t.board[2][0] == turn)
	return win
}

func (t *Game) GameIsOver() bool {
	return t.Win(int(t.player1)) || t.Win(int(t.player2)) || len(t.GetAvailableMoves()) == 0
}

func (t *Game) GetAvailableMoves() []int {
	moves := make([]int, 0)
	for i := range 3 {
		for j := range 3 {
			if t.board[i][j] == 0 {
				moves = append(moves, i*3+j)
			}
		}
	}
	return moves
}

func (t *Game) GetNewState(move int) Game {
	// Создаем копию игры с новым ходом
	newGame := *t
	row := move / 3
	col := move % 3
	newGame.board[row][col] = int(t.turn)

	// Меняем ход
	if t.turn == t.player1 {
		newGame.turn = t.player2
	} else {
		newGame.turn = t.player1
	}
	return newGame
}

func (t *Game) GameOverStatus() string {
	status := ""
	if t.Win(int(t.player1)) {
		status = string(S_PLAYER1_WIN)
	} else if t.Win(int(t.player2)) {
		status = string(S_PLAYER2_WIN)
	} else if len(t.GetAvailableMoves()) == 0 {
		status = string(S_DRAW)
	}
	return status
}
