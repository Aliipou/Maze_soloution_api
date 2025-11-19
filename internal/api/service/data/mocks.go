package data

import (
	"context"
	"goapi/internal/api/repository/models"
)

// * Mock implementation of DataService for testing purposes, always returns a successful response and Data object(s) *
type MockDataServiceSuccessful struct{}

func (m *MockDataServiceSuccessful) ReadMany(page int, rowsPerPage int, ctx context.Context) ([]*models.Data, error) {
	return []*models.Data{
		{
			ID:          1,
			DeviceID:    "device1",
			DeviceName:  "device1",
			Value:       1.0,
			Type:        "type1",
			DateTime:    "2021-01-01 00:00:00",
			Description: "description1",
		},
		{
			ID:          2,
			DeviceID:    "device2",
			DeviceName:  "device2",
			Value:       2.0,
			Type:        "type2",
			DateTime:    "2021-01-01 00:00:00",
			Description: "description2",
		},
	}, nil
}

func (m *MockDataServiceSuccessful) ReadOne(id int, ctx context.Context) (*models.Data, error) {
	return &models.Data{
		ID:          1,
		DeviceID:    "device1",
		DeviceName:  "device1",
		Value:       1.0,
		Type:        "type1",
		DateTime:    "2021-01-01 00:00:00",
		Description: "description1",
	}, nil
}

func (m *MockDataServiceSuccessful) Create(data *models.Data, ctx context.Context) error {
	return nil
}

func (m *MockDataServiceSuccessful) Update(data *models.Data, ctx context.Context) (int64, error) {
	return 1, nil
}

func (m *MockDataServiceSuccessful) Delete(data *models.Data, ctx context.Context) (int64, error) {
	return 1, nil
}

func (m *MockDataServiceSuccessful) ValidateData(data *models.Data) error {
	return nil
}

// * Mock implementation of DataService for testing purposes, always returns empty data *

type MockDataServiceNotFound struct{}

func (m *MockDataServiceNotFound) ReadMany(page int, rowsPerPage int, ctx context.Context) ([]*models.Data, error) {
	return []*models.Data{}, nil
}

func (m *MockDataServiceNotFound) ReadOne(id int, ctx context.Context) (*models.Data, error) {
	return nil, nil
}

func (m *MockDataServiceNotFound) Create(data *models.Data, ctx context.Context) error {
	return nil
}

func (m *MockDataServiceNotFound) Update(data *models.Data, ctx context.Context) (int64, error) {
	return 0, nil
}

func (m *MockDataServiceNotFound) Delete(data *models.Data, ctx context.Context) (int64, error) {
	return 0, nil
}

func (m *MockDataServiceNotFound) ValidateData(data *models.Data) error {
	return nil
}

// * Mock implementation of DataService for testing purposes, always returns an error *
type MockDataServiceError struct{}

func (m *MockDataServiceError) ReadMany(page int, rowsPerPage int, ctx context.Context) ([]*models.Data, error) {
	return nil, DataError{Message: "Error reading data."}
}

func (m *MockDataServiceError) ReadOne(id int, ctx context.Context) (*models.Data, error) {
	return nil, DataError{Message: "Error reading data."}
}

func (m *MockDataServiceError) Create(data *models.Data, ctx context.Context) error {
	return DataError{Message: "Error creating data."}
}

func (m *MockDataServiceError) Update(data *models.Data, ctx context.Context) (int64, error) {
	return 0, DataError{Message: "Error updating data."}
}

func (m *MockDataServiceError) Delete(data *models.Data, ctx context.Context) (int64, error) {
	return 0, DataError{Message: "Error deleting data."}
}

func (m *MockDataServiceError) ValidateData(data *models.Data) error {
	return nil
}
