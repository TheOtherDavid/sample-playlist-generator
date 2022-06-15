package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
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
	artistNames = append(artistNames, "VNV Nation")
	artistNames = append(artistNames, "Nova Twins")
	artistNames = append(artistNames, "Welle: Erdball")

	var artistIds []string

	for _, artistName := range artistNames {
		artistId := SearchForArtistId(artistName, accessToken)
		artistIds = append(artistIds, artistId)
	}

	var playlistName = "Sample Playlist"
	fmt.Println(playlistName)

	var playlistTrackIds []string

	for _, artistId := range artistIds {
		selectedTrackIds := GetTopTrackIds(artistId, accessToken)
		fmt.Println(selectedTrackIds)
		playlistTrackIds = append(playlistTrackIds, selectedTrackIds...)
	}

	//Create empty playlist
	playlistId := CreateEmptySpotifyPlaylist(playlistName, accessToken)
	fmt.Println(playlistId)
	//Add selected songs to playlist

	snapshotId := AddTracksToSpotifyPlaylist(playlistTrackIds, playlistId, accessToken)
	fmt.Println(snapshotId)
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
	}
	var responseBody SpotifyRefreshTokenResponse
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		fmt.Println("Oh no, error.")
	}

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
	qEncoded := &url.URL{Path: q}
	qEncodedString := qEncoded.String()
	fmt.Println(qEncodedString)

	url := "https://api.spotify.com/v1/search?q=" + qEncodedString + "&type=artist"

	client := http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Oh no, error.")
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
		fmt.Println("Oh no, error.")
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

	return selectedTrackIds[:3]
}

type SpotifyCreatePlaylistRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Public      bool   `json:"public"`
}

type SpotifyCreatePlaylistResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func CreateEmptySpotifyPlaylist(playlistName string, accessToken string) string {
	userId := os.Getenv("USER_ID")
	createPlaylistRequest := SpotifyCreatePlaylistRequest{
		Name:        playlistName,
		Description: "Generated automatically",
		//TODO: Change this to true later. Maybe have it env-specific so Dev playlists aren't public?
		Public: false,
	}
	body, _ := json.Marshal(createPlaylistRequest)
	url := "https://api.spotify.com/v1/users/" + userId + "/playlists"

	client := http.Client{}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("Oh no, error.")
	}

	req.Header = http.Header{
		"Authorization": {"Bearer " + accessToken},
		"Content-Type":  {"application/json"},
	}

	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Oh no, error.")
	}

	var responseBody SpotifyCreatePlaylistResponse
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		fmt.Println("Oh no, error.")
	}

	playlistId := responseBody.Id
	return playlistId
}

type SpotifyAddTrackToPlaylistResponse struct {
	SnapshotId string `json:"snapshot_id"`
}

func AddTracksToSpotifyPlaylist(trackIds []string, playlistId string, accessToken string) string {
	//convert trackId to full URI
	spotifyUris := []string{}
	for _, trackId := range trackIds {
		uri := "spotify:track:" + trackId
		spotifyUris = append(spotifyUris, uri)
	}
	uriString := strings.Join(spotifyUris[:], ",")
	fmt.Println(strings.Join(spotifyUris[:], ","))

	url := "https://api.spotify.com/v1/playlists/" + playlistId + "/tracks?uris=" + uriString

	client := http.Client{}

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		fmt.Println("Oh no, error.")
	}

	req.Header = http.Header{
		"Authorization": {"Bearer " + accessToken},
	}

	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Oh no, error.")
	}

	var responseBody SpotifyAddTrackToPlaylistResponse
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		fmt.Println("Oh no, error.")
	}

	snapshotId := responseBody.SnapshotId
	return snapshotId
}
