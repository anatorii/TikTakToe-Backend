package datasource

import "github.com/google/uuid"

type Player struct {
	Id    uuid.UUID `db:"id"`
	Wins  int       `db:"wins"`
	Fails int       `db:"fails"`
	Draws int       `db:"draws"`
}
