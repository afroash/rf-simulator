package tdma

import (
	"errors"
	"time"
)

// TDMAFrame represents a complete TDMA frame.
type TDMAFrame struct {
	Duration    time.Duration
	GuardTime   time.Duration
	TimeSlots   []*TimeSlot
	NumCarriers int
}

// NewTDMAFrame creates a new TDMA frame with the specified parameters
func NewTDMAFrame(frameDuration time.Duration, guardTime time.Duration, numCarriers int) (*TDMAFrame, error) {
	if frameDuration <= 0 || guardTime <= 0 || numCarriers <= 0 {
		return nil, errors.New("invalid parameters: all values must be positive")
	}

	frame := &TDMAFrame{
		Duration:    frameDuration,
		GuardTime:   guardTime,
		NumCarriers: numCarriers,
	}

	// Calculate slot duration (excluding guard time)
	slotDuration := (frameDuration - time.Duration(numCarriers)*guardTime) / time.Duration(numCarriers)

	// Create time slots with guard times
	currentTime := time.Duration(0)
	for i := 0; i < numCarriers; i++ {
		// Add main time slot
		frame.TimeSlots = append(frame.TimeSlots, &TimeSlot{
			StartTime:   currentTime,
			Duration:    slotDuration,
			IsGuardTime: false,
		})
		currentTime += slotDuration

		// Add guard time slot
		frame.TimeSlots = append(frame.TimeSlots, &TimeSlot{
			StartTime:   currentTime,
			Duration:    guardTime,
			IsGuardTime: true,
		})
		currentTime += guardTime
	}

	return frame, nil
}
