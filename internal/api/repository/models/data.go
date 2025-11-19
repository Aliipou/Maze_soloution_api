package models

import "context"

type Data struct {
	ID          int     `json:"id"`
	DeviceID    string  `json:"device_id"`
	DeviceName  string  `json:"device_name"`
	Value       float64 `json:"value"`
	Type        string  `json:"type"`
	DateTime    string  `json:"date_time"`
	Description string  `json:"description"`
}

type DataRepository interface {
	Create(Data *Data, ctx context.Context) error
	ReadOne(id int, ctx context.Context) (*Data, error)
	ReadMany(page int, rowsPerPage int, ctx context.Context) ([]*Data, error)
	Update(data *Data, ctx context.Context) (int64, error)
	Delete(data *Data, ctx context.Context) (int64, error)
}
