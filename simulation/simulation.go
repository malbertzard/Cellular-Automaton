 package simulation

import "cellular-autometon/common"

// Simulation represents a general simulation interface.
type Simulation interface {
	InitializeGrid()
	ComputeGridDiff()
	ApplyGridDiff()
	CheckEndState()
	IsEndState() bool
	GetConfig() common.Config
	GetGrid() [][]float64
	GetGridHistory() [][][]float64
	AddGridToHistory(grid [][]float64)
}

// SimulationType represents a general simulation type interface.
type SimulationType interface {
	NewSimulation(config common.Config) Simulation
}

