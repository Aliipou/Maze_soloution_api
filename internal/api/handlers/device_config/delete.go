package device_config

import (
	"context"
	"goapi/internal/api/repository/models"
	"goapi/internal/api/service/device_config"
	"log"
	"net/http"
	"strconv"
	"time"
)

// DeleteHandler handles DELETE requests to remove a device config
// curl -X DELETE http://127.0.0.1:8080/device/config/1 -u admin:password
func DeleteHandler(w http.ResponseWriter, r *http.Request, logger *log.Logger, service device_config.DeviceConfigService) {
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

	// Delete the config from the database
	config := &models.DeviceConfig{ID: id}
	rowsAffected, err := service.Delete(config, ctx)
	if err != nil {
		logger.Println("Error deleting device config:", err)
		http.Error(w, "Internal server error.", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "Device config not found."}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Device config deleted successfully."}`))
}
