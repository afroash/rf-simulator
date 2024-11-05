package main

import (
	"fmt"
	"log"
	"time"

	"github.com/afroash/rf-simulator/internal/tdma"
)

func main() {
	log.Println("Starting simulator...")

	// Create a new TDMA frame (1ms duration, 50µs guard time, 4 carriers)
	frame, err := tdma.NewTDMAFrame(2*time.Millisecond, 50*time.Microsecond, 4)
	if err != nil {
		log.Fatalf("Error creating TDMA frame: %v", err)
	}

	//Add some test bursts
	frame.AddBurst(0, []byte("Carrier 0 data"))
	frame.AddBurst(1, []byte("Carrier 1 data"))

	// Print the frame structure
	frame.PrintFrame()

	//create a new channel configuration
	ch := &tdma.Channel{
		CenterFreq:    1550.0, //MHz
		Bandwidth:     36.0,   //MHz
		Modulation:    tdma.QPSK,
		SymbolRate:    10.0, //Msymbols/s
		BitsPerSymbol: 2,
	}
	fmt.Printf("Channel Capacity: %.2f Mbps\n", ch.CalculateChannelCapacity())
}
