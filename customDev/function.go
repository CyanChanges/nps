package customDev

import (
	"math/rand"
	"time"
)

// Random one string with numbers and strings
func RandStr(length int) string {
	letterRunes := []rune("0123456789abcdefghijklmnopqrstuvwxyz")
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
