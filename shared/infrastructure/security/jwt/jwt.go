package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	Document       string  `json:"document"`
	Name           string  `json:"name"`
	AccessGroupIds []int16 `json:"accessGroupIds"`
	TokenType      string  `json:"token_type"`

	jwt.RegisteredClaims
}

type TokenPair struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	ExpiresIn        int64  `json:"expires_in"`
	RefreshExpiresIn int64  `json:"refresh_expires_in"`
}

type JwtManager struct {
	accessSecret  []byte
	refreshSecret []byte
}

func NewJwtManager(accessSecret, refreshSecret string) *JwtManager {
	return &JwtManager{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
	}
}

func (j *JwtManager) GenerateTokenPair(
	userID string,
	document string,
	accessGroupIds []int16,
	name string,
) (*TokenPair, error) {

	now := time.Now()

	accessExp := now.Add(1 * time.Hour)
	refreshExp := now.Add(7 * 24 * time.Hour)

	accessClaims := Claims{
		Document:       document,
		Name:           name,
		AccessGroupIds: accessGroupIds,
		TokenType:      "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(accessExp),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        uuid.NewString(),
		},
	}

	accessToken, err := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		accessClaims,
	).SignedString(j.accessSecret)
	if err != nil {
		return nil, err
	}

	refreshClaims := accessClaims
	refreshClaims.TokenType = "refresh"
	refreshClaims.ExpiresAt = jwt.NewNumericDate(refreshExp)
	refreshClaims.ID = uuid.NewString()

	refreshToken, err := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		refreshClaims,
	).SignedString(j.refreshSecret)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken,
		ExpiresIn:        accessExp.Unix(),
		RefreshExpiresIn: refreshExp.Unix(),
	}, nil
}

func (j *JwtManager) ValidateAccessToken(token string) (*Claims, error) {
	return j.validateToken(token, j.accessSecret, "access")
}

func (j *JwtManager) ValidateRefreshToken(token string) (*Claims, error) {
	return j.validateToken(token, j.refreshSecret, "refresh")
}

func (j *JwtManager) validateToken(
	tokenString string,
	secret []byte,
	expectedType string,
) (*Claims, error) {

	token, err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(t *jwt.Token) (interface{}, error) {
			return secret, nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	if claims.TokenType != expectedType {
		return nil, errors.New("invalid token type")
	}

	return claims, nil
}
