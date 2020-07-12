package id

import "math/rand"

var letterRunes = []rune("abcdefghijklmnpqrstuvwxyz1234567890")

func New(length int) string {
	b := make([]rune, length)

	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b)
}
