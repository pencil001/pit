package util

import (
	"bytes"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"strings"
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

func FindInBytes(bs []byte, c byte, start int) int {
	for i, r := range bs {
		if i < start {
			continue
		}
		if r == c {
			return i
		}
	}
	return -1
}

func BytesToHexStr(bs []byte) string {
	var sb strings.Builder
	for _, b := range bs {
		sb.WriteString(fmt.Sprintf("%02x", b))
	}
	return sb.String()
}

func HexStrToBytes(hexStr string) []byte {
	buf := new(bytes.Buffer)
	for i := 0; i < len(hexStr); i += 2 {
		ih, err := strconv.ParseUint(hexStr[i:i+2], 16, 8)
		if err != nil {
			log.Panic(err)
		}
		if err := binary.Write(buf, binary.LittleEndian, uint8(ih)); err != nil {
			log.Panic(err)
		}
	}
	return buf.Bytes()
}
