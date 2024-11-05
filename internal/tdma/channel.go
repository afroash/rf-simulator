package tdma

const (
	BPSK ModulationType = iota
	QPSK
	PSK8
)

type ModulationType int

// Channel represents an RF channel Configuration.
type Channel struct {
	//Frequency in MHz
	CenterFreq    float64 //MHz
	Bandwidth     float64 //MHz
	Modulation    ModulationType
	SymbolRate    float64 //Msymbols/s
	BitsPerSymbol int
}

// CalculateChannelCapacity calculates theoretical channel capacity
func (ch *Channel) CalculateChannelCapacity() float64 {
	return ch.SymbolRate * float64(ch.BitsPerSymbol) // Basic calculation in Mbps
}
