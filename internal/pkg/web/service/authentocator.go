package service

import (
	"fmt"
	"net/http"
	"strings"
)

type UserAuthenticator struct {
	service UserService
	jwt     JwtProvider
}

func NewAuthenticator(service UserService, jwt JwtProvider) UserAuthenticator {
	return UserAuthenticator{
		service: service,
		jwt:     jwt,
	}
}

func (a UserAuthenticator) ValidUser(login, pass string) bool {
	return a.service.VerifyUser(login, pass)
}

func (a UserAuthenticator) ValidJwt(token string) (*AuthClaims, error) {
	claims, err := a.jwt.ValidateAccessToken(token)
	if err != nil {
		return nil, err
	}
	id, err := jwtp.UuidFromToken(token)
	if err != nil {
		return nil, err
	}
	_, err = a.service.GetUser(id)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func (a UserAuthenticator) ExtractBearerToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header missing")
	}
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fmt.Errorf("authorization header is not Bearer token")
	}
	return strings.TrimPrefix(authHeader, "Bearer "), nil
}
