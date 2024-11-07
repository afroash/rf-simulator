package utils

import (
	"math"
	"time"

	"github.com/afroash/rf-simulator/internal/modulation"
	"github.com/afroash/rf-simulator/internal/tdma"
)

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

func CalculateBurstSize(terminal *tdma.Terminal, slotDuration time.Duration, snr float64) int {
	// Calculate maximum bytes that can be sent in this slot based on modulation
	mod := modulation.GetOptimalModulation(snr)
	scheme := modulation.GetModulationScheme(mod)

	// Calculate theoretical throughput for the slot duration
	bitsPerSymbol := scheme.BitsPerSymbol
	symbolRate := float64(slotDuration.Nanoseconds()) * 25.0 // 25 symbols per nanosecond
	maxBits := int64(symbolRate * bitsPerSymbol)
	maxBytes := maxBits / 8

	// Return either the maximum possible or remaining data, whichever is smaller
	if maxBytes > terminal.RemainingData {
		return int(terminal.RemainingData)
	}
	return int(maxBytes)
}
