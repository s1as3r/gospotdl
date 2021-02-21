package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/kkdai/youtube"
	"github.com/schollz/progressbar"
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

	fileName := strings.Join(artistNames, ", ") + " - " + s.Name
	file, err := os.Create(fileName + ".temp")
	if err != nil {
		return err
	}
	defer file.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		fmt.Sprintf("Downloading: %s", fileName),
	)
	_, err = io.Copy(io.MultiWriter(file, bar), resp.Body)
	if err != nil {
		return err
	}
	cmd := newCmd(fileName+".temp", fileName+".mp3", video.Formats[0].Bitrate)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Erro runnning ffmpeg command: %s", err)
	}
	if err := setId3Data(fileName+".mp3", s); err != nil {
		return fmt.Errorf("Error setting id3 data: %s", err)
	}
	if err := os.Remove(fileName + ".temp"); err != nil {
		return fmt.Errorf("Error removing temp file(%s.temp): %s", fileName, err)
	}
	return nil
}
