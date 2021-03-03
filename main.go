// A Golang implementation of spotDl (python)
package main

import (
	"flag"
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
	var dir string
	flag.StringVar(&dir, "output", ".", "Output Directory")
	flag.StringVar(&dir, "o", ".", "Output Directory")
	flag.Parse()

	if err := os.Chdir(dir); err != nil {
		log.Fatalf("Error changing directory: %s", err)
	}

	for _, url := range flag.Args() {
		t, id := parseArg(url)
		switch t {
		case "track":
			track, err := SongObjFromID(spotifyClient, spotify.ID(id))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			}
			if err := Download(track); err != nil {
				log.Fatalf("Error Downloading Track(%s): %s\n", track.Name, err)
			}

		case "album":
			tracks, err := GetAlbumTracks(spotifyClient, spotify.ID(id))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			}
			DownloadMulti(tracks)

		case "playlist":
			tracks, err := GetPlaylistTracks(spotifyClient, spotify.ID(id))
			if err != nil {
				log.Fatalf("Error Getting Tracks(%s): %s", id, err)
			}
			DownloadMulti(tracks)

		case "query":
			track, err := SongObjFromQuery(spotifyClient, id)
			if err != nil {
				log.Fatalf("Error Getting Song(%s): %s", id, err)
			}
			if err := Download(track); err != nil {
				log.Fatalf("Error Downloading Track(%s): %s\n", track.Name, err)
			}
		}
	}
}
