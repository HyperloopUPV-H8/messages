package main

import (
	"math/rand"

	vehicle_models "github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"
)

type MessageGenerator struct {
	boardIds []string
}

func NewMessageGenerator() MessageGenerator {
	return MessageGenerator{
		boardIds: []string{"LCU_MASTER", "LCU_SLAVE", "VCU", "BMSH", "BMSL"},
	}
}

func (gen *MessageGenerator) New() (any, int) {
	kinds := []string{"protection", "info"}

	kind := kinds[rand.Intn(len(kinds))]

	if kind == "protection" {
		return getProtectionMessage()

	} else {
		return getInfoMessage()
	}

}

func getInfoMessage() (InfoMessageAdapter, int) {
	return InfoMessageAdapter{
		BoardId:   2,
		Timestamp: randomTimestamp(),
		Msg:       "asdfljdsf asdlfj alsjd fasdk kflkjsda k kdsaj jsdlkfjalñsjdfajkdslfa klasjdflkad > 10",
	}, 1
}

func getProtectionMessage() (ProtectionMessageAdapter, int) {
	protection, id := randomProtection()

	return ProtectionMessageAdapter{
		BoardId:    2,
		Timestamp:  randomTimestamp(),
		Protection: protection,
	}, id
}

func randomProtection() (ProtectionAdapter, int) {
	messagesKinds := map[string]int{"warning": 2, "fault": 3}
	kind := RandKey(messagesKinds)

	protectionKinds := []string{"OUT_OF_BOUNDS", "LOWER_BOUND", "UPPER_BOUND", "NOT_EQUALS", "EQUALS", "TIME_ACCUMULATION", "ERROR_HANDLER"}

	protectionKind := protectionKinds[rand.Intn(len(protectionKinds))]

	return ProtectionAdapter{
		Name: "VCELL1",
		Type: protectionKind,
		Data: randomProtectionData(protectionKind),
	}, messagesKinds[kind]

}

func randomProtectionData(kind string) any {
	switch kind {
	case "OUT_OF_BOUNDS":
		return vehicle_models.OutOfBounds{
			Value:  rand.Float64() * 100,
			Bounds: [2]float64{rand.Float64() * 100, rand.Float64() * 100},
		}
	case "LOWER_BOUND":
		return vehicle_models.LowerBound{
			Value: rand.Float64() * 100,
			Bound: rand.Float64() * 100,
		}
	case "UPPER_BOUND":
		return vehicle_models.UpperBound{
			Value: rand.Float64() * 100,
			Bound: rand.Float64() * 100,
		}
	case "EQUALS":
		return vehicle_models.Equals{
			Value: rand.Float64() * 100,
		}
	case "NOT_EQUALS":
		return vehicle_models.NotEquals{
			Value: rand.Float64() * 100,
			Want:  rand.Float64() * 100,
		}
	case "TIME_ACCUMULATION":
		return vehicle_models.TimeLimit{
			Value:     rand.Float64() * 100,
			Bound:     rand.Float64() * 100,
			TimeLimit: rand.Float64() * 100,
		}
	case "ERROR_HANDLER":
		return "This is an error"
	default:
		return vehicle_models.NotEquals{
			Value: rand.Float64() * 100,
			Want:  rand.Float64() * 100,
		}
	}
}

func randomTimestamp() vehicle_models.Timestamp {
	return vehicle_models.Timestamp{
		Counter: uint16(rand.Int()),
		Second:  uint8(rand.Int()),
		Minute:  uint8(rand.Int()),
		Hour:    uint8(rand.Int()),
		Day:     uint8(rand.Int()),
		Month:   uint8(rand.Int()),
		Year:    uint16(rand.Int()),
	}
}
