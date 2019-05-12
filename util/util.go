package util

import (
	"crypto/sha1"
	"encoding/hex"
)

func CalcSHA(content []byte) string {
	sha := sha1.Sum(content)
	return hex.EncodeToString(sha[:])
}

func FindInRunes(rs []rune, c rune, start int) int {
	for i, r := range rs {
		if i < start {
			continue
		}
		if r == c {
			return i
		}
	}
	return -1
}
