package main

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/lastvoidtemplar/sabbac/internal/db"
	"github.com/lastvoidtemplar/sabbac/internal/fingerprint"
	"github.com/lastvoidtemplar/sabbac/internal/logger"
	"github.com/lastvoidtemplar/sabbac/internal/spectrogram"
	"github.com/lastvoidtemplar/sabbac/internal/youtube"
)

var custLogger *slog.Logger
var downloader youtube.YouTubeDownloader
var database *db.DB

type SongDTO struct {
	SongId string `json:"song_id"`
}

func main() {
	var err error

	custLogger = logger.New()

	downloader, err = youtube.NewYtDlpDownload("downloads", custLogger)
	if err != nil {
		panic(err.Error())
	}

	database, err = db.New("fingerprints.db", custLogger)

	if err != nil {
		panic(err.Error())
	}

	http.HandleFunc("POST /song", AddSongHandler)
	http.ListenAndServe(":3000", nil)
}

func AddSongHandler(w http.ResponseWriter, r *http.Request) {
	var dto SongDTO

	err := json.NewDecoder(r.Body).Decode(&dto)

	if err != nil {
		custLogger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid json format"))
		return
	}

	if !youtube.ValidateUrl(dto.SongId) {
		custLogger.With(slog.String("song_id", dto.SongId)).Warn("Invalid song id")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid song id"))
		return
	}

	go fingerPrintASong(downloader, database, dto.SongId, custLogger)
}

func fingerPrintASong(downloader youtube.YouTubeDownloader, db *db.DB, songId string, logger *slog.Logger) {
	wavePath, err := downloader.DownloadWav(songId)
	if err != nil {
		return
	}

	spectrogr, timePerColumn := spectrogram.STFT(wavePath, logger)

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

	for hash, videoTimestamps := range fingerprints {
		if len(videoTimestamps) > 50 {
			continue
		}
		db.InsertFingerprint(hash, videoTimestamps)
	}
}
