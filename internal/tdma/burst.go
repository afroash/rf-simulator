package tdma

import (
	"errors"
	"fmt"
	"time"
)

type BurstType int

const (
	DataBurst BurstType = iota
	ControlBurst
	MaintenanceBurst
)

type Burst struct {
	Data      []byte
	StartTime time.Duration
	Duration  time.Duration
	CarrierID int
}

type TimeSlot struct {
	StartTime   time.Duration
	Duration    time.Duration
	Burst       *Burst
	IsGuardTime bool
}

// AddBurst adds a burst to a specific carrier's time slot
func (frame *TDMAFrame) AddBurst(carrierID int, data []byte) error {
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

	// Create and assign the burst
	slot.Burst = &Burst{
		Data:      data,
		StartTime: slot.StartTime,
		Duration:  slot.Duration,
		CarrierID: carrierID,
	}

	return nil
}
