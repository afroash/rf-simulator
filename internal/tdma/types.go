package tdma

import (
	"time"

	"github.com/afroash/rf-simulator/internal/modulation"
)

// BurstType represents the type of burst
type BurstType int

const (
	DataBurst BurstType = iota
	ControlBurst
	MaintenanceBurst
)

// Terminal is a represents a user sending data.
type Terminal struct {
	ID            int
	TotalSize     int64
	RemainingData int64
	BaselineSNR   float64
	Priority      int
}

// v2 Additions
// FrameConfig holds configuration for TDMA frame
type FrameConfig struct {
	FrameDuration time.Duration
	GuardTime     time.Duration
	NumCarriers   int
	SlotDurations []time.Duration // Allow different slot durations
}

// Frame represents a sequence of TDMA frames
type Frame struct {
	Config      FrameConfig
	TimeSlots   []*TimeSlot
	FrameNumber int
	SNRProfiles map[int]*SNRProfile // Per-carrier SNR profiles
	Terminals   []*Terminal
}

// SNRProfile represents dynamic SNR characteristics for a carrier
type SNRProfile struct {
	BaselineSNR     float64
	Variation       float64 // Maximum SNR variation
	UpdateInterval  int     // Frames between SNR updates
	LastUpdateFrame int
}

// TDMAFrame represents a complete TDMA frame.
type TDMAFrame struct {
	Duration    time.Duration
	GuardTime   time.Duration
	TimeSlots   []*TimeSlot
	NumCarriers int
}

// Constants for calculations
const (
	// Typical symbol rate for satellite communications (30 MHz bandwidth)
	BaseSymbolRate = 25e6 // 25 MSymbols/second
)

type Burst struct {
	Data          []byte
	StartTime     time.Duration
	Duration      time.Duration
	CarrierID     int
	Type          BurstType
	Utilisation   float64 // Percentage of the time slot used by the burst
	Modulation    modulation.ModulationScheme
	SNR           float64
	BER           float64
	Datarate      float64
	SymbolsPacked int
}

type TimeSlot struct {
	StartTime   time.Duration
	Duration    time.Duration
	Burst       *Burst
	IsGuardTime bool
}

// TerminalStatus represents current terminal transmission status
type TerminalStatus struct {
	ID              int
	SNR             float64
	ModScheme       string
	DataRate        float64
	TotalData       int64
	RemainingData   int64
	DataThisFrame   int64
	SlotUtilization float64
}
