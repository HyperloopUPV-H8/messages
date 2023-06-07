package main

import (
	"bytes"
	"encoding/binary"
	"math/rand"
)

type OrderGenerator struct {
	addId       uint16
	removeId    uint16
	stateOrders map[string][]uint16
	boardToId   map[string]uint16
}

type StateOrders struct {
	Id      uint16
	BoardId uint16
	Len     uint8
	Orders  []uint16
}

func (o StateOrders) Bytes() ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, &o.Id)

	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.LittleEndian, &o.BoardId)

	if err != nil {
		return nil, err
	}

	err = binary.Write(buf, binary.LittleEndian, &o.Len)

	if err != nil {
		return nil, err
	}

	for _, order := range o.Orders {
		err = binary.Write(buf, binary.LittleEndian, &order)

		if err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func NewOrderGenerator(addId uint16, removeId uint16, stateOrders map[string][]uint16, boardToId map[string]uint16) OrderGenerator {
	return OrderGenerator{
		addId:       addId,
		removeId:    removeId,
		stateOrders: stateOrders,
		boardToId:   boardToId,
	}
}

func (generator OrderGenerator) New() StateOrders {
	boardName := RandKey(generator.stateOrders)
	boardId := generator.boardToId[boardName]

	orderNum := RandInt(len(generator.stateOrders[boardName]))
	orders := &Set[uint16]{}
	for i := 0; i < orderNum; i++ {
		orders.Add(generator.stateOrders[boardName][RandInt(len(generator.stateOrders[boardName]))])
	}

	var id uint16
	if rand.Intn(2) == 0 {
		id = generator.addId
	} else {
		id = generator.removeId
	}

	return StateOrders{
		Id:      id,
		BoardId: boardId,
		Len:     byte(len(orders.AsSlice())),
		Orders:  orders.AsSlice(),
	}

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
