package models

import "context"

// MazeDeviceStatus represents the status of a maze-based waking device
// The device uses a simple mechanical maze with metal balls and a Hall sensor
type MazeDeviceStatus struct {
	ID              int    `json:"id"`
	DeviceID        string `json:"device_id"`        // Hardware identifier of the Arduino
	AlarmActive     bool   `json:"alarm_active"`     // Whether the alarm is currently active
	MazeCompleted   bool   `json:"maze_completed"`   // Whether the maze has been completed
	HallSensorValue bool   `json:"hall_sensor_value"` // Hall sensor detection (true = ball detected at end)
	BatteryLevel    int    `json:"battery_level"`    // Battery level 0-100
	Timestamp       string `json:"timestamp"`        // Server timestamp in RFC3339 format
}

// MazeDeviceStatusRepository defines the interface for maze device status database operations
type MazeDeviceStatusRepository interface {
	Create(status *MazeDeviceStatus, ctx context.Context) error
	ReadOne(id int, ctx context.Context) (*MazeDeviceStatus, error)
	ReadMany(page int, rowsPerPage int, ctx context.Context) ([]*MazeDeviceStatus, error)
	ReadByDeviceID(deviceID string, ctx context.Context) ([]*MazeDeviceStatus, error)
	Update(status *MazeDeviceStatus, ctx context.Context) (int64, error)
	Delete(status *MazeDeviceStatus, ctx context.Context) (int64, error)
}
