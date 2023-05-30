package main

import (
	vehicle_models "github.com/HyperloopUPV-H8/Backend-H8/vehicle/models"
)

type InfoMessageAdapter struct {
	BoardId   uint16                   `json:"boardId"`
	Timestamp vehicle_models.Timestamp `json:"timestamp"`
	Msg       string                   `json:"msg"`
}

type ProtectionMessageAdapter struct {
	BoardId    uint16                   `json:"boardId"`
	Timestamp  vehicle_models.Timestamp `json:"timestamp"`
	Protection ProtectionAdapter        `json:"protection"`
}

type ProtectionAdapter struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Data any    `json:"data"`
}
