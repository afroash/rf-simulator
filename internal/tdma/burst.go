package tdma

import (
	"errors"
	"fmt"
	"time"

	"github.com/afroash/rf-simulator/internal/modulation"
)

type BurstType int

const (
	DataBurst BurstType = iota
	ControlBurst
	MaintenanceBurst
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

// NewBurst creates a new burst with the specified data.
func NewBurst(data []byte, carrierID int, burstType BurstType, mod modulation.ModulationType) *Burst {
	modScheme := modulation.GetModulationScheme(mod)

	//Default  duration for calculation purposes.
	defaultDuration := 450 * time.Millisecond

	utilisation, dataRate, symbolsPacked := calculateModulationBasedUtilization(len(data), *modScheme, defaultDuration)

	burst := &Burst{
		Data:          data,
		CarrierID:     carrierID,
		Type:          burstType,
		Utilisation:   utilisation,
		Modulation:    *modScheme,
		SNR:           20.0, // Default SNR value
		Datarate:      dataRate,
		SymbolsPacked: symbolsPacked,
	}
	burst.BER = modScheme.CalculateBER(burst.SNR)
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
// taking into account the modulation scheme's efficiency
func calculateModulationBasedUtilization(dataSize int, mod modulation.ModulationScheme, duration time.Duration) (float64, float64, int) {
	// Convert duration to seconds
	durationSec := duration.Seconds()

	// Calculate symbol rate (assuming 1 Hz per symbol as base rate)
	baseSymbolRate := 1000000.0 // 1 MHz base symbol rate
	maxSymbols := int(baseSymbolRate * durationSec)

	// Calculate how many bits we can send with this modulation
	bitsPerSymbol := mod.BitsPerSymbol
	maxBits := int(float64(maxSymbols) * bitsPerSymbol)

	// Calculate actual bits we're trying to send
	actualBits := dataSize * 8

	// Calculate utilization percentage
	utilization := (float64(actualBits) / float64(maxBits)) * 100
	if utilization > 100 {
		utilization = 100
	}

	// Calculate effective data rate
	dataRate := float64(actualBits) / durationSec // bits per second

	return utilization, dataRate, maxSymbols
}
