package main

import (
	"fmt"
	"log"
	"time"

	"github.com/afroash/rf-simulator/internal/tdma"
)

func main() {
	log.Println("Starting Dynamic TDMA Simulation...")

	// Create frame configuration with variable slot durations
	config := tdma.FrameConfig{
		FrameDuration: 2 * time.Millisecond,
		GuardTime:     50 * time.Microsecond,
		NumCarriers:   4,
		SlotDurations: []time.Duration{
			600 * time.Microsecond, // More time for high-priority carrier
			400 * time.Microsecond,
			300 * time.Microsecond,
			200 * time.Microsecond,
		},
	}

	frame, err := tdma.NewFrame(config)
	if err != nil {
		log.Fatal(err)
	}

	// Simulate multiple frames
	numFramesToSimulate := 5
	fmt.Printf("\nSimulating %d frames with dynamic SNR...\n", numFramesToSimulate)

	for frameNum := 0; frameNum < numFramesToSimulate; frameNum++ {
		fmt.Printf("\nFrame %d:\n", frameNum)
		fmt.Printf("════════════════════════════════════════════\n")

		// Get updated SNR values for this frame
		snrValues := frame.UpdateSNR()

		// Create bursts with current SNR values
		for carrierID := 0; carrierID < config.NumCarriers; carrierID++ {
			// Calculate data size based on slot duration
			slotDuration := config.SlotDurations[carrierID]
			baseDataSize := int(slotDuration.Microseconds() * 2400 / 450) // Scale data size with duration

			// Get current SNR for this carrier
			snr := snrValues[carrierID]
			if snr == 0 { // If no update this frame, skip
				snr = frame.SNRProfiles[carrierID].BaselineSNR
			}

			burst := tdma.NewBurstWithSNR(
				make([]byte, baseDataSize),
				carrierID,
				tdma.DataBurst,
				snr,
			)

			fmt.Printf("\nCarrier %d:\n", carrierID)
			fmt.Printf("  Slot Duration: %v\n", slotDuration)
			fmt.Printf("  SNR: %.1f dB\n", snr)
			fmt.Printf("  → Modulation: %s\n", burst.Modulation.Name)
			fmt.Printf("  → Data Rate: %.2f Mbps\n", burst.Datarate/1e6)
			fmt.Printf("  → Efficiency: %.1f%%\n", burst.Utilisation)

			err := frame.AddBurst(burst.CarrierID, burst)
			if err != nil {
				log.Printf("Error adding burst to carrier %d: %v", carrierID, err)
			}
		}

		frame.PrintDetailedFrameStructure()
		frame.AdvanceFrame()

		// Small delay to make output readable
		time.Sleep(time.Second)
	}
}
