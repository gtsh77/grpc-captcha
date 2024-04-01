package tools

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
)

func RandHex(bytes int) (string, error) {
	var (
		rr   *rand.Rand
		hash []byte
		err  error
	)

	hash = make([]byte, bytes)
	rr = rand.New(rand.NewSource(time.Now().UnixNano()))

	if _, err = rr.Read(hash); err != nil {
		return "", fmt.Errorf("rand.Read: %w", err)
	}

	return hex.EncodeToString(hash), nil
}
