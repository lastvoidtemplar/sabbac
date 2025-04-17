package main

import (
	"github.com/lastvoidtemplar/sabbac/internal/logger"
	"github.com/lastvoidtemplar/sabbac/internal/spectrogram"
)

func main() {
	logger := logger.New()

	// downloader, err := youtube.NewYtDlpDownload("downloads", logger)
	// if err != nil {
	// 	return
	// }

	// wavePath, err := downloader.DownloadWav("https://youtu.be/eAy2eVlvaZQ?si=s0Mt9ImSTTVqfva8")
	// if err != nil {
	// 	return
	// }

	wavePath := "downloads/SABATON - The Valley Of Death (Official Lyric Video).wav"
	spectrogram.NewSpectrogram(wavePath, logger)
}
