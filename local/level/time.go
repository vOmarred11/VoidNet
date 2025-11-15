package world

var (
	TimeDay      = 1000
	TimeNight    = 13000
	TimeNoon     = 6000
	TimeMidnight = 18000
	TimeSunrise  = 23000
	TimeSunset   = 12000
)

// Time is the time on the world
type Time struct {
	// Amount defines the time amount
	Amount uint8
}

func (t Time) timeAmount() uint8 {
	return t.Amount
}
