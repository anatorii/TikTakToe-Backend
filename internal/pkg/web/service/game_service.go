package service

import (
	"tiktaktoe/internal/pkg/domain"
	drepository "tiktaktoe/internal/pkg/domain/repository"
	"tiktaktoe/internal/pkg/domain/service"

	"github.com/google/uuid"
)

type GameService interface {
	GetGame(gameId uuid.UUID) (domain.Game, error)
	CreateGame(user domain.User, single int) (domain.Game, error)
	JoinGame(gameId uuid.UUID, id uuid.UUID) (domain.Game, error)
	NextMove(gameId uuid.UUID, move int, player uuid.UUID) (domain.Game, error)
	ValidateBoard(gameId uuid.UUID, newBoard domain.Board) error
	GameIsOver(gameId uuid.UUID) bool
	ListGames(user domain.User) []domain.Game
	FinishedGames(userId uuid.UUID) []domain.Game
	TopPlayers(numPlayers int) []domain.Player
}

func NewGameServ(r drepository.GameRepository) GameService {
	return service.NewGameService(r)
}
