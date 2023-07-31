package main

import (
	"fmt"
	"time"

	"cellular-autometon/common"
	"cellular-autometon/output"
	"cellular-autometon/simulation"
)

func main() {
	// Terminal dialog to set grid dimensions
	fmt.Println("Enter grid dimensions for the simulation:")
	var config common.Config
	fmt.Print("Width: ")
	fmt.Scan(&config.Width)
	fmt.Print("Height: ")
	fmt.Scan(&config.Height)

	// Set the configuration parameters to default values
	config.InfluenceRadius = 11.0
	config.AlphaN = 0.028
	config.AlphaM = 0.147
	config.ThresholdB1 = 0.278
	config.ThresholdB2 = 0.365
	config.ThresholdD1 = 0.267
	config.ThresholdD2 = 0.445
	config.Dt = 0.05

	// Create a map to associate simulation types with their constructor functions
	simulationTypes := map[string]simulation.SimulationType{
		"SmoothLife": simulation.SmoothLifeSimulationType{},
		"Conway":     simulation.ConwaySimulationType{},
	}

	// Get the selected simulation type from the user
	fmt.Println("Select a simulation type:")
	for key := range simulationTypes {
		fmt.Println("-", key)
	}
	var selectedType string
	fmt.Scan(&selectedType)

	// Check if the selected simulation type is valid
	selectedSimulationType, found := simulationTypes[selectedType]
	if !found {
		fmt.Println("Invalid simulation type selected.")
		return
	}

	// Initialize the selected simulation
	sim := selectedSimulationType.NewSimulation(config)
	sim.InitializeGrid()

	// Create instances of TerminalOutput and SVGOutput to use as the output interfaces
	terminalOutput := output.NewTerminalOutput()
	svgOutput := output.NewSVGOutput()

	done := make(chan struct{})
	go func() {
		for {
			sim.CheckEndState()
			if sim.IsEndState() {
				done <- struct{}{}
				break
			}
			time.Sleep(500 * time.Millisecond) // Adjust the check interval if needed
		}
	}()

	// Perform the simulation and display the grid
	for {
		sim.ComputeGridDiff()
		sim.ApplyGridDiff()

		// Store the current state of the grid in the history
		currentGrid := make([][]float64, sim.GetConfig().Height)
		for i := range currentGrid {
			currentGrid[i] = make([]float64, sim.GetConfig().Width)
			copy(currentGrid[i], sim.GetGrid()[i])
		}
		sim.AddGridToHistory(currentGrid)

		// Use the terminal output interface to display the grid
		terminalOutput.Display(sim.GetGrid())

		select {
		case <-done:
			// Notify the simulation completion through the terminal output interface
			terminalOutput.SimulationFinished()

			// Generate the SVG animation using the SVG output interface
			svgOutput.GenerateSVGAnimation(sim.GetGridHistory(), config)

			return
		default:
			// Continue with the simulation
		}

		// Add a short delay between simulation steps for smoother animation
		time.Sleep(100 * time.Millisecond)
	}
}
