package utils

import (
	"math/rand"
	"strings"
)

const allAlphabets = "abcdefghijklmnopqrstuvwxyz"

func GenerateRandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func GenerateRandomString(n int) string {
	var sb strings.Builder
	k := len(allAlphabets)

	for i := 0; i < n; i++ {
		c := allAlphabets[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func GenerateRandomAmount() int64 {
	return GenerateRandomInt(0, 1000)
}
