package util

import (
	"crypto/sha1"
	"encoding/hex"
)

func CalcSHA(content []byte) string {
	sha := sha1.Sum(content)
	return hex.EncodeToString(sha[:])
}
