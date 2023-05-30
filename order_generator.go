package main

import (
	"encoding/binary"
	"math/rand"
)

type OrderGenerator struct {
	id          uint16
	stateOrders map[string][]uint16
	boardToId   map[string]uint16
}

func NewOrderGenerator(id uint16, stateOrders map[string][]uint16, boardToId map[string]uint16) OrderGenerator {
	return OrderGenerator{
		id:          id,
		stateOrders: stateOrders,
		boardToId:   boardToId,
	}
}

func (generator OrderGenerator) New() ([]byte, error) {
	buffer := []byte{}
	buffer = binary.LittleEndian.AppendUint16(buffer, generator.id)

	board := RandKey(generator.stateOrders)
	buffer = binary.LittleEndian.AppendUint16(buffer, generator.boardToId[board])

	orderNum := rand.Intn(len(generator.stateOrders[board]))

	orders := &Set[uint16]{}
	for i := 0; i < orderNum; i++ {
		orders.Add(generator.stateOrders[board][rand.Intn(len(generator.stateOrders[board]))])
	}

	for _, order := range orders.AsSlice() {
		buffer = binary.LittleEndian.AppendUint16(buffer, order)
	}

	return buffer, nil
}

type Set[T comparable] map[T]struct{}

func (set *Set[T]) Add(value T) {
	(*set)[value] = struct{}{}
}

func (set *Set[T]) AsSlice() []T {
	elements := make([]T, 0, len(*set))
	for elem := range *set {
		elements = append(elements, elem)
	}
	return elements
}
