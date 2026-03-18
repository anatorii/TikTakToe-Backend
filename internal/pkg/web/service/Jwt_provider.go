package service

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var jwtSecret = []byte(os.Getenv("SECRET"))

var jwtp *JwtProvider

type session struct {
	UserId uuid.UUID
	Token  string
}

type AuthClaims struct {
	UserId string `json:"user_id"`
	Login  string `json:"login"`
	jwt.RegisteredClaims
}

type JwtProvider struct {
	Sessions map[string]session
}

func NewJwtProvider() JwtProvider {
	if jwtp == nil {
		jwtp = &JwtProvider{
			Sessions: make(map[string]session),
		}
	}
	return *jwtp
}

func (jwtp JwtProvider) GenRegClaims(subject string, mins int) jwt.RegisteredClaims {
	now := time.Now()
	return jwt.RegisteredClaims{
		Issuer:    "tiktaltoe",
		Subject:   subject,
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(mins) * time.Minute)),
	}
}

func (jwtp JwtProvider) GenAccessToken(userId uuid.UUID) (string, error) {
	claims := AuthClaims{
		UserId:           userId.String(),
		RegisteredClaims: jwtp.GenRegClaims("access", 5),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", errors.New("could not create access token")
	}
	return tokenString, nil
}

func (jwtp JwtProvider) GenRefreshToken(userId uuid.UUID) (string, error) {
	claims := AuthClaims{
		UserId:           userId.String(),
		RegisteredClaims: jwtp.GenRegClaims("refresh", 30),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", errors.New("could not create refresh token")
	}
	return tokenString, nil
}

func (jwtp JwtProvider) ValidateToken(tokenString string) (*AuthClaims, error) {
	claims := &AuthClaims{}
	token, _ := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return jwtSecret, nil
	})
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

func (jwtp JwtProvider) ValidateAccessToken(tokenString string) (*AuthClaims, error) {
	claims, err := jwtp.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}
	if claims.Subject != "access" {
		return nil, errors.New("wrong type of token")
	}
	return claims, nil
}

func (jwtp JwtProvider) ValidateRefreshToken(tokenString string) (*AuthClaims, error) {
	claims, err := jwtp.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}
	if claims.Subject != "refresh" {
		return nil, errors.New("wrong type of token")
	}
	return claims, nil
}

func (jwtp JwtProvider) UuidFromToken(tokenString string) (uuid.UUID, error) {
	claims := &AuthClaims{}
	token, _ := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return jwtSecret, nil
	})
	if !token.Valid {
		return uuid.Nil, errors.New("invalid token")
	}
	userId, err := uuid.Parse(claims.UserId)
	if err != nil {
		return uuid.Nil, err
	}
	return userId, nil
}

func (jwtp JwtProvider) findSession(tokenString string) (*session, error) {
	s, ok := jwtp.Sessions[tokenString]
	if !ok {
		return nil, errors.New("session not found")
	}
	return &s, nil
}

func (jwtp *JwtProvider) rotateRefreshToken(s *session, t string) error {
	newSession := session{
		Token:  t,
		UserId: s.UserId,
	}
	jwtp.saveSession(&newSession)
	jwtp.deleteSession(s)
	return nil
}

func (jwtp JwtProvider) saveSession(session *session) error {
	jwtp.Sessions[session.Token] = *session
	return nil
}

func (jwtp *JwtProvider) deleteSession(session *session) error {
	delete(jwtp.Sessions, session.Token)
	return nil
}
