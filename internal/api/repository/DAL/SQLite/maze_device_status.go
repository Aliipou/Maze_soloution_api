package SQLite

import (
	"context"
	"database/sql"
	"goapi/internal/api/repository/DAL"
	"goapi/internal/api/repository/models"
)

type MazeDeviceStatusRepository struct {
	sqlDB *sql.DB
	createStmt,
	readStmt,
	readManyStmt,
	readByDeviceIDStmt,
	updateStmt,
	deleteStmt *sql.Stmt
	ctx context.Context
}

func NewMazeDeviceStatusRepository(sqlDB DAL.SQLDatabase, ctx context.Context) (models.MazeDeviceStatusRepository, error) {

	repo := &MazeDeviceStatusRepository{
		sqlDB: sqlDB.Connection(),
		ctx:   ctx,
	}

	// Create the maze_device_status table if it doesn't exist
	if _, err := repo.sqlDB.Exec(`CREATE TABLE IF NOT EXISTS maze_device_status (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		device_id VARCHAR(50) NOT NULL,
		alarm_active BOOLEAN NOT NULL,
		maze_completed BOOLEAN NOT NULL,
		hall_sensor_value BOOLEAN NOT NULL,
		battery_level INTEGER NOT NULL CHECK(battery_level >= 0 AND battery_level <= 100),
		timestamp TIMESTAMP NOT NULL
	);`); err != nil {
		repo.sqlDB.Close()
		return nil, err
	}

	// Create index on device_id for faster queries
	if _, err := repo.sqlDB.Exec(`CREATE INDEX IF NOT EXISTS idx_maze_device_status_device_id ON maze_device_status(device_id);`); err != nil {
		repo.sqlDB.Close()
		return nil, err
	}

	// Prepare SQL statements
	createStmt, err := repo.sqlDB.Prepare(`INSERT INTO maze_device_status (device_id, alarm_active, maze_completed, hall_sensor_value, battery_level, timestamp) VALUES (?, ?, ?, ?, ?, ?)`)
	if err != nil {
		repo.sqlDB.Close()
		return nil, err
	}
	repo.createStmt = createStmt

	readStmt, err := repo.sqlDB.Prepare("SELECT id, device_id, alarm_active, maze_completed, hall_sensor_value, battery_level, timestamp FROM maze_device_status WHERE id = ?")
	if err != nil {
		repo.sqlDB.Close()
		return nil, err
	}
	repo.readStmt = readStmt

	readManyStmt, err := repo.sqlDB.Prepare("SELECT id, device_id, alarm_active, maze_completed, hall_sensor_value, battery_level, timestamp FROM maze_device_status LIMIT ? OFFSET ?")
	if err != nil {
		repo.sqlDB.Close()
		return nil, err
	}
	repo.readManyStmt = readManyStmt

	readByDeviceIDStmt, err := repo.sqlDB.Prepare("SELECT id, device_id, alarm_active, maze_completed, hall_sensor_value, battery_level, timestamp FROM maze_device_status WHERE device_id = ? ORDER BY timestamp DESC")
	if err != nil {
		repo.sqlDB.Close()
		return nil, err
	}
	repo.readByDeviceIDStmt = readByDeviceIDStmt

	updateStmt, err := repo.sqlDB.Prepare("UPDATE maze_device_status SET device_id = ?, alarm_active = ?, maze_completed = ?, hall_sensor_value = ?, battery_level = ?, timestamp = ? WHERE id = ?")
	if err != nil {
		repo.sqlDB.Close()
		return nil, err
	}
	repo.updateStmt = updateStmt

	deleteStmt, err := repo.sqlDB.Prepare("DELETE FROM maze_device_status WHERE id = ?")
	if err != nil {
		repo.sqlDB.Close()
		return nil, err
	}
	repo.deleteStmt = deleteStmt

	go CloseMazeDeviceStatus(ctx, repo)

	return repo, nil
}

func CloseMazeDeviceStatus(ctx context.Context, r *MazeDeviceStatusRepository) {
	<-ctx.Done()
	r.createStmt.Close()
	r.readStmt.Close()
	r.updateStmt.Close()
	r.deleteStmt.Close()
	r.readManyStmt.Close()
	r.readByDeviceIDStmt.Close()
	r.sqlDB.Close()
}

func (r *MazeDeviceStatusRepository) Create(status *models.MazeDeviceStatus, ctx context.Context) error {
	res, err := r.createStmt.ExecContext(ctx, status.DeviceID, status.AlarmActive, status.MazeCompleted, status.HallSensorValue, status.BatteryLevel, status.Timestamp)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	status.ID = int(id)
	return nil
}

func (r *MazeDeviceStatusRepository) ReadOne(id int, ctx context.Context) (*models.MazeDeviceStatus, error) {
	row := r.readStmt.QueryRowContext(ctx, id)
	var status models.MazeDeviceStatus
	err := row.Scan(&status.ID, &status.DeviceID, &status.AlarmActive, &status.MazeCompleted, &status.HallSensorValue, &status.BatteryLevel, &status.Timestamp)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &status, nil
}

func (r *MazeDeviceStatusRepository) ReadMany(page int, rowsPerPage int, ctx context.Context) ([]*models.MazeDeviceStatus, error) {
	if page < 1 {
		return r.ReadAll(ctx)
	}

	offset := rowsPerPage * (page - 1)
	rows, err := r.readManyStmt.QueryContext(ctx, rowsPerPage, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var statuses []*models.MazeDeviceStatus
	for rows.Next() {
		var s models.MazeDeviceStatus
		err := rows.Scan(&s.ID, &s.DeviceID, &s.AlarmActive, &s.MazeCompleted, &s.HallSensorValue, &s.BatteryLevel, &s.Timestamp)
		if err != nil {
			return nil, err
		}
		statuses = append(statuses, &s)
	}
	return statuses, nil
}

func (r *MazeDeviceStatusRepository) ReadByDeviceID(deviceID string, ctx context.Context) ([]*models.MazeDeviceStatus, error) {
	rows, err := r.readByDeviceIDStmt.QueryContext(ctx, deviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var statuses []*models.MazeDeviceStatus
	for rows.Next() {
		var s models.MazeDeviceStatus
		err := rows.Scan(&s.ID, &s.DeviceID, &s.AlarmActive, &s.MazeCompleted, &s.HallSensorValue, &s.BatteryLevel, &s.Timestamp)
		if err != nil {
			return nil, err
		}
		statuses = append(statuses, &s)
	}
	return statuses, nil
}

func (r *MazeDeviceStatusRepository) ReadAll(ctx context.Context) ([]*models.MazeDeviceStatus, error) {
	rows, err := r.sqlDB.QueryContext(ctx, "SELECT id, device_id, alarm_active, maze_completed, hall_sensor_value, battery_level, timestamp FROM maze_device_status ORDER BY timestamp DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var statuses []*models.MazeDeviceStatus
	for rows.Next() {
		var s models.MazeDeviceStatus
		err := rows.Scan(&s.ID, &s.DeviceID, &s.AlarmActive, &s.MazeCompleted, &s.HallSensorValue, &s.BatteryLevel, &s.Timestamp)
		if err != nil {
			return nil, err
		}
		statuses = append(statuses, &s)
	}
	return statuses, nil
}

func (r *MazeDeviceStatusRepository) Update(status *models.MazeDeviceStatus, ctx context.Context) (int64, error) {
	res, err := r.updateStmt.ExecContext(ctx, status.DeviceID, status.AlarmActive, status.MazeCompleted, status.HallSensorValue, status.BatteryLevel, status.Timestamp, status.ID)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}

func (r *MazeDeviceStatusRepository) Delete(status *models.MazeDeviceStatus, ctx context.Context) (int64, error) {
	res, err := r.deleteStmt.ExecContext(ctx, status.ID)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return rowsAffected, nil
}
