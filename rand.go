package main

import "math/rand"

func RandKey[K comparable, V any](m map[K]V) K {
	index := rand.Intn(len(m))

	i := 0
	for key := range m {
		if i == index {
			return key
		}

		i++
	}

	panic("should have returned")
}
