package maze_device

import (
	"context"
	"encoding/json"
	"goapi/internal/api/service/maze_device"
	"log"
	"net/http"
	"strconv"
	"time"
)

// GetByIDHandler handles GET requests to retrieve a specific maze device status by ID
// curl -X GET http://127.0.0.1:8080/device/status/1 -u admin:password
func GetByIDHandler(w http.ResponseWriter, r *http.Request, logger *log.Logger, service maze_device.MazeDeviceStatusService) {
	// Extract ID from URL path parameter
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Invalid ID format."}`))
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	// Retrieve the status from the database
	status, err := service.ReadOne(id, ctx)
	if err != nil {
		logger.Println("Error reading maze device status:", err)
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}

	if status == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "Maze device status not found."}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(status); err != nil {
		logger.Println("Error encoding maze device status:", err)
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
}
