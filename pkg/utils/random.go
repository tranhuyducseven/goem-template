package utils

import (
	"crypto/rand"
	"math/big"

	"github.com/rs/zerolog/log"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var numberRunes = []rune("0123456789")

// random string generator

func RandomString(n int) (string, error) {
	b := make([]rune, n)
	for i := range b {
		random, err := rand.Int(rand.Reader, big.NewInt(int64(len(letterRunes))))
		if err != nil {
			log.Err(err).Msg("failed to generate random string")
			return "", err
		}
		b[i] = letterRunes[random.Int64()]
	}
	return string(b), nil
}

func RandomNumber(n int) (string, error) {
	b := make([]rune, n)
	for i := range b {
		random, err := rand.Int(rand.Reader, big.NewInt(int64(len(numberRunes))))
		if err != nil {
			log.Err(err).Msg("failed to generate random number")
			return "", err
		}
		b[i] = numberRunes[random.Int64()]
	}
	return string(b), nil
}
