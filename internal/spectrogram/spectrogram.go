package spectrogram

import (
	"fmt"
	"log/slog"
)

func NewSpectrogram(wavePath string, logger *slog.Logger) [][]complex128 {
	wavParser, err := newWavParser(wavePath, logger)
	if err != nil {
		logger.Error("Couldn`t create wav parser")
		return nil
	}
	defer wavParser.Close()

	for i, sample := range wavParser.NewWindowIter(1024, 512) {
		fmt.Println(i, len(sample))
	}
}
