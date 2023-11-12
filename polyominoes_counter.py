import click
from collections import Counter
from typing import Tuple, List, Set


class Graph:
    def __init__(self):
        self.vertices = {}

    def __str__(self) -> str:
        return str(self.vertices)

    def add_vertex(self, node: Tuple[int, int]):
        self.vertices[node] = set()

    def add_edge(self, node_1: Tuple[int, int], node_2: Tuple[int, int]) -> None:
        self.vertices[node_1].add(node_2)
        self.vertices[node_2].add(node_1)

    def connect_adjacent_nodes(self) -> None:
        for key in self.vertices:
            x, y = key
            if self.vertices.get((x+1, y)) is not None:
                self.add_edge((x, y), (x+1, y))
            if self.vertices.get((x-1, y)) is not None:
                self.add_edge((x, y), (x-1, y))
            if self.vertices.get((x, y+1)) is not None:
                self.add_edge((x, y), (x, y+1))
            if self.vertices.get((x, y-1)) is not None:
                self.add_edge((x, y), (x, y-1))

    def get_neighbours(self, node: Tuple[int, int]) -> Set[Tuple[int, int]]:
        return self.vertices[node]


def generate_graph(size: int) -> Graph:
    """
    :param size: size of polyominoes to be calculated
    :return: A Graph where nodes represent lattice cells and edges represent adjacency.
    """
    lattice_graph = Graph()
    for x in range(0, size):
        lattice_graph.add_vertex((x, 0))
    for y in range(1, size):
        for x in range(-size + y + 1, size - y):
            lattice_graph.add_vertex((x, y))
    lattice_graph.connect_adjacent_nodes()
    return lattice_graph


def count_polyominoes(graph: Graph, depth: int, max_size: int, untried_set: List[Tuple[int, int]], cells_added: List[Tuple[int, int]]):
    """
    :param graph: square lattice - the graph on which we want to count the amount of sub-graphs up to given size
    :param depth: current recursion depth
    :param max_size: size of polyominoes to compute up to
    :param untried_set: list of elements that have not been in the square lattice
    :param cells_added: list of cells in square lattice
    :return: amount of fixed polyominoes up to give max_size
    """
    old_neighbours = set()
    if len(untried_set) != 0 and depth+1 < max_size:
        for cell in cells_added:
            old_neighbours.add(cell)
            old_neighbours.update(graph.get_neighbours(cell))
    element_count = Counter()
    while len(untried_set) != 0:
        random_element, *untried_set = untried_set  # Step 1
        cells_added.append(random_element)  # Step 2
        element_count[depth+1] += 1  # Step 3

        if depth + 1 < max_size:  # Step 4
            new_neighbours = []
            for neighbour in graph.get_neighbours(random_element):
                if neighbour not in old_neighbours:
                    new_neighbours.append(neighbour)

            untried_set.extend(new_neighbours)
            new_counts = count_polyominoes(graph, depth+1, max_size, untried_set, cells_added)
            element_count.update(new_counts)
            untried_set = [e for e in untried_set if e not in new_neighbours]

        cells_added.remove(random_element)  # Step 5
    return element_count


@click.command()
@click.argument("n", type=int)
def main(n):
    lattice_graph = generate_graph(n)
    count = count_polyominoes(lattice_graph, 0, n, [(0, 0)], [])
    print(count)
    return count


if __name__ == "__main__":
    main()
