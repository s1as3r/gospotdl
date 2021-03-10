// A Golang implementation of spotDl (python)
package main

import (
	"flag"
	"log"
	"os"
)

const (
	spotifyClientId     = "b9640826ab8949e4b34a4eeee5b26ce6"
	spotifyClientSecret = "c25b0399324d49a08d8210138cbf9e27"
)

func main() {
	spotifyClient, err := GetSpotifyClient(spotifyClientId, spotifyClientSecret)
	if err != nil {
		log.Fatalf("Error getting Spotify Client: %s\n", err)
	}
	var dir string
	flag.StringVar(&dir, "output", ".", "Output Directory")
	flag.StringVar(&dir, "o", ".", "Output Directory")
	flag.Parse()

	if err := os.Chdir(dir); err != nil {
		log.Fatalf("Error changing directory: %s", err)
	}

	for _, url := range flag.Args() {
		handleArg(url, spotifyClient)
	}
}
