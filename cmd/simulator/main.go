package main

import (
	"log"
	"time"

	"github.com/afroash/rf-simulator/internal/modulation"
	"github.com/afroash/rf-simulator/internal/tdma"
)

func main() {
	log.Println("Starting simulator...")

	// Create a new TDMA frame (2ms frame duration, 50µs guard time, 4 carriers)
	frame, err := tdma.NewTDMAFrame(2*time.Millisecond, 50*time.Microsecond, 4)
	if err != nil {
		log.Fatal(err)
	}

	// Calculate data sizes based on modulation efficiency
	// For 450µs timeslot with 25 MSps:
	// Base data size for BPSK (will be multiplied by modulation efficiency)
	baseDataSize := 1400 // bytes

	bursts := []*tdma.Burst{
		// BPSK: 1 bit per symbol
		tdma.NewBurst(make([]byte, baseDataSize*1), 0, tdma.DataBurst, modulation.BPSK),

		// QPSK: 2 bits per symbol
		tdma.NewBurst(make([]byte, baseDataSize*2), 1, tdma.DataBurst, modulation.QPSK),

		// 16-QAM: 4 bits per symbol
		tdma.NewBurst(make([]byte, baseDataSize*4), 2, tdma.DataBurst, modulation.QAM16),

		// 64-QAM: 6 bits per symbol
		tdma.NewBurst(make([]byte, baseDataSize*6), 3, tdma.DataBurst, modulation.QAM64),
	}

	// Add debug information
	for _, burst := range bursts {
		log.Printf("Debug - Carrier %d: %s, DataRate: %.2f Mbps, Utilization: %.2f%%, Data Size: %d bytes",
			burst.CarrierID,
			burst.Modulation.Name,
			burst.Datarate/1e6, // Convert to Mbps
			burst.Utilisation,
			len(burst.Data),
		)

		err := frame.AddBurst(burst.CarrierID, burst)
		if err != nil {
			log.Printf("Error adding burst to carrier %d: %v", burst.CarrierID, err)
		}
	}

	frame.PrintDetailedFrameStructure()
}
