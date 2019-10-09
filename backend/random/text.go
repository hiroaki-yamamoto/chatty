package random

import (
	"math/rand"
	"strings"
)

// GenerateRandomText generates random string.
func GenerateRandomText(charMap string, size int) string {
	var buf strings.Builder
	for i := 0; i < size; i++ {
		buf.WriteByte(charMap[rand.Intn(len(charMap))])
	}
	return buf.String()
}
