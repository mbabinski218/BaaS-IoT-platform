package types

import (
	"time"

	"github.com/google/uuid"
)

type Device struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	CreatedAt time.Time `json:"createdAt"`
}

type IotData struct {
	EncryptedData string `json:"encrypted_data"`
	PublicKey     string `json:"public_key"`
	Signature     string `json:"signature"`
	Type          string `json:"type"`
}

type DeviceStore interface {
	IsDeviceRegistered(id int) (*Device, error)
	GetDeviceByID(id int) (*Device, error)
	RegisterDevice(Device) error
	PushData(data IotData) error
}

type RegisterDevicePayload struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	Price       float64 `json:"price" validate:"required"`
	Quantity    int     `json:"quantity" validate:"required"`
}
