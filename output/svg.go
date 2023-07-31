package output

import (
	"cellular-autometon/common"
	"fmt"
	"os"
	"strconv"
)

// SVGOutput is a struct that implements the SVGOutputInterface for creating an SVG animation.
type SVGOutput struct{}

// NewSVGOutput creates a new instance of SVGOutput.
func NewSVGOutput() *SVGOutput {
	return &SVGOutput{}
}

// GenerateSVGAnimation implements the SVGOutputInterface to create an SVG animation.
func (s SVGOutput) GenerateSVGAnimation(gridHistory [][][]float64, config common.Config) {
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

		// Set a shorter duration for smoother animation
		duration := 100 // 100 milliseconds

		// Add the animate element to loop the frames seamlessly with shorter duration
		frame += fmt.Sprintf(`<animate attributeName="display" from="none" to="inline" begin="%dms" dur="%dms" repeatCount="indefinite" />`, frameIndex*100, duration)

		svgFrames += frame
	}

	// Create the SVG footer
	svgFooter := "</svg>"

	// Combine the header, frames, and footer to complete the SVG code
	svgCode := svgHeader + svgFrames + svgFooter

	// Save the minified SVG code to a file named "animation.svg"
	err := os.WriteFile("animation.svg", []byte(svgCode), 0644)
	if err != nil {
		fmt.Println("Error writing SVG animation to file:", err)
		return
	}

	fmt.Println("SVG animation generated and optimized. Check 'animation.svg' in the current directory.")
}

