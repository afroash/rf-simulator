package modulation

import (
	"math"
)

// ModulationType represents different modulation schemes.
type ModulationType int

const (
	BPSK ModulationType = iota
	QPSK
	PSK8
	QAM16
	QAM64
)

// ModulationScheme represents a modulation scheme with specific parameters.
type ModulationScheme struct {
	Type           ModulationType
	BitsPerSymbol  float64
	Name           string
	ReqSNR         float64
	MaxErrorRate   float64
	SpecEfficiency float64
}

// GetModulationScheme returns the modulation scheme for the specified type.
func GetModulationScheme(mt ModulationType) *ModulationScheme {
	switch mt {
	case BPSK:
		return &ModulationScheme{
			Type:           BPSK,
			BitsPerSymbol:  1,
			Name:           "BPSK",
			ReqSNR:         8.4,
			MaxErrorRate:   1e-6,
			SpecEfficiency: 1.0,
		}
	case QPSK:
		return &ModulationScheme{
			Type:           QPSK,
			BitsPerSymbol:  2,
			Name:           "QPSK",
			ReqSNR:         11.5,
			MaxErrorRate:   1e-6,
			SpecEfficiency: 2.0,
		}
	case PSK8:
		return &ModulationScheme{
			Type:           PSK8,
			BitsPerSymbol:  3,
			Name:           "8-PSK",
			ReqSNR:         14.0,
			MaxErrorRate:   1e-6,
			SpecEfficiency: 3.0,
		}
	case QAM16:
		return &ModulationScheme{
			Type:           QAM16,
			BitsPerSymbol:  4,
			Name:           "16-QAM",
			ReqSNR:         15.1,
			MaxErrorRate:   1e-6,
			SpecEfficiency: 4.0,
		}
	case QAM64:
		return &ModulationScheme{
			Type:           QAM64,
			BitsPerSymbol:  6,
			Name:           "64-QAM",
			ReqSNR:         19.5,
			MaxErrorRate:   1e-6,
			SpecEfficiency: 6.0,
		}
	default:
		return nil
	}
}

// CalculateSpecEfficiency calculates the spectral efficiency of a modulation scheme. (bps/Hz)
func (ms *ModulationScheme) CalculateSpecEfficiency() float64 {
	return ms.SpecEfficiency
}

// CalculateTheorecticalThroughput calculates the theoretical throughput of a modulation scheme. (bps)
func (ms *ModulationScheme) CalculateTheorecticalThroughput(bandwidthHz float64) float64 {
	return bandwidthHz * ms.CalculateSpecEfficiency()
}

// CalculateBER calculates the bit error rate for a modulation scheme given the SNR.
func (ms *ModulationScheme) CalculateBER(snrDb float64) float64 {
	snr := math.Pow(10, snrDb/10)

	switch ms.Type {
	case BPSK:
		return 0.5 * math.Erfc(math.Sqrt(snr))
	case QPSK:
		return math.Erfc(math.Sqrt(snr))
	case PSK8:
		return (2.0 / 3.0) * math.Erfc(math.Sqrt((3.0/2.0)*snr))
	case QAM16:
		return (3.0 / 4.0) * math.Erfc(math.Sqrt(snr/10.0))
	case QAM64:
		return (7.0 / 12.0) * math.Erfc(math.Sqrt(snr/42.0))
	default:
		return 1.0
	}
}

// GetOptimalModulation returns the optimal modulation scheme for the given SNR.
func GetOptimalModulation(snr float64) ModulationType {
	//SNR Thresholds for different modulation schemes
	// These are currently simplified values but are representive.
	switch {
	case snr >= 19.5:
		return QAM64 // Highest speed, requires excellent SNR.
	case snr >= 15.1:
		return QAM16 // High speed, requires good SNR.
	case snr >= 14.0:
		return PSK8 // Reliable with a moderate speed.
	case snr >= 11.5:
		return QPSK // Good balance of robustness and speed.
	default:
		return BPSK // Most robust, lowest data rate.
	}
}

// CalculateEffectiveDataRate calculates the effective data rate for a modulation scheme given the SNR.
func (ms ModulationScheme) CalculateEffectiveDataRate(snr float64, symbolRate float64) float64 {
	// Basic data rate (bits per second)
	baseDataRate := symbolRate * ms.BitsPerSymbol

	// Calculate SNR margin
	snrMargin := snr - ms.ReqSNR

	// Efficiency factor based on SNR Margin
	efficiency := 1.0
	if snrMargin > 0 {
		// More aggressive efficiency scaling
		efficiency = math.Min(2.0, 1.0+(snrMargin/10.0))
	} else if snrMargin < 0 {
		// Reduce efficiency when below required SNR
		efficiency = math.Max(0.5, 1.0+(snrMargin/20.0))
	}

	// Apply coding rate (assume 3/4 coding rate for error correction)
	codingRate := 0.75

	// Calculate final data rate
	effectiveRate := baseDataRate * efficiency * codingRate

	return effectiveRate
}
