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

	for y := 1; y < n; y++ {
		for x := -n + y + 1; x < n-y; x++ {
			node := Node{X: x, Y: y, Z: 0}
			latticeGraph.AddVertex(node)
		}
	}

	for z := 1; z < n; z++ {
		for x := -n + z + 1; x < n-z; x++ {
			for y := -n + 1; y < n; y++ {
				node := Node{X: x, Y: y, Z: z}
				latticeGraph.AddVertex(node)
			}
		}
	}

	latticeGraph.ConnectAdjacentNodes()
	return latticeGraph
}

func CountPolyominoes(graph *Graph, depth int, maxSize int, untriedSet []Node, cellsAdded []Node, oldNeighbours map[Node]int) []int {
	newUntriedSet := make([]Node, len(untriedSet))
	copy(newUntriedSet, untriedSet)

	elementCount := make([]int, maxSize)
	for len(newUntriedSet) != 0 {
		randomElement := newUntriedSet[len(newUntriedSet)-1] // step 1
		newUntriedSet = newUntriedSet[:len(newUntriedSet)-1] // step 2

		elementCount[depth]++ // Step 3

		if depth+1 < maxSize { // Step 4
			var newNeighbours []Node
			for _, neighbour := range graph.GetNeighbours(randomElement) {
				value, exists := oldNeighbours[neighbour]
				if !exists || value == 0 {
					newNeighbours = append(newNeighbours, neighbour)
				}
			}
			for _, neighbour := range graph.GetNeighbours(randomElement) {
				oldNeighbours[neighbour]++
			}
			newCounts := CountPolyominoes(graph, depth+1, maxSize, append(newUntriedSet, newNeighbours...), append(cellsAdded, randomElement), oldNeighbours)
			for i := range elementCount {
				elementCount[i] += newCounts[i]
			}
			for _, neighbour := range graph.GetNeighbours(randomElement) {
				oldNeighbours[neighbour]--
			}

		} else {
			//fmt.Printf("%v\n", append(cellsAdded, randomElement))
		}
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
	//fmt.Printf("Lattice: %v\n", latticeGraph)
	untriedSet := []Node{{X: 0, Y: 0, Z: 0}}
	oldNeighbours := make(map[Node]int)
	oldNeighbours[Node{X: 0, Y: 0, Z: 0}]++
	cellsAdded := make([]Node, 0, n)
	count := CountPolyominoes(latticeGraph, 0, n, untriedSet, cellsAdded, oldNeighbours)
	fmt.Println(count)
}
