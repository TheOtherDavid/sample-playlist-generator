package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/TheOtherDavid/sample-playlist-generator"
	"github.com/gorilla/mux"
)

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/generate", generatePlaylist()).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main() {
	fmt.Println("Listening on port 8080")
	handleRequests()
}

func generatePlaylist() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var body generatePlaylistRequest

		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&body); err != nil {
			//respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			fmt.Println(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		playlistName := body.PlaylistName
		artistNames := body.ArtistNames
		generate.GeneratePlaylist(artistNames, playlistName)
		response := generatePlaylistResponse{
			Success: true,
		}
		respondWithJSON(w, 200, response)
	}
}

type generatePlaylistRequest struct {
	PlaylistName string   `json:"playlistName"`
	ArtistNames  []string `json:"artistNames"`
}

type generatePlaylistResponse struct {
	Success bool `json:"success"`
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
