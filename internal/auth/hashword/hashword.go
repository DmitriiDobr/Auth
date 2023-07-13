package hashword

import (
	"crypto/sha1"
	"fmt"
)

type Hash struct {
	salt string
}

func NewHash(salt string) *Hash {
	return &Hash{salt}
}

func (h *Hash) Generate(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(h.salt)))
}
