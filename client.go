package main

import (
	"context"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
)

// GetSpotifyClient takes a clientId and a clientSecret key and returns
// a spotify.Client that can be used to make requests to the spotify API.
func GetSpotifyClient(clientId, clientSecret string) (spotify.Client, error) {
	config := &clientcredentials.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		TokenURL:     spotify.TokenURL,
	}

	token, err := config.Token(context.Background())
	if err != nil {
		return spotify.Client{}, err
	}

	client := spotify.Authenticator{}.NewClient(token)
	return client, nil
}
