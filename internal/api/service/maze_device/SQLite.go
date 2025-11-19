package maze_device

import (
	"context"
	"goapi/internal/api/repository/models"
	"time"
)

// MazeDeviceStatusServiceSQLite implements MazeDeviceStatusService for SQLite
type MazeDeviceStatusServiceSQLite struct {
	repo models.MazeDeviceStatusRepository
}

func NewMazeDeviceStatusServiceSQLite(repo models.MazeDeviceStatusRepository) *MazeDeviceStatusServiceSQLite {
	return &MazeDeviceStatusServiceSQLite{
		repo: repo,
	}
}

func (s *MazeDeviceStatusServiceSQLite) Create(status *models.MazeDeviceStatus, ctx context.Context) error {
	if err := s.ValidateStatus(status); err != nil {
		return MazeDeviceStatusError{Message: "Invalid maze device status: " + err.Error()}
	}
	return s.repo.Create(status, ctx)
}

func (s *MazeDeviceStatusServiceSQLite) ReadOne(id int, ctx context.Context) (*models.MazeDeviceStatus, error) {
	status, err := s.repo.ReadOne(id, ctx)
	if err != nil {
		return nil, err
	}
	return status, nil
}

func (s *MazeDeviceStatusServiceSQLite) ReadMany(page int, rowsPerPage int, ctx context.Context) ([]*models.MazeDeviceStatus, error) {
	return s.repo.ReadMany(page, rowsPerPage, ctx)
}

func (s *MazeDeviceStatusServiceSQLite) ReadByDeviceID(deviceID string, ctx context.Context) ([]*models.MazeDeviceStatus, error) {
	if deviceID == "" {
		return nil, MazeDeviceStatusError{Message: "device_id is required"}
	}
	return s.repo.ReadByDeviceID(deviceID, ctx)
}

func (s *MazeDeviceStatusServiceSQLite) Update(status *models.MazeDeviceStatus, ctx context.Context) (int64, error) {
	if err := s.ValidateStatus(status); err != nil {
		return 0, MazeDeviceStatusError{Message: "Invalid maze device status: " + err.Error()}
	}
	return s.repo.Update(status, ctx)
}

func (s *MazeDeviceStatusServiceSQLite) Delete(status *models.MazeDeviceStatus, ctx context.Context) (int64, error) {
	return s.repo.Delete(status, ctx)
}

// ValidateStatus validates the maze device status according to the requirements
func (s *MazeDeviceStatusServiceSQLite) ValidateStatus(status *models.MazeDeviceStatus) error {
	var errMsg string

	// Validate device_id (required, max 50 chars)
	if status.DeviceID == "" || len(status.DeviceID) > 50 {
		errMsg += "device_id is required and must be less than 50 characters. "
	}

	// Validate battery_level (must be 0-100)
	if status.BatteryLevel < 0 || status.BatteryLevel > 100 {
		errMsg += "battery_level must be between 0 and 100. "
	}

	// Validate timestamp format (RFC3339) and not in the future
	timestamp, err := time.Parse(time.RFC3339, status.Timestamp)
	if err != nil {
		errMsg += "timestamp must be in RFC3339 format (e.g., 2006-01-02T15:04:05Z07:00). "
	} else {
		// Check if timestamp is not in the future (allow 1 minute tolerance for clock skew)
		if timestamp.After(time.Now().Add(1 * time.Minute)) {
			errMsg += "timestamp must not be in the future. "
		}
	}

	// Business logic validation: if maze_completed is true, hall_sensor_value should be true
	if status.MazeCompleted && !status.HallSensorValue {
		errMsg += "maze_completed cannot be true when hall_sensor_value is false. "
	}

	if errMsg != "" {
		return MazeDeviceStatusError{Message: errMsg}
	}
	return nil
}
