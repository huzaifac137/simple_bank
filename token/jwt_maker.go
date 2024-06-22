package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTMaker struct {
	secretKey string
}

const keyMinSize int = 32

func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < keyMinSize {
		return nil, fmt.Errorf("secret key must be at least %d characters long", keyMinSize)
	}
	return &JWTMaker{secretKey}, nil

}

func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(maker.secretKey))
}

func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, InvalidTokenError
		}

		return []byte(maker.secretKey), nil
	}

	jwtToken, er := jwt.ParseWithClaims(token, &Payload{}, keyFunc)

	if er != nil {
		verr, ok := er.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrExpiredToken) {
			return nil, ErrExpiredToken
		}
		return nil, InvalidTokenError
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, InvalidTokenError
	}

	return payload, nil
}
