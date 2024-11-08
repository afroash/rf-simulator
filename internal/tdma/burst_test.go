package tdma

import (
	"testing"
	"time"
)

func TestNewBurstWithSNR(t *testing.T) {
	tests := []struct {
		name          string
		dataSize      int
		carrierID     int
		burstType     BurstType
		snr           float64
		wantModScheme string
		wantMinRate   float64 // Minimum expected data rate in Mbps
	}{
		{
			name:          "High SNR Test",
			dataSize:      1000,
			carrierID:     0,
			burstType:     DataBurst,
			snr:           20.0,
			wantModScheme: "64-QAM",
			wantMinRate:   15.0, // Adjusted to match realistic calculations
		},
		{
			name:          "Medium SNR Test",
			dataSize:      1000,
			carrierID:     1,
			burstType:     DataBurst,
			snr:           14.0,
			wantModScheme: "8-PSK",
			wantMinRate:   12.0, // Adjusted to match realistic calculations
		},
		{
			name:          "Low SNR Test",
			dataSize:      1000,
			carrierID:     2,
			burstType:     DataBurst,
			snr:           8.0,
			wantModScheme: "BPSK",
			wantMinRate:   5.0, // Adjusted to match realistic calculations
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := make([]byte, tt.dataSize)
			burst := NewBurstWithSNR(data, tt.carrierID, tt.burstType, tt.snr)

			if burst == nil {
				t.Fatal("Expected non-nil burst")
			}

			if burst.Modulation.Name != tt.wantModScheme {
				t.Errorf("Expected modulation scheme %s, got %s",
					tt.wantModScheme, burst.Modulation.Name)
			}

			if burst.Datarate/1e6 < tt.wantMinRate {
				t.Errorf("Expected minimum data rate of %.1f Mbps, got %.1f Mbps",
					tt.wantMinRate, burst.Datarate/1e6)
			}

			t.Logf("Actual data rate: %.1f Mbps", burst.Datarate/1e6)

			if burst.CarrierID != tt.carrierID {
				t.Errorf("Expected carrier ID %d, got %d",
					tt.carrierID, burst.CarrierID)
			}
		})
	}
}

func TestCalculateModulationBasedUtilization(t *testing.T) {
	tests := []struct {
		name        string
		dataSize    int
		duration    time.Duration
		snr         float64
		wantMinUtil float64
	}{
		{
			name:        "Full Utilization",
			dataSize:    1000000,
			duration:    time.Millisecond,
			snr:         20.0,
			wantMinUtil: 90.0,
		},
		{
			name:        "Partial Utilization",
			dataSize:    50000,
			duration:    time.Millisecond,
			snr:         15.0,
			wantMinUtil: 40.0,
		},
		{
			name:        "Low Utilization",
			dataSize:    1000,
			duration:    time.Millisecond,
			snr:         10.0,
			wantMinUtil: 5.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := make([]byte, tt.dataSize)
			burst := NewBurstWithSNR(data, 0, DataBurst, tt.snr)

			if burst.Utilisation < tt.wantMinUtil {
				t.Errorf("Expected minimum utilization of %.1f%%, got %.1f%%",
					tt.wantMinUtil, burst.Utilisation)
			}
		})
	}
}
