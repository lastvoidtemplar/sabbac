package spectrogram

import (
	"image"
	"image/color"
	"image/jpeg"
	"math/cmplx"
	"os"
)

func SpectrogramToImage(path string, spectrogram [][]complex128) {
	img := image.NewRGBA(image.Rect(0, 0, len(spectrogram), len(spectrogram[0])))

	maxMag := -1.0
	for _, bin := range spectrogram {
		for _, c := range bin {
			mag := cmplx.Abs(c)
			if mag > maxMag {
				maxMag = mag
			}
		}
	}

	// maxMag - 0
	// 0 - 255

	// (mag - curr)/maxmag*255

	con := 255 / maxMag

	for binInd, bin := range spectrogram {
		maxFr := len(bin)
		for freqInd, c := range bin {
			grey := uint8(cmplx.Abs(c) * con)
			img.Set(binInd, maxFr-freqInd, color.RGBA{R: grey, G: grey, B: grey, A: 255})
		}
	}

	imgFile, _ := os.Create(path)
	jpeg.Encode(imgFile, img, &jpeg.Options{Quality: 90})
	imgFile.Close()

}
