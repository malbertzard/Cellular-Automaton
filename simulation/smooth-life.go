package simulation

import (
	"cellular-autometon/common"
	"math"
	"math/rand"
)

// SmoothLifeSimulationType implements the SimulationType interface for SmoothLifeSimulation.
type SmoothLifeSimulationType struct{}

// NewSimulation creates a new instance of SmoothLifeSimulation.
func (t SmoothLifeSimulationType) NewSimulation(config common.Config) Simulation {
	return NewSmoothLifeSimulation(config)
}

// SmoothLifeSimulation represents the SmoothLife simulation.
type SmoothLifeSimulation struct {
	Config      common.Config
	Grid        [][]float64
	GridDiff    [][]float64
	GridHistory [][][]float64
	State       bool
}

// NewSmoothLifeSimulation creates a new instance of the SmoothLife simulation.
func NewSmoothLifeSimulation(config common.Config) Simulation {
	return &SmoothLifeSimulation{
		Config:      config,
		Grid:        nil,
		GridDiff:    nil,
		GridHistory: make([][][]float64, 0),
		State:       false,
	}
}

// InitializeGrid initializes the grid with random concentrations for the SmoothLife simulation.
func (s *SmoothLifeSimulation) InitializeGrid() {
	s.Grid = make([][]float64, s.Config.Height)
	s.GridDiff = make([][]float64, s.Config.Height)
	for i := range s.Grid {
		s.Grid[i] = make([]float64, s.Config.Width)
		s.GridDiff[i] = make([]float64, s.Config.Width)
	}

	w, h := s.Config.Width/3, s.Config.Height/3
	for dy := 0; dy < h; dy++ {
		for dx := 0; dx < w; dx++ {
			x := dx + s.Config.Width/2 - w/2
			y := dy + s.Config.Height/2 - h/2
			s.Grid[y][x] = rand.Float64()
		}
	}
}

// ComputeGridDiff computes the differences in concentration between neighboring cells for the SmoothLife simulation.
func (s *SmoothLifeSimulation) ComputeGridDiff() {
	for cy := 0; cy < s.Config.Height; cy++ {
		for cx := 0; cx < s.Config.Width; cx++ {
			m, M := 0.0, 0.0
			n, N := 0.0, 0.0

			for dy := -int(s.Config.InfluenceRadius - 1); dy <= int(s.Config.InfluenceRadius-1); dy++ {
				for dx := -int(s.Config.InfluenceRadius - 1); dx <= int(s.Config.InfluenceRadius-1); dx++ {
					x := emod(cx+dx, s.Config.Width)
					y := emod(cy+dy, s.Config.Height)
					if dx*dx+dy*dy <= int(s.Config.InfluenceRadius*s.Config.InfluenceRadius/9) {
						m += s.Grid[y][x]
						M += 1
					} else if dx*dx+dy*dy <= int(s.Config.InfluenceRadius*s.Config.InfluenceRadius) {
						n += s.Grid[y][x]
						N += 1
					}
				}
			}
			m /= M
			n /= N
			q := calculateFunctionS(n, m, s.Config.AlphaN, s.Config.AlphaM, s.Config.ThresholdB1, s.Config.ThresholdD1, s.Config.ThresholdB2, s.Config.ThresholdD2)
			s.GridDiff[cy][cx] = 2*q - 1
		}
	}
}

// ApplyGridDiff applies the computed differences to update the grid's concentration levels for the SmoothLife simulation.
func (s *SmoothLifeSimulation) ApplyGridDiff() {
	for y := 0; y < s.Config.Height; y++ {
		for x := 0; x < s.Config.Width; x++ {
			s.Grid[y][x] += s.Config.Dt * s.GridDiff[y][x]
			clamp(&s.Grid[y][x], 0, 1)
		}
	}
}

// CheckEndState checks if the grid has reached the end state (all cells have value 0) for the SmoothLife simulation.
func (s *SmoothLifeSimulation) CheckEndState() {
	hasEnded := true
	for y := 0; y < len(s.Grid); y++ {
		for x := 0; x < len(s.Grid[y]); x++ {
			if s.Grid[y][x] != 0 {
				hasEnded = false
				break
			}
		}
		if !hasEnded {
			break
		}
	}
	s.State = hasEnded
}

// GetConfig returns the configuration for the SmoothLife simulation.
func (s *SmoothLifeSimulation) GetConfig() common.Config {
	return s.Config
}

func (s *SmoothLifeSimulation) IsEndState() bool {
	return s.State
}

// GetGrid returns the current grid for the SmoothLife simulation.
func (s *SmoothLifeSimulation) GetGrid() [][]float64 {
	return s.Grid
}

// GetGridHistory returns the grid history for the SmoothLife simulation.
func (s *SmoothLifeSimulation) GetGridHistory() [][][]float64 {
	return s.GridHistory
}

// AddGridToHistory adds the current grid to the grid history for the SmoothLife simulation.
func (s *SmoothLifeSimulation) AddGridToHistory(grid [][]float64) {
	s.GridHistory = append(s.GridHistory, grid)
}

// Helper function to calculate the modulus correctly for negative numbers.
func emod(a, b int) int {
	return (a%b + b) % b
}

// Helper function to calculate the sigmoid function.
func calculateSigmoid(x, a, alpha float64) float64 {
	return 1.0 / (1.0 + math.Exp(-(x-a)*4/alpha))
}

// Helper function to calculate the 'n' component of the SmoothLife model.
func calculateComponentN(x, a, b, alpha float64) float64 {
	return calculateSigmoid(x, a, alpha) * (1 - calculateSigmoid(x, b, alpha))
}

// Helper function to calculate the 'm' component of the SmoothLife model.
func calculateComponentM(x, y, m, alpha float64) float64 {
	return x*(1-calculateSigmoid(m, 0.5, alpha)) + y*calculateSigmoid(m, 0.5, alpha)
}

// Helper function to calculate the function 's' used in the SmoothLife model.
func calculateFunctionS(n, m, alphaN, alphaM, b1, d1, b2, d2 float64) float64 {
	return calculateComponentN(n, calculateComponentM(b1, d1, m, alphaM), calculateComponentM(b2, d2, m, alphaM), alphaN)
}

// Helper function to clamp the value within a specified range.
func clamp(x *float64, min, max float64) {
	if *x < min {
		*x = min
	}
	if *x > max {
		*x = max
	}
}
