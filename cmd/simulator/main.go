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

	// Add different types of bursts
	bursts := []*tdma.Burst{
		tdma.NewBurst(make([]byte, 1024), 0, tdma.DataBurst, modulation.BPSK),  // Basic modulation
		tdma.NewBurst(make([]byte, 1024), 1, tdma.DataBurst, modulation.QPSK),  // Double efficiency
		tdma.NewBurst(make([]byte, 1024), 2, tdma.DataBurst, modulation.QAM16), // Quadruple efficiency
		tdma.NewBurst(make([]byte, 1024), 3, tdma.DataBurst, modulation.QAM64), // 6x efficiency
	}
	for _, burst := range bursts {
		log.Printf("Debug - Carrier %d: %s, DataRate: %.2f Mbps, Utilization: %.2f%%",
			burst.CarrierID,
			burst.Modulation.Name,
			burst.Datarate/1000000,
			burst.Utilisation,
		)
		err := frame.AddBurst(burst.CarrierID, burst)
		if err != nil {
			log.Printf("Error adding burst to carrier %d: %v", burst.CarrierID, err)
		}
	}

	frame.PrintDetailedFrameStructure()
}
