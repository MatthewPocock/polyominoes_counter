package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"strconv"
)

func CreateLattice(n int) *Graph {
	latticeGraph := NewGraph()

	for i := 0; i < n; i++ {
		node := Node{X: 0, Y: i}
		latticeGraph.AddVertex(node)
	}
	for i := 1; i < n; i++ {
		for j := -n + i + 1; j < n-i; j++ {
			node := Node{X: i, Y: j}
			latticeGraph.AddVertex(node)
		}
	}
	latticeGraph.ConnectAdjacentNodes()
	return latticeGraph
}

func CountPolyominoes(graph *Graph, depth int, maxSize int, untriedSet []Node, cellsAdded []Node) map[int]int {
	var oldNeighbours []Node
	if len(untriedSet) != 0 && depth+1 < maxSize {
		oldNeighbours = cellsAdded
		for _, cell := range cellsAdded {
			for k := range graph.GetNeighbours(cell) {
				oldNeighbours = append(oldNeighbours, k)
			}
		}
	}
	elementCount := make(map[int]int)
	for len(untriedSet) != 0 {
		randomElement := untriedSet[0] // Step 1
		untriedSet = untriedSet[1:]    // Step 2
		cellsAdded = append(cellsAdded, randomElement)

		elementCount[depth+1]++ // Step 3

		if depth+1 < maxSize { // Step 4
			var newNeighbours []Node
			for neighbour := range graph.GetNeighbours(randomElement) {
				if !contains(oldNeighbours, neighbour) {
					newNeighbours = append(newNeighbours, neighbour)
				}
			}
			newUntriedSet := append(untriedSet, newNeighbours...)
			newCounts := CountPolyominoes(graph, depth+1, maxSize, newUntriedSet, cellsAdded)
			for k, v := range newCounts {
				elementCount[k] = elementCount[k] + v
			}
		}
		cellsAdded = cellsAdded[:len(cellsAdded)-1]
	}
	return elementCount
}

func main() {
	cpuProfile, err := os.Create("cpu.pprof")
	if err != nil {
		fmt.Println("Could not create CPU profile: ", err)
		return
	}
	defer func(cpuProfile *os.File) {
		err := cpuProfile.Close()
		if err != nil {
		}
	}(cpuProfile)

	if err := pprof.StartCPUProfile(cpuProfile); err != nil {
		fmt.Println("Could not start CPU profile: ", err)
		return
	}
	defer pprof.StopCPUProfile()

	// Check if an argument is provided
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [n]")
		return
	}

	// Convert the argument from string to int
	n, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("Invalid input: %v\n", err)
		return
	}

	latticeGraph := CreateLattice(n)
	untriedSet := []Node{{X: 0, Y: 0}}
	var cellsAdded []Node
	count := CountPolyominoes(latticeGraph, 0, n, untriedSet, cellsAdded)
	fmt.Println(count)
}

func contains(slice []Node, value Node) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
