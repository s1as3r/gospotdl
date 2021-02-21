// A Golang implementation of spotDl (python)
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/zmb3/spotify"
)

const (
	spotifyClientId     = "b9640826ab8949e4b34a4eeee5b26ce6"
	spotifyClientSecret = "c25b0399324d49a08d8210138cbf9e27"
)

func main() {
	spotifyClient, err := GetSpotifyClient(spotifyClientId, spotifyClientSecret)
	if err != nil {
		log.Fatalf("Error getting Spotify Client: %s\n", err)
	}
	for _, url := range os.Args[1:] {
		t, id := parseUrl(url)
		switch t {
		case "track":
			track, err := SongObjFromID(spotifyClient, spotify.ID(id))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			}
			if err := Download(track); err != nil {
				log.Fatalf("Error Downloading Track(%s): %s\n", track.Name, err)
			}
			break
		case "album":
			tracks, err := GetAlbumTracks(spotifyClient, spotify.ID(id))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			}
			for _, track := range tracks {
				if err := Download(track); err != nil {
					fmt.Fprintf(os.Stderr, "Error Downloading Track(%s): %s\n", track.Name, err)
					continue
				}
			}
			break
		case "playlist":
			tracks, err := GetPlaylistTracks(spotifyClient, spotify.ID(id))
			if err != nil {
				log.Fatalf("Error Getting Tracks(%s): %s", id, err)
			}
			for _, track := range tracks {
				if err := Download(track); err != nil {
					fmt.Fprintf(os.Stderr, "Error Downloading Track(%s): %s\n", track.Name, err)
					continue
				}
			}
			break
		case "query":
			track, err := SongObjFromQuery(spotifyClient, id)
			if err != nil {
				log.Fatalf("Error Getting Song(%s): %s", id, err)
			}
			if err := Download(track); err != nil {
				log.Fatalf("Error Downloading Track(%s): %s\n", track.Name, err)
			}
			break
		}
	}
}
