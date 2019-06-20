package authorization

import (
	"log"
	"net/http"

	"github.com/thetetrabyte/go-spotify-txt/internal/app/spotify"
)

// Authorize handler for returning the authorization url for spotify
func Authorize(w http.ResponseWriter, r *http.Request) {
	// Get authorization url
	authorizationURL := spotify.GetAuthorizationURL()

	// Create response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"code": "authorization_required", "message": "` + authorizationURL + `"}`))
}

// Callback handler for handling the incoming callback from the authorization
func Callback(w http.ResponseWriter, r *http.Request) {
	authorizationCode := r.URL.Query().Get("code")

	_, tokenErr, err := spotify.RequestUserTokens(authorizationCode)
	if err != nil {
		log.Println(err)

		// Create response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"code": "failed", "message": "failed"}`))

		return
	}

	if (spotify.TokenError{}) != tokenErr {
		log.Println(tokenErr)

		// Create response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"code": "failed", "message": "failed"}`))

		return
	}

	// Create response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"code": "success", "message": "success"}`))
}
