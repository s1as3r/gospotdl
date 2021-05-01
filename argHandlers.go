package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/s1as3r/gospotdl/download"
	"github.com/s1as3r/gospotdl/search"
	"github.com/zmb3/spotify"
)

func parseArg(arg string) (string, string) {
	if strings.Contains(arg, "spotify.com") {
		url := strings.ReplaceAll(arg, "\\", "/")
		url = strings.TrimSuffix(url, "/")
		splitUrl := strings.Split(url, "/")
		id := splitUrl[len(splitUrl)-1]
		id = strings.Split(id, "?")[0]

		if strings.Contains(url, "track") {
			return "track", id
		} else if strings.Contains(url, "album") {
			return "album", id
		} else if strings.Contains(url, "playlist") {
			return "playlist", id
		}
	}
	return "query", arg
}

func handleArg(client *spotify.Client, arg string, max int) {
	argType, id := parseArg(arg)
	switch argType {
	case "track":
		handleTrack(client, id)
	case "album":
		handleAlbum(client, id, max)
	case "playlist":
		handlePlaylist(client, id, max)
	case "query":
		handleQuery(client, arg)
	}
}

func handleTrack(client *spotify.Client, id string) {
	track := &search.Song{}
	err := track.FromId(client, spotify.ID(id))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Downloading: %s", track.Name)
	if err := download.Download(track); err != nil {
		log.Fatalf("Error downloading track %s: %s", track.Name, err)
	}
}

func handleAlbum(client *spotify.Client, id string, max int) {
	tracks, err := search.GetPlaylistTracks(client, spotify.ID(id))
	if err != nil {
		log.Fatal(err)
	}
	download.AsyncDownload(tracks, max)
}

func handlePlaylist(client *spotify.Client, id string, max int) {
	tracks, err := search.GetPlaylistTracks(client, spotify.ID(id))
	if err != nil {
		log.Fatal(err)
	}
	download.AsyncDownload(tracks, max)
}

func handleQuery(client *spotify.Client, query string) {
	track := &search.Song{}
	if err := track.FromQuery(client, query); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Downloading: %s", track.Name)
	if err := download.Download(track); err != nil {
		log.Fatal(err)
	}
}
