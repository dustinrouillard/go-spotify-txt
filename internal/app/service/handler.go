package service

import (
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/thetetrabyte/go-spotify-txt/internal/app/spotify"
	"github.com/thetetrabyte/go-spotify-txt/internal/pkg/config"
)

// Initialize will pull all the accounts create their text files and queue them for updates
func Initialize(Account config.Account, User spotify.UserInfo) {
	// Check if the texts folder exists
	if _, err := os.Stat("texts"); os.IsNotExist(err) {
		log.Println("Texts directory does not exist, creating")
		os.Mkdir("texts", 0777)
	}

	// Log
	log.Println("Starting poller for " + User.DisplayName + " (" + User.ID + ")")

	// Check if song file exists for user and create if not
	if _, err := os.Stat("texts/" + User.ID + "-song.txt"); os.IsNotExist(err) {
		log.Println("Account text file does not exist, creating")
		songFile, err := os.Create("texts/" + User.ID + "-song.txt")
		if err != nil {
			log.Println("Failed to create text file for " + User.DisplayName)
			return
		}

		defer songFile.Close()
	}

	// Check if artist file exists for user and create if not
	if _, err := os.Stat("texts/" + User.ID + "-artist.txt"); os.IsNotExist(err) {
		log.Println("Account text file does not exist, creating")
		artistFile, err := os.Create("texts/" + User.ID + "-artist.txt")
		if err != nil {
			log.Println("Failed to create text file for " + User.DisplayName)
			return
		}

		defer artistFile.Close()
	}

	// Write song text
	songWriteErr := ioutil.WriteFile("texts/"+User.ID+"-song.txt", []byte("Loading..."), 0777)
	if songWriteErr != nil {
		log.Println("Failed to read data for " + User.DisplayName)
		return
	}

	// Write artist text
	artistWriteErr := ioutil.WriteFile("texts/"+User.ID+"-artist.txt", []byte("Loading..."), 0777)
	if artistWriteErr != nil {
		log.Println("Failed to read data for " + User.DisplayName)
		return
	}

	// Create loop
	tick := time.Tick(5 * time.Second)
	for range tick {
		// Define the slice for the user tokens
		Usr := config.GetJSON()

		// Get user details and update
		for index, a := range Usr {
			if a.Account == Account.Account {
				playerInf, playerErr := spotify.PullPlayerInfo(config.GetJSON()[index])
				if playerErr != nil {
					log.Println("Failed to pull player data for " + User.DisplayName)
					continue
				}

				if playerInf.IsPlaying {
					// Read the song file
					songFile, songErr := ioutil.ReadFile("texts/" + User.ID + "-song.txt")
					if songErr != nil {
						log.Println("Failed to read data for " + User.DisplayName)
						return
					}

					// Read the artist file
					artistFile, artistErr := ioutil.ReadFile("texts/" + User.ID + "-artist.txt")
					if artistErr != nil {
						log.Println("Failed to read data for " + User.DisplayName)
						return
					}

					// Define new texts
					songText := playerInf.Item.Name
					artistText := playerInf.Item.Artists[0].Name

					// Check data is not the same as what we're setting it to
					if string(songFile) != songText {
						// Update text files with current song details
						log.Println("Updating song " + User.ID + " to " + songText)

						// Write song text
						songWriteErr := ioutil.WriteFile("texts/"+User.ID+"-song.txt", []byte(songText), 0777)
						if songWriteErr != nil {
							log.Println("Failed to read data for " + User.DisplayName)
							return
						}
					}

					// Check data is not the same as what we're setting it to
					if string(artistFile) != artistText {
						// Update text files with current song details
						log.Println("Updating artist " + User.ID + " to " + artistText)

						// Write artist text
						artistWriteErr := ioutil.WriteFile("texts/"+User.ID+"-artist.txt", []byte(artistText), 0777)
						if artistWriteErr != nil {
							log.Println("Failed to read data for " + User.DisplayName)
							return
						}
					}
				} else {
					// Read the song file
					songFile, songErr := ioutil.ReadFile("texts/" + User.ID + "-song.txt")
					if songErr != nil {
						log.Println("Failed to read data for " + User.DisplayName)
						return
					}

					// Read the artist file
					artistFile, artistErr := ioutil.ReadFile("texts/" + User.ID + "-artist.txt")
					if artistErr != nil {
						log.Println("Failed to read data for " + User.DisplayName)
						return
					}

					// Define new texts
					songText := "Nothing playing"
					artistText := "N/A"

					// Check data is not the same as what we're setting it to
					if string(songFile) != songText {
						// Update text files with current song details
						log.Println("Updating song " + User.ID + " to " + songText)

						// Write song text
						songWriteErr := ioutil.WriteFile("texts/"+User.ID+"-song.txt", []byte(songText), 0777)
						if songWriteErr != nil {
							log.Println("Failed to read data for " + User.DisplayName)
							return
						}
					}

					// Check data is not the same as what we're setting it to
					if string(artistFile) != artistText {
						// Update text files with current song details
						log.Println("Updating artist " + User.ID + " to " + artistText)

						// Write artist text
						artistWriteErr := ioutil.WriteFile("texts/"+User.ID+"-artist.txt", []byte(artistText), 0777)
						if artistWriteErr != nil {
							log.Println("Failed to read data for " + User.DisplayName)
							return
						}
					}
				}
			}
		}

	}
}
