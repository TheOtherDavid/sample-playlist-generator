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

	fmt.Println("Refreshing access token.")
	accessToken := RefreshSpotifyAuth()
	fmt.Println("Access token refreshed.")
	fmt.Println(accessToken)

	//We need to take in some artists, but let's just hardcode one for now.
	var artistNames []string
	artistNames = append(artistNames, "Rammstein")
	fmt.Println(artistNames[0])
	//How do we get from the string Rammstein to the ID?
	var artistIds []string
	//Just hardcode this for now.
	for _, artistName := range artistNames {
		artistId := SearchForArtistId(artistName, accessToken)
		artistIds = append(artistIds, artistId)
	}

	var playlistName = "Sample Playlist"
	fmt.Println(playlistName)

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

type SpotifySearchResponse struct {
	Artists SpotifyArtist `json:"artists"`
}

type SpotifyArtist struct {
	Items []SpotifyArtistItem `json:"items"`
}

type SpotifyArtistItem struct {
	Genres     []string `json:"genres"`
	Id         string   `json:"id"`
	Name       string   `json:"name"`
	Popularity int      `json:"popularity"`
	Type       string   `json:"type"`
}

func SearchForArtistId(artistName string, accessToken string) string {

	q := artistName
	url := "https://api.spotify.com/v1/search?q=" + q + "&type=artist"

	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		//Handle Error
	}

	req.Header = http.Header{
		"Authorization": {"Bearer " + accessToken},
	}

	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Oh no, error.")
	}

	var responseBody SpotifySearchResponse
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		fmt.Println("Oh no, error.")
	}
	//Assume the first result is correct and take that ID.
	artistId := responseBody.Artists.Items[0].Id
	return artistId
}

type SpotifyTopTracksResponse struct {
	Tracks []SpotifyTrack `json:"tracks"`
}

type SpotifyTrack struct {
	Album       SpotifyAlbum        `json:"album"`
	Artists     []SpotifyArtistItem `json:"artists"`
	Id          string              `json:"id"`
	Name        string              `json:"name"`
	Popularity  int                 `json:"popularity"`
	TrackNumber int                 `json:"track_number"`
}

type SpotifyAlbum struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	ReleaseDate string `json:"release_date"`
}

func GetTopTrackIds(artistId string, accessToken string) []string {
	//Call Top Tracks for artist ID
	url := "https://api.spotify.com/v1/artists/" + artistId + "/top-tracks?market=US"

	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		//Handle Error
	}

	req.Header = http.Header{
		"Authorization": {"Bearer " + accessToken},
	}

	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Oh no, error.")
	}
	var responseBody SpotifyTopTracksResponse
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		fmt.Println("Oh no, error.")
	}
	tracks := responseBody.Tracks

	//Create a map with album ID as the key, and the Track object as the value
	selectedTracks := make(map[string]SpotifyTrack)
	for _, track := range tracks {
		if val, ok := selectedTracks[track.Album.Id]; !ok {
			//If the album isn't already in the map, add this one!
			selectedTracks[track.Album.Id] = track
		} else {
			//If the album IS in the map, check the probability
			oldTrack := val
			oldPopularity := oldTrack.Popularity
			newPopularity := track.Popularity
			if newPopularity > oldPopularity {
				//If the new track is more popular than the old track, replace it
				selectedTracks[track.Album.Id] = track
			}
		}
	}

	var selectedTrackIds []string

	//Get the values out of the map
	for _, track := range selectedTracks {
		selectedTrackIds = append(selectedTrackIds, track.Id)
	}

	return selectedTrackIds
}
