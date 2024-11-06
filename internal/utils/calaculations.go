package utils

import "math"

type RFParameters struct {
	FrequencyGHz    float64
	PathLengthKm    float64
	AtmosphericLoss float64
	Rainloss        float64
}

// CalculatePathLoss calculates the path loss based on the RF parameters.
func CalculatePathLoss(params RFParameters) float64 {
	// Free space path loss calculation : 20 * log10(f) + 20 * log10(d) + 92.45
	return 20*math.Log10(params.PathLengthKm) + 20*math.Log10(params.FrequencyGHz*1000) + 92.45
}

// CalculateTotalLoss calculates the total link loss based on the RF parameters.
func CalculateTotalLoss(params RFParameters) float64 {
	fspl := CalculatePathLoss(params)
	return fspl + params.AtmosphericLoss + params.Rainloss
}

// CalculateUtilisation returns the percentage of the time slot used by the burst. (Simple calculation for demonstration purposes)
func CalculateUtilisation(dataSize int) float64 {
	maxSize := 1024 // Assume 1024 bytes is the maximum data size
	if dataSize > maxSize {
		return 100.0
	}
	return float64(dataSize) / float64(maxSize) * 100.0 // Calculate the percentage of the time slot used by the burst
}
