package service

import (
	"errors"
	"tiktaktoe/internal/pkg/domain"
	drepository "tiktaktoe/internal/pkg/domain/repository"

	"github.com/google/uuid"
)

type userService struct {
	repo drepository.UserRepository
}

func NewUserService(r drepository.UserRepository) userService {
	return userService{
		repo: r,
	}
}

func (s userService) GetUser(id uuid.UUID) (domain.User, error) {
	return s.repo.GetUser(id)
}

func (s userService) GetUserByLogin(login string) (domain.User, error) {
	return s.repo.GetUserByLogin(login)
}

func (s userService) VerifyUser(login, password string) bool {
	user, err := s.repo.GetUserByLogin(login)
	if err == nil && user.Password == password {
		return true
	}

	return false
}

func (s userService) Register(login, pass string) (domain.User, error) {
	user, err := s.repo.GetUserByLogin(login)
	if err == nil {
		return domain.User{}, errors.New("User exists")
	}

	user = domain.User{
		Id:       uuid.New(),
		Login:    login,
		Password: pass,
	}

	err = s.repo.CreateUser(user)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (s userService) Authenticate(login, pass string) (uuid.UUID, error) {
	user, err := s.repo.GetUserByLogin(login)
	if err == nil && user.Password == pass {
		return user.Id, nil
	}

	return uuid.UUID{}, errors.New("User does not exist")
}
