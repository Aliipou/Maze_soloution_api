package maze_device

import (
	"context"
	"goapi/internal/api/repository/models"
)

// MazeDeviceStatusService defines the interface for maze device status business logic
type MazeDeviceStatusService interface {
	Create(status *models.MazeDeviceStatus, ctx context.Context) error
	ReadOne(id int, ctx context.Context) (*models.MazeDeviceStatus, error)
	ReadMany(page int, rowsPerPage int, ctx context.Context) ([]*models.MazeDeviceStatus, error)
	ReadByDeviceID(deviceID string, ctx context.Context) ([]*models.MazeDeviceStatus, error)
	Update(status *models.MazeDeviceStatus, ctx context.Context) (int64, error)
	Delete(status *models.MazeDeviceStatus, ctx context.Context) (int64, error)
	ValidateStatus(status *models.MazeDeviceStatus) error
}

// MazeDeviceStatusError represents a business logic error
type MazeDeviceStatusError struct {
	Message string
}

func (e MazeDeviceStatusError) Error() string {
	return e.Message
}
