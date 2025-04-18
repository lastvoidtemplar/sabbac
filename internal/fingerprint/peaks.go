package fingerprint

import (
	"math/cmplx"
)

type Peak struct {
	Freq uint16
	Time float64
}

var peakRanges = []struct {
	min int
	max int
}{{1, 10}, {10, 20}, {20, 40}, {40, 80}, {80, 160}, {160, 512}}

func FilterPeaks(spectrogram [][]complex128, timePerColumn float64) []Peak {
	peaks := make([]Peak, 0)

	maxMagPerRange := make([]float64, len(peakRanges))
	maxFreqPerRange := make([]int, len(peakRanges))
	for colmInd, colm := range spectrogram {
		for peakRangeInd, peakRange := range peakRanges {
			maxFreq := -1
			maxMag := -1.0

			for freq, c := range colm[peakRange.min:peakRange.max] {
				mag := cmplx.Abs(c)
				if mag > maxMag {
					maxFreq = peakRange.min + freq
					maxMag = mag
				}
			}

			maxFreqPerRange[peakRangeInd] = maxFreq
			maxMagPerRange[peakRangeInd] = maxMag
		}

		avgMaxMag := 0.0
		for _, mag := range maxMagPerRange {
			avgMaxMag += mag
		}
		avgMaxMag /= float64(len(maxMagPerRange))

		for peakRangeInd, mag := range maxMagPerRange {
			if mag > avgMaxMag {
				peaks = append(peaks, Peak{
					Freq: uint16(maxFreqPerRange[peakRangeInd]),
					Time: timePerColumn * float64(colmInd),
				})
			}
		}
	}

	return peaks
}
