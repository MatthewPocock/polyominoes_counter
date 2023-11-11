package main

type Node struct {
	X int
	Y int
}

type Graph struct {
	vertices map[Node]map[Node]bool
}

func (g *Graph) AddVertex(node Node) {
	if g.vertices[node] == nil {
		g.vertices[node] = make(map[Node]bool)
	}
}

func (g *Graph) AddEdge(node1, node2 Node) {
	g.vertices[node1][node2] = true
	g.vertices[node2][node1] = true
}

func (g *Graph) ConnectAdjacentNodes() {
	for key := range g.vertices {
		x, y := key.X, key.Y
		adjacent := []Node{{x + 1, y}, {x - 1, y}, {x, y + 1}, {x, y - 1}}
		for _, node := range adjacent {
			if _, exists := g.vertices[node]; exists {
				g.AddEdge(key, node)
			}
		}
	}
}

func (g *Graph) GetNeighbours(node Node) map[Node]bool {
	return g.vertices[node]
}

func NewGraph() *Graph {
	return &Graph{
		vertices: make(map[Node]map[Node]bool),
	}
}
