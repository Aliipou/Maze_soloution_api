package service

import (
	"context"
	"goapi/internal/api/repository/DAL"
	"goapi/internal/api/repository/DAL/SQLite"
	service "goapi/internal/api/service/data"
	"goapi/internal/api/service/device_config"
	"goapi/internal/api/service/maze_device"
	"log"
)

type DataServiceType int

const (
	SQLiteDataService DataServiceType = iota
)

type ServiceFactory struct {
	db     DAL.SQLDatabase
	logger *log.Logger
	ctx    context.Context
}

// * Factory for creating data service *
func NewServiceFactory(db DAL.SQLDatabase, logger *log.Logger, ctx context.Context) *ServiceFactory {
	return &ServiceFactory{
		db:     db,
		logger: logger,
		ctx:    ctx,
	}
}

func (sf *ServiceFactory) CreateDataService(serviceType DataServiceType) (*service.DataServiceSQLite, error) {

	switch serviceType {

	case SQLiteDataService:
		repo, err := SQLite.NewDataRepository(sf.db, sf.ctx)
		if err != nil {
			return nil, err
		}
		ds := service.NewDataServiceSQLite(repo)
		return ds, nil
	default:
		return nil, service.DataError{Message: "Invalid data service type."}
	}
}

func (sf *ServiceFactory) CreateMazeDeviceStatusService(serviceType DataServiceType) (*maze_device.MazeDeviceStatusServiceSQLite, error) {

	switch serviceType {

	case SQLiteDataService:
		repo, err := SQLite.NewMazeDeviceStatusRepository(sf.db, sf.ctx)
		if err != nil {
			return nil, err
		}
		service := maze_device.NewMazeDeviceStatusServiceSQLite(repo)
		return service, nil
	default:
		return nil, maze_device.MazeDeviceStatusError{Message: "Invalid service type."}
	}
}

func (sf *ServiceFactory) CreateDeviceConfigService(serviceType DataServiceType) (*device_config.DeviceConfigServiceSQLite, error) {

	switch serviceType {

	case SQLiteDataService:
		repo, err := SQLite.NewDeviceConfigRepository(sf.db, sf.ctx)
		if err != nil {
			return nil, err
		}
		service := device_config.NewDeviceConfigServiceSQLite(repo)
		return service, nil
	default:
		return nil, device_config.DeviceConfigError{Message: "Invalid service type."}
	}
}
