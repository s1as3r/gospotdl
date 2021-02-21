package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/bogem/id3v2"
)

// well..., GO doesnt have an an abs function. :-(
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// parseDuartion is used ot parse the length of a yt video.
func parseDuration(durationStr string) (int, error) {
	duraionRe := regexp.MustCompile(`PT(\d+H)?(\d+M)?(\d+S)?`)
	match := duraionRe.FindStringSubmatch(durationStr)
	if len(match) == 0 {
		return 0, fmt.Errorf("Error Parsing Duration: %s", durationStr)
	}
	var seconds int
	for _, i := range match[1:] {
		if i == "" {
			continue
		}
		p, _ := strconv.Atoi(i[:len(i)-1])
		if strings.Contains(i, "H") {
			seconds += p * 3600
		} else if strings.Contains(i, "M") {
			seconds += p * 60
		} else {
			seconds += p
		}
	}
	return seconds, nil
}

// parseUrl parses a spotify url and returns its ID.
func parseUrl(url string) (string, string) {
	if strings.Contains(url, "spotify.com") {
		url = strings.ReplaceAll(url, "\\", "/")
		list := strings.Split(url, "/")
		id := list[len(list)-1]
		id = strings.Split(id, "?")[0]
		if strings.Contains(url, "track") {
			return "track", id
		} else if strings.Contains(url, "album") {
			return "album", id
		} else if strings.Contains(url, "playlist") {
			return "playlist", id
		}
	}
	return "query", url
}

// newCmd is a helper function used to parse an ffmpeg command.
func newCmd(inputFile, outFile string, bitrate int) *exec.Cmd {
	return exec.Command("ffmpeg",
		"-v", "quiet",
		"-y",
		"-i", inputFile,
		"-acodec", "libmp3lame",
		"-b:a", fmt.Sprintf("%d", bitrate),
		"-af", "apad=pad_dur=2, dynaudnorm, loudnorm=I=-17",
		outFile,
	)
}

// setId3Data is used to set the metadata.
func setId3Data(filePath string, s SongObj) error {
	mp3File, err := id3v2.Open(filePath, id3v2.Options{Parse: true})
	if err != nil {
		return fmt.Errorf("Error while opening file(%s): %s\n", filePath, err)
	}
	defer mp3File.Close()
	mp3File.SetTitle(s.Name)
	mp3File.SetArtist(s.Artists[0].Name)
	mp3File.SetAlbum(s.Album.Name)
	mp3File.SetYear(s.Album.ReleaseDate)
	resp, err := http.Get(s.Album.Images[0].URL)
	if err != nil {
		return fmt.Errorf("Error Getting Cover Image: %s\n", err)
	}
	defer resp.Body.Close()
	img, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error reading cover image: %s", err)
	}
	albumCover := id3v2.PictureFrame{
		Encoding:    id3v2.EncodingUTF8,
		MimeType:    "image/jpeg",
		PictureType: id3v2.PTFrontCover,
		Description: "Front Cover",
		Picture:     img,
	}
	mp3File.AddAttachedPicture(albumCover)

	trackNumberTag := id3v2.NewEmptyTag()
	trackNumberFrame := id3v2.TextFrame{
		Encoding: id3v2.EncodingUTF8,
		Text:     strconv.Itoa(s.TrackNumber),
	}
	trackNumberTag.AddFrame("TRCK", trackNumberFrame)

	if err = mp3File.Save(); err != nil {
		return fmt.Errorf("Error while saving a tag: %s", err)
	}
	return nil
}
