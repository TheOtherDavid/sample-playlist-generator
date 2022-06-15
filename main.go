package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func main() {
	GeneratePlaylist()
}

func GeneratePlaylist() {
	//We need to take in some artists, but let's just hardcode one for now.
	var artists []string
	artists = append(artists, "Rammstein")
	fmt.Println(artists[0])
	//How do we get from the string Rammstein to the ID?
	var artistIds []string
	//Just hardcode this for now.
	artistIds = append(artistIds, "6wWVKhxIU2cEi0K81v7HvP")

	var playlistName = "Sample Playlist"
	fmt.Println(playlistName)

	fmt.Println("Refreshing access token.")
	accessToken := RefreshSpotifyAuth()
	fmt.Println("Access token refreshed.")
	fmt.Println(accessToken)

	for _, artistId := range artistIds {
		//Call top tracks. Put tracks into a map by album, to get the most popular track per album from the top tracks.
		selectedTrackIds := GetTopTrackIds(artistId, accessToken)
		fmt.Println(selectedTrackIds)
		//Create empty playlist
	}

	//Add selected songs to playlist
}

type SpotifyRefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

func RefreshSpotifyAuth() string {
	clientId := os.Getenv("SPOTIFY_ID")
	clientSecret := os.Getenv("SPOTIFY_SECRET")
	refreshToken := os.Getenv("REFRESH_TOKEN")
	grantType := "refresh_token"
	url := "https://accounts.spotify.com/api/token?client_id=" +
		clientId + "&client_secret=" + clientSecret + "&refresh_token=" +
		refreshToken + "&grant_type=" + grantType

	response, err := http.Post(url, "application/x-www-form-urlencoded", nil)
	if err != nil {
		fmt.Println("Oh no, error.")
	} else {
		fmt.Println(response)
	}
	var responseBody SpotifyRefreshTokenResponse
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		fmt.Println("Oh no, error.")
	}

	fmt.Println(responseBody)

	accessToken := responseBody.AccessToken
	return accessToken
}

func GetTopTrackIds(artistId string, accessToken string) []string {
	var trackIds []string
	trackIds = append(trackIds, "0C9p8YMtbdOkcXPPlEmZvY")
	return trackIds
}
