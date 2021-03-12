package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBestMatch(t *testing.T) {
	var tests = []struct {
		title   string
		artists []string
		dur     int
		link    string
	}{
		{"CLOUDS", []string{"NF"}, 243, "https://youtube.com/watch?v=vdR5ZeCD4Vk"},
		{"Hard On Yourself", []string{"Charlie Puth", "blackbear"}, 160, "https://youtube.com/watch?v=O7uc5Yqjhsg"},
		{"Toosie Slide", []string{"Drake"}, 247, "https://youtube.com/watch?v=eqMj9DTQcAQ"},
		{"Sun Is Shining", []string{"Axwell /\\ Ingrosso", "Axwell", "Sebastian Ingrosso"}, 255, "https://youtube.com/watch?v=7e-A7y9WesI"},
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
