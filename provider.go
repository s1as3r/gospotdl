package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/tidwall/gjson"
)

func getYtmResponse(query string) (string, error) {
	payload_str := `{context:{client:{clientName:"WEB_REMIX",clientVersion:"0.1",}},` +
		`query: "%s",` +
		`params:"Eg-KAQwIARAAGAAgACgAMABCAggBagoQBBADEAkQBRAK"}`
	payload_str = fmt.Sprintf(payload_str, query)
	payload := strings.NewReader(payload_str)
	url := "https://music.youtube.com/youtubei/v1/search?alt=json&key=AIzaSyC9XL3ZjWddXya6X74dJoCTL-WEYFDNX30"
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return "", fmt.Errorf("Error creating request: %s", err)
	}
	req.Header.Add("Referer", "https://music.youtube.com/search")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error getting data from ytmusic: %s", err)
	}
	defer res.Body.Close()

	jsonBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("Error reading response body:%s", err)
	}
	jsonString := string(jsonBytes)

	return jsonString, nil
}

func getYtmResults(response string) []string {
	parsedJson := gjson.Parse(response)
	ytIds := parsedJson.Get("contents" +
		".sectionListRenderer" +
		".contents.#" +
		".musicShelfRenderer" +
		".contents.#" +
		".musicResponsiveListItemRenderer" +
		".overlay" +
		".musicItemThumbnailOverlayRenderer" +
		".content" +
		".musicPlayButtonRenderer" +
		".playNavigationEndpoint" +
		".watchEndpoint" +
		".videoId|@flatten")

	var results []string
	for _, id := range ytIds.Array() {
		ytLink := fmt.Sprintf("https://youtube.com/watch?v=%s", id)
		results = append(results, ytLink)
	}
	return results
}

// GetBestMatch gets the best youtube match for a song.
func GetBestMatch(songName string, songArtists []string) (string, error) {
	query := songName + " " + strings.Join(songArtists, ", ")
	ytmResponse, err := getYtmResponse(query)
	if err != nil {
		return "", fmt.Errorf("Error getting ytm response: %s", err)
	}

	results := getYtmResults(ytmResponse)

	if len(results) == 0 {
		return "", fmt.Errorf("0 Matches Found")
	}

	return results[0], err
}
