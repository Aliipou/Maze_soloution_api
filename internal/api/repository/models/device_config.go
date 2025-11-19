package models

import "context"

// DeviceConfig represents configuration settings for a device
type DeviceConfig struct {
	ID               int    `json:"id"`
	DeviceID         string `json:"device_id"`         // Hardware identifier of the Arduino
	AlarmTimeout     int    `json:"alarm_timeout"`     // Alarm timeout in seconds
	SensitivityLevel int    `json:"sensitivity_level"` // Hall sensor sensitivity level (1-10)
	UpdatedAt        string `json:"updated_at"`        // Last update timestamp in RFC3339 format
}

// DeviceConfigRepository defines the interface for device config database operations
type DeviceConfigRepository interface {
	Create(config *DeviceConfig, ctx context.Context) error
	ReadOne(id int, ctx context.Context) (*DeviceConfig, error)
	ReadByDeviceID(deviceID string, ctx context.Context) (*DeviceConfig, error)
	ReadMany(page int, rowsPerPage int, ctx context.Context) ([]*DeviceConfig, error)
	Update(config *DeviceConfig, ctx context.Context) (int64, error)
	Delete(config *DeviceConfig, ctx context.Context) (int64, error)
}
