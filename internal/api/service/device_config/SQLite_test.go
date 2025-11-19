package device_config

import (
	"goapi/internal/api/repository/models"
	"testing"
	"time"
)

func TestValidateConfig(t *testing.T) {
	service := &DeviceConfigServiceSQLite{repo: nil}

	tests := []struct {
		name        string
		config      models.DeviceConfig
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid config",
			config: models.DeviceConfig{
				DeviceID:         "ARD001",
				AlarmTimeout:     300,
				SensitivityLevel: 5,
				UpdatedAt:        time.Now().Add(-1 * time.Minute).Format(time.RFC3339),
			},
			expectError: false,
		},
		{
			name: "Valid config with minimum values",
			config: models.DeviceConfig{
				DeviceID:         "ARD002",
				AlarmTimeout:     1,
				SensitivityLevel: 1,
				UpdatedAt:        time.Now().Format(time.RFC3339),
			},
			expectError: false,
		},
		{
			name: "Valid config with maximum values",
			config: models.DeviceConfig{
				DeviceID:         "ARD003",
				AlarmTimeout:     3600,
				SensitivityLevel: 10,
				UpdatedAt:        time.Now().Format(time.RFC3339),
			},
			expectError: false,
		},
		{
			name: "Empty device_id",
			config: models.DeviceConfig{
				DeviceID:         "",
				AlarmTimeout:     300,
				SensitivityLevel: 5,
				UpdatedAt:        time.Now().Format(time.RFC3339),
			},
			expectError: true,
			errorMsg:    "device_id is required",
		},
		{
			name: "Device_id too long",
			config: models.DeviceConfig{
				DeviceID:         "ThisIsAVeryLongDeviceIDThatExceedsTheFiftyCharacterLimitDefinedInTheSchema",
				AlarmTimeout:     300,
				SensitivityLevel: 5,
				UpdatedAt:        time.Now().Format(time.RFC3339),
			},
			expectError: true,
			errorMsg:    "device_id is required and must be less than 50 characters",
		},
		{
			name: "Alarm timeout zero",
			config: models.DeviceConfig{
				DeviceID:         "ARD004",
				AlarmTimeout:     0,
				SensitivityLevel: 5,
				UpdatedAt:        time.Now().Format(time.RFC3339),
			},
			expectError: true,
			errorMsg:    "alarm_timeout must be between 1 and 3600 seconds",
		},
		{
			name: "Alarm timeout negative",
			config: models.DeviceConfig{
				DeviceID:         "ARD005",
				AlarmTimeout:     -100,
				SensitivityLevel: 5,
				UpdatedAt:        time.Now().Format(time.RFC3339),
			},
			expectError: true,
			errorMsg:    "alarm_timeout must be between 1 and 3600 seconds",
		},
		{
			name: "Alarm timeout too large",
			config: models.DeviceConfig{
				DeviceID:         "ARD006",
				AlarmTimeout:     5000,
				SensitivityLevel: 5,
				UpdatedAt:        time.Now().Format(time.RFC3339),
			},
			expectError: true,
			errorMsg:    "alarm_timeout must be between 1 and 3600 seconds",
		},
		{
			name: "Sensitivity level zero",
			config: models.DeviceConfig{
				DeviceID:         "ARD007",
				AlarmTimeout:     300,
				SensitivityLevel: 0,
				UpdatedAt:        time.Now().Format(time.RFC3339),
			},
			expectError: true,
			errorMsg:    "sensitivity_level must be between 1 and 10",
		},
		{
			name: "Sensitivity level negative",
			config: models.DeviceConfig{
				DeviceID:         "ARD008",
				AlarmTimeout:     300,
				SensitivityLevel: -5,
				UpdatedAt:        time.Now().Format(time.RFC3339),
			},
			expectError: true,
			errorMsg:    "sensitivity_level must be between 1 and 10",
		},
		{
			name: "Sensitivity level too high",
			config: models.DeviceConfig{
				DeviceID:         "ARD009",
				AlarmTimeout:     300,
				SensitivityLevel: 15,
				UpdatedAt:        time.Now().Format(time.RFC3339),
			},
			expectError: true,
			errorMsg:    "sensitivity_level must be between 1 and 10",
		},
		{
			name: "Invalid updated_at format",
			config: models.DeviceConfig{
				DeviceID:         "ARD010",
				AlarmTimeout:     300,
				SensitivityLevel: 5,
				UpdatedAt:        "2024-01-15 10:30:00",
			},
			expectError: true,
			errorMsg:    "updated_at must be in RFC3339 format",
		},
		{
			name: "UpdatedAt in the future",
			config: models.DeviceConfig{
				DeviceID:         "ARD011",
				AlarmTimeout:     300,
				SensitivityLevel: 5,
				UpdatedAt:        time.Now().Add(10 * time.Minute).Format(time.RFC3339),
			},
			expectError: true,
			errorMsg:    "updated_at must not be in the future",
		},
		{
			name: "Multiple validation errors",
			config: models.DeviceConfig{
				DeviceID:         "",
				AlarmTimeout:     5000,
				SensitivityLevel: 20,
				UpdatedAt:        "invalid-timestamp",
			},
			expectError: true,
			errorMsg:    "device_id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ValidateConfig(&tt.config)

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
