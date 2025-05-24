package types

import (
	"github.com/google/uuid"
)

type DeviceType string

// type Device struct {
// 	Id        uuid.UUID  `json:"id"`
// 	Type      DeviceType `json:"type"`
// 	PublicKey string     `json:"public_key"`
// }

// type IotData struct {
// 	DeviceId   uuid.UUID      `json:"device_id"`
// 	DeviceType DeviceType     `json:"device_type"`
// 	Region     string         `json:"region"`
// 	Data       map[string]any `json:"data"`
// 	PublicKey  string         `json:"public_key"`
// 	Signature  string         `json:"signature"`
// }

// type Client interface {
// 	Register(Device) error
// 	IsDeviceRegistered(id uuid.UUID) (*Device, error)
// 	GetDevice(id uuid.UUID) (*Device, error)
// 	Send(data IotData) error
// 	GetData(deviceType DeviceType, region string) (*IotData, error)
// }

// type RegisterDevicePayload struct {
// 	DeviceType DeviceType `json:"device_type"`
// 	Region     string     `json:"region"`
// 	PublicKey  string     `json:"public_key"`
// }

type NewDataPayload struct {
	DeviceId uuid.UUID      `json:"device_id"`
	Data     map[string]any `json:"data"`
	Hash     string         `json:"hash"`
	DataId   uuid.UUID      `json:"data_id"`
	// Signature string         `json:"signature"`
}
