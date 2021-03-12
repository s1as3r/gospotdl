package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type YtMusicResponse struct {
	Contents1 struct {
		SectionListRenderer struct {
			Contents2 []struct {
				MusicShelfRenderer struct {
					Contents3 []struct {
						MusicResponsiveListItemRenderer struct {
							Overlay struct {
								MusicItemThumbnailOverlayRenderer struct {
									Content4 struct {
										MusicPlayButtonRenderer struct {
											PlayNavigationEndpoint struct {
												WatchEndpoint struct {
													VideoId string `json:"videoId"`
												} `json:"watchEndpoint"`
											} `json:"playNavigationEndpoint"`
										} `json:"musicPlayButtonRenderer"`
									} `json:"content"`
								} `json:"musicItemThumbnailOverlayRenderer"`
							} `json:"overlay"`
						} `json:"musicResponsiveListItemRenderer"`
					} `json:"contents"`
				} `json:"musicShelfRenderer"`
			} `json:"contents"`
		} `json:"sectionListRenderer"`
	} `json:"contents"`
}

func getYtmResponse(query string) (YtMusicResponse, error) {
	payload_str := `{context:{client:{clientName:"WEB_REMIX",clientVersion:"0.1",}},query: "%s", params:"Eg-KAQwIARAAGAAgACgAMABCAggBagoQBBADEAkQBRAK"}`
	payload_str = fmt.Sprintf(payload_str, query)
	payload := strings.NewReader(payload_str)
	url := "https://music.youtube.com/youtubei/v1/search?alt=json&key=AIzaSyC9XL3ZjWddXya6X74dJoCTL-WEYFDNX30"
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return YtMusicResponse{}, fmt.Errorf("Error creating request: %s", err)
	}
	req.Header.Add("Referer", "https://music.youtube.com/search")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return YtMusicResponse{}, fmt.Errorf("Error getting data from ytmusic: %s", err)
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	var response YtMusicResponse
	if err := decoder.Decode(&response); err != nil {
		return YtMusicResponse{}, fmt.Errorf("Error decoding ytmusic json: %s", err)
	}
	return response, nil
}

func getYtmResults(response *YtMusicResponse) []string {
	var results []string
	contentBlocks := response.
		Contents1.
		SectionListRenderer.
		Contents2
	for _, i := range contentBlocks {
		for _, block := range i.MusicShelfRenderer.Contents3 {
			ytId := block.
				MusicResponsiveListItemRenderer.
				Overlay.
				MusicItemThumbnailOverlayRenderer.
				Content4.
				MusicPlayButtonRenderer.
				PlayNavigationEndpoint.
				WatchEndpoint.
				VideoId
			ytLink := fmt.Sprintf("https://youtube.com/watch?v=%s", ytId)
			results = append(results, ytLink)
		}
	}

	return results
}

// GetBestMatch gets the best youtube match for a song.
func GetBestMatch(songName string, songArtists []string, songDuration int) (string, error) {
	query := songName + " " + strings.Join(songArtists, ", ")
	ytmResponse, err := getYtmResponse(query)
	if err != nil {
		return "", fmt.Errorf("Error getting ytm response: %s", err)
	}

	results := getYtmResults(&ytmResponse)

	if len(results) == 0 {
		return "", fmt.Errorf("0 Matches Found")
	}

	return results[0], err
}
