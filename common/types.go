package common

// Config holds the configuration parameters for the reaction-diffusion simulation
type Config struct {
	Width           int
	Height          int
	InfluenceRadius float64
	AlphaN          float64
	AlphaM          float64
	ThresholdB1     float64
	ThresholdB2     float64
	ThresholdD1     float64
	ThresholdD2     float64
	Dt              float64
}

// LevelString is the string used to represent different concentration levels.
const LevelString = " .-=coaA@#"

// LevelCount is the number of levels in the LevelString.
const LevelCount = len(LevelString) - 1
