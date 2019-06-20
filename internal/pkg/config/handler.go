package config

import (
	"encoding/json"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var (
	// config is the env variables
	config Config

	// JSONFile config for accounts
	JSONFile JSONConfig
)

// Initialize will load the config variables from the .env file
func Initialize() {
	godotenv.Load()
	config = Config{}

	if err := envconfig.Process("app", &config); err != nil {
		log.Fatalln(err)
	}
}

// Get config
func Get() Config {
	return config
}

// InitializeJSON will load the json data for spotify credentials
func InitializeJSON() {
	// Check if accounts file exists and create if not
	if _, err := os.Stat("accounts.json"); os.IsNotExist(err) {
		log.Println("Accounts file does not exist, creating")
		accountsFile, err := os.Create("accounts.json")
		if err != nil {
			log.Println("Failed to create accounts file")
			return
		}

		defer accountsFile.Close()
	}

	File, err := os.Open("accounts.json")
	defer File.Close()
	if err != nil {
		log.Println(err.Error())
	}

	jsonParser := json.NewDecoder(File)
	jsonParser.Decode(&JSONFile)
}

// GetJSON will return the json config struct
func GetJSON() JSONConfig {
	return JSONFile
}
