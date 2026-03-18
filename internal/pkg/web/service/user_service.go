package service

import (
	"tiktaktoe/internal/pkg/domain"
	drepository "tiktaktoe/internal/pkg/domain/repository"
	"tiktaktoe/internal/pkg/domain/service"

	"github.com/google/uuid"
)

type UserService interface {
	GetUser(id uuid.UUID) (domain.User, error)
	GetUserByLogin(login string) (domain.User, error)
	VerifyUser(login, password string) bool
	Register(login, pass string) (domain.User, error)
	Authenticate(login, pass string) (uuid.UUID, error)
}

func NewUserServ(r drepository.UserRepository) UserService {
	return service.NewUserService(r)
}
