package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseArg(t *testing.T) {
	var tests = []struct {
		url, argtype, id string
	}{
		{"https://open.spotify.com/track/18czZN7uruOjftj71Kt8oj/", "track", "18czZN7uruOjftj71Kt8oj"},
		{"https://open.spotify.com/track/18czZN7uruOjftj71Kt8oj?si=0tf31swlTz2_V6uQVF_r1A&", "track", "18czZN7uruOjftj71Kt8oj"},
		{"https://open.spotify.com/playlist/37i9dQZEVXbMDoHDwVN2tF?si=NKDZ03KdRrSBbkUgUGpJiw&utm_source=copy-link", "playlist", "37i9dQZEVXbMDoHDwVN2tF"},
		{"https://open.spotify.com/album/46xdC4Qcvscfs3Ai2RIHcv?si=24utlQHbSjqaCd2a45977Q&utm_source=copy-link", "album", "46xdC4Qcvscfs3Ai2RIHcv"},
		{"NF Clouds", "query", "NF Clouds"},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s,%s", tt.argtype, tt.id)
		t.Run(testname, func(t *testing.T) {
			argtype, id := parseArg(tt.url)
			assert.Equal(t, tt.argtype, argtype)
			assert.Equal(t, tt.id, id)
		})
	}
}
