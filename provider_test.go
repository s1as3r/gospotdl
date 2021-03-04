package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSeconds(t *testing.T) {
	var tests = []struct {
		videoId string
		dur     int
	}{
		{"5wivUfSS-T0", 8},
		{"7Hlb8YX2-W8", 39005},
		{"iZnLZFRylbs", 311},
	}

	for _, tt := range tests {
		t.Run(tt.videoId, func(t *testing.T) {
			got, err := getSeconds(tt.videoId)
			if assert.Nil(t, err) {
				assert.Equal(t, tt.dur, got)
			}
		})
	}
}

func TestGetBestMatch(t *testing.T) {
	var tests = []struct {
		title   string
		artists []string
		dur     int
		link    string
	}{
		{"Clouds", []string{"NF"}, 243, "https://youtube.com/watch?v=Z0Wc2-qDdn0"},
		{"Hard On Yourself", []string{"Charlie Puth", "blackbear"}, 160, "https://youtube.com/watch?v=L_RkUodwy6Y"},
		{"Toosie Slide", []string{"Drake"}, 247, "https://youtube.com/watch?v=EGr-fMU6Kd4"},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			got, err := GetBestMatch(tt.title, tt.artists, tt.dur)
			if assert.Nil(t, err) {
				assert.Equal(t, tt.link, got)
			}
		})
	}
}
