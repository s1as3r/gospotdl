package download

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"

	"github.com/kkdai/youtube/v2"
	"github.com/s1as3r/gospotdl/search"
	"github.com/schollz/progressbar/v3"
)

// Download downloads a Song and sets the underlying metadata to the
// downloaded file.
func Download(s *search.Song) error {
	ytClient := youtube.Client{}

	video, err := ytClient.GetVideo(s.YoutubeLink)
	if err != nil {
		return err
	}

	formats := video.Formats
	sort.Slice(formats, func(i, j int) bool {
		return formats[i].Bitrate > formats[i].Bitrate
	})
	filteredFormats := youtube.FormatList{}
	for _, f := range formats {
		if f.AudioChannels > 0 {
			filteredFormats = append(filteredFormats, f)
		}
	}
	bestFormat := filteredFormats[0]

	stream, err := ytClient.GetStream(video, &bestFormat)
	if err != nil {
		return err
	}
	defer stream.Body.Close()

	baseFileName := s.Artists[0].Name + " - " + s.Name
	if runtime.GOOS == "windows" {
		baseFileName_ := ""
		for _, r := range baseFileName {
			if strings.Contains("*<>\\/[];|", string(r)) {
				continue
			} else if r == ':' {
				baseFileName_ += "-"
			} else if r == '"' {
				baseFileName_ += "'"
			} else {
				baseFileName_ += string(r)
			}
		}
		baseFileName = baseFileName_
	}
	mp3FileName := baseFileName + ".mp3"
	tempFileName := baseFileName + ".temp"
	tempFile, err := os.Create(tempFileName)
	defer os.Remove(tempFileName)
	defer tempFile.Close()

	_, err = io.Copy(tempFile, stream.Body)
	if err != nil {
		return err
	}

	cmd := newCmd(tempFileName, mp3FileName, bestFormat.Bitrate)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Error running FFMpeg: %s", err)
	}

	err = setId3Tags(mp3FileName, s)
	if err != nil {
		return err
	}

	return nil
}

// AsyncDownload is a wrapper around `Download` and downloads
// `max` songs concurrently
func AsyncDownload(tracks []*search.Song, max int) {
	var wg sync.WaitGroup
	tokens := make(chan struct{}, max)
	bar := progressbar.New(len(tracks))
	bar.Describe("Downloading ")
	for _, track := range tracks {
		wg.Add(1)
		go func(track *search.Song, wg *sync.WaitGroup) {
			defer wg.Done()
			defer bar.Add(1)

			tokens <- struct{}{}
			if err := Download(track); err != nil {
				fmt.Fprintf(os.Stderr, "Error downloading %s: %s", track.Name, err)
			}
			<-tokens
		}(track, &wg)
	}
	wg.Wait()
}
