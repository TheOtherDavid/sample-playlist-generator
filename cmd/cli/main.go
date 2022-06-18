package main

import "os"
import "github.com/TheOtherDavid/sample-playlist-generator"

func main() {
	playlistName := os.Args[1]
	artistNames := os.Args[2:]
	generate.GeneratePlaylist(artistNames, playlistName)
}
