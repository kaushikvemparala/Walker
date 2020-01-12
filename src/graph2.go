package main

import "fmt"

type Node2 struct {
	id       int
	read     string
	prefkey  string
	suffkey  string
	innodes  []Node2
	outnodes []Node2
	visited  bool
}

type overlap struct {
	preNode       Node2
	sufNode       Node2
	overlapString string
}

type Graph2 struct {
	Nodes               map[int]Node2
	Edges               []*overlap
	connectedComponents [][]Node2
	connectCompCount    int
}

func MakeGraph2() Graph2 {
	return Graph2{make(map[int]Node2), make([]*overlap, 0), make([][]Node2, 0), 0}
}

func (graph *Graph2) addNode2(id int, read string) {
	graph.Nodes[id] = Node2{id, read, "", "", make([]Node2, 0), make([]Node2, 0), false}
}

func (graph *Graph2) addEdge2(id1, id2 int, overlapString string) {
	overlap1 := overlap{preNode: graph.Nodes[id1], sufNode: graph.Nodes[id2], overlapString: overlapString}
	graph.Edges = append(graph.Edges, &overlap1)
	pointerToNode := graph.Nodes[id1]
	fmt.Println("outnode before:", len(pointerToNode.outnodes))
	(&pointerToNode).addOutNode(graph.Nodes[id2])
	fmt.Println("outnode after:", len(pointerToNode.outnodes))
	pointerToNode2 := graph.Nodes[id2]
	fmt.Println("innode before:", len(pointerToNode2.innodes))
	(&pointerToNode2).addInNode(graph.Nodes[id1])
	fmt.Println("innode after:", len(pointerToNode2.innodes))
	graph.Nodes[id1] = pointerToNode
	graph.Nodes[id2] = pointerToNode2
}

func (graph *Graph2) removeNode(id int) {
	delete(graph.Nodes, id)
}

func (graph Graph2) PrintGraph2() {

	for _, node := range graph.Nodes {
		fmt.Println("node id", node.id)
		fmt.Println("		number of innodes:", len(node.innodes))
		fmt.Println("		number of outnodes:", len(node.outnodes))
	}

	fmt.Println(len(graph.Edges))
}

func (node Node2) getID2() int {
	return node.id
}

func (node *Node2) setprefkey(key string) {
	node.prefkey = key
}

func (node *Node2) setsuffkey(key string) {
	node.suffkey = key
}

func (node Node2) getprefkey() string {
	return node.prefkey
}

func (node Node2) getsuffkey() string {
	return node.suffkey
}

/*
func (graph Graph2) getNonBranchingPaths() [][]int {
	paths := make([][]int, 0)
	for node := range graph.Nodes {
		if !(node.isOnetoOne()) {
			if len(node.outnodes) > 0 {
			}
		}
	}
	return paths
}
*/

func (node Node2) isOnetoOne() bool {
	if (len(node.innodes) == 1) && (len(node.outnodes) == 1) {
		return true
	} else {
		return false
	}

}

func (node *Node2) addInNode(node2 Node2) {
	node.innodes = append(node.innodes, node2)
}

func (node *Node2) addOutNode(node2 Node2) {
	node.outnodes = append(node.outnodes, node2)
}

func (graph *Graph2) dfs(node Node2, connectedComp []Node2) {
	fmt.Println("entering dfs")
	fmt.Println("innodes:", len(node.innodes))
	fmt.Println("outnodes:", len(node.outnodes))
	fmt.Println("node of dfs:", node.id)
	//if node.visited {
	//fmt.Println("visited")
	//} else {
	//fmt.Println("not visited")
	//}
	//fmt.Println("how many accumulatd components : ",len(connectedComp))
	//fmt.Println("How many neighbors:", len(node.edges))
	//fmt.Println("accumulated comps:", len(connectedComp))
	if len(node.outnodes) == 0 {
		fmt.Println("Found a connected Component")
		//fmt.Print("{ ")
		//for i := range connectedComp {
		//	fmt.Print(connectedComp[i].id, " ")
		//}
		//fmt.Println("}")
		graph.addConnectedComp(connectedComp)
	}
	for _, outnode := range node.outnodes {
		fmt.Println("edge found:", outnode.id)
		if !(graph.Nodes[outnode.id].isInConnectedComp(connectedComp)) {
			node := graph.Nodes[outnode.id]
			node.visited = true
			graph.Nodes[outnode.id] = node

			//fmt.Println("node that is about to be appended:", graph.Nodes[id].id)
			connectedComp = append(connectedComp, graph.Nodes[outnode.id])
			nodes := graph.Nodes
			delete(nodes, outnode.id)
			graph.Nodes = nodes
			//fmt.Println("accumulated comps:", len(connectedComp))

			graph.dfs(graph.Nodes[outnode.id], connectedComp)
		} else if !(isAConnectedComp2(graph, connectedComp)) {
			fmt.Println("Found a connected Component")
			graph.addConnectedComp(connectedComp)
		}
	}
}

func (node Node2) isInConnectedComp(connComp []Node2) bool {
	for _, neigh := range connComp {
		if node.id == neigh.id {
			return true
		}
	}
	return false
}

func isAConnectedComp2(graph *Graph2, connectedComp []Node2) bool {
	for _, comps := range graph.connectedComponents {
		if comps[0].id == connectedComp[0].id {
			return true
		}
	}
	return false
}
func (node Node2) isANeighbor(node1 Node2) bool {
	for _, neighbor := range node.innodes {
		if neighbor.id == node1.id {
			return true
		}
	}
	for _, neighbor := range node.outnodes {
		if neighbor.id == node1.id {
			return true
		}
	}
	return false
}

func (graph *Graph2) FindConnectedComponents() {
	for i := range graph.Nodes {
		//node := graph.Nodes[i]
		//node.visited = true
		//graph.Nodes[i] = node
		if !graph.Nodes[i].visited {

			node := graph.Nodes[i]
			node.visited = true
			graph.Nodes[i] = node

			connectedComp := make([]Node2, 0)
			connectedComp = append(connectedComp, graph.Nodes[i])

			temp := connectedComp
			graph.dfs(graph.Nodes[i], temp)
		}
	}
}

func (graph Graph2) getConnectedComponents() [][]Node2 {
	return graph.connectedComponents
}

func (graph *Graph2) addConnectedComp(cc []Node2) {
	//fmt.Println("length of connected comps before adding", len(graph.getConnectedComponents()))
	graph.connectedComponents = append(graph.connectedComponents, cc)
	graph.connectCompCount++
	//fmt.Println("length of connected comps after adding", len(graph.getConnectedComponents()))
}

func (graph *Graph2) removeConnectedComp(index int) {
	//fmt.Println("length of connected components before removal:", len(graph.connectedComponents))
	newSlice := make([][]Node2, 0)
	for i := range graph.connectedComponents {
		if i != index {
			newSlice = append(newSlice, graph.connectedComponents[i])
		}
	}
	graph.connectedComponents = newSlice
	graph.connectCompCount--
	//fmt.Println("length of connected components after removal:", len(graph.connectedComponents))
}

func (graph *Graph2) PruneConnectedComps() {
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
	fmt.Println("number of connencted components:", len(graph.connectedComponents))
}

func Graph2Statistics(graph Graph2) {
	totalNodes := make([]int, 0)
	for _, con := range graph.connectedComponents {
		numNodes := 0
		for i := 0; i < len(con); i++ {
			numNodes++
		}
		totalNodes = append(totalNodes, numNodes)
	}
	fmt.Println("total connected comps:", len(totalNodes))
	tot := 0
	for _, lens := range totalNodes {
		tot += lens
	}
	avgLen := tot / (len(totalNodes))
	fmt.Println("avg connected comp size:", avgLen)
}
