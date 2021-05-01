package main

import (
	"flag"
	"log"
	"os"

	"github.com/s1as3r/gospotdl/search"
)

const (
	clientId     = "b9640826ab8949e4b34a4eeee5b26ce6"
	clientSecret = "c25b0399324d49a08d8210138cbf9e27"
)

func main() {
	client, err := search.GetSpotifyClient(clientId, clientSecret)
	if err != nil {
		log.Fatalf("Error getting spotify client: %s", err)
	}

	var dir string
	var max int
	flag.StringVar(&dir, "output", ".", "output directory")
	flag.StringVar(&dir, "o", ".", "output directory")
	flag.IntVar(&max, "max", 4, "max concurrent downloads")
	flag.IntVar(&max, "m", 4, "max concurrent downloads")
	flag.Parse()

	if err := os.Chdir(dir); err != nil {
		log.Fatalf("Error changing directory: %s", err)
	}

	for _, arg := range flag.Args() {
		handleArg(client, arg, max)
	}
}
