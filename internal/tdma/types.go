package tdma

// BurstType represents the type of burst
type BurstType int

const (
	DataBurst BurstType = iota
	ControlBurst
	MaintenanceBurst
)
