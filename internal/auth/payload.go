package auth

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

type TokenType byte

const (
	AccessToken  TokenType = 'A'
	RefreshToken TokenType = 'R'
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	Type      TokenType `json:"token_type"`
	Username  string    `json:"username"`
	Role      string    `json:"user_role"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expire_at"`
}

func NewPayload(username string, role string, duration time.Duration, tokenType TokenType) (*Payload, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	paylaod := &Payload{
		ID:        tokenID,
		Type:      tokenType,
		Username:  username,
		Role:      role,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}
	return paylaod, nil
}

func (payload *Payload) Valid(tokenType TokenType) error {
	if payload.Type != tokenType {
		return ErrInvalidToken
	}
	if time.Now().After(payload.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
