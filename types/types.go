package types

import (
	"github.com/google/uuid"
)

type BlockchainMode int

type NewDataPayload struct {
	DeviceId uuid.UUID      `json:"device_id"`
	Data     map[string]any `json:"data"`
	Hash     string         `json:"hash"`
	DataId   uuid.UUID      `json:"data_id"`
}

type DocData struct {
	Id    uuid.UUID      `json:"_id"`
	Data  map[string]any `json:"data"`
	Proof [][]byte       `json:"proof,omitempty"`
}
