package search

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetYoutubeLink(t *testing.T) {
	tests := []struct {
		title   string
		artists []string
		link    string
	}{
		{"CLOUDS", []string{"NF"}, "https://youtube.com/watch?v=JXOYZXb0no4"},
		{"Blinding Lights", []string{"The Weeknd"}, "https://youtube.com/watch?v=J7p4bzqLvCw"},
		{"Closer", []string{"The Chainsmokers", "Halsey"}, "https://youtube.com/watch?v=u-YGV5xt-jk"},
		{"Fearless Pt. II", []string{"Lost Sky", "Chris Linton"}, "https://youtube.com/watch?v=JTjmZZ1W2ew"},
	}

	for _, tCase := range tests {
		t.Run(tCase.title, func(t *testing.T) {
			ytLink, err := GetYoutubeLink(tCase.title, tCase.artists)
			if assert.Nil(t, err) {
				assert.Equal(t, tCase.link, ytLink)
			}
		})
	}
}
