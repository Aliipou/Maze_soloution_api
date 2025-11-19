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

// PutHandler handles PUT requests to update device config
// curl -X PUT http://127.0.0.1:8080/device/config -u admin:password -H "Content-Type: application/json" -d '{"id":1,"device_id":"ARD001","alarm_timeout":600,"sensitivity_level":7,"updated_at":"2024-01-15T10:35:00Z"}'
func PutHandler(w http.ResponseWriter, r *http.Request, logger *log.Logger, service device_config.DeviceConfigService) {
	var config models.DeviceConfig

	// Decode the JSON payload from the request body
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Invalid request data. Please check your input."}`))
		return
	}

	// Validate that ID is provided
	if config.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "ID is required for update."}`))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	// Try to update the config in the database
	rowsAffected, err := service.Update(&config, ctx)
	if err != nil {
		switch err.(type) {
		case device_config.DeviceConfigError:
			// Client error: validation failed
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "` + err.Error() + `"}`))
			return
		default:
			// Server error
			logger.Println("Error updating device config:", err, config)
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
			return
		}
	}

	if rowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "Device config not found."}`))
		return
	}

	// Return the updated config with 200 OK
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(config); err != nil {
		logger.Println("Error encoding device config:", err, config)
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
}
