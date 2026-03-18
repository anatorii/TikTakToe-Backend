package domain

import "github.com/google/uuid"

type Player struct {
	Id    uuid.UUID
	Wins  int
	Loses int
	Draws int
}
