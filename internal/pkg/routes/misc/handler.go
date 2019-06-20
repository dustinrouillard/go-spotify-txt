package misc

import (
	"net/http"
)

// NotFound is the handler for a 404
func NotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"code": "endpoint_not_found", "message": "Endpoint does not exist"}`))
}
