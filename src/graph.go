package main

import "fmt"

// Node is a node in the graph; it has a (unique) ID and a sequence of
// edges to other nodes.
type Node struct {
	id      int
	read    string
	edges   []int
	visited bool
}

// Graph contains a set of Nodes, uniquely identified by numeric IDs.
type Graph struct {
	Nodes               map[int]Node
	connectedComponents [][]Node
	connectCompCount    int
}

// MakeGraph makes an empty Graph
func MakeGraph() Graph {
	return Graph{make(map[int]Node), make([][]Node, 0), 0}
}

func (graph *Graph) addNode(id int, read string) {
	graph.Nodes[id] = Node{id, read, make([]int, 0), false}
}

func (graph *Graph) addEdge(id, id2 int) {
	var node = graph.Nodes[id2]
	node.edges = append(node.edges, id)
	graph.Nodes[id2] = node
}

func (node Node) getID() int {
	return node.id
}

func (graph *Graph) removeNode(id int) {
	delete(graph.Nodes, id)
}

func (graph Graph) PrintGraph() {
	for i := range graph.Nodes {
		fmt.Println(i, ": ")
		//fmt.Println("Read: ", graph.Nodes[i].read)
		fmt.Print("Neighbors: ")
		for j := range graph.Nodes[i].edges {
			fmt.Println(graph.Nodes[i].edges[j])
		}
		fmt.Println()
		fmt.Println()

	}
	fmt.Println("CC's: ")

	for i := range graph.connectedComponents {
		fmt.Print("{ ")
		for j := range graph.connectedComponents[i] {
			fmt.Print(graph.connectedComponents[i][j].id, " ")
		}
		fmt.Println("}")
	}

	fmt.Println(graph.connectCompCount)

}

func (graph *Graph) dfs(node Node, connectedComp []Node) {
	fmt.Println("entering dfs")
	fmt.Println("node of dfs:", node.id)
	if node.visited {
		//fmt.Println("visited")
	} else {
		//fmt.Println("not visited")
	}
	//fmt.Println("how many accumulatd components : ",len(connectedComp))
	//fmt.Println("How many neighbors:", len(node.edges))
	//fmt.Println("accumulated comps:", len(connectedComp))
	if len(node.edges) == 0 {
		fmt.Println("Found a connected Component")
		//fmt.Print("{ ")
		//for i := range connectedComp {
		//	fmt.Print(connectedComp[i].id, " ")
		//}
		//fmt.Println("}")
		graph.addConnectedComp(connectedComp)
		return
	}

	for _, id := range node.edges {
		fmt.Println("edge found:", id)
		if !(graph.Nodes[id].isInConnectedComp(connectedComp)) {
			node := graph.Nodes[id]
			node.visited = true
			graph.Nodes[id] = node

			//fmt.Println("node that is about to be appended:", graph.Nodes[id].id)
			connectedComp = append(connectedComp, graph.Nodes[id])
			//fmt.Println("accumulated comps:", len(connectedComp))

			graph.dfs(graph.Nodes[id], connectedComp)
		} else if !(isAConnectedComp(graph, connectedComp)) {
			fmt.Println("Found a connected Component")
			graph.addConnectedComp(connectedComp)
		}
	}
}

func (node Node) isInConnectedComp(connComp []Node) bool {
	for _, neigh := range connComp {
		if node.id == neigh.id {
			return true
		}
	}
	return false
}

func isAConnectedComp(graph *Graph, connectedComp []Node) bool {
	for _, comps := range graph.connectedComponents {
		if comps[0].id == connectedComp[0].id {
			return true
		}
	}
	return false
}
func (node Node) isANeighbor(node1 Node) bool {
	for _, neighbor := range node.edges {
		if neighbor == node1.id {
			return true
		}
	}
	return false
}

func (graph *Graph) FindConnectedComponents() {

	for i := range graph.Nodes {
		if graph.Nodes[i].visited {
			//fmt.Println("visited")
		} else {
			//fmt.Println("not visited")
		}
		//fmt.Println()
	}
	for i := range graph.Nodes {
		//node := graph.Nodes[i]
		//node.visited = true
		//graph.Nodes[i] = node
		if !graph.Nodes[i].visited {

			node := graph.Nodes[i]
			node.visited = true
			graph.Nodes[i] = node

			connectedComp := make([]Node, 0)
			connectedComp = append(connectedComp, graph.Nodes[i])

			temp := connectedComp
			graph.dfs(graph.Nodes[i], temp)
		}
	}
}

func (graph Graph) getConnectedComponents() [][]Node {
	return graph.connectedComponents
}

func (graph *Graph) addConnectedComp(cc []Node) {
	//fmt.Println("length of connected comps before adding", len(graph.getConnectedComponents()))
	graph.connectedComponents = append(graph.connectedComponents, cc)
	graph.connectCompCount++
	//fmt.Println("length of connected comps after adding", len(graph.getConnectedComponents()))
}

func (graph *Graph) removeConnectedComp(index int) {
	//fmt.Println("length of connected components before removal:", len(graph.connectedComponents))
	newSlice := make([][]Node, 0)
	for i := range graph.connectedComponents {
		if i != index {
			newSlice = append(newSlice, graph.connectedComponents[i])
		}
	}
	graph.connectedComponents = newSlice
	graph.connectCompCount--
	//fmt.Println("length of connected components after removal:", len(graph.connectedComponents))
}

func (graph *Graph) PruneConnectedComps() {
	fmt.Println("pruning cc's")
	fmt.Println("length of connectedComponents:", len(graph.connectedComponents))
	for i := len(graph.connectedComponents) - 1; i >= 0; i-- {
		fmt.Println("Entered the", i, "for loop")
		//fmt.Println("i:", i)
		for j := i - 1; j >= 0; j-- {
			//fmt.Println("j:", j)
			if len(graph.connectedComponents[i]) == len(graph.connectedComponents[j]) {
				id1 := graph.connectedComponents[j][0].id
				id2 := graph.connectedComponents[i][0].id
				if id1 == id2 {
					graph.removeConnectedComp(j)
					i--
				}
			} else if len(graph.connectedComponents[i]) > len(graph.connectedComponents[j]) {
				id1 := graph.connectedComponents[j][0].id
				//fmt.Println(id1)
				for _, node := range graph.connectedComponents[i] {
					if node.id == id1 {
						//fmt.Println("removing connected comp and i > j")
						graph.removeConnectedComp(j)
						i--
					}
				}
			} else if len(graph.connectedComponents[j]) > len(graph.connectedComponents[i]) {
				id1 := graph.connectedComponents[i][0].id
				//fmt.Println(id1)
				for _, node := range graph.connectedComponents[j] {
					if node.id == id1 {
						//fmt.Println("removing connected comp and i < j")
						graph.removeConnectedComp(i)
						i--
					}
				}
			}
		}
	}
}
