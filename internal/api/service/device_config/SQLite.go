package device_config

import (
	"context"
	"goapi/internal/api/repository/models"
	"time"
)

// DeviceConfigServiceSQLite implements DeviceConfigService for SQLite
type DeviceConfigServiceSQLite struct {
	repo models.DeviceConfigRepository
}

func NewDeviceConfigServiceSQLite(repo models.DeviceConfigRepository) *DeviceConfigServiceSQLite {
	return &DeviceConfigServiceSQLite{
		repo: repo,
	}
}

func (s *DeviceConfigServiceSQLite) Create(config *models.DeviceConfig, ctx context.Context) error {
	if err := s.ValidateConfig(config); err != nil {
		return DeviceConfigError{Message: "Invalid device config: " + err.Error()}
	}
	return s.repo.Create(config, ctx)
}

func (s *DeviceConfigServiceSQLite) ReadOne(id int, ctx context.Context) (*models.DeviceConfig, error) {
	config, err := s.repo.ReadOne(id, ctx)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (s *DeviceConfigServiceSQLite) ReadByDeviceID(deviceID string, ctx context.Context) (*models.DeviceConfig, error) {
	if deviceID == "" {
		return nil, DeviceConfigError{Message: "device_id is required"}
	}
	config, err := s.repo.ReadByDeviceID(deviceID, ctx)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (s *DeviceConfigServiceSQLite) ReadMany(page int, rowsPerPage int, ctx context.Context) ([]*models.DeviceConfig, error) {
	return s.repo.ReadMany(page, rowsPerPage, ctx)
}

func (s *DeviceConfigServiceSQLite) Update(config *models.DeviceConfig, ctx context.Context) (int64, error) {
	if err := s.ValidateConfig(config); err != nil {
		return 0, DeviceConfigError{Message: "Invalid device config: " + err.Error()}
	}
	return s.repo.Update(config, ctx)
}

func (s *DeviceConfigServiceSQLite) Delete(config *models.DeviceConfig, ctx context.Context) (int64, error) {
	return s.repo.Delete(config, ctx)
}

// ValidateConfig validates the device config according to the requirements
func (s *DeviceConfigServiceSQLite) ValidateConfig(config *models.DeviceConfig) error {
	var errMsg string

	// Validate device_id (required, max 50 chars)
	if config.DeviceID == "" || len(config.DeviceID) > 50 {
		errMsg += "device_id is required and must be less than 50 characters. "
	}

	// Validate alarm_timeout (must be positive, reasonable range 1-3600 seconds = 1 hour)
	if config.AlarmTimeout < 1 || config.AlarmTimeout > 3600 {
		errMsg += "alarm_timeout must be between 1 and 3600 seconds. "
	}

	// Validate sensitivity_level (must be 1-10)
	if config.SensitivityLevel < 1 || config.SensitivityLevel > 10 {
		errMsg += "sensitivity_level must be between 1 and 10. "
	}

	// Validate updated_at format (RFC3339) and not in the future
	updatedAt, err := time.Parse(time.RFC3339, config.UpdatedAt)
	if err != nil {
		errMsg += "updated_at must be in RFC3339 format (e.g., 2006-01-02T15:04:05Z07:00). "
	} else {
		// Check if updated_at is not in the future (allow 1 minute tolerance for clock skew)
		if updatedAt.After(time.Now().Add(1 * time.Minute)) {
			errMsg += "updated_at must not be in the future. "
		}
	}

	if errMsg != "" {
		return DeviceConfigError{Message: errMsg}
	}
	return nil
}
