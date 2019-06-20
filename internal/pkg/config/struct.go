package config

// Config defines the structure of the config variables
type Config struct {
	Port    string `envconfig:"PORT" default:"9090"`
	Spotify Spotify
}

// Spotify defines the structure of the spotify credentials
type Spotify struct {
	Client   string `envconfig:"SPOTIFY_CLIENT_ID"`
	Secret   string `envconfig:"SPOTIFY_CLIENT_SECRET"`
	Callback string `envconfig:"SPOTIFY_CALLBACK"`
	Scopes   string `envconfig:"SPOTIFY_SCOPES"`
}

// JSONConfig defines the structure of the spotify accounts stored locally
type JSONConfig []Account

// Account defines an account in the json array
type Account struct {
	Email   string `json:"email"`
	Account string `json:"account"`
	Token   string `json:"token"`
	Refresh string `json:"refresh"`
}
