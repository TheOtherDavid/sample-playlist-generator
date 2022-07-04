package generate

import (
	"fmt"
)

func GeneratePlaylist(artistNames []string, playlistName string) {

	fmt.Println("Refreshing access token.")
	accessToken, err := RefreshSpotifyAuth()
	if err != nil {
		panic(0)
	}
	fmt.Println("Access token refreshed.")
	fmt.Println(accessToken)

	var artists []SpotifyArtistItem

	for _, artistName := range artistNames {
		artist := SearchForArtist(artistName, accessToken)
		artists = append(artists, artist)
	}

	var playlistTrackIds []string

	for _, artist := range artists {
		selectedTrackIds := GetTopTrackIds(artist.Id, accessToken)
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
