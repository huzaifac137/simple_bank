package token

import (
	"fmt"
	"time"

	paseto "github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) < keyMinSize {
		return nil, fmt.Errorf("secret key must be at least %d characters long", keyMinSize)
	}
	return &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}, nil
}

func (maker *PasetoMaker) createToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	pasetToken, err := paseto.NewV2().Encrypt(maker.symmetricKey, payload, nil)
	if err != nil {
		return "", err
	}

	return pasetToken, nil
}
func (maker *PasetoMaker) verifyToken(token string) (*Payload, error) {
	payload := &Payload{}
	payloadErr := paseto.NewV2().Decrypt(token, maker.symmetricKey, payload, nil)
	if payloadErr != nil {
		return nil, fmt.Errorf("invalid token")
	}

	return payload, nil
}
