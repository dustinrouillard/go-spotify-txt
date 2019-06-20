package main

import (
	"log"
	"net/http"

	"github.com/thetetrabyte/go-spotify-txt/internal/app/rest"
	"github.com/thetetrabyte/go-spotify-txt/internal/app/service"
	"github.com/thetetrabyte/go-spotify-txt/internal/app/spotify"
	"github.com/thetetrabyte/go-spotify-txt/internal/pkg/config"
)

func main() {
	// Initialize Config
	config.Initialize()

	// Initialize JSON config
	config.InitializeJSON()

	// Initialize Routes
	router := rest.Initialize()

	// Check the user details and initialize their data
	Users := config.GetJSON()

	// Check if the user details already exist
	for _, a := range Users {
		userInf, userErr := spotify.PullUserInfo(a)
		if userErr != nil {
			log.Println("Erroring accessing data for "+a.Email, userErr)
			continue
		}

		log.Println("Tracking data for " + userInf.DisplayName + " (" + userInf.ID + ")")

		// Initialize the polling for each accounts player
		go func() { service.Initialize(a, userInf) }()
	}

	// Log about being ready
	log.Println("Running on " + config.Get().Port)
	log.Fatalln(http.ListenAndServe(":"+config.Get().Port, router))
}
