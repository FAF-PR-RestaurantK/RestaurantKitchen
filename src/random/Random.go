package random

import "math/rand"

func Range(min int, max int) int {
	diff := max - min
	value := rand.Intn(diff) + min

	return value
}
