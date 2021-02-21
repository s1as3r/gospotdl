package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/zmb3/spotify"
)

// GetAlbumTracks gets all the tracks in in album.
func GetAlbumTracks(spotifyClient spotify.Client, spotifyID spotify.ID) ([]SongObj, error) {
	page, err := spotifyClient.GetAlbumTracks(spotifyID)
	if err != nil {
		return []SongObj{}, err
	}
	var tracks []SongObj

	album, err := spotifyClient.GetAlbum(spotifyID)
	if err != nil {
		return []SongObj{}, err
	}
	fmt.Printf("Fetching Songs from: %s\n", album.Name)
	for {
		for _, track := range page.Tracks {
			songObj, err := SongObjFromSimpleTrack(spotifyClient, &track)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Skipping %s as no match was found on Yoututbe", track.Name)
				continue
			}
			fmt.Printf("Found: %s\n", track.Name)
			tracks = append(tracks, songObj)
		}

		if strings.Contains(page.Next, "spotify") {
			opt := spotify.Options{
				Offset: &page.Offset,
			}
			page, err = spotifyClient.GetAlbumTracksOpt(spotifyID, &opt)
			if err != nil {
				return []SongObj{}, err
			}
		} else {
			break
		}
	}
	return tracks, nil
}

// GetPlaylistTracks gets all the songs of a playlist.
func GetPlaylistTracks(spotifyClient spotify.Client, spotifyID spotify.ID) ([]SongObj, error) {
	page, err := spotifyClient.GetPlaylistTracks(spotifyID)
	if err != nil {
		return []SongObj{}, err
	}
	var tracks []SongObj

	playlist, err := spotifyClient.GetPlaylist(spotifyID)
	if err != nil {
		return []SongObj{}, err
	}
	fmt.Printf("Fetching Songs from: %s\n", playlist.Name)
	for {
		for _, track := range page.Tracks {
			songObj, err := SongObjFromPlaylistTrack(spotifyClient, &track)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Skipping %q: %s\n", track.Track.Name, err)
				continue
			}
			fmt.Printf("Found: %s\n", track.Track.Name)
			tracks = append(tracks, songObj)
		}

		if strings.Contains(page.Next, "spotify") {
			opt := spotify.Options{
				Offset: &page.Offset,
			}
			page, err = spotifyClient.GetPlaylistTracksOpt(spotifyID, &opt, "")
			if err != nil {
				return []SongObj{}, err
			}
		} else {
			break
		}
	}
	return tracks, nil
}
