package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"strconv"
)

func CreateLattice(n int) *Graph {
	latticeGraph := NewGraph()

	for x := 0; x < n; x++ {
		node := Node{X: x, Y: 0, Z: 0}
		latticeGraph.AddVertex(node)
	}

	for z := 1; z < n; z++ {
		for x := -n + z + 1; x < n-z; x++ {
			node := Node{X: x, Y: 0, Z: z}
			latticeGraph.AddVertex(node)
		}
	}

	for y := 1; y < n; y++ {
		for x := -n + y + 1; x < n-y; x++ {
			for z := -n + 1; z < n; z++ {
				node := Node{X: x, Y: y, Z: z}
				latticeGraph.AddVertex(node)
			}

		}
	}
	latticeGraph.ConnectAdjacentNodes()
	return latticeGraph
}

func CountPolyominoes(graph *Graph, depth int, maxSize int, untriedSet []Node, cellsAdded []Node, path []Node) []int {
	newUntriedSet := make([]Node, len(untriedSet))
	copy(newUntriedSet, untriedSet)

	var oldNeighbours []Node
	if len(newUntriedSet) != 0 && depth+1 < maxSize {
		oldNeighbours = cellsAdded
		for _, cell := range cellsAdded {
			for _, neighbour := range graph.GetNeighbours(cell) {
				oldNeighbours = append(oldNeighbours, neighbour)
			}
		}
	}
	elementCount := make([]int, maxSize)
	for len(newUntriedSet) != 0 {
		randomElement := newUntriedSet[len(newUntriedSet)-1] // step 1
		newUntriedSet = newUntriedSet[:len(newUntriedSet)-1] // step 2
		cellsAdded = append(cellsAdded, randomElement)

		elementCount[depth]++ // Step 3

		if depth+1 < maxSize { // Step 4
			var newNeighbours []Node
			for _, neighbour := range graph.GetNeighbours(randomElement) {
				if !contains(oldNeighbours, neighbour) {
					newNeighbours = append(newNeighbours, neighbour)
				}
			}
			newCounts := CountPolyominoes(graph, depth+1, maxSize, append(newUntriedSet, newNeighbours...), cellsAdded, append(path, randomElement))
			for i := range elementCount {
				elementCount[i] += newCounts[i]
			}
		} else {
			//fmt.Printf("%v\n", append(path, randomElement))
		}
		cellsAdded = cellsAdded[:len(cellsAdded)-1]
	}
	return elementCount
}

func contains(slice []Node, value Node) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
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
	//fmt.Printf("Lattice: %v\n", latticeGraph)
	untriedSet := []Node{{X: 0, Y: 0, Z: 0}}
	var cellsAdded []Node
	var path []Node
	count := CountPolyominoes(latticeGraph, 0, n, untriedSet, cellsAdded, path)
	fmt.Println(count)
}
