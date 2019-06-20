package spotify

// PlayerResponse is the incoming response when asking for player data
type PlayerResponse struct {
	Device               PlayerDevice  `json:"device"`
	ShuffleState         bool          `json:"shuffle_state"`
	RepeatState          string        `json:"repeat_state"`
	Timestamp            int64         `json:"timestamp"`
	Context              PlayerContext `json:"context"`
	ProgressMS           int64         `json:"progress_ms"`
	Item                 PlayerItem    `json:"item"`
	CurrentlyPlayingType string        `json:"currently_playing_type"`
	IsPlaying            bool          `json:"is_playing"`
}

// PlayerDevice is the device information that is playing a song
type PlayerDevice struct {
	ID               string `json:"id"`
	IsActive         bool   `json:"is_active"`
	IsPrivateSession bool   `json:"is_private_session"`
	IsRestricted     bool   `json:"is_restricted"`
	Name             string `json:"name"`
	Type             string `json:"type"`
	VolumePercent    int64  `json:"volume_percent"`
}

// PlayerContext is the context object of the current song
type PlayerContext struct {
	ExternalUrls ExternalUrls `json:"external_urls"`
	Href         string       `json:"href"`
	Type         string       `json:"type"`
	URI          string       `json:"uri"`
}

// ExternalUrls is the struct for all external urls
type ExternalUrls struct {
	Spotify string `json:"spotify"`
}

// PlayerItem is the struct for the current item that is playing
type PlayerItem struct {
	Album            PlayerAlbum    `json:"album"`
	Artists          []PlayerArtist `json:"artists"`
	AvailableMarkets []string       `json:"available_markets"`
	DiscNumber       int64          `json:"disc_number"`
	DurationMS       int64          `json:"duration_ms"`
	Explicit         bool           `json:"explicit"`
	ExternalUrls     ExternalUrls   `json:"external_urls"`
	Href             string         `json:"href"`
	ID               string         `json:"id"`
	IsLocal          bool           `json:"is_local"`
	Name             string         `json:"name"`
	Popularity       int64          `json:"popularity"`
	PreviewURL       string         `json:"preview_url"`
	TrackNumber      int64          `json:"track_number"`
	Type             string         `json:"type"`
	URI              string         `json:"uri"`
}

// PlayerAlbum is the struct for the current playing songs album
type PlayerAlbum struct {
	AlbumType            string         `json:"album_type"`
	Artists              []PlayerArtist `json:"artists"`
	AvailableMarkets     []string       `json:"available_markets"`
	ExternalUrls         ExternalUrls   `json:"external_urls"`
	Href                 string         `json:"href"`
	ID                   string         `json:"id"`
	Images               []PlayerImage  `json:"images"`
	Name                 string         `json:"name"`
	ReleaseDate          string         `json:"release_date"`
	ReleaseDatePrecision string         `json:"release_date_precision"`
	TotalTracks          int64          `json:"total_tracks"`
	Type                 string         `json:"type"`
	URI                  string         `json:"uri"`
}

// PlayerArtist is the struct for the song or album artist
type PlayerArtist struct {
	ExternalUrls ExternalUrls `json:"external_urls"`
	Href         string       `json:"href"`
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Type         string       `json:"type"`
	URI          string       `json:"uri"`
}

// PlayerImage is the struct for all song/slbum/artist images
type PlayerImage struct {
	Height int64  `json:"height"`
	URL    string `json:"url"`
	Width  int64  `json:"width"`
}

// Tokens is the struct for an authorization request
type Tokens struct {
	AccessToken  string `json:"access_token,omitempty"`
	ExpiresIn    int64  `json:"expires_in,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// UserInfo is the struct for a user profile
type UserInfo struct {
	DisplayName  string        `json:"display_name"`
	Email        string        `json:"email"`
	ExternalUrls ExternalUrls  `json:"external_urls"`
	Followers    Followers     `json:"followers"`
	Href         string        `json:"href"`
	ID           string        `json:"id"`
	Images       []interface{} `json:"images"`
	Type         string        `json:"type"`
	URI          string        `json:"uri"`
}

// Followers is the struct for a user followers
type Followers struct {
	Href  interface{} `json:"href"`
	Total int64       `json:"total"`
}

// TokenError defines the struct for an error when requesting tokens
type TokenError struct {
	Error            string `json:"error,omitempty"`
	ErrorDescription string `json:"error_description,omitempty"`
}

// AuthError is the error returned by spotify when requesting data as a user
type AuthError struct {
	Error Error `json:"error"`
}

// Error defines the error struct for AuthError
type Error struct {
	Status  int64  `json:"status"`
	Message string `json:"message"`
}
