package search

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zmb3/spotify"
)

func TestFromId(t *testing.T) {
	tests := []struct {
		title string
		id    spotify.ID
	}{
		{"DRIFTING", spotify.ID("2hRlHXzOf14ArYmOPeAXsa")},
		{"The Search", spotify.ID("3oLe5ZILASG8vU5dxIMfLY")},
		{"Blinding Lights", spotify.ID("0VjIjW4GlUZAMYd2vXMi3b")},
		{"Humnava Mere", spotify.ID("0loZn1c5heXie7OAtvK6nH")},
		{"Afterglow", spotify.ID("0E4Y1XIbs8GrAT1YqVy6dq")},
	}

	for _, tCase := range tests {
		t.Run(tCase.title, func(t *testing.T) {
			song := Song{}
			err := song.FromId(testClient, tCase.id)
			if assert.Nil(t, err) {
				assert.Equal(t, tCase.title, song.Name)
			}
		})
	}
}

func TestFromQuery(t *testing.T) {
	tests := []struct {
		title string
		query string
	}{
		{"DRIFTING", "DRIFTING NF"},
		{"The Search", "The Search NF"},
		{"Blinding Lights", "Blinding Lights The Weeknd"},
		{"Humnava Mere", "Humnava Mere Jubin Nautiyal"},
		{"Afterglow", "Afterglow Ed Sheeran"},
	}

	for _, tCase := range tests {
		t.Run(tCase.title, func(t *testing.T) {
			song := Song{}
			err := song.FromQuery(testClient, tCase.query)
			if assert.Nil(t, err) {
				assert.Equal(t, tCase.title, song.Name)
			}
		})
	}
}
