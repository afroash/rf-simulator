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

func createProgressBar(percent float64, width int) string {
	// Ensure percent is between 0 and 100
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}

	filled := int(percent * float64(width) / 100)
	// Ensure filled is between 0 and width
	if filled < 0 {
		filled = 0
	}
	if filled > width {
		filled = width
	}

	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	return fmt.Sprintf("[%s] %5.1f%%", bar, percent)
}

func PrintFrameSummary(frameNum int, statuses []TerminalStatus, frameDuration time.Duration) {
	// Clear screen (ANSI escape sequence)
	fmt.Print("\033[H\033[2J")

	width := 64 // Total width of the box

	// Print frame header
	fmt.Printf("╔%s╗\n", strings.Repeat("═", width-2))
	fmt.Printf("║ TDMA Frame #%-4d              Duration: %-12v ║\n",
		frameNum, frameDuration)
	fmt.Printf("╠%s╣\n", strings.Repeat("═", width-2))

	// Print column headers with fixed widths
	fmt.Printf("║ Term │ SNR  │ Mod   │ Rate    │ Progress                    ║\n")
	fmt.Printf("╟──────┼──────┼───────┼─────────┼────────────────────────────╢\n")

	// Print each terminal's status with proper spacing
	for _, status := range statuses {
		// Calculate progress
		var progress float64
		if status.TotalData > 0 {
			progress = 100.0 * float64(status.TotalData-status.RemainingData) /
				float64(status.TotalData)
		}

		progressBar := createProgressBar(progress, 20)

		fmt.Printf("║ %4d │ %4.1f │ %-5s │ %5.1fMb │ %-28s ║\n",
			status.ID,
			status.SNR,
			status.ModScheme,
			status.DataRate,
			progressBar,
		)
	}

	fmt.Printf("╚%s╝\n", strings.Repeat("═", width-2))

	// Print detailed statistics with proper spacing
	fmt.Printf("\nDetailed Statistics:\n")
	fmt.Printf("══════════════════\n")

	activeTransmissions := false
	for _, status := range statuses {
		if status.RemainingData > 0 {
			activeTransmissions = true
			fmt.Printf("Terminal %d:\n", status.ID)
			fmt.Printf("  ├── Data Transferred: %d bytes\n",
				status.DataThisFrame)
			fmt.Printf("  ├── Remaining: %d bytes\n",
				status.RemainingData)
			fmt.Printf("  └── Slot Utilization: %.1f%%\n",
				status.SlotUtilization)
		}
	}

	if !activeTransmissions {
		fmt.Println("\nAll transmissions complete!")
	}
}
