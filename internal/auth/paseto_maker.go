package auth

import (
	"fmt"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/google/uuid"
)

// PasetoMaker is a PASETO token maker
type PasetoMaker struct {
	symmetricKey paseto.V4SymmetricKey
}

// NewPasetoMaker creates a new PasetoMaker
func NewPasetoMaker(symmetricKey string) (Maker, error) {
	key, err := paseto.V4SymmetricKeyFromBytes([]byte(symmetricKey))
	if err != nil {
		return nil, fmt.Errorf("invalid key: %w", err)
	}

	maker := &PasetoMaker{
		symmetricKey: key,
	}

	return maker, nil
}

// CreateToken creates a new token for a specific username and duration
func (maker *PasetoMaker) CreateToken(
	username string,
	role string,
	duration time.Duration,
	tokenType TokenType,
) (string, *Payload, error) {

	payload, err := NewPayload(username, role, duration, tokenType)
	if err != nil {
		return "", nil, err
	}

	token := paseto.NewToken()
	token.SetExpiration(payload.ExpiredAt)
	token.SetIssuedAt(payload.IssuedAt)
	token.SetNotBefore(payload.IssuedAt)
	token.SetString("username", payload.Username)
	token.SetString("role", payload.Role)
	token.SetString("token_type", string(payload.Type))
	token.SetString("id", payload.ID.String())

	signed := token.V4Encrypt(maker.symmetricKey, nil)

	return signed, payload, nil
}

// VerifyToken checks if a token is valid or not
func (maker *PasetoMaker) VerifyToken(token string, tokenType TokenType) (*Payload, error) {
	parser := paseto.NewParser()
	parser.AddRule(paseto.ValidAt(time.Now()))

	parsedToken, err := parser.ParseV4Local(maker.symmetricKey, token, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	payload, err := extractPayload(parsedToken)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.Valid(tokenType)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func extractPayload(token *paseto.Token) (*Payload, error) {
	payload := &Payload{}

	idStr, err := token.GetString("id")
	if err != nil {
		return nil, err
	}
	payload.ID, err = uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}

	payload.Username, err = token.GetString("username")
	if err != nil {
		return nil, err
	}

	payload.Role, err = token.GetString("role")
	if err != nil {
		return nil, err
	}

	typeStr, err := token.GetString("token_type")
	if err != nil {
		return nil, err
	}
	if len(typeStr) > 0 {
		payload.Type = TokenType(typeStr[0])
	}

	exp, err := token.GetExpiration()
	if err != nil {
		return nil, err
	}
	payload.ExpiredAt = exp

	iat, err := token.GetIssuedAt()
	if err != nil {
		return nil, err
	}
	payload.IssuedAt = iat

	return payload, nil
}
