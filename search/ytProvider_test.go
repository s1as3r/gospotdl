package search

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetYoutubeLink(t *testing.T) {
	tests := []struct {
		title    string
		artists  []string
		duration int
		link     string
	}{
		{"CLOUDS", []string{"NF"}, 244, "https://youtube.com/watch?v=JXOYZXb0no4"},
		{"Blinding Lights", []string{"The Weeknd"}, 201, "https://youtube.com/watch?v=2ru92T7Y5z0"},
		{"Fearless Pt. II", []string{"Lost Sky", "Chris Linton"}, 194, "https://youtube.com/watch?v=JTjmZZ1W2ew"},
	}

	for _, tCase := range tests {
		t.Run(tCase.title, func(t *testing.T) {
			ytLink, err := GetYoutubeLink(tCase.title, tCase.artists, tCase.duration)
			if assert.Nil(t, err) {
				assert.Equal(t, tCase.link, ytLink)
			}
		})
	}
}
