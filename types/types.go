package types

import (
	"time"

	"github.com/google/uuid"
)

type Device struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"createdAt"`
}

type DeviceType string

type IotData struct {
	DeviceId   uuid.UUID      `json:"device_id"`
	DeviceType DeviceType     `json:"device_type"`
	Region     string         `json:"region"`
	Data       map[string]any `json:"data"`
	PublicKey  string         `json:"public_key"`
	Signature  string         `json:"signature"`
}

type Client interface {
	Register(Device) error
	IsDeviceRegistered(id uuid.UUID) (*Device, error)
	GetDevice(id uuid.UUID) (*Device, error)
	Send(data IotData) error
	GetData(deviceType DeviceType, region string) (*IotData, error)
}

type RegisterDevicePayload struct {
}

type NewDataPayload struct {
}
