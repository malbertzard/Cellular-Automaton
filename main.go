package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

const (
	LevelString = " .-=coaA@#"
	LevelCount  = len(LevelString) - 1
)

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

// Display the grid with appropriate characters based on concentration levels
func displayGrid(grid [][]float64) {
	for y := 0; y < len(grid); y++ {
		for x := 0; x < len(grid[y]); x++ {
			c := LevelString[int(grid[y][x]*float64(LevelCount))]
			fmt.Printf("%c%c", c, c)
		}
		fmt.Println()
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
	rand.Seed(time.Now().UnixNano())

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

	// Perform the simulation and display the grid
	for {
		computeGridDiff(config)
		applyGridDiff(config)

		displayGrid(Grid)

		select {
		case <-done:
			fmt.Println("Simulation finished. All cells have reached the end state (value 0).")
			return
		default:
			// Continue with the simulation
		}

		// Add a short delay between simulation steps for smoother animation
		time.Sleep(100 * time.Millisecond)
	}
}
