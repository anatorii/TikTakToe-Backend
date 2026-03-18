package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tiktaktoe/internal/pkg/web/auth"
	wservice "tiktaktoe/internal/pkg/web/service"
)

type AuthHandler struct {
	authService *wservice.AuthenticationService
}

func NewAuthHandler(auth *wservice.AuthenticationService) *AuthHandler {
	return &AuthHandler{
		authService: auth,
	}
}

func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		SendJsonResponse(w, http.StatusBadRequest, ErrorResponse("Failed to parse form"))
		return
	}

	login := r.Form.Get("login")
	password := r.Form.Get("password")

	if login == "" || password == "" {
		SendJsonResponse(w, http.StatusBadRequest, ErrorResponse("Login and password are required"))
		return
	}

	if len(login) < 3 || len(login) > 64 {
		SendJsonResponse(w, http.StatusBadRequest, ErrorResponse("Login must be 3–64 characters"))
		return
	}

	req := auth.SignUpRequest{Login: login, Password: password}

	_, err := h.authService.Register(req)
	if err != nil {
		SendJsonResponse(w, http.StatusInternalServerError, ErrorResponse(err.Error()))
		return
	}
	SendJsonResponse(w, http.StatusOK, "Registration succeded")
}

func (h *AuthHandler) AuthHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		http.Error(w, `{"error": "Authorization header required"}`, http.StatusUnauthorized)
		return
	}

	uuid, err := h.authService.Authenticate(authHeader)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusUnauthorized)
		return
	}

	id := map[string]interface{}{"uuid": uuid}

	SendJsonResponse(w, http.StatusOK, id)
}

func (h *AuthHandler) JwtAuthHandler(w http.ResponseWriter, r *http.Request) {
	// extract credentials
	var request auth.JwtRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		SendJsonResponse(w, http.StatusInternalServerError, ErrorResponse("Invalid request body"))
		return
	}
	// authenticate
	response, err := h.authService.JwtAuthenticate(request)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusUnauthorized)
		return
	}
	// send responce with tokens
	SendJsonResponse(w, http.StatusOK, response)
}

func (h *AuthHandler) UpdateAccessTokenHandler(w http.ResponseWriter, r *http.Request) {
	// extract refresh token
	var request auth.RefreshJwtRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		SendJsonResponse(w, http.StatusInternalServerError, ErrorResponse("Invalid request body"))
	}
	// update access token
	response, err := h.authService.UpdateAccessToken(request.RefreshToken)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusUnauthorized)
		return
	}
	// send responce with tokens
	SendJsonResponse(w, http.StatusOK, response)
}

func (h *AuthHandler) UpdateRefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	// extract refresh token
	var request auth.RefreshJwtRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		SendJsonResponse(w, http.StatusInternalServerError, ErrorResponse("Invalid request body"))
	}
	// update access token
	response, err := h.authService.UpdateRefreshToken(request.RefreshToken)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "%s"}`, err.Error()), http.StatusUnauthorized)
		return
	}
	// send responce with tokens
	SendJsonResponse(w, http.StatusOK, response)
}
