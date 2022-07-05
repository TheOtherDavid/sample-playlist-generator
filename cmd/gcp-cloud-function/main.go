package gcp

import (
	"encoding/json"
	"fmt"
	"github.com/TheOtherDavid/sample-playlist-generator"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("HelloHTTP", helloHTTP)
}

// helloHTTP is an HTTP Cloud Function with a request parameter.
func helloHTTP(w http.ResponseWriter, r *http.Request) {
	var event struct {
		PlaylistName string   `json:"playlistName"`
		Artists      []string `json:"artistNames"`
	}
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		fmt.Fprint(w, "Hello, World!")
		return
	}
	artists := event.Artists
	playlistName := event.PlaylistName
	generate.GeneratePlaylist(artists, playlistName)
}
