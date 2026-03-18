package repository

import (
	"tiktaktoe/internal/pkg/datasource/repository"
	"tiktaktoe/internal/pkg/domain"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(user domain.User) error
	GetUser(id uuid.UUID) (domain.User, error)
	GetUserByLogin(login string) (domain.User, error)
	UpdateUser(user domain.User) (domain.User, error)
	DeleteUser(id uuid.UUID) error
}

func NewUserRepo() UserRepository {
	return repository.NewUserRepository()
}
