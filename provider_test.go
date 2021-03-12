package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBestMatch(t *testing.T) {
	var tests = []struct {
		title   string
		artists []string
		link    string
	}{
		{"CLOUDS", []string{"NF"}, "https://youtube.com/watch?v=JXOYZXb0no4"},
		{"Hard On Yourself", []string{"Charlie Puth", "blackbear"}, "https://youtube.com/watch?v=O7uc5Yqjhsg"},
		{"Toosie Slide", []string{"Drake"}, "https://youtube.com/watch?v=eqMj9DTQcAQ"},
		{"Sun Is Shining", []string{"Axwell /\\ Ingrosso", "Axwell", "Sebastian Ingrosso"}, "https://youtube.com/watch?v=7e-A7y9WesI"},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			got, err := GetBestMatch(tt.title, tt.artists)
			if assert.Nil(t, err) {
				assert.Equal(t, tt.link, got)
			}
		})
	}
}
