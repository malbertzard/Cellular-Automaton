package simulation

import (
	"cellular-autometon/common"
	"math/rand"
)

// ConwaySimulationType represents Conway's Game of Life simulation type.
type ConwaySimulationType struct{}

func (c ConwaySimulationType) NewSimulation(config common.Config) Simulation {
	return NewConwaySimulation(config)
}

// ConwaySimulation represents the Conway's Game of Life simulation.
type ConwaySimulation struct {
	Config      common.Config
	Grid        [][]float64
	NextGrid    [][]float64
	GridHistory [][][]float64
	State       bool
}

// NewConwaySimulation creates a new instance of Conway's Game of Life simulation.
func NewConwaySimulation(config common.Config) Simulation {
	return &ConwaySimulation{
		Config:      config,
		Grid:        nil,
		NextGrid:    nil,
		GridHistory: make([][][]float64, 0),
		State:       false,
	}
}

// InitializeGrid initializes the grid with random values for Conway's Game of Life simulation.
func (c *ConwaySimulation) InitializeGrid() {
	c.Grid = make([][]float64, c.Config.Height)
	c.NextGrid = make([][]float64, c.Config.Height)
	for i := range c.Grid {
		c.Grid[i] = make([]float64, c.Config.Width)
		c.NextGrid[i] = make([]float64, c.Config.Width)
		for j := range c.Grid[i] {
			// Randomly assign values 0 or 1
			c.Grid[i][j] = float64(rand.Intn(2))
		}
	}
}

// ComputeGridDiff computes the differences in concentration between neighboring cells for Conway's Game of Life.
func (c *ConwaySimulation) ComputeGridDiff() {
	for y := 0; y < c.Config.Height; y++ {
		for x := 0; x < c.Config.Width; x++ {
			neighbors := 0
			for dy := -1; dy <= 1; dy++ {
				for dx := -1; dx <= 1; dx++ {
					if dy == 0 && dx == 0 {
						continue
					}
					ny, nx := y+dy, x+dx
					if ny >= 0 && ny < c.Config.Height && nx >= 0 && nx < c.Config.Width && c.Grid[ny][nx] == 1 {
						neighbors++
					}
				}
			}
			if c.Grid[y][x] == 1 {
				if neighbors < 2 || neighbors > 3 {
					c.NextGrid[y][x] = 0
				} else {
					c.NextGrid[y][x] = 1
				}
			} else {
				if neighbors == 3 {
					c.NextGrid[y][x] = 1
				} else {
					c.NextGrid[y][x] = 0
				}
			}
		}
	}
}

// ApplyGridDiff applies the computed differences to update the grid's concentration levels for Conway's Game of Life.
func (c *ConwaySimulation) ApplyGridDiff() {
	for y := 0; y < c.Config.Height; y++ {
		for x := 0; x < c.Config.Width; x++ {
			c.Grid[y][x] = c.NextGrid[y][x]
		}
	}
}

// CheckEndState checks if the grid has reached the end state (all cells have the same value) for Conway's Game of Life.
func (c *ConwaySimulation) CheckEndState() {
	// For Conway's Game of Life, we will consider reaching a stable state (grid stops changing) as the end state.
	if len(c.GridHistory) >= 2 {
		latestGrid := c.GridHistory[len(c.GridHistory)-1]
		previousGrid := c.GridHistory[len(c.GridHistory)-2]

		isStableState := true
		for y := 0; y < c.Config.Height; y++ {
			for x := 0; x < c.Config.Width; x++ {
				if latestGrid[y][x] != previousGrid[y][x] {
					isStableState = false
					break
				}
			}
			if !isStableState {
				break
			}
		}

		c.State = isStableState
	}
}

// GetConfig returns the configuration for Conway's Game of Life simulation.
func (c *ConwaySimulation) GetConfig() common.Config {
	return c.Config
}

func (c *ConwaySimulation) IsEndState() bool {
	return c.State
}

// GetGrid returns the current grid for Conway's Game of Life simulation.
func (c *ConwaySimulation) GetGrid() [][]float64 {
	return c.Grid
}

// GetGridHistory returns the grid history for Conway's Game of Life simulation.
func (c *ConwaySimulation) GetGridHistory() [][][]float64 {
	return c.GridHistory
}

// AddGridToHistory adds the current grid to the grid history for Conway's Game of Life simulation.
func (c *ConwaySimulation) AddGridToHistory(grid [][]float64) {
	c.GridHistory = append(c.GridHistory, grid)
}
