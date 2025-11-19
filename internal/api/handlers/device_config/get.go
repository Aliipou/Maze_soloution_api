package device_config

import (
	"context"
	"encoding/json"
	"goapi/internal/api/service/device_config"
	"log"
	"net/http"
	"strconv"
	"time"
)

// GetHandler handles GET requests to retrieve device configs
// Supports pagination: GET /device/config?page=1&rows_per_page=10
// Supports filtering by device_id: GET /device/config?device_id=ARD001
// curl -X GET "http://127.0.0.1:8080/device/config?page=1&rows_per_page=10" -u admin:password
// curl -X GET "http://127.0.0.1:8080/device/config?device_id=ARD001" -u admin:password
func GetHandler(w http.ResponseWriter, r *http.Request, logger *log.Logger, service device_config.DeviceConfigService) {
	// Parse query parameters
	pageStr := r.URL.Query().Get("page")
	rowsPerPageStr := r.URL.Query().Get("rows_per_page")
	deviceID := r.URL.Query().Get("device_id")

	page, _ := strconv.Atoi(pageStr)
	rowsPerPage, _ := strconv.Atoi(rowsPerPageStr)

	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	// If device_id is provided, return that specific config
	if deviceID != "" {
		config, err := service.ReadByDeviceID(deviceID, ctx)
		if err != nil {
			switch err.(type) {
			case device_config.DeviceConfigError:
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"error": "` + err.Error() + `"}`))
				return
			default:
				logger.Println("Error reading device config by device_id:", err)
				http.Error(w, "Internal server error.", http.StatusInternalServerError)
				return
			}
		}
		if config == nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"error": "Device config not found."}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(config); err != nil {
			logger.Println("Error encoding device config:", err)
			http.Error(w, "Internal server error.", http.StatusInternalServerError)
			return
		}
		return
	}

	// Otherwise, return paginated results
	configs, err := service.ReadMany(page, rowsPerPage, ctx)
	if err != nil {
		logger.Println("Error reading device configs:", err)
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(configs); err != nil {
		logger.Println("Error encoding device configs:", err)
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}
}
