package search

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

type YtmResult struct {
	Url      string
	Duration int
}

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
// a list of YtmResult
func parseYtmResults(json string) []YtmResult {
	gjson.AddModifier("seconds", func(json, arg string) string {
		json = strings.ReplaceAll(json, "\"", "")
		dur := strings.Split(json, ":")
		if len(dur) != 2 {
			return "0"
		}
		minutes, _ := strconv.Atoi(dur[0])
		seconds, _ := strconv.Atoi(dur[1])
		return strconv.Itoa(minutes*60 + seconds)
	})

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
		".videoId|@flatten").Array()

	durations := parsedJson.Get("contents" +
		".sectionListRenderer" +
		".contents.#" +
		".musicShelfRenderer" +
		".contents.#" +
		".musicResponsiveListItemRenderer" +
		".flexColumns.1" +
		".musicResponsiveListItemFlexColumnRenderer" +
		".text" +
		".runs.4" +
		".text.@seconds|@flatten").Array()

	var resutls []YtmResult
	for k, id := range ytIds {
		ytLink := fmt.Sprintf("https://youtube.com/watch?v=%s", id)
		duration := durations[k].Int()
		resutls = append(resutls, YtmResult{
			Url:      ytLink,
			Duration: int(duration),
		})
	}
	return resutls
}

// getBestYtMatch returns the best match song for a given query
// and duration (in seconds). It sorts the results by calculating
// the difference in duration between a match and given duration
// and returns the match with the least duartion difference.
func getBestYtMatch(query string, duration int) (string, error) {
	ytmJson, err := getYtmJson(query)
	if err != nil {
		return "", fmt.Errorf("[getBestYtMatch] Error getting json: %s", err)
	}
	results := parseYtmResults(ytmJson)
	if len(results) == 0 {
		return "", fmt.Errorf("0 matches found on YTm for: %s", query)
	}
	sort.Slice(results, func(i, j int) bool {
		diff1 := abs(results[i].Duration - duration)
		diff2 := abs(results[j].Duration - duration)
		return diff1 < diff2
	})

	return results[0].Url, nil
}

// GetYoutubeLink takes a song's name, its artists' names and
// it's duration in seconds and returns it's youtube link
func GetYoutubeLink(songName string, songArtists []string, duration int) (string, error) {
	query := songName + " " + strings.Join(songArtists, ", ")
	ytLink, err := getBestYtMatch(query, duration)
	if err != nil {
		return "", err
	}
	return ytLink, nil
}
