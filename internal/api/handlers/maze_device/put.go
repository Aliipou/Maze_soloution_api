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

// PutHandler handles PUT requests to update maze device status
// curl -X PUT http://127.0.0.1:8080/device/status -u admin:password -H "Content-Type: application/json" -d '{"id":1,"device_id":"ARD001","alarm_active":false,"maze_completed":true,"hall_sensor_value":true,"battery_level":80,"timestamp":"2024-01-15T10:35:00Z"}'
func PutHandler(w http.ResponseWriter, r *http.Request, logger *log.Logger, service maze_device.MazeDeviceStatusService) {
	var status models.MazeDeviceStatus

	// Decode the JSON payload from the request body
	if err := json.NewDecoder(r.Body).Decode(&status); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Invalid request data. Please check your input."}`))
		return
	}

	// Validate that ID is provided
	if status.ID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "ID is required for update."}`))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	// Try to update the status in the database
	rowsAffected, err := service.Update(&status, ctx)
	if err != nil {
		switch err.(type) {
		case maze_device.MazeDeviceStatusError:
			// Client error: validation failed
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "` + err.Error() + `"}`))
			return
		default:
			// Server error
			logger.Println("Error updating maze device status:", err, status)
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
			return
		}
	}

	if rowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "Maze device status not found."}`))
		return
	}

	// Return the updated status with 200 OK
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(status); err != nil {
		logger.Println("Error encoding maze device status:", err, status)
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
}
