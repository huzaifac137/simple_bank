package token

import "time"

// manages creating and verifying tokens
type Maker interface {
	createToken(username string, duration time.Duration) (string, error)
	verifyToken(token string) (*Payload, error)
}
