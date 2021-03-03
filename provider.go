package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
)

const (
	ytApiKey = "AIzaSyC05R5TtMtuyUFZKASO2H1Cvv6N2la6e60"
)

type songData struct {
	Link   string
	Length int
}

type diffResult struct {
	timeDiff int
	Link     string
}

// getSeconds gets the length of a yt video in seconds.
func getSeconds(videoId string) (int, error) {
	videoUrl := "https://youtube.googleapis.com/youtube/v3/videos"
	parsedUrl := fmt.Sprintf("%s?part=contentDetails&key=%s&id=%s",
		videoUrl, ytApiKey, videoId)
	response, err := http.Get(parsedUrl)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	var duration struct {
		Items []struct {
			ContentDetails struct {
				DurationStr string `json:"duration"`
			} `json:"contentDetails"`
		} `json:"items"`
	}

	if err := json.NewDecoder(response.Body).Decode(&duration); err != nil {
		return 0, err
	}
	durationStr := duration.Items[0].ContentDetails.DurationStr
	seconds, err := parseDuration(durationStr)
	if err != nil {
		return 0, err
	}
	return seconds, nil
}

// mapIdToSongData is used to get a songData struct.
func mapIdToSongData(id string) (songData, error) {
	var song songData
	song.Link = fmt.Sprintf("https://youtube.com/watch?v=%s", id)
	length, err := getSeconds(id)
	if err != nil {
		return song, err
	}
	song.Length = length
	return song, nil
}

// queryAndSimplify queries youtube and returns a list of simplified matches.
func queryAndSimplify(query string) ([]songData, error) {
	query = strings.Join(strings.Split(query, " "), "%20")
	searchUrl := "https://www.googleapis.com/youtube/v3/search"
	parsedUrl := fmt.Sprintf("%s?part=snippet&maxResults=5&type=video&q=%s&key=%s",
		searchUrl, query, ytApiKey)
	response, err := http.Get(parsedUrl)
	if err != nil {
		return []songData{}, err
	}
	defer response.Body.Close()

	var parsedJson struct {
		Items []struct {
			Id struct {
				VideoId string `json:"videoId"`
			} `json:"id"`
		} `json:"items"`
	}
	if err := json.NewDecoder(response.Body).Decode(&parsedJson); err != nil {
		return []songData{}, err
	}
	var results []songData
	for _, i := range parsedJson.Items {
		data, err := mapIdToSongData(i.Id.VideoId)
		if err != nil {
			return []songData{}, err
		}
		results = append(results, data)
	}
	return results, nil
}

// getYtResults takes gets yt results and also the difference in duration
// between the result and the actual spotify track.
func getYtResults(songName string, songArtists []string, songDuration int) ([]diffResult, error) {
	artistStr := strings.Join(songArtists, ", ")
	query := fmt.Sprintf("%s - %s", songName, artistStr)
	searchResults, err := queryAndSimplify(query)
	if err != nil {
		return []diffResult{}, err
	}
	var results []diffResult
	for _, result := range searchResults {
		timeDiff := abs(result.Length - songDuration)
		results = append(results, diffResult{
			timeDiff: timeDiff,
			Link:     result.Link,
		})
	}
	return results, nil
}

// GetBestMatch gets the best youtube match for a song.
func GetBestMatch(songName string, songArtists []string, songDuration int) (string, error) {
	results, err := getYtResults(songName, songArtists, songDuration)
	if err != nil {
		return "", err
	}
	if len(results) == 0 {
		return "", fmt.Errorf("0 Matches Found")
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].timeDiff < results[j].timeDiff
	})

	return results[0].Link, err
}
