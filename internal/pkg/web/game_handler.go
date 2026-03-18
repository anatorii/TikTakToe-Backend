package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	wservice "tiktaktoe/internal/pkg/web/service"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type GameHandler struct {
	gameService wservice.GameService
	userService wservice.UserService
}

func NewGameHandler(gs wservice.GameService, us wservice.UserService) *GameHandler {
	return &GameHandler{
		gameService: gs,
		userService: us,
	}
}

// обработчик POST /game
func (h *GameHandler) CreateGameHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendJsonResponse(w, http.StatusInternalServerError, ErrorResponse("No input params absense"))
		return
	}

	single, err := GameRequestToCreate(&req)
	if err != nil {
		SendJsonResponse(w, http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	// Вызов сервисного слоя
	usewrId, _ := uuid.Parse(r.Context().Value("user_id").(string))
	user, _ := h.userService.GetUser(usewrId)
	game, err := h.gameService.CreateGame(user, single)
	if err != nil {
		SendJsonResponse(w, http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	jsonMessage := DomainToResponse(game)
	SendJsonResponse(w, http.StatusOK, jsonMessage)
}

// обработчик POST /game/{uuid}/move
func (h *GameHandler) MakeMoveHandler(w http.ResponseWriter, r *http.Request) {
	// Извлечение GameID из URL
	vars := mux.Vars(r)
	gameID, err := uuid.Parse(vars["uuid"])
	if err != nil {
		SendJsonResponse(w, http.StatusInternalServerError, ErrorResponse("Invalid game ID"))
		return
	}

	// Парсинг тела запроса, request json -> request struct
	var req MoveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendJsonResponse(w, http.StatusInternalServerError, ErrorResponse("Invalid request body"))
		return
	}

	// Конвертация запроса, request struct -> int
	move, err := RequestToMove(&req)
	if err != nil {
		SendJsonResponse(w, http.StatusInternalServerError, ErrorResponse("Board conversion error"))
		return
	}

	// Вызов сервисного слоя для обработки хода
	usewrId, _ := uuid.Parse(r.Context().Value("user_id").(string))
	user, _ := h.userService.GetUser(usewrId)
	updatedGame, err := h.gameService.NextMove(gameID, move, user.Id)
	if err != nil {
		SendJsonResponse(w, http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	jsonMessage := DomainToResponse(updatedGame)
	SendJsonResponse(w, http.StatusOK, jsonMessage)
}

// обработчик GET /game/{uuid}
func (h *GameHandler) GetGameHandler(w http.ResponseWriter, r *http.Request) {
	// Извлечение GameID из URL
	vars := mux.Vars(r)
	gameId, err := uuid.Parse(vars["uuid"])
	if err != nil {
		SendJsonResponse(w, http.StatusBadRequest, ErrorResponse(fmt.Sprint("Invalid game ID:", gameId)))
		return
	}

	// Вызов сервисного слоя
	game, err := h.gameService.GetGame(gameId)
	if err != nil {
		SendJsonResponse(w, http.StatusBadRequest, ErrorResponse(err.Error()))
		return
	}

	jsonMessage := DomainToResponse(game)
	SendJsonResponse(w, http.StatusOK, jsonMessage)
}

// обработчик POST /game/{game_UUID}/join
func (h *GameHandler) JoinGameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameId, err := uuid.Parse(vars["uuid"])
	if err != nil {
		SendJsonResponse(w, http.StatusBadRequest, ErrorResponse(fmt.Sprint("Invalid game ID:", gameId)))
		return
	}

	usewrId, _ := uuid.Parse(r.Context().Value("user_id").(string))
	user, _ := h.userService.GetUser(usewrId)
	game, err := h.gameService.JoinGame(gameId, user.Id)
	if err != nil {
		SendJsonResponse(w, http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	jsonMessage := DomainToResponse(game)
	SendJsonResponse(w, http.StatusOK, jsonMessage)
}

// обработчик GET /game
func (h *GameHandler) ListGamesHandler(w http.ResponseWriter, r *http.Request) {
	usewrId, _ := uuid.Parse(r.Context().Value("user_id").(string))
	user, _ := h.userService.GetUser(usewrId)
	gameList := h.gameService.ListGames(user)
	response := make([]GameResponse, 0)
	for _, game := range gameList {
		response = append(response, *DomainToResponse(game))
	}
	SendJsonResponse(w, http.StatusOK, response)
}

// обработчик GET /game/finished
func (h *GameHandler) FinishedGamesHandler(w http.ResponseWriter, r *http.Request) {
	usewrId, _ := uuid.Parse(r.Context().Value("user_id").(string))

	gameList := h.gameService.FinishedGames(usewrId)
	response := make([]GameResponse, 0)
	for _, game := range gameList {
		response = append(response, *DomainToResponse(game))
	}
	SendJsonResponse(w, http.StatusOK, response)
}

// обработчик GET /game/top
func (h *GameHandler) TopPlayersHandler(w http.ResponseWriter, r *http.Request) {
	var req TopPlayersRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendJsonResponse(w, http.StatusInternalServerError, ErrorResponse("No input params absense"))
		return
	}
	num, err := TopPlayersRequestToNumPlayers(&req)
	if err != nil {
		SendJsonResponse(w, http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}

	players := h.gameService.TopPlayers(num)
	response := make([]PlayerResponse, 0)
	for _, player := range players {
		pu, _ := h.userService.GetUser(player.Id)
		response = append(response, *PlayerToPlayerResponse(player, pu.Login))
	}
	SendJsonResponse(w, http.StatusOK, response)
}

// обработчик GET /user/{uuid}, GET /user
func (h *GameHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if id, ok := vars["uuid"]; ok {
		userId, err := uuid.Parse(id)
		if err != nil {
			SendJsonResponse(w, http.StatusBadRequest, ErrorResponse(fmt.Sprint("Invalid user id: ", id)))
			return
		}
		user, err := h.userService.GetUser(userId)
		if err != nil {
			SendJsonResponse(w, http.StatusBadRequest, ErrorResponse("User not found"))
			return
		}
		SendJsonResponse(w, http.StatusOK, user)
	} else {
		usewrId, _ := uuid.Parse(r.Context().Value("user_id").(string))
		user, _ := h.userService.GetUser(usewrId)
		SendJsonResponse(w, http.StatusOK, user)
	}
}
