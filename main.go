package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const (
	LevelString = " .-=coaA@#"
	LevelCount  = len(LevelString) - 1
)

// OutputInterface defines the interface for displaying the grid and handling simulation completion.
type OutputInterface interface {
	Display(grid [][]float64)
	SimulationFinished()
}

// TerminalOutput is a struct that implements the OutputInterface for displaying the grid in the terminal.
type TerminalOutput struct{}

// Display implements the OutputInterface to show the grid in the terminal.
func (t TerminalOutput) Display(grid [][]float64) {
	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[y]); x++ {
			c := LevelString[int(grid[y][x]*float64(LevelCount))]
			fmt.Printf("%c%c", c, c)
		}
		fmt.Println()
	}
}

// SimulationFinished implements the OutputInterface to notify the end of the simulation in the terminal.
func (t TerminalOutput) SimulationFinished() {
	fmt.Println("Simulation finished. All cells have reached the end state (value 0).")
}

// SVGOutputInterface defines the interface for creating an SVG animation.
type SVGOutputInterface interface {
	GenerateSVGAnimation(gridHistory [][][]float64, config Config)
}

// SVGOutput is a struct that implements the SVGOutputInterface for creating an SVG animation.
type SVGOutput struct{}

// GenerateSVGAnimation implements the SVGOutputInterface to create an SVG animation.
func (s SVGOutput) GenerateSVGAnimation(gridHistory [][][]float64, config Config) {
	// Create the SVG header
	svgHeader := fmt.Sprintf(
		`<svg xmlns="http://www.w3.org/2000/svg" width="%d" height="%d">`,
		config.Width, config.Height,
	)

	// Calculate the maximum value in the grid history for normalization
	maxValue := 0.0
	for _, grid := range gridHistory {
		for y := 0; y < config.Height; y++ {
			for x := 0; x < config.Width; x++ {
				if grid[y][x] > maxValue {
					maxValue = grid[y][x]
				}
			}
		}
	}

	// Create the SVG frames for each time step in gridHistory
	var svgFrames string
	for frameIndex, grid := range gridHistory {
		frame := fmt.Sprintf(`<g display="%s" fill="none">`, strconv.Itoa(frameIndex*100))

		for y := 0; y < config.Height; y++ {
			for x := 0; x < config.Width; x++ {
				// Calculate the normalized value for color
				normalizedValue := grid[y][x] / maxValue
				// Convert the normalized value to a color (here, we'll use a blue-to-red gradient)
				greyValue := int(255 * normalizedValue)
				color := fmt.Sprintf("#%02x%02x%02x", greyValue, greyValue, greyValue)
				// Define the SVG rectangle for the cell with the calculated color
				rect := fmt.Sprintf(`<rect x="%d" y="%d" width="1" height="1" fill="%s" />`, x, y, color)
				frame += rect
			}
		}

		frame += "</g>"
		svgFrames += frame
	}

	// Create the SVG footer
	svgFooter := "</svg>"

	// Combine the header, frames, and footer to complete the SVG code
	svgCode := svgHeader + svgFrames + svgFooter

	// Save the SVG code to a file named "animation.svg"
	err := os.WriteFile("animation.svg", []byte(svgCode), 0644)
	if err != nil {
		fmt.Println("Error writing SVG animation to file:", err)
		return
	}

	fmt.Println("SVG animation generated. Check 'animation.svg' in the current directory.")
}

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

var (
	Grid          [][]float64
	GridDiff      [][]float64
	DefaultConfig = Config{
		Width:           100,
		Height:          100,
		InfluenceRadius: 11.0,
		AlphaN:          0.028,
		AlphaM:          0.147,
		ThresholdB1:     0.278,
		ThresholdB2:     0.365,
		ThresholdD1:     0.267,
		ThresholdD2:     0.445,
		Dt:              0.05,
	}
)

// Helper function to generate a random float between 0 and 1
func randFloat() float64 {
	return rand.Float64()
}

// Initialize the grid with concentrations in the central region
func initializeGrid(config Config) {
	Grid = make([][]float64, config.Height)
	GridDiff = make([][]float64, config.Height)
	for i := range Grid {
		Grid[i] = make([]float64, config.Width)
		GridDiff[i] = make([]float64, config.Width)
	}

	w, h := config.Width/3, config.Height/3
	for dy := 0; dy < h; dy++ {
		for dx := 0; dx < w; dx++ {
			x := dx + config.Width/2 - w/2
			y := dy + config.Height/2 - h/2
			Grid[y][x] = randFloat()
		}
	}
}

// Modulo function that handles negative numbers correctly
func emod(a, b int) int {
	return (a%b + b) % b
}

// Calculate the sigmoid function
func calculateSigmoid(x, a, alpha float64) float64 {
	return 1.0 / (1.0 + math.Exp(-(x-a)*4/alpha))
}

// Calculate the 'n' component of the reaction-diffusion model
func calculateComponentN(x, a, b, alpha float64) float64 {
	return calculateSigmoid(x, a, alpha) * (1 - calculateSigmoid(x, b, alpha))
}

// Calculate the 'm' component of the reaction-diffusion model
func calculateComponentM(x, y, m, alpha float64) float64 {
	return x*(1-calculateSigmoid(m, 0.5, alpha)) + y*calculateSigmoid(m, 0.5, alpha)
}

// Calculate the function 's' used in the reaction-diffusion model
func calculateFunctionS(n, m, alphaN, alphaM, b1, d1, b2, d2 float64) float64 {
	return calculateComponentN(n, calculateComponentM(b1, d1, m, alphaM), calculateComponentM(b2, d2, m, alphaM), alphaN)
}

// Compute the differences in concentration between neighboring cells
func computeGridDiff(config Config) {
	for cy := 0; cy < config.Height; cy++ {
		for cx := 0; cx < config.Width; cx++ {
			m, M := 0.0, 0.0
			n, N := 0.0, 0.0

			for dy := -int(config.InfluenceRadius - 1); dy <= int(config.InfluenceRadius-1); dy++ {
				for dx := -int(config.InfluenceRadius - 1); dx <= int(config.InfluenceRadius-1); dx++ {
					x := emod(cx+dx, config.Width)
					y := emod(cy+dy, config.Height)
					if dx*dx+dy*dy <= int(config.InfluenceRadius*config.InfluenceRadius/9) {
						m += Grid[y][x]
						M += 1
					} else if dx*dx+dy*dy <= int(config.InfluenceRadius*config.InfluenceRadius) {
						n += Grid[y][x]
						N += 1
					}
				}
			}
			m /= M
			n /= N
			q := calculateFunctionS(n, m, config.AlphaN, config.AlphaM, config.ThresholdB1, config.ThresholdD1, config.ThresholdB2, config.ThresholdD2)
			GridDiff[cy][cx] = 2*q - 1
		}
	}
}

// Clamp the value within a specified range
func clamp(x *float64, min, max float64) {
	if *x < min {
		*x = min
	}
	if *x > max {
		*x = max
	}
}

// Apply the computed differences to update the grid's concentration levels
func applyGridDiff(config Config) {
	for y := 0; y < config.Height; y++ {
		for x := 0; x < config.Width; x++ {
			Grid[y][x] += config.Dt * GridDiff[y][x]
			clamp(&Grid[y][x], 0, 1)
		}
	}
}

// Check if the grid has reached the end state (all cells have value 0)
func isEndState(grid [][]float64, done chan<- struct{}) {
	hasEnded := true
	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[y]); x++ {
			if grid[y][x] != 0 {
				hasEnded = false
				break
			}
		}
		if !hasEnded {
			break
		}
	}
	if hasEnded {
		done <- struct{}{}
	}
}

func main() {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	// Terminal dialog to set grid dimensions
	fmt.Println("Enter grid dimensions for the reaction-diffusion simulation:")
	var config Config
	fmt.Print("Width: ")
	fmt.Scan(&config.Width)
	fmt.Print("Height: ")
	fmt.Scan(&config.Height)

	// Set the configuration parameters to default values
	config.InfluenceRadius = DefaultConfig.InfluenceRadius
	config.AlphaN = DefaultConfig.AlphaN
	config.AlphaM = DefaultConfig.AlphaM
	config.ThresholdB1 = DefaultConfig.ThresholdB1
	config.ThresholdB2 = DefaultConfig.ThresholdB2
	config.ThresholdD1 = DefaultConfig.ThresholdD1
	config.ThresholdD2 = DefaultConfig.ThresholdD2
	config.Dt = DefaultConfig.Dt

	// Initialize the grid
	initializeGrid(config)

	done := make(chan struct{})
	go func() {
		for {
			isEndState(Grid, done)
			time.Sleep(500 * time.Millisecond) // Adjust the check interval if needed
		}
	}()

	// Create an instance of TerminalOutput to use as the output interface
	terminalOutput := TerminalOutput{}
	svgOutput := SVGOutput{}

	// Perform the simulation and display the grid
	var gridHistory [][][]float64
	for {
		computeGridDiff(config)
		applyGridDiff(config)

		// Store the current state of the grid in the history
		currentGrid := make([][]float64, config.Height)
		for i := range currentGrid {
			currentGrid[i] = make([]float64, config.Width)
			copy(currentGrid[i], Grid[i])
		}
		gridHistory = append(gridHistory, currentGrid)

		// Use the terminal output interface to display the grid
		terminalOutput.Display(Grid)

		select {
		case <-done:
			// Notify the simulation completion through the terminal output interface
			terminalOutput.SimulationFinished()

			// Generate the SVG animation using the SVG output interface
			svgOutput.GenerateSVGAnimation(gridHistory, config)

			return
		default:
			// Continue with the simulation
		}

		// Add a short delay between simulation steps for smoother animation
		time.Sleep(100 * time.Millisecond)
	}
}

