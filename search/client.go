package search

import (
	"context"
	"fmt"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2/clientcredentials"
)

// GetSpotifyClient returns a spotify client that's used for searching
// and getting metadata from spotify.
func GetSpotifyClient(clientId, clientSecret string) (*spotify.Client, error) {
	config := &clientcredentials.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		TokenURL:     spotify.TokenURL,
	}

	token, err := config.Token(context.Background())
	if err != nil {
		return &spotify.Client{}, fmt.Errorf("[GetSpotifyClient] Error getting token: %s", err)
	}

	client := spotify.Authenticator{}.NewClient(token)
	return &client, nil
}
