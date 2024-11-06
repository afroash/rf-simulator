package tdma

import (
	"fmt"
	"strings"
)

// PrintDetailedFrameStructure prints an enhanced visual representation of the frame
func (frame *TDMAFrame) PrintDetailedFrameStructure() {
	fmt.Printf("\nTDMA Frame Structure (Duration: %v)\n", frame.Duration)
	fmt.Println("═══════════════════════════════════════════════════════════════════════════")

	fmt.Println("Time →")
	fmt.Printf("%-10s %-15s %-45s %-20s\n",
		"Slot",
		"Duration",
		"Burst Info [Type] Modulation (Data Rate)",
		"Utilization",
	)
	fmt.Println("───────────────────────────────────────────────────────────────────────────")

	for i, slot := range frame.TimeSlots {
		if slot.IsGuardTime {
			printGuardTimeVisual()
		} else {
			printSlotVisual(i/2, slot)
		}
	}

	fmt.Println("═══════════════════════════════════════════════════════════════════════════")
	printLegend()
	printModulationEfficiency()
}

func printModulationEfficiency() {
	fmt.Println("\nModulation Efficiency:")
	fmt.Println("BPSK:  1 bit/symbol  - Most robust, lowest data rate")
	fmt.Println("QPSK:  2 bits/symbol - Good balance of robustness and speed")
	fmt.Println("8PSK:  3 bits/symbol - Higher speed, needs better SNR")
	fmt.Println("16QAM: 4 bits/symbol - High speed, requires good SNR")
	fmt.Println("64QAM: 6 bits/symbol - Highest speed, requires excellent SNR")
}

func printSlotVisual(slotNum int, slot *TimeSlot) {
	var burstInfo, utilizationBar, modInfo, dataRateInfo string

	if slot.Burst != nil {
		burstType := GetBurstsType(slot.Burst.Type)
		modInfo = slot.Burst.Modulation.Name
		dataRateInfo = fmt.Sprintf("%.2f Mbps", slot.Burst.Datarate/1000000)

		burstInfo = fmt.Sprintf("Carrier %d [%s] %s (%s)",
			slot.Burst.CarrierID,
			burstType,
			modInfo,
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
