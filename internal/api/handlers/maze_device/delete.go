package maze_device

import (
	"context"
	"goapi/internal/api/repository/models"
	"goapi/internal/api/service/maze_device"
	"log"
	"net/http"
	"strconv"
	"time"
)

// DeleteHandler handles DELETE requests to remove a maze device status
// curl -X DELETE http://127.0.0.1:8080/device/status/1 -u admin:password
func DeleteHandler(w http.ResponseWriter, r *http.Request, logger *log.Logger, service maze_device.MazeDeviceStatusService) {
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

	// Delete the status from the database
	status := &models.MazeDeviceStatus{ID: id}
	rowsAffected, err := service.Delete(status, ctx)
	if err != nil {
		logger.Println("Error deleting maze device status:", err)
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "Maze device status not found."}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Maze device status deleted successfully."}`))
}
