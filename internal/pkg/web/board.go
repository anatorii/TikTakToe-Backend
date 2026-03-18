package web

import (
	"errors"
	"math"
	"tiktaktoe/internal/pkg/domain"

	"github.com/google/uuid"
)

type WebBoard [3][3]int

type TopPlayersRequest struct {
	NumPlayers int `json:"num_players"`
}

type CreateGameRequest struct {
	Single int `json:"single"`
}

type MoveRequest struct {
	Move int `json:"move"`
}

type GameResponse struct {
	GameID    uuid.UUID `json:"game_id"`
	Status    string    `json:"status"`
	Player1   uuid.UUID `json:"player1"`
	Player2   uuid.UUID `json:"player2"`
	Board     WebBoard  `json:"board"`
	Single    int       `json:"single"`
	CreatedAt int64     `json:"created_at"`
}

type PlayerResponse struct {
	Player  uuid.UUID `json:"player"`
	Login   string    `json:"login"`
	Percent int       `json:"percent"`
}

func TopPlayersRequestToNumPlayers(request *TopPlayersRequest) (int, error) {
	if request.NumPlayers < 1 {
		return 0, errors.New("Incorrect players number")
	}
	return request.NumPlayers, nil
}

func GameRequestToCreate(request *CreateGameRequest) (int, error) {
	if request.Single < 1 || request.Single > 2 {
		return 0, errors.New("Incorrect players number")
	}
	return request.Single, nil
}

func DomainToResponse(game domain.Game) *GameResponse {
	return &GameResponse{
		Board:     convertBoardToWeb(game.GetBoard()),
		GameID:    game.GetId(),
		Status:    game.GetStatus(),
		Player1:   game.GetPlayer(1),
		Player2:   game.GetPlayer(2),
		Single:    game.GetSingle(),
		CreatedAt: game.GetCreatedAt().Unix(),
	}
}

func PlayerToPlayerResponse(player domain.Player, login string) *PlayerResponse {
	return &PlayerResponse{
		Player:  player.Id,
		Login:   login,
		Percent: int(math.Round(float64(player.Wins) / (float64(player.Draws) + float64(player.Loses) + float64(player.Wins)) * 100)),
	}
}

func RequestToMove(request *MoveRequest) (int, error) {
	return request.Move, nil
}

func convertBoardToDomain(webBoard WebBoard) (domain.Board, error) {
	return domain.Board(webBoard), nil
}

func convertBoardToWeb(domainBoard domain.Board) WebBoard {
	return WebBoard(domainBoard)
}
