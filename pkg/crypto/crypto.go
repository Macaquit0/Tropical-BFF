package crypto

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
)

type Crypto struct {
}

func New() *Crypto {
	return &Crypto{}
}

func (c *Crypto) Hash(ctx context.Context, text string) (string, error) {
	hash := sha512.New()
	hash.Write([]byte(text))
	hashedBytes := hash.Sum(nil)
	hashedString := hex.EncodeToString(hashedBytes)
	return hashedString, nil
}
