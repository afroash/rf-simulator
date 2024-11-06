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
	Type          ModulationType
	BitsPerSymbol float64
	Name          string
	ReqSNR        float64
	MaxErrorRate  float64
}

// GetModulationScheme returns the modulation scheme for the specified type.
func GetModulationScheme(mt ModulationType) *ModulationScheme {
	switch mt {
	case BPSK:
		return &ModulationScheme{
			Type:          BPSK,
			BitsPerSymbol: 1,
			Name:          "BPSK",
			ReqSNR:        8.4,
			MaxErrorRate:  1e-6,
		}
	case QPSK:
		return &ModulationScheme{
			Type:          QPSK,
			BitsPerSymbol: 2,
			Name:          "QPSK",
			ReqSNR:        11.5,
			MaxErrorRate:  1e-6,
		}
	case PSK8:
		return &ModulationScheme{
			Type:          PSK8,
			BitsPerSymbol: 3,
			Name:          "8-PSK",
			ReqSNR:        14.0,
			MaxErrorRate:  1e-6,
		}
	case QAM16:
		return &ModulationScheme{
			Type:          QAM16,
			BitsPerSymbol: 4,
			Name:          "16-QAM",
			ReqSNR:        15.1,
			MaxErrorRate:  1e-6,
		}
	case QAM64:
		return &ModulationScheme{
			Type:          QAM64,
			BitsPerSymbol: 6,
			Name:          "64-QAM",
			ReqSNR:        19.5,
			MaxErrorRate:  1e-6,
		}
	default:
		return nil
	}
}

// CalculateSpecEfficiency calculates the spectral efficiency of a modulation scheme. (bps/Hz)
func (ms *ModulationScheme) CalculateSpecEfficiency() float64 {
	return ms.BitsPerSymbol
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
		return (3.0 / 8.0) * math.Erfc(math.Sqrt((4.0/10.0)*snr))
	case QAM64:
		return (7.0 / 24.0) * math.Erfc(math.Sqrt((1.0/42.0)*snr))
	default:
		return 0
	}
}
