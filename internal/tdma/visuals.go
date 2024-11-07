// internal/tdma/visualization.go
package tdma

import (
	"fmt"
	"strings"
	"time"
)

// PrintDetailedFrameStructure prints a detailed view of the frame
func (f *Frame) PrintDetailedFrameStructure() {
	fmt.Printf("\nTDMA Frame Structure (Frame #%d, Duration: %v)\n",
		f.FrameNumber,
		f.Config.FrameDuration,
	)
	fmt.Println("═══════════════════════════════════════════════════════════════════════════")

	fmt.Println("Time →")
	fmt.Printf("%-10s %-15s %-45s %-20s\n",
		"Slot",
		"Duration",
		"Burst Info [Type] Modulation (Data Rate)",
		"Utilization",
	)
	fmt.Println("───────────────────────────────────────────────────────────────────────────")

	slotNum := 0
	for _, slot := range f.TimeSlots {
		if slot.IsGuardTime {
			printGuardTimeVisual()
		} else {
			printSlotVisual(slotNum, slot)
			slotNum++
		}
	}

	fmt.Println("═══════════════════════════════════════════════════════════════════════════")
	printLegend()
}

func printSlotVisual(slotNum int, slot *TimeSlot) {
	var burstInfo, utilizationBar string

	if slot.Burst != nil {
		burstType := getBurstTypeSymbol(slot.Burst.Type)
		dataRateInfo := fmt.Sprintf("%.2f Mbps", slot.Burst.Datarate/1e6)
		snrInfo := fmt.Sprintf("SNR: %.1fdB", slot.Burst.SNR)

		burstInfo = fmt.Sprintf("Carrier %d [%s] %s %s (%s)",
			slot.Burst.CarrierID,
			burstType,
			slot.Burst.Modulation.Name,
			snrInfo,
			dataRateInfo,
		)
		utilizationBar = createUtilizationBar(slot.Burst.Utilisation)
	} else {
		burstInfo = "Empty"
		utilizationBar = createUtilizationBar(0)
	}

	fmt.Printf("%-10d %-15v %-45s %s\n",
		slotNum,
		slot.Duration,
		burstInfo,
		utilizationBar,
	)
}

func printGuardTimeVisual() {
	fmt.Printf("%s\n", strings.Repeat("-", 80))
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

func getBurstTypeSymbol(burstType BurstType) string {
	switch burstType {
	case DataBurst:
		return "D"
	case ControlBurst:
		return "C"
	case MaintenanceBurst:
		return "M"
	default:
		return "?"
	}
}

func printLegend() {
	fmt.Println("\nLegend:")
	fmt.Println("D - Data Burst")
	fmt.Println("C - Control Burst")
	fmt.Println("M - Maintenance Burst")
	fmt.Println("█ - Utilized capacity")
	fmt.Println("░ - Available capacity")
	fmt.Printf("Guard Time: %v between slots\n", time.Duration(50*time.Microsecond))
}
