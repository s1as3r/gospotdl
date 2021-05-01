// Package search provides the required metadata and youtube
// link to download a song.
package search

import (
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
		return err
	}
	songName := trackMeta.Name
	var songArtsits []string
	for _, artist := range trackMeta.Artists {
		songArtsits = append(songArtsits, artist.Name)
	}

	ytLink, err := GetYoutubeLink(songName, songArtsits)
	if err != nil {
		return err
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
		return err
	}
	tracks := searchResults.Tracks.Tracks

	trackMeta := tracks[0]
	songName := trackMeta.Name
	var songArtsits []string
	for _, artist := range trackMeta.Artists {
		songArtsits = append(songArtsits, artist.Name)
	}

	ytLink, err := GetYoutubeLink(songName, songArtsits)
	if err != nil {
		return err
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
		return err
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

	ytLink, err := GetYoutubeLink(songName, songArtsits)
	if err != nil {
		return err
	}

	s.FullTrack = &trackMeta
	s.YoutubeLink = ytLink

	return nil
}
