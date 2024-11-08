// internal/tdma/frame_test.go
package tdma

import (
	"testing"
	"time"
)

func TestNewFrame(t *testing.T) {
	tests := []struct {
		name    string
		config  FrameConfig
		wantErr bool
	}{
		{
			name: "Valid Configuration",
			config: FrameConfig{
				FrameDuration: 2 * time.Millisecond,
				GuardTime:     50 * time.Microsecond,
				NumCarriers:   4,
				SlotDurations: []time.Duration{
					600 * time.Microsecond,
					400 * time.Microsecond,
					300 * time.Microsecond,
					200 * time.Microsecond,
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid Slot Duration Count",
			config: FrameConfig{
				FrameDuration: 2 * time.Millisecond,
				GuardTime:     50 * time.Microsecond,
				NumCarriers:   4,
				SlotDurations: []time.Duration{
					600 * time.Microsecond,
					400 * time.Microsecond,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			frame, err := NewFrame(tt.config)

			if (err != nil) != tt.wantErr {
				t.Errorf("NewFrame() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && frame == nil {
				t.Error("Expected non-nil frame for valid configuration")
			}

			if !tt.wantErr {
				// Check if time slots were created correctly
				expectedSlots := tt.config.NumCarriers * 2 // Including guard slots
				if len(frame.TimeSlots) != expectedSlots {
					t.Errorf("Expected %d time slots, got %d",
						expectedSlots, len(frame.TimeSlots))
				}

				// Verify SNR profiles initialization
				if len(frame.SNRProfiles) != tt.config.NumCarriers {
					t.Errorf("Expected %d SNR profiles, got %d",
						tt.config.NumCarriers, len(frame.SNRProfiles))
				}
			}
		})
	}
}

func TestUpdateSNR(t *testing.T) {
	config := FrameConfig{
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

	frame, err := NewFrame(config)
	if err != nil {
		t.Fatalf("Failed to create frame: %v", err)
	}

	// Test SNR updates over multiple frames
	for i := 0; i < 10; i++ {
		snrValues := frame.UpdateSNR()

		// Verify SNR values are within expected ranges
		for carrierID, snr := range snrValues {
			profile := frame.SNRProfiles[carrierID]
			minSNR := profile.BaselineSNR - profile.Variation
			maxSNR := profile.BaselineSNR + profile.Variation

			if snr != 0 && (snr < minSNR || snr > maxSNR) {
				t.Errorf("Frame %d, Carrier %d: SNR %.1f outside expected range [%.1f, %.1f]",
					i, carrierID, snr, minSNR, maxSNR)
			}
		}

		frame.AdvanceFrame()
	}
}

func TestAddBurst(t *testing.T) {
	config := FrameConfig{
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

	frame, _ := NewFrame(config)

	tests := []struct {
		name      string
		carrierID int
		burst     *Burst
		wantErr   bool
	}{
		{
			name:      "Valid Burst Addition",
			carrierID: 0,
			burst:     NewBurstWithSNR(make([]byte, 1000), 0, DataBurst, 15.0),
			wantErr:   false,
		},
		{
			name:      "Invalid Carrier ID",
			carrierID: 10,
			burst:     NewBurstWithSNR(make([]byte, 1000), 10, DataBurst, 15.0),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := frame.AddBurst(tt.carrierID, tt.burst)

			if (err != nil) != tt.wantErr {
				t.Errorf("AddBurst() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				// Verify burst was added correctly
				slot := frame.TimeSlots[tt.carrierID*2] // Multiply by 2 for guard slots
				if slot.Burst != tt.burst {
					t.Error("Burst was not properly added to time slot")
				}
			}
		})
	}
}
