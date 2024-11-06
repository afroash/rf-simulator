package main

import (
	"fmt"
	"log"
	"time"

	"github.com/afroash/rf-simulator/internal/tdma"
)

func main() {
	log.Println("Starting simulator...")
	log.Println("SNR to Data Rate Relationship Demo...")

	frame, err := tdma.NewTDMAFrame(2*time.Millisecond, 50*time.Microsecond, 4)
	if err != nil {
		log.Fatal(err)
	}

	// Test cases showing SNR → Modulation → Data Rate relationship
	testCases := []struct {
		snr      float64
		dataSize int
		expected string
	}{
		{21.0, 1400, "Excellent - Expect 64-QAM"},
		{16.0, 1400, "Good - Expect 16-QAM"},
		{13.0, 1400, "Moderate - Expect QPSK"},
		{8.0, 1400, "Poor - Expect BPSK"},
	}

	fmt.Println("\nSNR to Modulation to Data Rate Relationship:")
	fmt.Println("════════════════════════════════════════════")

	for i, tc := range testCases {
		burst := tdma.NewBurstWithSNR(
			make([]byte, tc.dataSize),
			i,
			tdma.DataBurst,
			tc.snr,
		)

		fmt.Printf("\nCarrier %d:\n", i)
		fmt.Printf("  SNR: %.1f dB (%s)\n", tc.snr, tc.expected)
		fmt.Printf("  → Selected Modulation: %s (%g bits/symbol)\n",
			burst.Modulation.Name,
			burst.Modulation.BitsPerSymbol)
		fmt.Printf("  → Achieved Data Rate: %.2f Mbps\n",
			burst.Datarate/1e6)
		fmt.Printf("  → Bit Error Rate: %.2e\n",
			burst.BER)

		err := frame.AddBurst(burst.CarrierID, burst)
		if err != nil {
			log.Printf("Error adding burst to carrier %d: %v", burst.CarrierID, err)
		}
	}

	fmt.Println("\nFull TDMA Frame Structure:")
	fmt.Println("════════════════════════════════════════════")
	frame.PrintDetailedFrameStructure()
}
