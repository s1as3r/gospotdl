package search

import (
	"fmt"
	"os"

	"github.com/zmb3/spotify"
)

// GetAlbumTracks gets all the tracks in in album using its spotify id.
func GetAlbumTracks(client *spotify.Client, spotifyId spotify.ID) ([]*Song, error) {
	page, err := client.GetAlbumTracks(spotifyId)
	if err != nil {
		return []*Song{}, err
	}
	var tracks []*Song

	album, err := client.GetAlbum(spotifyId)
	if err != nil {
		return []*Song{}, err
	}
	fmt.Printf("Fetching Songs from: %s\n", album.Name)
	for {
		for _, track := range page.Tracks {
			song := &Song{}
			err := song.FromSimpleTrack(client, &track)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Skipping %s as no match was found on Yoututbe", track.Name)
				continue
			}
			fmt.Printf("Found: %s\n", track.Name)
			tracks = append(tracks, song)
		}

		if page.Total > len(tracks) {
			opt := spotify.Options{
				Offset: &page.Limit,
			}
			page, err = client.GetAlbumTracksOpt(spotifyId, &opt)
			if err != nil {
				return []*Song{}, err
			}
		} else {
			break
		}
	}
	return tracks, nil
}

// GetPlaylistTracks gets all the songs of a playlist using its spotify id.
func GetPlaylistTracks(client *spotify.Client, spotifyId spotify.ID) ([]*Song, error) {
	page, err := client.GetPlaylistTracks(spotifyId)
	if err != nil {
		return []*Song{}, err
	}
	var tracks []*Song

	playlist, err := client.GetPlaylist(spotifyId)
	if err != nil {
		return []*Song{}, err
	}
	fmt.Printf("Fetching Songs from: %s\n", playlist.Name)
	for {
		for _, track := range page.Tracks {
			song := &Song{}
			err := song.FromPlaylistTrack(&track)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Skipping %q: %s\n", track.Track.Name, err)
				continue
			}
			fmt.Printf("Found: %s\n", track.Track.Name)
			tracks = append(tracks, song)
		}

		if page.Total > len(tracks) {
			opt := spotify.Options{
				Offset: &page.Limit,
			}
			page, err = client.GetPlaylistTracksOpt(spotifyId, &opt, "")
			if err != nil {
				return []*Song{}, err
			}
		} else {
			break
		}
	}
	return tracks, nil
}
