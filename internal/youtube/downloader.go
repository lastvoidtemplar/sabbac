package youtube

import (
	"bytes"
	"errors"
	"log/slog"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
)

var ErrInvalidDownloadUrl = errors.New("invalid download youtube url")
var ErrInvalidDirPath = errors.New("invalid path to dir")
var ErrUnSuccessfulDownload = errors.New("unsuccessful download")

type YouTubeDownloader interface {
	DownloadWav(url string) (string, error)
}

type ytdlpDownloader struct {
	logger    *slog.Logger
	outputDir string
}

func NewYtDlpDownload(outputDir string, logger *slog.Logger) (YouTubeDownloader, error) {
	st, err := os.Stat(outputDir)

	if err != nil {
		logger.With(slog.String("output_dir", outputDir), slog.String("err", err.Error())).Error("Invalid path")
		return nil, ErrInvalidDirPath
	}

	if !st.IsDir() {
		logger.With(slog.String("output_dir", outputDir)).Error("The file is not a dir")
		return nil, ErrInvalidDirPath
	}

	logger.With(slog.String("output_dir", outputDir)).Info("YtDlp downloader is created successfully")
	return &ytdlpDownloader{
		logger:    logger,
		outputDir: outputDir,
	}, nil
}

func ValidateUrl(rawUrl string) bool {
	u, err := url.ParseRequestURI(rawUrl)

	if err != nil {
		return false
	}

	if u.Host != "youtu.be" {
		return false
	}

	return true
}

func (downloader *ytdlpDownloader) DownloadWav(rawUrl string) (string, error) {
	if !ValidateUrl(rawUrl) {
		downloader.logger.With(slog.String("song_id", rawUrl)).Error("Invalid url")
		return "", ErrInvalidDownloadUrl
	}

	cmd := exec.Command(
		"venv/bin/yt-dlp",
		"-x",
		"--audio-format", "wav",
		"-o", filepath.Join(downloader.outputDir, "%(title)s.%(ext)s"),
		"--postprocessor-args", "-ac 1",
		rawUrl,
	)

	cmdOutput, err := cmd.Output()
	if err != nil {
		downloader.logger.With(slog.String("song_id", rawUrl), slog.String("err", err.Error())).Error("YtDlp failed")
		return "", err
	}

	ind1 := bytes.LastIndexByte(cmdOutput[:len(cmdOutput)-1], '\n')
	if ind1 == -1 {
		downloader.logger.With(slog.String("song_id", rawUrl), slog.String("ytdlp_output", string(cmdOutput))).Error("No new line found")
		return "", ErrUnSuccessfulDownload
	}

	ind2 := bytes.LastIndexByte(cmdOutput[:ind1], '\n')
	if ind2 == -1 {
		downloader.logger.With(slog.String("song_id", rawUrl), slog.String("ytdlp_output", string(cmdOutput))).Error("No new line found")
		return "", ErrUnSuccessfulDownload
	}

	destinationOutput := cmdOutput[ind2+1 : ind1]

	ind3 := bytes.IndexByte(destinationOutput, ':')
	if ind3 == -1 {
		downloader.logger.With(slog.String("song_id", rawUrl), slog.String("dest_output", string(destinationOutput))).Error("No colon found")
		return "", ErrUnSuccessfulDownload
	}

	outputPath := string(destinationOutput[ind3+2:])
	downloader.logger.With(slog.String("song_id", rawUrl), slog.String("output_path", outputPath)).Info("Successful audio download")

	return outputPath, nil
}
