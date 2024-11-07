package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/afroash/rf-simulator/internal/tdma"
)

// Creates a log file for detailed analysis
func LogFrameDetails(frameNum int, statuses []tdma.TerminalStatus) error {
	// Open log file in append mode
	f, err := os.OpenFile(fmt.Sprintf("tdma_simulation_%s.log",
		time.Now().Format("2006-01-02")),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write frame details
	fmt.Fprintf(f, "\nFrame %d\n", frameNum)
	fmt.Fprintf(f, "Time: %s\n", time.Now().Format("15:04:05.000"))

	for _, status := range statuses {
		fmt.Fprintf(f, "Terminal %d: SNR=%.1f, Mod=%s, Rate=%.1fMb, Sent=%d bytes\n",
			status.ID, status.SNR, status.ModScheme, status.DataRate, status.DataThisFrame)
	}

	return nil
}
