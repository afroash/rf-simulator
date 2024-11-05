package tdma

import (
	"errors"
	"fmt"
	"time"
)

// TDMAFrame represents a complete TDMA frame.
type TDMAFrame struct {
	Duration    time.Duration
	GuardTime   time.Duration
	TimeSlots   []*TimeSlot
	NumCarriers int
}

// NewTDMAFrame creates a new TDMA frame with the specified duration and guard time.
func NewTDMAFrame(duration, guardTime time.Duration, numCarriers int) (*TDMAFrame, error) {
	if duration <= 0 || guardTime <= 0 || numCarriers <= 0 {
		return nil, errors.New("invalid parameters: All parameters must be greater than 0")
	}
	frame := &TDMAFrame{
		Duration:    duration,
		GuardTime:   guardTime,
		NumCarriers: numCarriers,
	}

	slotDuration := (duration - time.Duration(numCarriers)*guardTime) / time.Duration(numCarriers)

	currentTime := time.Duration(0)
	for i := 0; i < numCarriers; i++ {
		// Add main time slot
		frame.TimeSlots = append(frame.TimeSlots, &TimeSlot{
			StartTime:   currentTime,
			Duration:    slotDuration,
			IsGuardTime: false,
		})
		currentTime += slotDuration

		// Add guard time
		frame.TimeSlots = append(frame.TimeSlots, &TimeSlot{
			StartTime:   currentTime,
			Duration:    guardTime,
			IsGuardTime: true,
		})
		currentTime += guardTime
	}
	return frame, nil
}

// PrintFrame prints the TDMA frame to the console.
func (f *TDMAFrame) PrintFrame() {
	fmt.Printf("TDMA Frame Structure (Duration: %v)\n", f.Duration)
	fmt.Println("=====================================")

	for i, slot := range f.TimeSlots {
		if slot.IsGuardTime {
			fmt.Printf("Guard Time: %v|\n", slot.Duration)
		} else {
			burstInfo := "Empty"
			if slot.Burst != nil {
				burstInfo = fmt.Sprintf("Carrier: %d", slot.Burst.CarrierID)
			}
			fmt.Printf("Slot %d: %v, - %s\n", i/2, slot.Duration, burstInfo)
		}
	}
	fmt.Println("=====================================")
}
