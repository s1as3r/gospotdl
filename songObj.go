package main

import (
	"fmt"

	"github.com/zmb3/spotify"
)

// SongObj is just a spotify.FullTrack along with a YT Link.
type SongObj struct {
	*spotify.FullTrack
	YoutubeLink string
}

// SongObjFromUrl gets a SongObj from a spotify url.
func SongObjFromID(spotifyClient spotify.Client, spotifyID spotify.ID) (SongObj, error) {
	trackMeta, err := spotifyClient.GetTrack(spotifyID)
	if err != nil {
		return SongObj{}, fmt.Errorf("Error Getting Track(%s): %s", spotifyID, err)
	}
	songName := trackMeta.Name
	var songArtists []string
	for _, i := range trackMeta.Artists {
		songArtists = append(songArtists, i.Name)
	}
	songDuration := trackMeta.Duration / 1000
	ytLink, err := GetBestMatch(songName, songArtists, songDuration)
	if err != nil {
		return SongObj{}, fmt.Errorf("Error Getting Youtube Link: %s", err)
	}
	return SongObj{
		FullTrack:   trackMeta,
		YoutubeLink: ytLink,
	}, nil
}

// SongObjFromQuery gets a SongObj from a query.
func SongObjFromQuery(spotifyClient spotify.Client, query string) (SongObj, error) {
	result, err := spotifyClient.Search(query, spotify.SearchTypeTrack)
	if err != nil {
		return SongObj{}, err
	}
	trackResults := result.Tracks
	if len(trackResults.Tracks) == 0 {
		return SongObj{}, fmt.Errorf("Found 0 Tracks Matching: %s", query)
	}

	trackMeta := trackResults.Tracks[0]
	songName := trackMeta.Name
	var songArtists []string
	for _, i := range trackMeta.Artists {
		songArtists = append(songArtists, i.Name)
	}
	songDuration := trackMeta.Duration / 1000
	ytLink, err := GetBestMatch(songName, songArtists, songDuration)
	if err != nil {
		return SongObj{}, err
	}
	return SongObj{
		FullTrack:   &trackMeta,
		YoutubeLink: ytLink,
	}, nil
}

// SongObjFromSimpleTrack gets a SongObj from a spotify.SimpleTrack
func SongObjFromSimpleTrack(spotifyClient spotify.Client, simpleTrack *spotify.SimpleTrack) (SongObj, error) {
	trackID := simpleTrack.ID
	song, err := SongObjFromID(spotifyClient, trackID)
	if err != nil {
		return SongObj{}, err
	}
	return song, nil
}

// SongObjFromPlaylistTrack gets a SongObj from a spotify.PlaylistTrack
func SongObjFromPlaylistTrack(spotifyClient spotify.Client, playlistTrack *spotify.PlaylistTrack) (SongObj, error) {
	trackID := playlistTrack.Track.ID
	song, err := SongObjFromID(spotifyClient, trackID)
	if err != nil {
		return SongObj{}, err
	}
	return song, nil
}
