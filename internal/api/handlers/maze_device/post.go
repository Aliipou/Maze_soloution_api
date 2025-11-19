package maze_device

import (
	"context"
	"encoding/json"
	"goapi/internal/api/repository/models"
	"goapi/internal/api/service/maze_device"
	"log"
	"net/http"
	"time"
)

// PostHandler handles POST requests to create new maze device status
// curl -X POST http://127.0.0.1:8080/device/status -u admin:password -H "Content-Type: application/json" -d '{"device_id":"ARD001","alarm_active":true,"maze_completed":false,"hall_sensor_value":false,"battery_level":85,"timestamp":"2024-01-15T10:30:00Z"}'
func PostHandler(w http.ResponseWriter, r *http.Request, logger *log.Logger, service maze_device.MazeDeviceStatusService) {
	var status models.MazeDeviceStatus

	// Decode the JSON payload from the request body
	if err := json.NewDecoder(r.Body).Decode(&status); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Invalid request data. Please check your input."}`))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	// Try to create the status in the database
	if err := service.Create(&status, ctx); err != nil {
		switch err.(type) {
		case maze_device.MazeDeviceStatusError:
			// Client error: validation failed
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "` + err.Error() + `"}`))
			return
		default:
			// Server error
			logger.Println("Error creating maze device status:", err, status)
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
			return
		}
	}

	// Return the created status with 201 Created
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(status); err != nil {
		logger.Println("Error encoding maze device status:", err, status)
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
}
