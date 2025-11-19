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

// GetHandler handles GET requests to retrieve multiple maze device statuses
// Supports pagination: GET /device/status?page=1&rows_per_page=10
// curl -X GET "http://127.0.0.1:8080/device/status?page=1&rows_per_page=10" -u admin:password
func GetHandler(w http.ResponseWriter, r *http.Request, logger *log.Logger, service maze_device.MazeDeviceStatusService) {
	// Parse query parameters for pagination
	pageStr := r.URL.Query().Get("page")
	rowsPerPageStr := r.URL.Query().Get("rows_per_page")
	deviceID := r.URL.Query().Get("device_id")

	page, _ := strconv.Atoi(pageStr)
	rowsPerPage, _ := strconv.Atoi(rowsPerPageStr)

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	// If device_id is provided, filter by device_id
	if deviceID != "" {
		statuses, err := service.ReadByDeviceID(deviceID, ctx)
		if err != nil {
			switch err.(type) {
			case maze_device.MazeDeviceStatusError:
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"error": "` + err.Error() + `"}`))
				return
			default:
				logger.Println("Error reading maze device statuses by device_id:", err)
				http.Error(w, "Internal server error.", http.StatusInternalServerError)
				return
			}
		}
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(statuses); err != nil {
			logger.Println("Error encoding maze device statuses:", err)
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
			return
		}
		return
	}

	// Otherwise, return paginated results
	statuses, err := service.ReadMany(page, rowsPerPage, ctx)
	if err != nil {
		logger.Println("Error reading maze device statuses:", err)
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(statuses); err != nil {
		logger.Println("Error encoding maze device statuses:", err)
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
}
