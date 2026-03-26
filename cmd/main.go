package main

import (
	"fmt"
	"lem-in/internal/optimizer"
	"lem-in/internal/parser"
	"lem-in/internal/simulator"
	"lem-in/internal/solver"
	"lem-in/utils"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run ./cmd/ <filename>")
		os.Exit(1)
	}

	filename := os.Args[1]

	// Print original input (required by subject)
	lines, err := utils.ReadFile(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: cannot read file:", err)
		os.Exit(1)
	}
	for _, line := range lines {
		fmt.Println(line)
	}
	fmt.Println() // blank line before moves

	// Parse the file
	colony, err := parser.ParseFile(filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Verify a path exists
	shortestPath := solver.FindShortestPath(colony)
	if shortestPath == nil {
		fmt.Fprintln(os.Stderr, "ERROR: no path found from start to end")
		os.Exit(1)
	}

	// Find vertex-disjoint paths via max-flow
	disjointPaths := solver.FindDisjointPaths(colony)
	if len(disjointPaths) == 0 {
		fmt.Fprintln(os.Stderr, "ERROR: no paths found")
		os.Exit(1)
	}

	// Find optimal ant distribution across those paths
	opt := optimizer.NewOptimizer(disjointPaths, colony.AntCount)
	optimalAssignment := opt.FindOptimalAssignment()
	if optimalAssignment == nil {
		fmt.Fprintln(os.Stderr, "ERROR: cannot find optimal path assignment")
		os.Exit(1)
	}

	// Run simulation
	sim := simulator.NewSimulator(colony, optimalAssignment)
	sim.RunSimulation()
}