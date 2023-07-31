# Cellular Automaton

<!--toc:start-->
- [Cellular Automaton](#cellular-automaton)
  - [Installation](#installation)
  - [Usage](#usage)
  - [Available Simulation Types](#available-simulation-types)
<!--toc:end-->

Cellular Automaton is a simulation project that allows users to explore and visualize various cellular automaton models. The project includes different simulation types, such as the classic Conway's Game of Life and the continuous variant known as SmoothLife. Users can select the simulation type and interact with the simulation through a terminal interface or generate animated visualizations in SVG format.

## Installation

1. Clone the repository to your local machine:

```
git clone https://github.com/malbertzard/cellular-automaton.git
```

2. Change into the project directory:

```
cd cellular-automaton
```

3. Build the project (if required):

```
go build
```

## Usage

1. Run the main executable:

```
./cellular-automaton
```

2. The application will prompt you to enter grid dimensions for the simulation.

3. Select a simulation type from the available options (e.g., SmoothLife, Conway's Game of Life).

4. The simulation will start, and you can observe the changes in the grid over time.

5. If the simulation reaches the end state (all cells have value 0), the program will display a message indicating the completion.

6. The simulation output will be displayed in the terminal in real-time.

7. An SVG animation will be generated and saved to the output directory upon completion of the simulation.

## Available Simulation Types

1. **SmoothLife**: A continuous variant of Conway's Game of Life that allows for smoother transitions between cell states.

2. **Conway's Game of Life**: A classic cellular automaton that follows simple rules to determine the next state of each cell based on its neighbors.
