package spotify

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/thetetrabyte/go-spotify-txt/internal/pkg/config"
)

// GetAuthorizationURL will send the user back the authorization url to authenticate
func GetAuthorizationURL() string {
	return "https://accounts.spotify.com/authorize?response_type=code&client_id=" + config.Get().Spotify.Client + "&scope=" + url.QueryEscape(config.Get().Spotify.Scopes) + "&redirect_uri=" + url.QueryEscape(config.Get().Spotify.Callback)
}

// RequestUserTokens will make a request to spotify accounts api and return the tokens needed to authenticate as the user
func RequestUserTokens(Code string) (config.Account, TokenError, error) {
	var tokens Tokens
	var account config.Account
	var tokenError TokenError

	// Create base64 string with client id and client secret
	Token := base64.StdEncoding.EncodeToString([]byte(config.Get().Spotify.Client + `:` + config.Get().Spotify.Secret))

	// Define the body and create the request
	body := strings.NewReader(`grant_type=authorization_code&code=` + Code + `&redirect_uri=` + url.QueryEscape(config.Get().Spotify.Callback))
	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", body)
	if err != nil {
		return account, tokenError, err
	}

	// Set the needed headers
	req.Header.Set("Authorization", "Basic "+Token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create the client and send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return account, tokenError, err
	}

	// Close the client
	defer resp.Body.Close()

	// Create decoder
	decoder := json.NewDecoder(resp.Body)

	switch resp.StatusCode {
	case http.StatusBadRequest:
		if err := decoder.Decode(&tokenError); err != nil {
			return account, tokenError, err
		}

		break
	case http.StatusOK:
		if err := decoder.Decode(&tokens); err != nil {
			return account, tokenError, err
		}

		break
	default:
		err := errors.New("failed_to_get_tokens")
		return account, tokenError, err
	}

	userInf, userErr := PullUserInfo(config.Account{Token: tokens.AccessToken, Refresh: tokens.RefreshToken, Email: "temp@temp.com"})
	if userErr != nil {
		log.Println(userErr)

		return account, tokenError, userErr
	}

	// Define the slice for the user tokens
	User := config.GetJSON()

	// Check if the user details already exist
	for _, a := range User {
		if a.Email == userInf.Email {
			return account, tokenError, errors.New("user_exists")
		}
	}

	// Append the new record to the current json
	User = append(User, config.Account{Email: userInf.Email, Account: userInf.ID, Token: tokens.AccessToken, Refresh: tokens.RefreshToken})

	result, err := json.MarshalIndent(User, "", "\t")
	if err != nil {
		log.Println(err)
	}

	_ = ioutil.WriteFile("accounts.json", result, 0644)

	// Reinitialize JSON config
	config.InitializeJSON()

	// Return the tokens
	return account, tokenError, nil
}

// RefreshAccessToken will make a request to spotify accounts api with the refresh token and return the tokens needed to authenticate as the user
func RefreshAccessToken(Refresh string, Email string) (config.Account, TokenError, error) {
	var tokens Tokens
	var tokenError TokenError
	var account config.Account

	// Create base64 string with client id and client secret
	Token := base64.StdEncoding.EncodeToString([]byte(config.Get().Spotify.Client + `:` + config.Get().Spotify.Secret))

	// Define the body and create the request
	body := strings.NewReader(`grant_type=refresh_token&refresh_token=` + Refresh)
	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", body)
	if err != nil {
		return account, tokenError, err
	}

	// Set the needed headers
	req.Header.Set("Authorization", "Basic "+Token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create the client and send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return account, tokenError, err
	}

	// Close the client
	defer resp.Body.Close()

	// Create decoder
	decoder := json.NewDecoder(resp.Body)

	switch resp.StatusCode {
	case http.StatusBadRequest:
		if err := decoder.Decode(&tokenError); err != nil {
			return account, tokenError, err
		}

		break
	case http.StatusOK:
		if err := decoder.Decode(&tokens); err != nil {
			return account, tokenError, err
		}

		break
	default:
		err := errors.New("failed_to_get_tokens")
		return account, tokenError, err
	}

	// Define the slice for the user tokens
	User := config.GetJSON()

	// Get user details and update
	for index, a := range User {
		if a.Email == Email {
			config.GetJSON()[index].Token = tokens.AccessToken
			result, err := json.MarshalIndent(User, "", "\t")
			if err != nil {
				log.Println(err)
			}

			_ = ioutil.WriteFile("accounts.json", result, 0644)

			// Reinitialize JSON config
			config.InitializeJSON()

			account = config.GetJSON()[index]
		}
	}

	// Return the tokens
	return account, tokenError, nil
}

// PullUserInfo will pull the users info from the spotify api
func PullUserInfo(Account config.Account) (UserInfo, error) {
	var User UserInfo

	// Create the request
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
	if err != nil {
		return User, err
	}

	// Set the required headers
	req.Header.Set("Authorization", "Bearer "+Account.Token)

	// Create the client and send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return User, err
	}

	// Close the client
	defer resp.Body.Close()

	// Create decoder
	decoder := json.NewDecoder(resp.Body)

	// Switch case for status code forbidden
	switch resp.StatusCode {
	case http.StatusOK:
		if err := decoder.Decode(&User); err != nil {
			return User, err
		}

		break
	case http.StatusUnauthorized:
		var authError AuthError
		if err := decoder.Decode(&authError); err != nil {
			log.Println("Error unmarshaling json for auth error")
		}

		if authError.Error.Message == "The access token expired" {
			log.Println("Refreshing access token for " + Account.Email)

			account, tokenErr, err := RefreshAccessToken(Account.Refresh, Account.Email)
			if err != nil {
				return User, err
			}

			if (TokenError{}) != tokenErr {
				return User, errors.New("token_fetch_error")
			}

			user, err := PullUserInfo(account)
			if err != nil {
				return User, err
			}

			User = user
		} else if authError.Error.Message == "Invalid access token" {
			log.Println("Invalid access token for " + Account.Email)
			return User, errors.New("Invalid access token for " + Account.Email)
		}

		break
	case http.StatusBadRequest:
		return User, errors.New("bad_request")
	}

	// Return the user
	return User, nil
}

// PullPlayerInfo will pull the users current player from the spotify api
func PullPlayerInfo(Account config.Account) (PlayerResponse, error) {
	var Player PlayerResponse

	// Create the request
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/player", nil)
	if err != nil {
		return Player, err
	}

	// Set the required headers
	req.Header.Set("Authorization", "Bearer "+Account.Token)

	// Create the client and send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return Player, err
	}

	// Close the client
	defer resp.Body.Close()

	// Create decoder
	decoder := json.NewDecoder(resp.Body)

	// Switch case for status code forbidden
	switch resp.StatusCode {
	case http.StatusOK:
		if err := decoder.Decode(&Player); err != nil {
			return Player, err
		}

		break
	case http.StatusUnauthorized:
		var authError AuthError
		if err := decoder.Decode(&authError); err != nil {
			log.Println("Error unmarshaling json for auth error")
		}

		if authError.Error.Message == "The access token expired" {
			log.Println("Refreshing access token for " + Account.Email)

			account, tokenErr, err := RefreshAccessToken(Account.Refresh, Account.Email)
			if err != nil {
				return Player, err
			}

			if (TokenError{}) != tokenErr {
				return Player, errors.New("token_fetch_error")
			}

			player, err := PullPlayerInfo(account)
			if err != nil {
				return Player, err
			}

			Player = player
		} else if authError.Error.Message == "Invalid access token" {
			log.Println("Invalid access token for " + Account.Email)
			return Player, errors.New("Invalid access token for " + Account.Email)
		}

		break
	case http.StatusBadRequest:
		return Player, errors.New("bad_request")
	}

	// Return the player
	return Player, nil
}
