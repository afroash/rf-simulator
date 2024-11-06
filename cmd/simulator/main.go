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
		tdma.NewBurst(make([]byte, 1024), 0, tdma.DataBurst, modulation.QAM64),      // Full data burst
		tdma.NewBurst(make([]byte, 512), 1, tdma.DataBurst, modulation.QPSK),        // Half data burst
		tdma.NewBurst(make([]byte, 256), 2, tdma.ControlBurst, modulation.BPSK),     // Control burst
		tdma.NewBurst(make([]byte, 128), 3, tdma.MaintenanceBurst, modulation.QPSK), // Maintenance burst
	}
	for _, burst := range bursts {
		err := frame.AddBurst(burst.CarrierID, burst)
		if err != nil {
			log.Printf("Error adding burst to carrier %d: %v", burst.CarrierID, err)
		}
	}

	frame.PrintDetailedFrameStructure()
}
