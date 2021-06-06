package search

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/tidwall/gjson"
)

// getYtmJson gets the JSON String of a youtube music search
func getYtmJson(query string) (string, error) {
	payloadStr := `{context:{client:{clientName:"WEB_REMIX",clientVersion:"0.1",}},` +
		`query: "%s",` +
		`params:"Eg-KAQwIARAAGAAgACgAMABCAggBagoQBBADEAkQBRAK"}`
	payloadStr = fmt.Sprintf(payloadStr, query)
	payload := strings.NewReader(payloadStr)

	url := "https://music.youtube.com/youtubei/v1/" +
		"search?alt=json&key=AIzaSyC9XL3ZjWddXya6X74dJoCTL-WEYFDNX30"

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return "", fmt.Errorf("[getYtmJson] Error creating request: %s", err)
	}

	req.Header.Add("Referer", "https://music.youtube.com/search")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("[getYtmJson] Error making request: %s", err)
	}
	defer resp.Body.Close()

	jsonBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("[getYtmJson] Error reading response body: %s", err)
	}
	jsonStr := string(jsonBytes)

	return jsonStr, nil
}

// parseYtmResults parses the ytm JSON string and provides
// the yt links to the results
func parseYtmResults(json string) []string {
	parsedJson := gjson.Parse(json)
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

	var resutls []string
	for _, id := range ytIds.Array() {
		ytLink := fmt.Sprintf("https://youtube.com/watch?v=%s", id)
		resutls = append(resutls, ytLink)
	}
	return resutls
}

// GetYoutubeLink takes a song's name and its artists' names and
// returns it's youtube link
func GetYoutubeLink(songName string, songArtists []string) (string, error) {
	query := songName + " " + strings.Join(songArtists, ", ")
	ytmJson, err := getYtmJson(query)
	if err != nil {
		return "", fmt.Errorf("[GetYoutubeLink] Error getting json: %s", err)
	}

	results := parseYtmResults(ytmJson)
	if len(results) == 0 {
		return "", fmt.Errorf("0 matches found on youtube")
	}

	return results[0], nil
}
