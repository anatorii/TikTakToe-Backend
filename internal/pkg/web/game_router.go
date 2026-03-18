package web

import (
	"context"
	"net/http"
	"tiktaktoe/internal/pkg/web/service"

	"github.com/gorilla/mux"
)

var AppRouter *mux.Router

func NewRouter() *mux.Router {
	if AppRouter == nil {
		AppRouter = mux.NewRouter()
	}
	return AppRouter
}

type Handler func(w http.ResponseWriter, r *http.Request)

func JwtMiddleware(a service.UserAuthenticator, next Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := a.ExtractBearerToken(r)
		if err != nil {
			SendJsonResponse(w, http.StatusUnauthorized, ErrorResponse(err.Error()))
			return
		}
		claims, err := a.ValidJwt(token)
		if err != nil {
			SendJsonResponse(w, http.StatusUnauthorized, ErrorResponse(err.Error()))
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", claims.UserId)
		next(w, r.WithContext(ctx))
	}
}

func RegisterAuthRoutes(router *mux.Router, handler *AuthHandler, authenticator service.UserAuthenticator) {
	router.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		handler.RegisterHandler(w, r)
	}).Methods("POST")

	router.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		handler.JwtAuthHandler(w, r)
	}).Methods("POST")

	router.HandleFunc("/access", func(w http.ResponseWriter, r *http.Request) {
		handler.UpdateAccessTokenHandler(w, r)
	}).Methods("POST")

	router.HandleFunc("/refresh", func(w http.ResponseWriter, r *http.Request) {
		handler.UpdateRefreshTokenHandler(w, r)
	}).Methods("POST")
}

func RegisterGameRoutes(router *mux.Router, handler *GameHandler, authenticator service.UserAuthenticator) {
	router.HandleFunc("/game", JwtMiddleware(authenticator,
		func(w http.ResponseWriter, r *http.Request) {
			handler.ListGamesHandler(w, r)
		},
	)).Methods("GET")

	router.HandleFunc("/game/finished", JwtMiddleware(authenticator,
		func(w http.ResponseWriter, r *http.Request) {
			handler.FinishedGamesHandler(w, r)
		},
	)).Methods("GET")

	router.HandleFunc("/game/top", JwtMiddleware(authenticator,
		func(w http.ResponseWriter, r *http.Request) {
			handler.TopPlayersHandler(w, r)
		},
	)).Methods("GET")

	router.HandleFunc("/game/{uuid}", JwtMiddleware(authenticator,
		func(w http.ResponseWriter, r *http.Request) {
			handler.GetGameHandler(w, r)
		},
	)).Methods("GET")

	router.HandleFunc("/game/create", JwtMiddleware(authenticator,
		func(w http.ResponseWriter, r *http.Request) {
			handler.CreateGameHandler(w, r)
		},
	)).Methods("POST")

	router.HandleFunc("/game/{uuid}/join", JwtMiddleware(authenticator,
		func(w http.ResponseWriter, r *http.Request) {
			handler.JoinGameHandler(w, r)
		},
	)).Methods("POST")

	router.HandleFunc("/game/{uuid}/move", JwtMiddleware(authenticator,
		func(w http.ResponseWriter, r *http.Request) {
			handler.MakeMoveHandler(w, r)
		},
	)).Methods("POST")

	router.HandleFunc("/user/{uuid}", JwtMiddleware(authenticator,
		func(w http.ResponseWriter, r *http.Request) {
			handler.GetUserHandler(w, r)
		},
	)).Methods("GET")

	router.HandleFunc("/user", JwtMiddleware(authenticator,
		func(w http.ResponseWriter, r *http.Request) {
			handler.GetUserHandler(w, r)
		},
	)).Methods("GET")
}
