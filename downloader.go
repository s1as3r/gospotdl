package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/kkdai/youtube/v2"
	"github.com/schollz/progressbar/v3"
)

// Download downlaods a song from Youtube
func Download(s SongObj) error {
	ytClient := youtube.Client{}

	video, err := ytClient.GetVideo(s.YoutubeLink)
	if err != nil {
		return err
	}

	format := video.Formats.FindByItag(251) // Best quality audio available?
	if format == nil {
		format = video.Formats.FindByItag(140) // Fallback Quality
	}
	resp, err := ytClient.GetStream(video, format)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var artistNames []string
	for _, i := range s.Artists {
		artistNames = append(artistNames, i.Name)
	}

	baseFileName := strings.Join(artistNames, ", ") + " - " + s.Name
	tempFileName := baseFileName + ".temp"
	mp3FileName := baseFileName + ".mp3"

	tempFile, err := os.Create(tempFileName)
	if err != nil {
		return err
	}
	defer os.Remove(tempFileName)
	defer tempFile.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		fmt.Sprintf("Downloading: %s", baseFileName),
	)

	_, err = io.Copy(io.MultiWriter(tempFile, bar), resp.Body)
	if err != nil {
		return err
	}

	cmd := newCmd(tempFileName, mp3FileName, video.Formats[0].Bitrate)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Erro runnning ffmpeg command: %s", err)
	}
	if err := SetId3Data(mp3FileName, s); err != nil {
		return fmt.Errorf("Error setting id3 data: %s", err)
	}
	return nil
}

// DownloadMulti downlaods multiple songs in parallel.
func DownloadMulti(tracks []SongObj) {
	var wg sync.WaitGroup
	tokens := make(chan struct{}, 4)
	for _, track := range tracks {
		wg.Add(1)
		go func(track SongObj, wg *sync.WaitGroup) {
			defer wg.Done()
			tokens <- struct{}{}
			if err := Download(track); err != nil {
				fmt.Fprintf(os.Stderr, "Error Downloading %s: %s", track.Name, err)
			}
			<-tokens
		}(track, &wg)
	}
	wg.Wait()
}
