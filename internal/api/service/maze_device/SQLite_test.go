package maze_device

import (
	"goapi/internal/api/repository/models"
	"testing"
	"time"
)

func TestValidateStatus(t *testing.T) {
	service := &MazeDeviceStatusServiceSQLite{repo: nil}

	tests := []struct {
		name        string
		status      models.MazeDeviceStatus
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid status",
			status: models.MazeDeviceStatus{
				DeviceID:        "ARD001",
				AlarmActive:     true,
				MazeCompleted:   false,
				HallSensorValue: false,
				BatteryLevel:    85,
				Timestamp:       time.Now().Add(-1 * time.Minute).Format(time.RFC3339),
			},
			expectError: false,
		},
		{
			name: "Valid status with maze completed and sensor true",
			status: models.MazeDeviceStatus{
				DeviceID:        "ARD002",
				AlarmActive:     false,
				MazeCompleted:   true,
				HallSensorValue: true,
				BatteryLevel:    90,
				Timestamp:       time.Now().Add(-5 * time.Minute).Format(time.RFC3339),
			},
			expectError: false,
		},
		{
			name: "Empty device_id",
			status: models.MazeDeviceStatus{
				DeviceID:        "",
				AlarmActive:     true,
				MazeCompleted:   false,
				HallSensorValue: false,
				BatteryLevel:    50,
				Timestamp:       time.Now().Format(time.RFC3339),
			},
			expectError: true,
			errorMsg:    "device_id is required",
		},
		{
			name: "Device_id too long",
			status: models.MazeDeviceStatus{
				DeviceID:        "ThisIsAVeryLongDeviceIDThatExceedsTheFiftyCharacterLimitDefinedInTheSchema",
				AlarmActive:     true,
				MazeCompleted:   false,
				HallSensorValue: false,
				BatteryLevel:    50,
				Timestamp:       time.Now().Format(time.RFC3339),
			},
			expectError: true,
			errorMsg:    "device_id is required and must be less than 50 characters",
		},
		{
			name: "Battery level negative",
			status: models.MazeDeviceStatus{
				DeviceID:        "ARD003",
				AlarmActive:     true,
				MazeCompleted:   false,
				HallSensorValue: false,
				BatteryLevel:    -10,
				Timestamp:       time.Now().Format(time.RFC3339),
			},
			expectError: true,
			errorMsg:    "battery_level must be between 0 and 100",
		},
		{
			name: "Battery level above 100",
			status: models.MazeDeviceStatus{
				DeviceID:        "ARD004",
				AlarmActive:     true,
				MazeCompleted:   false,
				HallSensorValue: false,
				BatteryLevel:    150,
				Timestamp:       time.Now().Format(time.RFC3339),
			},
			expectError: true,
			errorMsg:    "battery_level must be between 0 and 100",
		},
		{
			name: "Battery level at 0 (edge case - valid)",
			status: models.MazeDeviceStatus{
				DeviceID:        "ARD005",
				AlarmActive:     true,
				MazeCompleted:   false,
				HallSensorValue: false,
				BatteryLevel:    0,
				Timestamp:       time.Now().Format(time.RFC3339),
			},
			expectError: false,
		},
		{
			name: "Battery level at 100 (edge case - valid)",
			status: models.MazeDeviceStatus{
				DeviceID:        "ARD006",
				AlarmActive:     true,
				MazeCompleted:   false,
				HallSensorValue: false,
				BatteryLevel:    100,
				Timestamp:       time.Now().Format(time.RFC3339),
			},
			expectError: false,
		},
		{
			name: "Invalid timestamp format",
			status: models.MazeDeviceStatus{
				DeviceID:        "ARD007",
				AlarmActive:     true,
				MazeCompleted:   false,
				HallSensorValue: false,
				BatteryLevel:    50,
				Timestamp:       "2024-01-15 10:30:00",
			},
			expectError: true,
			errorMsg:    "timestamp must be in RFC3339 format",
		},
		{
			name: "Timestamp in the future",
			status: models.MazeDeviceStatus{
				DeviceID:        "ARD008",
				AlarmActive:     true,
				MazeCompleted:   false,
				HallSensorValue: false,
				BatteryLevel:    50,
				Timestamp:       time.Now().Add(10 * time.Minute).Format(time.RFC3339),
			},
			expectError: true,
			errorMsg:    "timestamp must not be in the future",
		},
		{
			name: "Maze completed but sensor false (inconsistent state)",
			status: models.MazeDeviceStatus{
				DeviceID:        "ARD009",
				AlarmActive:     false,
				MazeCompleted:   true,
				HallSensorValue: false,
				BatteryLevel:    70,
				Timestamp:       time.Now().Format(time.RFC3339),
			},
			expectError: true,
			errorMsg:    "maze_completed cannot be true when hall_sensor_value is false",
		},
		{
			name: "Multiple validation errors",
			status: models.MazeDeviceStatus{
				DeviceID:        "",
				AlarmActive:     true,
				MazeCompleted:   false,
				HallSensorValue: false,
				BatteryLevel:    200,
				Timestamp:       "invalid-timestamp",
			},
			expectError: true,
			errorMsg:    "device_id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ValidateStatus(&tt.status)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error containing '%s', got nil", tt.errorMsg)
				} else if tt.errorMsg != "" {
					errStr := err.Error()
					found := false
					if len(tt.errorMsg) > 0 {
						for i := 0; i <= len(errStr)-len(tt.errorMsg); i++ {
							if errStr[i:i+len(tt.errorMsg)] == tt.errorMsg {
								found = true
								break
							}
						}
					}
					if !found {
						t.Errorf("Expected error containing '%s', got '%s'", tt.errorMsg, errStr)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}
			}
		})
	}
}
