# Smooth Life - Reaction-Diffusion Simulation

This is a Go program that implements the "Smooth Life" simulation, a variant of Conway's Game of Life using the reaction-diffusion model. The reaction-diffusion model is a mathematical model used to describe the behavior of certain chemical reactions in space and time. The simulation is displayed on the terminal using characters to represent different concentrations.

## Prerequisites

To run the simulation, you need to have Go (Golang) installed on your system. If you haven't installed Go yet, you can download and install it from the official Go website: https://golang.org/

## Running the Simulation

1. Clone or download this repository to your local machine.

2. Open a terminal or command prompt and navigate to the project's root directory.

3. Run the simulation by executing the following command:

```
go run main.go
```

4. The program will prompt you to enter the grid dimensions for the simulation. You can specify the width and height of the grid. The larger the dimensions, the longer the simulation will take to complete.

5. The simulation will start, and the terminal will display the evolving patterns based on the reaction-diffusion model using characters.

6. The simulation will continue until all cells in the grid have reached the end state (value 0), at which point it will print a message indicating the simulation has finished.

## Customizing the Simulation

You can modify the default configuration parameters in the `DefaultConfig` variable to experiment with different simulations. The configuration parameters control various aspects of the reaction-diffusion model, such as the influence radius, alpha values, thresholds, and time step (dt).

Note: Modifying the configuration parameters may result in different patterns and behaviors in the simulation.

## Exiting the Simulation

The simulation will automatically exit once all cells in the grid have reached the end state (value 0). Alternatively, you can interrupt the simulation by pressing `Ctrl + C` in the terminal.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

The simulation code is based on the reaction-diffusion model and its implementation, inspired by various resources and research on reaction-diffusion systems. The initial version of the code might have been adapted from a specific source, but it has been modified and generalized for this repository.
