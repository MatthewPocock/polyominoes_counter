package main

import (
	"fmt"
	"os"
	"runtime/pprof"
	"strconv"
	"sync"
	"time"
)

const branchDepth = 6

func CreateLattice(n int) *Graph {
	latticeGraph := NewGraph()

	for x := 0; x < n; x++ {
		node := Node{X: x, Y: 0}
		latticeGraph.AddVertex(node)
	}
	for y := 1; y < n; y++ {
		for x := -n + y + 1; x < n-y; x++ {
			node := Node{X: x, Y: y}
			latticeGraph.AddVertex(node)
		}
	}
	latticeGraph.ConnectAdjacentNodes()
	return latticeGraph
}

func copyMap(originalMap map[Node]int) map[Node]int {
	newMap := make(map[Node]int)
	for key, value := range originalMap {
		newMap[key] = value
	}
	return newMap
}

func compareCodes(slice1, slice2 []bool) int {
	// returns -1 when slice 1 is less than slice 2
	for i := 0; i < len(slice1); i++ {
		if slice1[i] == true && slice2[i] != true {
			return -1
		} else if slice2[i] == true && slice1[i] != true {
			return 1
		}
	}
	return 0
}

func rotateNodes(nodes []Node) {
	for i := range nodes {
		nodes[i].X, nodes[i].Y = nodes[i].Y, -nodes[i].X
	}
}

func getCode(nodes []Node) []bool {
	minW, maxW, maxH, minH := 0, 0, 0, 0
	for _, node := range nodes {
		if node.X < minW {
			minW = node.X
		}
		if node.X > maxW {
			maxW = node.X
		}
		if node.Y > maxH {
			maxH = node.Y
		}
		if node.Y < minH {
			minH = node.Y
		}
	}

	squareSize := max(maxW-minW+1, maxH-minH+1)
	code := make([]bool, squareSize*squareSize)
	for _, node := range nodes {
		index := node.X - minW + ((maxH - node.Y) * squareSize)
		code[index] = true
	}
	return code
}

func isCanonical(nodes []Node) bool {
	// get main code
	code := getCode(nodes)
	// get other codes
	rotatedNodes := make([]Node, len(nodes))
	copy(rotatedNodes, nodes)

	for i := 0; i < 3; i++ {
		//minWidth, maxHeight, maxWidth, minHeight = minHeight, -minWidth, maxHeight, -maxWidth
		rotateNodes(rotatedNodes)
		if compareCodes(code, getCode(rotatedNodes)) == 1 {
			return false
		}
	}
	return true
}

// CountPolyominoes implements Redelmeier's algorithm to enumerate fixed polycubes w/ an additional isCanonical check
// to determine if polycube is unique when considering rotations.
func CountPolyominoes(
	graph *Graph,
	depth int,
	maxSize int,
	untriedSet []Node,
	cellsAdded []Node,
	oldNeighbours map[Node]int,
	ch chan []int,
	wg *sync.WaitGroup,
) []int {
	newUntriedSet := make([]Node, len(untriedSet))
	copy(newUntriedSet, untriedSet)

	elementCount := make([]int, maxSize)
	for len(newUntriedSet) != 0 {
		randomElement := newUntriedSet[len(newUntriedSet)-1] // step 1
		newUntriedSet = newUntriedSet[:len(newUntriedSet)-1] // step 2

		if isCanonical(append(cellsAdded, randomElement)) {
			elementCount[depth]++ // Step 3
		}
		//elementCount[depth]++ // Step 3

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
			if depth == branchDepth { // parallelDepth is the depth at which to start parallelization
				wg.Add(1)
				newOldNeighbours := copyMap(oldNeighbours)
				untriedSetCopy := make([]Node, len(newUntriedSet), len(newUntriedSet)+len(newNeighbours))
				copy(untriedSetCopy, newUntriedSet)
				cellsAddedCopy := make([]Node, len(cellsAdded), len(cellsAdded)+1)
				copy(cellsAddedCopy, cellsAdded)
				go CountPolyominoes(graph, depth+1, maxSize, append(untriedSetCopy, newNeighbours...), append(cellsAddedCopy, randomElement), newOldNeighbours, ch, wg)
			} else {
				newCounts := CountPolyominoes(graph, depth+1, maxSize, append(newUntriedSet, newNeighbours...), append(cellsAdded, randomElement), oldNeighbours, ch, wg)
				for i := range elementCount {
					elementCount[i] += newCounts[i]
				}
			}
			for _, neighbour := range graph.GetNeighbours(randomElement) {
				oldNeighbours[neighbour]--
			}

		} else {
			//fmt.Printf("%v\n", append(cellsAdded, randomElement))
		}
	}
	if depth == branchDepth+1 {
		ch <- elementCount
		wg.Done()
	}
	if depth == 0 {
		go func() {
			wg.Wait()
			close(ch)
		}()
		for result := range ch {
			for i := range elementCount {
				elementCount[i] += result[i]
			}
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
	startTime := time.Now()

	latticeGraph := CreateLattice(n)
	//fmt.Printf("Lattice: %v\n", latticeGraph)

	untriedSet := []Node{{X: 0, Y: 0}}
	oldNeighbours := make(map[Node]int)
	oldNeighbours[Node{X: 0, Y: 0}]++
	cellsAdded := make([]Node, 0, n)

	ch := make(chan []int)
	var wg sync.WaitGroup

	count := CountPolyominoes(latticeGraph, 0, n, untriedSet, cellsAdded, oldNeighbours, ch, &wg)
	fmt.Println(count)
	fmt.Println("Completed in:", time.Since(startTime))
}
