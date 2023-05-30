package main

import (
	"math/rand"

	vehicle_models "github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"
)

type MessageGenerator struct {
	boardIds map[string]uint16
	kinds    map[string]uint16
}

func NewMessageGenerator(kinds, boardIds map[string]uint16) MessageGenerator {
	return MessageGenerator{
		boardIds: boardIds,
		kinds:    kinds,
	}
}

func (generator MessageGenerator) New() (any, uint16) {
	kind := RandKey(generator.kinds)

	if kind == "info" {
		return generator.getInfoMessage()
	}

	return generator.getProtectionMessage(kind)
}

func (generator MessageGenerator) getInfoMessage() (InfoMessageAdapter, uint16) {
	return InfoMessageAdapter{
		BoardId:   RandVal(generator.boardIds),
		Timestamp: randomTimestamp(),
		Msg:       "We are about to win, if you are a member of the femenine genere, please be aware of Juan. You have been informed.",
	}, generator.kinds["info"]
}

func (generator MessageGenerator) getProtectionMessage(kind string) (ProtectionMessageAdapter, uint16) {
	protection := generator.randomProtection(kind)

	return ProtectionMessageAdapter{
		BoardId:    RandVal(generator.boardIds),
		Timestamp:  randomTimestamp(),
		Protection: protection,
	}, generator.kinds[kind]
}

var protectionKinds = []string{"OUT_OF_BOUNDS", "LOWER_BOUND", "UPPER_BOUND", "NOT_EQUALS", "EQUALS", "TIME_ACCUMULATION", "ERROR_HANDLER"}

func (generator MessageGenerator) randomProtection(kind string) ProtectionAdapter {
	protectionKind := protectionKinds[rand.Intn(len(protectionKinds))]

	return ProtectionAdapter{
		Name: "VCELL1",
		Type: protectionKind,
		Data: randomProtectionData(protectionKind),
	}
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
		return "Starting booster!"
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
