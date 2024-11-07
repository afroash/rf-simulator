package main

import (
	"fmt"
	"log"
	"time"

	"github.com/afroash/rf-simulator/internal/tdma"
	"github.com/afroash/rf-simulator/internal/utils"
)

func main() {
	log.Println("Starting Enhanced TDMA Simulation...")

	// Initialize terminals with different data requirements
	terminals := []*tdma.Terminal{
		{
			ID:            0,
			TotalSize:     1000000, // 1MB
			RemainingData: 1000000,
			BaselineSNR:   18.0, // Good SNR - can use QAM64
			Priority:      1,
		},
		{
			ID:            1,
			TotalSize:     500000, // 500KB
			RemainingData: 500000,
			BaselineSNR:   14.0, // Moderate SNR - will use QPSK/8PSK
			Priority:      2,
		},
		{
			ID:            2,
			TotalSize:     750000, // 750KB
			RemainingData: 750000,
			BaselineSNR:   12.0, // Lower SNR - will use QPSK
			Priority:      3,
		},
		{
			ID:            3,
			TotalSize:     250000, // 250KB
			RemainingData: 250000,
			BaselineSNR:   9.0, // Poor SNR - will use BPSK
			Priority:      4,
		},
	}

	config := tdma.FrameConfig{
		FrameDuration: 2 * time.Millisecond,
		GuardTime:     50 * time.Microsecond,
		NumCarriers:   4,
		SlotDurations: []time.Duration{
			600 * time.Microsecond,
			400 * time.Microsecond,
			300 * time.Microsecond,
			200 * time.Microsecond,
		},
	}

	frame, err := tdma.NewFrame(config)
	if err != nil {
		log.Fatal(err)
	}

	// Add terminals to frame
	frame.Terminals = terminals

	// Simulate until all data is transmitted
	frameNum := 0
	allDataTransmitted := false

	for !allDataTransmitted {
		// Collect status for all terminals
		statuses := make([]tdma.TerminalStatus, len(frame.Terminals))
		allDataTransmitted = true

		snrValues := frame.UpdateSNR()

		for idx, terminal := range frame.Terminals {
			// Initialize status with default values
			status := tdma.TerminalStatus{
				ID:            terminal.ID,
				SNR:           terminal.BaselineSNR,
				TotalData:     terminal.TotalSize,
				RemainingData: terminal.RemainingData,
			}

			if terminal.RemainingData > 0 {
				allDataTransmitted = false

				// Get current SNR
				snr := snrValues[terminal.ID]
				if snr == 0 {
					snr = terminal.BaselineSNR
				}
				status.SNR = snr

				// Create burst and calculate transmission
				slotDuration := config.SlotDurations[terminal.ID]
				burst := tdma.NewBurstWithSNR(
					make([]byte, utils.CalculateBurstSize(terminal, slotDuration, snr)),
					terminal.ID,
					tdma.DataBurst,
					snr,
				)

				// Calculate data transmitted
				dataTransmitted := int64(float64(len(burst.Data)) * burst.Utilisation / 100)
				terminal.RemainingData -= dataTransmitted

				if terminal.RemainingData < 0 {
					terminal.RemainingData = 0
				}

				// Update status
				status.ModScheme = burst.Modulation.Name
				status.DataRate = burst.Datarate / 1e6
				status.DataThisFrame = dataTransmitted
				status.SlotUtilization = burst.Utilisation
				status.RemainingData = terminal.RemainingData

				// Add burst to frame
				frame.AddBurst(burst.CarrierID, burst)
			} else {
				// For completed terminals, set final status
				status.ModScheme = "-"
				status.DataRate = 0
				status.DataThisFrame = 0
				status.SlotUtilization = 0
				status.RemainingData = 0
			}

			statuses[idx] = status
		}

		// Display frame summary
		tdma.PrintFrameSummary(frameNum, statuses, config.FrameDuration)

		// Log detailed information
		utils.LogFrameDetails(frameNum, statuses)

		if allDataTransmitted {
			break
		}

		frame.AdvanceFrame()
		frameNum++

		// Small delay for visualization
		time.Sleep(200 * time.Millisecond)
	}

	// Print final summary
	fmt.Printf("\nSimulation Complete\n")
	fmt.Printf("Total Frames: %d\n", frameNum+1)
	fmt.Printf("Total Time: %v\n", time.Duration(frameNum+1)*config.FrameDuration)
	fmt.Printf("\nDetailed logs available in: tdma_simulation_%s.log\n",
		time.Now().Format("2006-01-02"))
}
