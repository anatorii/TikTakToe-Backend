package domain

import "github.com/google/uuid"

type User struct {
	Id       uuid.UUID
	Login    string
	Password string
}
