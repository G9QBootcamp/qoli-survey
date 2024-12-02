package util

import (
	"time"

	"golang.org/x/exp/rand"
)

func GenerateNumericString(length int) string {
	rand.Seed(uint64(time.Now().UnixNano())) // Seed the random number generator
	numericString := make([]byte, length)

	for i := range numericString {
		numericString[i] = '0' + byte(rand.Intn(10)) // Generate a random digit between 0 and 9
	}

	return string(numericString)
}
