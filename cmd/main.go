package main

import (
	"fmt"
	"log/slog"

	"github.com/lastvoidtemplar/sabbac/internal/fingerprint"
	"github.com/lastvoidtemplar/sabbac/internal/logger"
	"github.com/lastvoidtemplar/sabbac/internal/spectrogram"
)

func main() {
	logger := logger.New()

	// downloader, err := youtube.NewYtDlpDownload("downloads", logger)
	// if err != nil {
	// 	return
	// }

	// wavePath, err := downloader.DownloadWav("https://youtu.be/xMNz5QqHsF0?si=WF4hnOdFqPEudmDd")
	// if err != nil {
	// 	return
	// }

	wavePath := "downloads/Demi Lovato - Heart Attack (Lyrics).wav"
	spectrogr, timePerColumn := spectrogram.STFT(wavePath, logger)

	spectrogram.SpectrogramToImage("ouptput3.jpg", spectrogr)

	logger.With(
		slog.String("wav_path", wavePath),
	).Info("The spectrogram is generate successfully")

	peaks := fingerprint.FilterPeaks(spectrogr, timePerColumn)

	logger.With(
		slog.String("wav_path", wavePath),
	).Info("The peaks are filter successfully")

	fingerprints := fingerprint.GenerateFingerprints(wavePath, peaks)

	logger.With(
		slog.String("wav_path", wavePath),
	).Info("The fingerprints are generated successfully")

	for hash, timestamp := range fingerprints {
		fmt.Printf("%x %d\n", hash, len(timestamp))
	}

}
