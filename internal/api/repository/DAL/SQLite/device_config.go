package SQLite

import (
	"context"
	"database/sql"
	"goapi/internal/api/repository/DAL"
	"goapi/internal/api/repository/models"
)

type DeviceConfigRepository struct {
	sqlDB *sql.DB
	createStmt,
	readStmt,
	readByDeviceIDStmt,
	readManyStmt,
	updateStmt,
	deleteStmt *sql.Stmt
	ctx context.Context
}

func NewDeviceConfigRepository(sqlDB DAL.SQLDatabase, ctx context.Context) (models.DeviceConfigRepository, error) {

	repo := &DeviceConfigRepository{
		sqlDB: sqlDB.Connection(),
		ctx:   ctx,
	}

	// Create the device_config table if it doesn't exist
	if _, err := repo.sqlDB.Exec(`CREATE TABLE IF NOT EXISTS device_config (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		device_id VARCHAR(50) NOT NULL UNIQUE,
		alarm_timeout INTEGER NOT NULL,
		sensitivity_level INTEGER NOT NULL CHECK(sensitivity_level >= 1 AND sensitivity_level <= 10),
		updated_at TIMESTAMP NOT NULL
	);`); err != nil {
		repo.sqlDB.Close()
		return nil, err
	}

	// Create index on device_id for faster queries
	if _, err := repo.sqlDB.Exec(`CREATE INDEX IF NOT EXISTS idx_device_config_device_id ON device_config(device_id);`); err != nil {
		repo.sqlDB.Close()
		return nil, err
	}

	// Prepare SQL statements
	createStmt, err := repo.sqlDB.Prepare(`INSERT INTO device_config (device_id, alarm_timeout, sensitivity_level, updated_at) VALUES (?, ?, ?, ?)`)
	if err != nil {
		repo.sqlDB.Close()
		return nil, err
	}
	repo.createStmt = createStmt

	readStmt, err := repo.sqlDB.Prepare("SELECT id, device_id, alarm_timeout, sensitivity_level, updated_at FROM device_config WHERE id = ?")
	if err != nil {
		repo.sqlDB.Close()
		return nil, err
	}
	repo.readStmt = readStmt

	readByDeviceIDStmt, err := repo.sqlDB.Prepare("SELECT id, device_id, alarm_timeout, sensitivity_level, updated_at FROM device_config WHERE device_id = ?")
	if err != nil {
		repo.sqlDB.Close()
		return nil, err
	}
	repo.readByDeviceIDStmt = readByDeviceIDStmt

	readManyStmt, err := repo.sqlDB.Prepare("SELECT id, device_id, alarm_timeout, sensitivity_level, updated_at FROM device_config LIMIT ? OFFSET ?")
	if err != nil {
		repo.sqlDB.Close()
		return nil, err
	}
	repo.readManyStmt = readManyStmt

	updateStmt, err := repo.sqlDB.Prepare("UPDATE device_config SET device_id = ?, alarm_timeout = ?, sensitivity_level = ?, updated_at = ? WHERE id = ?")
	if err != nil {
		repo.sqlDB.Close()
		return nil, err
	}
	repo.updateStmt = updateStmt

	deleteStmt, err := repo.sqlDB.Prepare("DELETE FROM device_config WHERE id = ?")
	if err != nil {
		repo.sqlDB.Close()
		return nil, err
	}
	repo.deleteStmt = deleteStmt

	go CloseDeviceConfig(ctx, repo)

	return repo, nil
}

func CloseDeviceConfig(ctx context.Context, r *DeviceConfigRepository) {
	<-ctx.Done()
	r.createStmt.Close()
	r.readStmt.Close()
	r.readByDeviceIDStmt.Close()
	r.updateStmt.Close()
	r.deleteStmt.Close()
	r.readManyStmt.Close()
	r.sqlDB.Close()
}

func (r *DeviceConfigRepository) Create(config *models.DeviceConfig, ctx context.Context) error {
	res, err := r.createStmt.ExecContext(ctx, config.DeviceID, config.AlarmTimeout, config.SensitivityLevel, config.UpdatedAt)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	config.ID = int(id)
	return nil
}

func (r *DeviceConfigRepository) ReadOne(id int, ctx context.Context) (*models.DeviceConfig, error) {
	row := r.readStmt.QueryRowContext(ctx, id)
	var config models.DeviceConfig
	err := row.Scan(&config.ID, &config.DeviceID, &config.AlarmTimeout, &config.SensitivityLevel, &config.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &config, nil
}

func (r *DeviceConfigRepository) ReadByDeviceID(deviceID string, ctx context.Context) (*models.DeviceConfig, error) {
	row := r.readByDeviceIDStmt.QueryRowContext(ctx, deviceID)
	var config models.DeviceConfig
	err := row.Scan(&config.ID, &config.DeviceID, &config.AlarmTimeout, &config.SensitivityLevel, &config.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &config, nil
}

func (r *DeviceConfigRepository) ReadMany(page int, rowsPerPage int, ctx context.Context) ([]*models.DeviceConfig, error) {
	if page < 1 {
		return r.ReadAll(ctx)
	}

	offset := rowsPerPage * (page - 1)
	rows, err := r.readManyStmt.QueryContext(ctx, rowsPerPage, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []*models.DeviceConfig
	for rows.Next() {
		var c models.DeviceConfig
		err := rows.Scan(&c.ID, &c.DeviceID, &c.AlarmTimeout, &c.SensitivityLevel, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		configs = append(configs, &c)
	}
	return configs, nil
}

func (r *DeviceConfigRepository) ReadAll(ctx context.Context) ([]*models.DeviceConfig, error) {
	rows, err := r.sqlDB.QueryContext(ctx, "SELECT id, device_id, alarm_timeout, sensitivity_level, updated_at FROM device_config")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []*models.DeviceConfig
	for rows.Next() {
		var c models.DeviceConfig
		err := rows.Scan(&c.ID, &c.DeviceID, &c.AlarmTimeout, &c.SensitivityLevel, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		configs = append(configs, &c)
	}
	return configs, nil
}

func (r *DeviceConfigRepository) Update(config *models.DeviceConfig, ctx context.Context) (int64, error) {
	res, err := r.updateStmt.ExecContext(ctx, config.DeviceID, config.AlarmTimeout, config.SensitivityLevel, config.UpdatedAt, config.ID)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

func (r *DeviceConfigRepository) Delete(config *models.DeviceConfig, ctx context.Context) (int64, error) {
	res, err := r.deleteStmt.ExecContext(ctx, config.ID)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}
