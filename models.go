package main

import (
	"encoding/binary"
	"encoding/json"

	vehicle_models "github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"
)

type Message interface {
	Bytes() ([]byte, error)
}

type InfoMessageAdapter struct {
	id        uint16
	BoardId   uint16                   `json:"boardId"`
	Timestamp vehicle_models.Timestamp `json:"timestamp"`
	Msg       string                   `json:"msg"`
}

func (adapter InfoMessageAdapter) Bytes() ([]byte, error) {
	buf, err := json.Marshal(adapter)

	if err != nil {
		return nil, err
	}

	return addHeader(int(adapter.id), buf), nil
}

type ProtectionMessageAdapter struct {
	id         uint16
	BoardId    uint16                   `json:"boardId"`
	Timestamp  vehicle_models.Timestamp `json:"timestamp"`
	Protection ProtectionAdapter        `json:"protection"`
}

func (adapter ProtectionMessageAdapter) Bytes() ([]byte, error) {
	buf, err := json.Marshal(adapter)

	if err != nil {
		return nil, err
	}

	return addHeader(int(adapter.id), buf), nil
}

type ProtectionAdapter struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Data any    `json:"data"`
}

func addHeader(id int, payload []byte) []byte {
	buffer := make([]byte, 2)
	binary.LittleEndian.PutUint16(buffer, uint16(id))
	buffer = append(buffer, '\n', '\n')
	buffer = append(buffer, payload...)
	buffer = append(buffer, 0x00)

	return buffer
}
