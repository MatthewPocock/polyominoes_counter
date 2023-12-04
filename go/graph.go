package main

type Node struct {
	X int
	Y int
}

type Graph struct {
	vertices map[Node][]Node
}

func (g *Graph) AddVertex(node Node) {
	g.vertices[node] = []Node{}
}

func (g *Graph) AddEdge(node1, node2 Node) {
	g.vertices[node1] = append(g.vertices[node1], node2)
	//g.vertices[node2] = append(g.vertices[node2], node1)
}

func (g *Graph) ConnectAdjacentNodes() {
	for key := range g.vertices {
		x, y := key.X, key.Y
		adjacent := []Node{{X: x + 1, Y: y}, {X: x - 1, Y: y}, {X: x, Y: y + 1}, {X: x, Y: y - 1}}
		for _, node := range adjacent {
			if _, exists := g.vertices[node]; exists {
				g.AddEdge(key, node)
			}
		}
	}
}

func (g *Graph) GetNeighbours(node Node) []Node {
	return g.vertices[node]
}

func NewGraph() *Graph {
	return &Graph{
		vertices: make(map[Node][]Node),
	}
}
