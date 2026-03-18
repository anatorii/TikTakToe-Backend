package service

import (
	"encoding/base64"
	"errors"
	"strings"
	"tiktaktoe/internal/pkg/web/auth"

	"github.com/google/uuid"
)

type AuthenticationService struct {
	userService UserService
	jwtProvider JwtProvider
}

func NewAuthenticationService(us UserService, jwtp JwtProvider) *AuthenticationService {
	return &AuthenticationService{
		userService: us,
		jwtProvider: jwtp,
	}
}

func (s *AuthenticationService) Register(req auth.SignUpRequest) (uuid.UUID, error) {
	user, err := s.userService.Register(req.Login, req.Password)
	return user.Id, err
}

func (s *AuthenticationService) Authenticate(authHeader string) (uuid.UUID, error) {
	if !strings.HasPrefix(authHeader, "Basic ") {
		err := errors.New("Invalid authentication format")
		return uuid.Nil, err
	}

	encoded := strings.TrimPrefix(authHeader, "Basic ")
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		err := errors.New("Invalid base64 encoding")
		return uuid.Nil, err
	}

	parts := strings.SplitN(string(decoded), ":", 2)
	if len(parts) != 2 {
		err := errors.New("Invalid credentials format")
		return uuid.Nil, err
	}

	login, password := parts[0], parts[1]

	return s.userService.Authenticate(login, password)
}

func (s *AuthenticationService) JwtAuthenticate(req auth.JwtRequest) (auth.JwtResponse, error) {
	userId, err := s.userService.Authenticate(req.Login, req.Password)
	if err != nil {
		return auth.JwtResponse{}, errors.New("User not found")
	}

	accessString, err := jwtp.GenAccessToken(userId)
	if err != nil {
		return auth.JwtResponse{}, errors.New("Generate token failed")
	}
	refreshToken, err := jwtp.GenRefreshToken(userId)
	if err != nil {
		return auth.JwtResponse{}, errors.New("Generate refresh token failed")
	}

	jwtp.saveSession(&session{Token: refreshToken, UserId: userId})

	response := auth.JwtResponse{
		AccessToken:  accessString,
		RefreshToken: refreshToken,
	}

	return response, nil
}

func (s *AuthenticationService) UpdateAccessToken(refreshToken string) (auth.JwtResponse, error) {
	_, err := jwtp.ValidateRefreshToken(refreshToken)
	if err != nil {
		return auth.JwtResponse{}, err
	}
	userId, err := jwtp.UuidFromToken(refreshToken)
	if err != nil {
		return auth.JwtResponse{}, err
	}
	_, err = s.userService.GetUser(userId)
	if err != nil {
		return auth.JwtResponse{}, err
	}
	_, err = jwtp.findSession(refreshToken)
	if err != nil {
		return auth.JwtResponse{}, err
	}

	accessString, err := jwtp.GenAccessToken(userId)
	if err != nil {
		return auth.JwtResponse{}, errors.New("Generate token failed")
	}

	response := auth.JwtResponse{
		AccessToken:  accessString,
		RefreshToken: refreshToken,
	}

	return response, nil
}

func (s *AuthenticationService) UpdateRefreshToken(refreshToken string) (auth.JwtResponse, error) {
	_, err := jwtp.ValidateRefreshToken(refreshToken)
	if err != nil {
		return auth.JwtResponse{}, err
	}
	userId, err := jwtp.UuidFromToken(refreshToken)
	if err != nil {
		return auth.JwtResponse{}, err
	}
	_, err = s.userService.GetUser(userId)
	if err != nil {
		return auth.JwtResponse{}, err
	}
	ses, err := jwtp.findSession(refreshToken)
	if err != nil {
		return auth.JwtResponse{}, err
	}

	accessString, err := jwtp.GenAccessToken(userId)
	if err != nil {
		return auth.JwtResponse{}, errors.New("Generate token failed")
	}
	newRefreshToken, err := jwtp.GenRefreshToken(userId)
	if err != nil {
		return auth.JwtResponse{}, errors.New("Generate refresh token failed")
	}

	jwtp.rotateRefreshToken(ses, newRefreshToken)

	response := auth.JwtResponse{
		AccessToken:  accessString,
		RefreshToken: newRefreshToken,
	}

	return response, nil
}
