package device_config

import (
	"context"
	"goapi/internal/api/repository/models"
)

// DeviceConfigService defines the interface for device config business logic
type DeviceConfigService interface {
	Create(config *models.DeviceConfig, ctx context.Context) error
	ReadOne(id int, ctx context.Context) (*models.DeviceConfig, error)
	ReadByDeviceID(deviceID string, ctx context.Context) (*models.DeviceConfig, error)
	ReadMany(page int, rowsPerPage int, ctx context.Context) ([]*models.DeviceConfig, error)
	Update(config *models.DeviceConfig, ctx context.Context) (int64, error)
	Delete(config *models.DeviceConfig, ctx context.Context) (int64, error)
	ValidateConfig(config *models.DeviceConfig) error
}

// DeviceConfigError represents a business logic error
type DeviceConfigError struct {
	Message string
}

func (e DeviceConfigError) Error() string {
	return e.Message
}
