package service

import (
	"errors"
	"tiktaktoe/internal/pkg/domain"
	drepository "tiktaktoe/internal/pkg/domain/repository"

	"github.com/google/uuid"
)

type gameService struct {
	repo drepository.GameRepository
}

func NewGameService(r drepository.GameRepository) gameService {
	return gameService{
		repo: r,
	}
}

func (s gameService) CreateGame(user domain.User, single int) (domain.Game, error) {
	game := domain.NewGame()
	game.SetPlayer(1, user.Id)
	game.SetSingle(single)
	game.SetActiveTurn(domain.PLAYER1_CHAR)
	if single == 1 {
		game.SetStatus(string(domain.S_PLAYER1_TURN))
	} else {
		game.SetStatus(string(domain.S_WAITING))
	}
	s.repo.CreateGame(*game)
	return *game, nil
}

func (s gameService) GetGame(gameId uuid.UUID) (domain.Game, error) {
	return s.repo.GetGame(gameId)
}

func (s gameService) JoinGame(gameId uuid.UUID, playerId uuid.UUID) (domain.Game, error) {
	game, err := s.repo.GetGame(gameId)
	if err != nil {
		return domain.Game{}, errors.New("Join game. Game not found")
	}
	if game.GetPlayer(1) == playerId || game.GetPlayer(2) == playerId {
		return domain.Game{}, errors.New("Join game. Already joined")
	}
	if game.GetSingle() == 1 {
		return domain.Game{}, errors.New("Join game. No vacancies")
	} else {
		if game.GetPlayer(2) != (uuid.UUID{}) {
			return domain.Game{}, errors.New("Join game. No vacancies")
		} else {
			game.SetPlayer(2, playerId)
			game.SetActiveTurn(domain.PLAYER1_CHAR)
			game.SetStatus(string(domain.S_PLAYER1_TURN))
			s.repo.UpdateGame(game)
		}
	}

	return game, nil
}

func (s gameService) NextMove(gameId uuid.UUID, move int, playerId uuid.UUID) (domain.Game, error) {
	game, err := s.repo.GetGame(gameId)
	if err != nil {
		return domain.Game{}, errors.New("Game not found")
	}

	if game.GetStatus() == game.GameOverStatus() {
		return domain.Game{}, errors.New("Game is over")
	}

	if game.GetStatus() == string(domain.S_WAITING) {
		return domain.Game{}, errors.New("Game is not ready")
	}

	if move < 0 || move > 8 {
		return domain.Game{}, errors.New("Wrong move")
	}

	var player int
	if playerId == game.GetPlayer(1) {
		player = 1
	} else if playerId == game.GetPlayer(2) {
		player = 2
	}

	if (player == 1 && game.GetStatus() == string(domain.S_PLAYER2_TURN)) ||
		(player == 2 && game.GetStatus() == string(domain.S_PLAYER1_TURN)) {
		return domain.Game{}, errors.New("Wrong turn")
	}

	newGame := game.GetNewState(move)
	if err := s.ValidateBoard(game.GetId(), newGame.GetBoard()); err != nil {
		return domain.Game{}, err
	}
	if game.GetSingle() == 1 {
		_, choice := (domain.MinimaxAlg{}).GetNextMove(newGame)
		newGame = newGame.GetNewState(choice)
	}
	if game.GetSingle() != 1 {
		if player == 1 {
			newGame.SetStatus(string(domain.S_PLAYER2_TURN))
		} else {
			newGame.SetStatus(string(domain.S_PLAYER1_TURN))
		}
	}

	if st := newGame.GameOverStatus(); st != "" {
		newGame.SetStatus(string(st))
	}

	s.repo.UpdateGame(newGame)

	return newGame, nil
}

func (s gameService) ValidateBoard(gameId uuid.UUID, newBoard domain.Board) error {
	game, err := s.repo.GetGame(gameId)
	if err != nil {
		return err
	}

	currentBoard := game.GetBoard()
	diffCount := 0
	for row := range 3 {
		for col := range 3 {
			if currentBoard[row][col] != newBoard[row][col] {
				if currentBoard[row][col] != 0 {
					return errors.New("Can not modify existing cells")
				}
				if newBoard[row][col] != int(game.GetActiveTurn()) {
					return errors.New("Wrong player move")
				}
				diffCount++
			}
		}
	}
	if diffCount == 0 {
		return errors.New("Cell already filled")
	}
	if diffCount != 1 {
		return errors.New("Exactly one move must be made")
	}

	return nil
}

func (s gameService) GameIsOver(gameId uuid.UUID) bool {
	game, _ := s.repo.GetGame(gameId)
	return game.GameIsOver()
}

func (s gameService) ListGames(user domain.User) []domain.Game {
	gms := make([]domain.Game, 0)
	games, _ := s.repo.ListGames(user)
	for _, game := range games {
		gms = append(gms, game)
	}
	return games
}

func (s gameService) FinishedGames(userId uuid.UUID) []domain.Game {
	gms := make([]domain.Game, 0)
	games, _ := s.repo.FinishedGames(userId)
	for _, game := range games {
		gms = append(gms, game)
	}
	return gms
}

func (s gameService) TopPlayers(numPlayers int) []domain.Player {
	players := make([]domain.Player, 0)
	users, _ := s.repo.TopPlayers(numPlayers)
	for _, user := range users {
		players = append(players, user)
	}
	return players
}
