package output

import (
	"cellular-autometon/common"
	"fmt"
)

// TerminalOutput is a struct that implements the OutputInterface for displaying the grid in the terminal.
type TerminalOutput struct{}

// NewTerminalOutput creates a new instance of TerminalOutput.
func NewTerminalOutput() *TerminalOutput {
	return &TerminalOutput{}
}

// Display implements the OutputInterface to show the grid in the terminal.
func (t TerminalOutput) Display(grid [][]float64) {
	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[y]); x++ {
			c := common.LevelString[int(grid[y][x]*float64(common.LevelCount))]
			fmt.Printf("%c%c", c, c)
		}
		fmt.Println()
	}
}

// SimulationFinished implements the OutputInterface to notify the end of the simulation in the terminal.
func (t TerminalOutput) SimulationFinished() {
	fmt.Println("Simulation finished. All cells have reached the end state (value 0).")
}
