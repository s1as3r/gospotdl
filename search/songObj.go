// Package search provides the required metadata and youtube
// link to download a song.
package search

import (
	"fmt"

	"github.com/zmb3/spotify"
)

type Song struct {
	*spotify.FullTrack
	YoutubeLink string
}

// FromId sets the metadata and the yt link of the underlying Song `s`
// from the spotify id of a track.
func (s *Song) FromId(client *spotify.Client, spotifyId spotify.ID) error {
	trackMeta, err := client.GetTrack(spotifyId)
	if err != nil {
		return fmt.Errorf("[Song.FromId] Error getting track: %s", err)
	}
	songName := trackMeta.Name
	var songArtsits []string
	for _, artist := range trackMeta.Artists {
		songArtsits = append(songArtsits, artist.Name)
	}

	duration := trackMeta.Duration / 1000
	ytLink, err := GetYoutubeLink(songName, songArtsits, duration)
	if err != nil {
		return fmt.Errorf("[Song.FromId] Erro getting Youtube link: %s", err)
	}

	s.FullTrack = trackMeta
	s.YoutubeLink = ytLink
	return nil
}

// FromQuery gets search for the `query` and get's sets the metadata
// and yt link of the underlying Song `s`
func (s *Song) FromQuery(client *spotify.Client, query string) error {
	searchResults, err := client.Search(query, spotify.SearchTypeTrack)
	if err != nil {
		return fmt.Errorf("[Song.FromQuery] Error while searching: %s", err)
	}
	tracks := searchResults.Tracks.Tracks

	trackMeta := tracks[0]
	songName := trackMeta.Name
	var songArtsits []string
	for _, artist := range trackMeta.Artists {
		songArtsits = append(songArtsits, artist.Name)
	}

	duration := trackMeta.Duration / 1000
	ytLink, err := GetYoutubeLink(songName, songArtsits, duration)
	if err != nil {
		return fmt.Errorf("[Song.FromQuery] Error getting Youtube link: %s", err)
	}

	s.FullTrack = &trackMeta
	s.YoutubeLink = ytLink
	return nil
}

// FromSimpleTrack sets the underlying Songs `FullTrack` and yt link by using a
// SimpleTrack
func (s *Song) FromSimpleTrack(client *spotify.Client, simpleTrack *spotify.SimpleTrack) error {
	trackId := simpleTrack.ID
	err := s.FromId(client, trackId)
	if err != nil {
		return fmt.Errorf("[Song.FromSimpleTrack] : %s", err)
	}
	return nil
}

// FromPlaylistTrack sets the underlying Songs `FullTrack` and yt link by using a
// PlaylistTrack
func (s *Song) FromPlaylistTrack(playlistTrack *spotify.PlaylistTrack) error {
	trackMeta := playlistTrack.Track
	songName := trackMeta.Name
	var songArtsits []string
	for _, artist := range trackMeta.Artists {
		songArtsits = append(songArtsits, artist.Name)
	}

	duration := trackMeta.Duration / 1000
	ytLink, err := GetYoutubeLink(songName, songArtsits, duration)
	if err != nil {
		return fmt.Errorf("[Song.FromPlaylistTrack] Error getting Youtube link: %s", err)
	}

	s.FullTrack = &trackMeta
	s.YoutubeLink = ytLink

	return nil
}
