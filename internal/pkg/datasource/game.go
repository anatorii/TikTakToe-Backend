package datasource

import (
	"github.com/google/uuid"
)

type Game struct {
	Id         uuid.UUID `db:"id"`
	Board      [3][3]int `db:"board"`
	Turn       int       `db:"turn"`
	Status     string    `db:"status"`
	Single     int       `db:"single"`
	Player1_id uuid.UUID `db:"player1_id"`
	Player2_id uuid.UUID `db:"player2_id"`
	Created_at int64     `db:"created_at"`
}

func (g Game) BoardJson() (interface{}, error) {
	b := make([][]int, 3)
	for i := range b {
		b[i] = make([]int, 3)
		copy(b[i], g.Board[i][:])
	}
	return b, nil
}
