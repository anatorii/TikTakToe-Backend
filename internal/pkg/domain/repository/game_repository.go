package repository

import (
	"tiktaktoe/internal/pkg/datasource/repository"
	"tiktaktoe/internal/pkg/domain"

	"github.com/google/uuid"
)

type GameRepository interface {
	CreateGame(user domain.Game) error
	GetGame(id uuid.UUID) (domain.Game, error)
	UpdateGame(user domain.Game) (domain.Game, error)
	DeleteGame(id uuid.UUID) error
	ListGames(user domain.User) ([]domain.Game, error)
	FinishedGames(userId uuid.UUID) ([]domain.Game, error)
	TopPlayers(numPlayers int) ([]domain.Player, error)
}

func NewGameRepo() GameRepository {
	return repository.NewGameRepository()
}
