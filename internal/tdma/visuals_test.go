// internal/tdma/visuals_test.go
package tdma

import (
	"strings"
	"testing"
	"time"
)

func TestCreateProgressBar(t *testing.T) {
	tests := []struct {
		name       string
		percent    float64
		width      int
		wantFilled int
		wantEmpty  int
	}{
		{
			name:       "Zero Percent",
			percent:    0,
			width:      10,
			wantFilled: 0,
			wantEmpty:  10,
		},
		{
			name:       "Full Progress",
			percent:    100,
			width:      10,
			wantFilled: 10,
			wantEmpty:  0,
		},
		{
			name:       "Half Progress",
			percent:    50,
			width:      10,
			wantFilled: 5,
			wantEmpty:  5,
		},
		{
			name:       "Negative Progress",
			percent:    -10,
			width:      10,
			wantFilled: 0,
			wantEmpty:  10,
		},
		{
			name:       "Over 100 Percent",
			percent:    150,
			width:      10,
			wantFilled: 10,
			wantEmpty:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := createProgressBar(tt.percent, tt.width)

			// Count filled and empty blocks
			filled := strings.Count(result, "█")
			empty := strings.Count(result, "░")

			if filled != tt.wantFilled {
				t.Errorf("Expected %d filled blocks, got %d", tt.wantFilled, filled)
			}
			if empty != tt.wantEmpty {
				t.Errorf("Expected %d empty blocks, got %d", tt.wantEmpty, empty)
			}
		})
	}
}

func TestPrintFrameSummary(t *testing.T) {
	statuses := []TerminalStatus{
		{
			ID:              0,
			SNR:             15.0,
			ModScheme:       "QPSK",
			DataRate:        30.0,
			TotalData:       1000,
			RemainingData:   500,
			DataThisFrame:   100,
			SlotUtilization: 50.0,
		},
	}

	// This is mostly a smoke test since we're dealing with output
	PrintFrameSummary(0, statuses, 2*time.Millisecond)
}
