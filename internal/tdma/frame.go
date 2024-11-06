package tdma

import (
	"errors"
	"fmt"
	"strings"
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

// PrintDetailedFrameStructure prints an enhanced visual representation of the frame
func (frame *TDMAFrame) PrintDetailedFrameStructure() {
	fmt.Printf("\nTDMA Frame Structure (Duration: %v)\n", frame.Duration)
	fmt.Println("═══════════════════════════════════════════════")

	// Frame header
	fmt.Println("Time →")
	fmt.Printf("%-10s %-15s %-20s %-15s\n", "Slot", "Duration", "Burst Info", "Utilization")
	fmt.Println("───────────────────────────────────────────────")

	for i, slot := range frame.TimeSlots {
		if slot.IsGuardTime {
			printGuardTimeVisual()
		} else {
			printSlotVisual(i/2, slot)
		}
	}

	fmt.Println("═══════════════════════════════════════════════")
	printLegend()
}

func printSlotVisual(slotNum int, slot *TimeSlot) {
	var burstInfo, utilizationBar, modInfo string

	if slot.Burst != nil {
		burstType := GetBurstsType(slot.Burst.Type)

		modInfo = fmt.Sprintf("%s, SNR: %.1fdB",
			slot.Burst.Modulation.Name,
			slot.Burst.SNR,
		)

		burstInfo = fmt.Sprintf("Carrier %d [%s] %s",
			slot.Burst.CarrierID,
			burstType,
			modInfo,
		)
		utilizationBar = createUtilizationBar(slot.Burst.Utilisation)
	} else {
		burstInfo = "Empty"
		utilizationBar = createUtilizationBar(0)
	}

	fmt.Printf("%-10d %-15v %-35s %s\n",
		slotNum,
		slot.Duration,
		burstInfo,
		utilizationBar,
	)
}
func printGuardTimeVisual() {
	fmt.Printf("%s\n", strings.Repeat("-", 50))
}

func createUtilizationBar(utilization float64) string {
	const barLength = 10
	filledSlots := int((utilization / 100.0) * float64(barLength))

	bar := strings.Builder{}
	bar.WriteString("[")
	bar.WriteString(strings.Repeat("█", filledSlots))
	bar.WriteString(strings.Repeat("░", barLength-filledSlots))
	bar.WriteString("]")

	return fmt.Sprintf("%s %.1f%%", bar.String(), utilization)
}

func printLegend() {
	fmt.Println("\nLegend:")
	fmt.Println("D - Data Burst")
	fmt.Println("C - Control Burst")
	fmt.Println("M - Maintenance Burst")
	fmt.Println("█ - Utilized capacity")
	fmt.Println("░ - Available capacity")
}
