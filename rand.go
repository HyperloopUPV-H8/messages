package main

import "math/rand"

func RandKey[K comparable, V any](m map[K]V) K {
	index := RandInt(len(m))

	i := 0
	for key := range m {
		if i == index {
			return key
		}

		i++
	}

	panic("failed to generate random key")
}

func RandVal[K comparable, V any](m map[K]V) V {
	index := rand.Intn(len(m))

	i := 0
	for _, val := range m {
		if i == index {
			return val
		}

		i++
	}

	panic("failed to generate random value")
}

func RandInt(n int) int {
	if n == 0 {
		return 0
	}

	return rand.Intn(n)
}
