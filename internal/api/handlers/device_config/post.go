package device_config

import (
	"context"
	"encoding/json"
	"goapi/internal/api/repository/models"
	"goapi/internal/api/service/device_config"
	"log"
	"net/http"
	"time"
)

// PostHandler handles POST requests to create new device config
// curl -X POST http://127.0.0.1:8080/device/config -u admin:password -H "Content-Type: application/json" -d '{"device_id":"ARD001","alarm_timeout":300,"sensitivity_level":5,"updated_at":"2024-01-15T10:30:00Z"}'
func PostHandler(w http.ResponseWriter, r *http.Request, logger *log.Logger, service device_config.DeviceConfigService) {
	var config models.DeviceConfig

	// Decode the JSON payload from the request body
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Invalid request data. Please check your input."}`))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	// Try to create the config in the database
	if err := service.Create(&config, ctx); err != nil {
		switch err.(type) {
		case device_config.DeviceConfigError:
			// Client error: validation failed
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "` + err.Error() + `"}`))
			return
		default:
			// Server error (could be duplicate device_id due to UNIQUE constraint)
			logger.Println("Error creating device config:", err, config)
			// Check if it's a unique constraint error
			if err.Error() == "UNIQUE constraint failed: device_config.device_id" {
				w.WriteHeader(http.StatusConflict)
				w.Write([]byte(`{"error": "Device config already exists for this device_id."}`))
				return
			}
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
			return
		}
	}

	// Return the created config with 201 Created
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(config); err != nil {
		logger.Println("Error encoding device config:", err, config)
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
}
