package main

import (
	"fmt"
	"os"
)

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
	LNBPs               [][]Node2
}

func MakeGraph2() Graph2 {
	return Graph2{make(map[int]Node2), make([]*overlap, 0), make([][]Node2, 0), 0, make([][]Node2, 0)}
}

func (graph *Graph2) addNode2(id int, read string) {
	graph.Nodes[id] = Node2{id, read, "", "", make([]Node2, 0), make([]Node2, 0), false}
}

func (graph *Graph2) addEdge2(id1, id2 int, overlapString string) {
	overlap1 := overlap{preNode: graph.Nodes[id1], sufNode: graph.Nodes[id2], overlapString: overlapString}
	graph.Edges = append(graph.Edges, &overlap1)
	pointerToNode := graph.Nodes[id1]
	//fmt.Println("outnode before:", len(pointerToNode.outnodes))
	(&pointerToNode).addOutNode(graph.Nodes[id2])
	//fmt.Println("outnode after:", len(pointerToNode.outnodes))
	pointerToNode2 := graph.Nodes[id2]
	//fmt.Println("innode before:", len(pointerToNode2.innodes))
	(&pointerToNode2).addInNode(graph.Nodes[id1])
	//fmt.Println("innode after:", len(pointerToNode2.innodes))
	graph.Nodes[id1] = pointerToNode
	graph.Nodes[id2] = pointerToNode2
}

func (graph *Graph2) removeNode(id int) {
	delete(graph.Nodes, id)
}

func (graph Graph2) PrintGraph2(logFile *os.File) {

	for _, node := range graph.Nodes {
		fmt.Fprintln(logFile, "node id", node.id)
		fmt.Fprintln(logFile, "		number of innodes:", len(node.innodes))
		fmt.Fprintln(logFile, "		number of outnodes:", len(node.outnodes))
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

func (graph Graph2) getNonBranchingPaths() [][]Node2 {
	paths := make([][]Node2, 0)
	for _, node := range graph.Nodes {
		//if !(node.isOnetoOne()) {
		if len(node.outnodes) > 0 {
			for _, outnode := range node.outnodes {
				nonBranchingPath := make([]Node2, 0)
				nonBranchingPath = append(nonBranchingPath, node)
				nonBranchingPath = append(nonBranchingPath, outnode)
				for outnode.isOnetoOne() {
					fmt.Println(outnode.id)
					outnode2 := outnode.outnodes[0]
					nonBranchingPath = append(nonBranchingPath, outnode2)
					outnode = outnode2
				}
				paths = append(paths, nonBranchingPath)
			}
		}
		//}
	}
	return paths
}

func (node Node2) isOnetoOne() bool {
	if (len(node.innodes) == 1) && (len(node.outnodes) == 1) {
		return true
	}
	return false
}

func (node *Node2) addInNode(node2 Node2) {
	node.innodes = append(node.innodes, node2)
}

func (node *Node2) addOutNode(node2 Node2) {
	node.outnodes = append(node.outnodes, node2)
}

func (graph *Graph2) dfs(node Node2, connectedComp []Node2, logFile *os.File) {
	fmt.Fprintln(logFile, "entering dfs")
	fmt.Fprintln(logFile, "innodes:", len(node.innodes))
	fmt.Fprintln(logFile, "outnodes:", len(node.outnodes))
	fmt.Fprintln(logFile, "node of dfs:", node.id)
	//if node.visited {
	//fmt.Println("visited")
	//} else {
	//fmt.Println("not visited")
	//}
	//fmt.Println("how many accumulatd components : ",len(connectedComp))
	//fmt.Println("How many neighbors:", len(node.edges))
	//fmt.Println("accumulated comps:", len(connectedComp))
	if len(node.outnodes) == 0 {
		fmt.Fprintln(logFile, "Found a connected Component")
		//fmt.Print("{ ")
		//for i := range connectedComp {
		//	fmt.Print(connectedComp[i].id, " ")
		//}
		//fmt.Println("}")
		graph.addConnectedComp(connectedComp)
		//graph.FindLongestPath(connectedComp, logFile)
	}
	for _, outnode := range node.outnodes {
		fmt.Fprintln(logFile, "edge found:", outnode.id)
		if !(graph.Nodes[outnode.id].isInConnectedComp(connectedComp)) {
			node := graph.Nodes[outnode.id]
			node.visited = true
			graph.Nodes[outnode.id] = node

			//fmt.Println("node that is about to be appended:", graph.Nodes[id].id)
			connectedComp = append(connectedComp, graph.Nodes[outnode.id])
			//nodes := graph.Nodes
			//delete(nodes, outnode.id)
			//graph.Nodes = nodes
			fmt.Fprintln(logFile, "accumulated comps:", len(connectedComp))

			graph.dfs(graph.Nodes[outnode.id], connectedComp, logFile)
		} else if !(isAConnectedComp2(graph, connectedComp)) {
			fmt.Fprintln(logFile, "Found a connected Component")
			graph.addConnectedComp(connectedComp)
			//graph.FindLongestPath(connectedComp, logFile)
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

func (node Node2) isInConnectedCompPoint(connComp *[]Node2) bool {
	for _, neigh := range *connComp {
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

func (graph *Graph2) FindConnectedComponents(logFile *os.File, summaryFile *os.File) {
	//nodesMaster := graph.Nodes
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
			graph.dfs(graph.Nodes[i], temp, logFile)
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

func (graph *Graph2) PruneConnectedComps(logFile *os.File, summaryFile *os.File) {
	fmt.Println("pruning cc's")
	fmt.Println("length of connectedComponents:", len(graph.connectedComponents))
	for i := len(graph.connectedComponents) - 1; i >= 0; i-- {
		fmt.Fprintln(logFile, "Entered the", i, "for loop")
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
	fmt.Fprintln(logFile, "number of connencted components:", len(graph.connectedComponents))
}

func Graph2Statistics(graph Graph2, summaryFile *os.File) {
	totalNodes := make([]int, 0)
	for _, con := range graph.connectedComponents {
		numNodes := 0
		for i := 0; i < len(con); i++ {
			numNodes++
		}
		totalNodes = append(totalNodes, numNodes)
	}
	fmt.Fprintln(summaryFile, "\t\ttotal connected comps:", len(totalNodes))
	fmt.Fprintln(summaryFile, "")
	tot := 0
	for _, lens := range totalNodes {
		tot += lens
	}
	avgLen := tot / (len(totalNodes))
	min := len(graph.connectedComponents[0])
	max := len(graph.connectedComponents[0])
	for _, cc := range graph.connectedComponents {
		if len(cc) < min {
			min = len(cc)
		}
		if len(cc) > max {
			max = len(cc)
		}
	}
	fmt.Fprintln(summaryFile, "\t\tmax cc size:", max)
	fmt.Fprintln(summaryFile, "")
	fmt.Fprintln(summaryFile, "\t\tmin cc size:", min)
	fmt.Fprintln(summaryFile, "")
	fmt.Fprintln(summaryFile, "\t\tavg connected comp size:", avgLen)
	fmt.Fprintln(summaryFile, "")
	fmt.Fprintln(summaryFile, "\t\tnumber of edges:", len(graph.Edges))
	fmt.Fprintln(summaryFile, "")
	fmt.Fprintln(summaryFile, "\t\tnumber of lnbps:", len(graph.LNBPs))
	fmt.Fprintln(summaryFile, "")
	for _, lnbp := range graph.LNBPs {
		fmt.Fprint(summaryFile, "[ ")
		for _, node := range lnbp {
			fmt.Fprint(summaryFile, node.id, " ")
		}
		fmt.Fprintln(summaryFile, "]")
	}
	//for i, edge := range graph.Edges {
	//	fmt.Fprintln(summaryFile, "edge ",i,": From node ",edge.preNode.id," to edge ", edge.sufNode.id)
	//	fmt.Fprintln(summaryFile, "")
	//}
}

func (graph *Graph2) removeNodeNBP(nbp *[]Node2, index int) {
	//fmt.Println("length of connected components before removal:", len(graph.connectedComponents))
	newSlice := make([]Node2, 0)
	for i := range *nbp {
		if i != index {
			newSlice = append(newSlice, (*nbp)[i])
		}
	}
	nbp = &newSlice
	//fmt.Println("length of connected components after removal:", len(graph.connectedComponents))
}

func (graph *Graph2) maxNBP(nbps [][]Node2) []Node2 {
	lnbp := make([]Node2, 0)
	for _, nbp := range nbps {
		if len(nbp) > len(lnbp) {
			lnbp = nbp
		}
	}
	return lnbp
}

func (graph *Graph2) FindLongestPath(conComp []Node2, logFile *os.File) {
	longestPath := make([]Node2, 0)
	pathsNBP := make([][]Node2, 0)
	for _, node := range conComp {
		tempPath := make([]Node2, 0)
		tempPath = append(tempPath, graph.Nodes[node.id])
		graph.dfs1(node, &tempPath, logFile, &pathsNBP)
	}
	longestPath = graph.maxNBP(pathsNBP)
	fmt.Fprintln(logFile, "length of longest path", len(longestPath))
	graph.addLNBP(longestPath)
}

func (graph *Graph2) dfs1(node Node2, nbp *[]Node2, logFile *os.File, pathsNBP *[][]Node2) {
	fmt.Fprintln(logFile, "entering dfs")
	fmt.Fprintln(logFile, "innodes:", len(node.innodes))
	fmt.Fprintln(logFile, "outnodes:", len(node.outnodes))
	fmt.Fprintln(logFile, "node of dfs:", node.id)
	//if node.visited {
	//fmt.Println("visited")
	//} else {
	//fmt.Println("not visited")
	//}
	//fmt.Println("how many accumulatd components : ",len(connectedComp))
	//fmt.Println("How many neighbors:", len(node.edges))
	//fmt.Println("accumulated comps:", len(connectedComp))
	if len(node.outnodes) == 0 {
		fmt.Fprintln(logFile, "path found")
		fmt.Fprintln(logFile, "done with", node.id)
		// clone the nbp array  , save it to a global data structure
		pathsNBPtemp := *pathsNBP
		pathsNBPtemp = append(pathsNBPtemp, *nbp)
		pathsNBP = &pathsNBPtemp
		graph.removeNodeNBP(nbp, len(*nbp)-1)
		return
	}
	for _, outnode := range node.outnodes {
		fmt.Fprintln(logFile, "edge found:", outnode.id)
		inCC := graph.Nodes[outnode.id].isInConnectedCompPoint(nbp)
		fmt.Fprintln(logFile, "is it in the cc?", graph.Nodes[outnode.id].isInConnectedCompPoint(nbp))

		if inCC {
			fmt.Fprintln(logFile, "done with", node.id)
			pathsNBPtemp := *pathsNBP
			pathsNBPtemp = append(pathsNBPtemp, *nbp)
			pathsNBP = &pathsNBPtemp
			graph.removeNodeNBP(nbp, len(*nbp)-1)
			return
		} else {
			//fmt.Println("node that is about to be appended:", graph.Nodes[id].id)
			nbpTemp := *nbp
			nbpTemp = append(nbpTemp, graph.Nodes[outnode.id])
			nbp = &nbpTemp
			fmt.Fprintln(logFile, "accumulated comps:", len(*nbp))

			graph.dfs1(graph.Nodes[outnode.id], nbp, logFile, pathsNBP)
		}
	}
	fmt.Fprintln(logFile, "done with", node.id)
}

func (graph *Graph2) addLNBP(lnbp []Node2) {
	graph.LNBPs = append(graph.LNBPs, lnbp)
}

func (graph *Graph2) LNBPfinder(logFile *os.File) {
	for _, cc := range graph.connectedComponents {
		graph.FindLongestPath(cc, logFile)
	}
}
