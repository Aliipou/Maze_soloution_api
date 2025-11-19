package device_config

import "net/http"

// OptionsHandler handles OPTIONS requests for CORS preflight
func OptionsHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
