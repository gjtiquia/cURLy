package random

import "math/rand/v2"

func Range(minInclusive, maxExclusive int) int {

	// in-case users swapped the order
	if minInclusive > maxExclusive {
		minInclusive, maxExclusive = maxExclusive, minInclusive
	}

	return rand.N(maxExclusive-minInclusive) + minInclusive

}
