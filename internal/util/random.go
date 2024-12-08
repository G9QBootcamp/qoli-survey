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

func ShuffleSlice[T any](slice []T) []T {
	rand.Seed(uint64(time.Now().UnixNano()))
	shuffled := make([]T, len(slice))
	copy(shuffled, slice)

	for i := len(shuffled) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}

	return shuffled
}
