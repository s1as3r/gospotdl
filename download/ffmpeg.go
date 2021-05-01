package download

import (
	"fmt"
	"os/exec"
)

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
