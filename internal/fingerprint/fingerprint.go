package fingerprint

type VideoTimestamp struct {
	SongId     string
	AnchorTime uint32
}

const targetZoneSize = 8

func GenerateFingerprints(songId string, peaks []Peak) map[uint32][]VideoTimestamp {
	fingerprints := make(map[uint32][]VideoTimestamp)

	for anchorInd, anchor := range peaks {
		for i := anchorInd + 1; i <= anchorInd+targetZoneSize && i < len(peaks); i++ {
			hash := GenerateHash(anchor, peaks[i])
			_, ok := fingerprints[hash]
			if ok {
				fingerprints[hash] = append(fingerprints[hash], VideoTimestamp{
					SongId:     songId,
					AnchorTime: uint32(anchor.Time * 1000),
				})
			} else {
				t := make([]VideoTimestamp, 0)
				t = append(t, VideoTimestamp{
					SongId:     songId,
					AnchorTime: uint32(anchor.Time * 1000),
				})
				fingerprints[hash] = t
			}
		}
	}

	return fingerprints
}

func GenerateHash(anchor Peak, target Peak) uint32 {
	deltaTime := uint32(target.Time*1000) - uint32(anchor.Time*1000)
	return uint32(anchor.Freq)<<23 | uint32(target.Freq) | deltaTime
}
