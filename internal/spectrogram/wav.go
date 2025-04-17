package spectrogram

import (
	"encoding/binary"
	"iter"
	"log/slog"
	"os"
)

type wavHeader struct {
	sampleRate    uint32
	bitsPerSample uint16
	dataSize      uint32
}

type wavParser struct {
	logger *slog.Logger

	wavPath string
	wavFile *os.File

	wavHeader wavHeader
}

func newWavParser(wavPath string, logger *slog.Logger) (*wavParser, error) {
	wavFile, err := os.Open(wavPath)
	if err != nil {
		logger.With(slog.String("wav_path", wavPath), slog.String("err", err.Error())).Error("Couldn`t open wav file")
		return nil, err
	}

	parser := &wavParser{
		logger:  logger,
		wavPath: wavPath,
		wavFile: wavFile,
	}

	parser.parseHeader()

	return parser, nil
}

func (parser *wavParser) Close() error {
	parser.logger.With(slog.String("wav_path", parser.wavPath)).Info("The file is closed")
	return parser.wavFile.Close()
}

// http://soundfile.sapp.org/doc/WaveFormat/

func (parser *wavParser) parseHeader() bool {
	buf := make([]byte, 78)
	parser.wavFile.Read(buf)

	riff := string(buf[:4])
	if riff != "RIFF" {
		parser.logger.With(slog.String("wav_path", parser.wavPath), slog.String("riff", riff)).Error("The first 4 bytes must be 'RIFF'")
		return false
	}

	format := string(buf[8:12])
	if format != "WAVE" {
		parser.logger.With(slog.String("wav_path", parser.wavPath), slog.String("format", format)).Error("The bytes between 8-12 must be 'WAVE'")
		return false
	}

	subId1 := string(buf[12:16])
	if subId1 != "fmt " {
		parser.logger.With(slog.String("wav_path", parser.wavPath), slog.String("subId1", subId1)).Error("The bytes between 12-16 must be 'fmt '")
		return false
	}

	subSize1 := binary.LittleEndian.Uint16(buf[16:20])
	audioFormat := binary.LittleEndian.Uint16(buf[20:22])
	numChannels := binary.LittleEndian.Uint16(buf[22:24])
	if subSize1 != 16 || audioFormat != 1 || numChannels != 1 {
		parser.logger.With(
			slog.String("wav_path", parser.wavPath),
			slog.Uint64("sub_size_1", uint64(subSize1)),
			slog.Uint64("audio_format", uint64(audioFormat)),
			slog.Uint64("number_of_channels", uint64(numChannels)),
		).Error("The file isn`t PCM")
		return false
	}

	sampleRate := binary.LittleEndian.Uint32(buf[24:28])
	//byte rate(28:32) = sample_rate * num_channels * bits_per_sample/8
	//block align (32:34) = num_channels * bits_per_sample/8
	bitsPerSample := binary.LittleEndian.Uint16(buf[34:36])

	//LIST CHUNK between 36 and 70

	data := string(buf[70:74])
	if data != "data" {
		parser.logger.With(slog.String("wav_path", parser.wavPath), slog.String("data", data)).Error("The bytes between 70-74 must be 'data'")
		return false
	}

	dataSize := binary.LittleEndian.Uint32(buf[74:78])

	parser.wavHeader = wavHeader{
		sampleRate:    sampleRate,
		bitsPerSample: bitsPerSample,
		dataSize:      dataSize,
	}

	parser.logger.With(
		slog.String("wav_path", parser.wavPath),
		slog.Uint64("sample_rate", uint64(parser.wavHeader.sampleRate)),
		slog.Uint64("bit_per_sample", uint64(parser.wavHeader.bitsPerSample)),
		slog.Uint64("data_size", uint64(parser.wavHeader.dataSize)),
	).Info("The wav header is parse successfully")

	return true
}

func (parser *wavParser) WindowCount(windowSize int, step int) int {
	return 1 + (int(parser.wavHeader.dataSize)-2*windowSize)/(2*step)

}

func (parser *wavParser) NewWindowIter(windowSize int, step int) iter.Seq2[int, []float64] {
	if windowSize > int(parser.wavHeader.dataSize) {
		panic("window size can`t be bigger than the data size")
	}

	return func(yield func(int, []float64) bool) {
		numWindows := parser.WindowCount(windowSize, step)

		buf := make([]byte, 2*windowSize)
		window := make([]float64, windowSize)
		temp := make([]float64, step)
		pos := 0
		for i := 0; i < numWindows; i++ {
			_, err := parser.wavFile.Read(buf[2*pos:])
			if err != nil {
				parser.logger.With(slog.String("wav_path", parser.wavPath), slog.String("err", err.Error())).Error("Error while reading the wav file")
				panic("Error while reading the wav file")
			}

			for j := pos; j < windowSize; j++ {
				sample := int16(binary.LittleEndian.Uint16(buf[2*j : 2*j+2]))
				window[j] = float64(sample)
			}

			if !yield(i, window) {
				break
			}
			pos = copy(temp, window[windowSize-step:])
			pos = copy(window[:pos], temp)
		}
	}
}
