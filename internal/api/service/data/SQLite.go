package data

import (
	"context"
	"goapi/internal/api/repository/models"
	"time"
)

// * Implementation of DataService for SQLite database *
type DataServiceSQLite struct {
	repo models.DataRepository
}

func NewDataServiceSQLite(repo models.DataRepository) *DataServiceSQLite {
	return &DataServiceSQLite{
		repo: repo,
	}
}

func (ds *DataServiceSQLite) Create(data *models.Data, ctx context.Context) error {

	if err := ds.ValidateData(data); err != nil {
		return DataError{Message: "InvalMockDataServiceSuccessfulid data."}
	}
	return ds.repo.Create(data, ctx)
}

func (ds *DataServiceSQLite) ReadOne(id int, ctx context.Context) (*models.Data, error) {

	data, err := ds.repo.ReadOne(id, ctx)
	if err != nil {
		return nil, err
	}

	_ = data

	// Tehdään datalle jotain, päätellään datasta jotain!!!
	// Tämä ohjaa toimintaa älykkäästi, esim. jos data on tietynlaista, niin tehdään jotain

	return data, nil
}

func (ds *DataServiceSQLite) ReadMany(page int, rowsPerPage int, ctx context.Context) ([]*models.Data, error) {
	return ds.repo.ReadMany(page, rowsPerPage, ctx)
}

func (ds *DataServiceSQLite) Update(data *models.Data, ctx context.Context) (int64, error) {

	if err := ds.ValidateData(data); err != nil {
		return 0, DataError{Message: "Invalid data: " + err.Error()}
	}
	return ds.repo.Update(data, ctx)
}

func (ds *DataServiceSQLite) Delete(data *models.Data, ctx context.Context) (int64, error) {
	return ds.repo.Delete(data, ctx)
}

func (ds *DataServiceSQLite) ValidateData(data *models.Data) error {
	var errMsg string
	if data.DeviceID == "" || len(data.DeviceID) > 50 {
		errMsg += "DeviceID is required and must be less than 50 characters. "
	}
	if len(data.DeviceName) > 50 {
		errMsg += "DeviceName must be less than 50 characters. "
	}
	if len(data.Type) > 20 {
		errMsg += "Type must be less than 20 characters. "
	}
	if len(data.Description) > 100 {
		errMsg += "Description must be less than 100 characters. "
	}
	_, err := time.Parse("2006-01-02T15:04:05Z", data.DateTime)
	if err != nil {
		errMsg += "DateTime must be in the format: 2021-01-01T12:00:00Z. "
	}
	if errMsg != "" {
		return DataError{Message: errMsg}
	}
	return nil
}
