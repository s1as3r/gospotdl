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

	resp, err := ytClient.GetStream(video, &video.Formats[0])
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

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		fmt.Sprintf("Downloading: %s", baseFileName),
	)

	_, err = io.Copy(io.MultiWriter(tempFile, bar), resp.Body)
	if err != nil {
		tempFile.Close()
		return err
	}

	cmd := newCmd(tempFileName, mp3FileName, video.Formats[0].Bitrate)
	if err := cmd.Run(); err != nil {
		tempFile.Close()
		return fmt.Errorf("Erro runnning ffmpeg command: %s", err)
	}
	tempFile.Close()
	if err := setId3Data(mp3FileName, s); err != nil {
		return fmt.Errorf("Error setting id3 data: %s", err)
	}
	if err := os.Remove(tempFileName); err != nil {
		return fmt.Errorf("Error removing temp file(%s): %s", tempFileName, err)
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
