package main

import (
	"fmt"
	"log/slog"

	"github.com/lastvoidtemplar/sabbac/internal/fingerfrint"
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
	spectrogr, timePerColumn := spectrogram.STFT(wavePath, logger)

	logger.With(
		slog.String("wav_path", wavePath),
	).Info("The spectrogram is generate successfully")

	peaks := fingerfrint.FilterPeaks(spectrogr, timePerColumn)

	logger.With(
		slog.String("wav_path", wavePath),
	).Info("The peaks are filter successfully")

	fmt.Println(peaks)
}
