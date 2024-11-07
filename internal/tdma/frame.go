package tdma

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"time"
)

// NewFrame creates a new TDMA frame with variable slot durations
func NewFrame(config FrameConfig) (*Frame, error) {
	if len(config.SlotDurations) > 0 && len(config.SlotDurations) != config.NumCarriers {
		return nil, fmt.Errorf("number of slot durations (%d) must match number of carriers (%d)",
			len(config.SlotDurations), config.NumCarriers)
	}

	frame := &Frame{
		Config:      config,
		FrameNumber: 0,
		SNRProfiles: make(map[int]*SNRProfile),
	}

	// Initialize SNR profiles for each carrier
	for i := 0; i < config.NumCarriers; i++ {
		frame.SNRProfiles[i] = &SNRProfile{
			BaselineSNR:     15.0, // Default baseline SNR
			Variation:       3.0,  // ±3dB variation
			UpdateInterval:  5,    // Update every 5 frames
			LastUpdateFrame: 0,
		}
	}

	err := frame.createTimeSlots()
	if err != nil {
		return nil, err
	}

	return frame, nil
}

func (f *Frame) createTimeSlots() error {
	var currentTime time.Duration
	defaultSlotDuration := (f.Config.FrameDuration -
		time.Duration(f.Config.NumCarriers)*f.Config.GuardTime) /
		time.Duration(f.Config.NumCarriers)

	f.TimeSlots = nil // Clear existing slots

	for i := 0; i < f.Config.NumCarriers; i++ {
		// Get slot duration (either from config or default)
		slotDuration := defaultSlotDuration
		if len(f.Config.SlotDurations) > i {
			slotDuration = f.Config.SlotDurations[i]
		}

		// Add main time slot
		f.TimeSlots = append(f.TimeSlots, &TimeSlot{
			StartTime:   currentTime,
			Duration:    slotDuration,
			IsGuardTime: false,
		})
		currentTime += slotDuration

		// Add guard time slot
		f.TimeSlots = append(f.TimeSlots, &TimeSlot{
			StartTime:   currentTime,
			Duration:    f.Config.GuardTime,
			IsGuardTime: true,
		})
		currentTime += f.Config.GuardTime
	}

	return nil
}

// UpdateSNR updates SNR values for each carrier based on their profiles
func (f *Frame) UpdateSNR() map[int]float64 {
	currentSNRs := make(map[int]float64)

	for carrierID, profile := range f.SNRProfiles {
		if f.FrameNumber-profile.LastUpdateFrame >= profile.UpdateInterval {
			// Generate random variation within defined range
			variation := (rand.Float64()*2 - 1) * profile.Variation
			currentSNRs[carrierID] = profile.BaselineSNR + variation
			profile.LastUpdateFrame = f.FrameNumber
		}
	}

	return currentSNRs
}

// AdvanceFrame moves to the next frame and updates SNR values
func (f *Frame) AdvanceFrame() {
	f.FrameNumber++
}

func (f *Frame) AddBurst(carrierID int, burst *Burst) error {
	if carrierID >= f.Config.NumCarriers {
		return fmt.Errorf("carrier ID %d exceeds number of carriers", carrierID)
	}

	// Find the corresponding time slot for this carrier
	slotIndex := carrierID * 2 // Multiply by 2 because we have guard slots
	if slotIndex >= len(f.TimeSlots) {
		return errors.New("slot index out of range")
	}

	slot := f.TimeSlots[slotIndex]
	if slot.Burst != nil {
		return fmt.Errorf("slot for carrier %d already contains a burst", carrierID)
	}

	slot.Burst = burst
	return nil
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
