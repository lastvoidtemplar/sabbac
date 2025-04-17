package spectrogram

import (
	"log/slog"
	"math"
	"math/cmplx"
)

const (
	windowSize int = 1024
	step       int = 512
)

// Spectogram slice of frequencies for window
func STFT(wavePath string, logger *slog.Logger) [][]complex128 {
	wavParser, err := newWavParser(wavePath, logger)
	if err != nil {
		logger.Error("Couldn`t create wav parser")
		return nil
	}
	defer wavParser.Close()

	numWindows := wavParser.WindowCount(windowSize, step)
	stftRes := make([][]complex128, numWindows)
	windowFunction := hammingWindow(windowSize)

	for i, sample := range wavParser.NewWindowIter(windowSize, step) {
		applyWindowFunction(sample, windowFunction)
		stftRes[i] = fft(sample, 0, len(sample)-1, 1)[:windowSize/2]
	}

	return stftRes
}

// https://en.wikipedia.org/wiki/Window_function#Hann_and_Hamming_windows
func hammingWindow(windowSize int) []float64 {
	const a0 float64 = 25.0 / 46
	const a1 = 1 - a0
	const piTimes2 float64 = 2 * math.Pi
	var windowSizeFloat64 float64 = float64(windowSize)

	window := make([]float64, windowSize)
	for i := 0; i < windowSize; i++ {
		window[i] = a0 - a1*math.Cos(piTimes2*float64(i)/windowSizeFloat64)
	}

	return window
}

func applyWindowFunction(sample []float64, windowFunction []float64) {
	for i, v := range windowFunction {
		sample[i] *= v
	}
}

func dft(sample []float64) []float64 {
	frequencies := make([]float64, 1+len(sample)/2)

	for fr := 0; fr < len(frequencies); fr++ {
		w := cmplx.Exp(complex(0, -2*float64(math.Pi)*float64(fr)/float64(len(sample))))
		e := complex(1.0, 0)
		sum := complex(0, 0)
		for j := 0; j < len(sample); j++ {
			sum += complex(sample[j], 0) * e
			e *= w
		}
		frequencies[fr] = cmplx.Abs(sum)
	}

	return frequencies
}

func fft(sample []float64, st int, end int, step int) []complex128 {
	n := (end-st)/(step) + 1

	if n == 1 {
		return []complex128{complex(sample[st], 0)}
	}

	frequencies := make([]complex128, n)

	frequenciesEven := fft(sample, st, end-step, 2*step)
	frequenciesOdd := fft(sample, st+step, end, 2*step)

	w := cmplx.Exp(complex(0, -2*float64(math.Pi)/float64(n)))
	e := complex(1, 0)
	for i := 0; i < n/2; i++ {
		frequencies[i] = frequenciesEven[i] + e*frequenciesOdd[i]
		frequencies[i+n/2] = frequenciesEven[i] - e*frequenciesOdd[i]
		e *= w
	}

	return frequencies
}
