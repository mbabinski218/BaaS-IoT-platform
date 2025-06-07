package types

import (
	"github.com/google/uuid"
)

type NewDataPayload struct {
	DeviceId uuid.UUID      `json:"device_id"`
	Data     map[string]any `json:"data"`
	Hash     string         `json:"hash"`
	DataId   uuid.UUID      `json:"data_id"`
	// Signature string         `json:"signature"`
}

type AuditData struct {
	Id   uuid.UUID      `json:"_id"`
	Data map[string]any `json:"data"`
}
