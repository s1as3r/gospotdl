package main

import (
	"fmt"
	"log"
	"os"

	"github.com/zmb3/spotify"
)

func handleArg(arg string, client spotify.Client) {
	t, id := parseArg(arg)
	switch t {
	case "track":
		handleTrack(client, id)
	case "album":
		handleAlbum(client, id)
	case "playlist":
		handlePlaylist(client, id)
	case "query":
		handleQuery(client, id)
	}
}

func handleTrack(client spotify.Client, id string) {
	track, err := SongObjFromID(client, spotify.ID(id))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	}
	if err := Download(track); err != nil {
		log.Fatalf("Error Downloading Track(%s): %s\n", track.Name, err)
	}
}

func handleAlbum(client spotify.Client, id string) {
	tracks, err := GetAlbumTracks(client, spotify.ID(id))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
	}
	DownloadMulti(tracks)
}

func handlePlaylist(client spotify.Client, id string) {
	tracks, err := GetPlaylistTracks(client, spotify.ID(id))
	if err != nil {
		log.Fatalf("Error Getting Tracks(%s): %s", id, err)
	}
	DownloadMulti(tracks)
}

func handleQuery(client spotify.Client, id string) {
	track, err := SongObjFromQuery(client, id)
	if err != nil {
		log.Fatalf("Error Getting Song(%s): %s", id, err)
	}
	if err := Download(track); err != nil {
		log.Fatalf("Error Downloading Track(%s): %s\n", track.Name, err)
	}
}
