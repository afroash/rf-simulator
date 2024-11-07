package tdma

import (
	"errors"
	"fmt"
	"time"

	"github.com/afroash/rf-simulator/internal/modulation"
)

// Constants for calculations
const (
	// Typical symbol rate for satellite communications (30 MHz bandwidth)
	BaseSymbolRate = 25e6 // 25 MSymbols/second
)

type Burst struct {
	Data          []byte
	StartTime     time.Duration
	Duration      time.Duration
	CarrierID     int
	Type          BurstType
	Utilisation   float64 // Percentage of the time slot used by the burst
	Modulation    modulation.ModulationScheme
	SNR           float64
	BER           float64
	Datarate      float64
	SymbolsPacked int
}

type TimeSlot struct {
	StartTime   time.Duration
	Duration    time.Duration
	Burst       *Burst
	IsGuardTime bool
}

// NewBurstWithSNR creates a burst with specific SNR value
func NewBurstWithSNR(data []byte, carrierID int, burstType BurstType, snr float64) *Burst {
	// Automatically select optimal modulation for given SNR
	optimalMod := modulation.GetOptimalModulation(snr)
	modScheme := modulation.GetModulationScheme(optimalMod)

	defaultDuration := 450 * time.Microsecond

	utilization, dataRate, symbolsPacked := calculateModulationBasedUtilization(
		len(data),
		*modScheme,
		defaultDuration,
		snr,
	)

	burst := &Burst{
		Data:          data,
		CarrierID:     carrierID,
		Type:          burstType,
		Modulation:    *modScheme,
		SNR:           snr,
		Utilisation:   utilization,
		Datarate:      dataRate,
		SymbolsPacked: symbolsPacked,
	}

	burst.BER = modScheme.CalculateBER(snr)

	return burst
}

// AddBurst adds a burst to a specific carrier's time slot
func (frame *TDMAFrame) AddBurst(carrierID int, burst *Burst) error {
	if carrierID >= frame.NumCarriers {
		return fmt.Errorf("carrier ID %d exceeds number of carriers", carrierID)
	}

	// Find the corresponding time slot for this carrier
	slotIndex := carrierID * 2 // Multiply by 2 because we have guard slots
	if slotIndex >= len(frame.TimeSlots) {
		return errors.New("slot index out of range")
	}

	slot := frame.TimeSlots[slotIndex]
	if slot.Burst != nil {
		return fmt.Errorf("slot for carrier %d already contains a burst", carrierID)
	}

	// Update burst timing information
	burst.StartTime = slot.StartTime
	burst.Duration = slot.Duration

	// Assign the burst to the slot
	slot.Burst = burst

	return nil
}

// getBurstsType returns the bursts of a specific type
func GetBurstsType(burstType BurstType) string {
	switch burstType {
	case DataBurst:
		return "D"
	case ControlBurst:
		return "C"
	case MaintenanceBurst:
		return "M"
	default:
		return "?" // Unknown burst type
	}
}

// calculateModulationBasedUtilization calculates how much of the slot is utilized
// taking into account the modulation scheme's efficiency and the SNR.
func calculateModulationBasedUtilization(dataSize int, mod modulation.ModulationScheme, duration time.Duration, snr float64) (float64, float64, int) {
	durationSec := duration.Seconds()

	// Calculate effective data rate considering SNR
	effectiveDataRate := mod.CalculateEffectiveDataRate(BaseSymbolRate, snr)

	// Maximum bits that could be transmitted in this time slot
	maxBitsInSlot := int(effectiveDataRate * durationSec)

	// Actual bits we're trying to send
	actualBits := dataSize * 8

	// Calculate utilization as percentage of maximum capacity
	utilization := (float64(actualBits) / float64(maxBitsInSlot)) * 100
	if utilization > 100 {
		utilization = 100
	}

	// Actual achieved data rate
	achievedDataRate := float64(actualBits) / durationSec
	if achievedDataRate > effectiveDataRate {
		achievedDataRate = effectiveDataRate
	}

	symbolsPacked := int(BaseSymbolRate * durationSec)

	return utilization, achievedDataRate, symbolsPacked
}
