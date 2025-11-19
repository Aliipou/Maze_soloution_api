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

// GetByIDHandler handles GET requests to retrieve a specific device config by ID
// curl -X GET http://127.0.0.1:8080/device/config/1 -u admin:password
func GetByIDHandler(w http.ResponseWriter, r *http.Request, logger *log.Logger, service device_config.DeviceConfigService) {
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

	// Retrieve the config from the database
	config, err := service.ReadOne(id, ctx)
	if err != nil {
		logger.Println("Error reading device config:", err)
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
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
}
