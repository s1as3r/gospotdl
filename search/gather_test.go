package search

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zmb3/spotify"
)

const (
	spotifyClientId     = "b9640826ab8949e4b34a4eeee5b26ce6"
	spotifyClientSecret = "c25b0399324d49a08d8210138cbf9e27"
)

var testClient, _ = GetSpotifyClient(spotifyClientId, spotifyClientSecret)

func TestGetAlbumTracks(t *testing.T) {
	tests := []struct {
		title   string
		id      spotify.ID
		nTracks int
	}{
		{"CLOUDS", spotify.ID("7eQGtkzCgrIWDOe76E9F8t"), 11},
		{"After Hours", spotify.ID("4yP0hdKOZPNshxUOjY0cZj"), 14},
	}

	for _, tCase := range tests {
		t.Run(tCase.title, func(t *testing.T) {
			tracks, err := GetAlbumTracks(testClient, tCase.id)
			if assert.Nil(t, err) {
				assert.Equal(t, tCase.nTracks, len(tracks))
			}
		})
	}
}

func TestGetPlaylistTracks(t *testing.T) {
	var tests = []struct {
		title   string
		id      spotify.ID
		nTracks int
	}{
		{"Test-Playlist-1", spotify.ID("6wuR62qIWHxIO4xyln6AKV"), 9},
	}

	for _, tCase := range tests {
		t.Run(tCase.title, func(t *testing.T) {
			tracks, err := GetPlaylistTracks(testClient, tCase.id)
			if assert.Nil(t, err) {
				assert.Equal(t, tCase.nTracks, len(tracks))
			}
		})
	}
}
