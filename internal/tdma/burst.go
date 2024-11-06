package tdma

import (
	"errors"
	"fmt"
	"time"

	"github.com/afroash/rf-simulator/internal/modulation"
	"github.com/afroash/rf-simulator/internal/utils"
)

type BurstType int

const (
	DataBurst BurstType = iota
	ControlBurst
	MaintenanceBurst
)

type Burst struct {
	Data        []byte
	StartTime   time.Duration
	Duration    time.Duration
	CarrierID   int
	Type        BurstType
	Utilisation float64 // Percentage of the time slot used by the burst
	Modulation  modulation.ModulationScheme
	SNR         float64
	BER         float64
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

	burst := &Burst{
		Data:        data,
		CarrierID:   carrierID,
		Type:        burstType,
		Utilisation: utils.CalculateUtilisation(len(data)),
		Modulation:  *modScheme,
		SNR:         10.0, // Default SNR value
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
