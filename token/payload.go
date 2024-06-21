package token

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Payload struct {
	Id         uuid.UUID `json:"id"`
	Username   string    `json:"username"`
	IssuedAt   time.Time `json:"issued_at"`
	ExpiredAt  time.Time `json:"expired_at"`
	Expiration time.Time `json:"expiration"`
}

var (
	ErrExpiredToken   error = errors.New("token has expired!!")
	InvalidTokenError error = errors.New("token is invalid")
)

func NewPayload(username string, duration time.Duration) (*Payload, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	return &Payload{
		Id:         id,
		Username:   username,
		IssuedAt:   time.Now(),
		ExpiredAt:  time.Now().Add(duration),
		Expiration: time.Now().Add(duration),
	}, nil
}

func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return ErrExpiredToken
	}
	return nil
}
