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
	connectedComponents []connectedComponent
	connectCompCount    int
	LNBPs               [][]Node2
}

type connectedComponent struct {
	nodes []Node2
	nbps  [][]Node2
}

func MakeGraph2() Graph2 {
	return Graph2{make(map[int]Node2), make([]*overlap, 0), make([]connectedComponent, 0), 0, make([][]Node2, 0)}
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

func (graph Graph2) getNonBranchingPaths(logFile *os.File) [][]Node2 {
	paths := make([][]Node2, 0)
	for _, node := range graph.Nodes {
		//if !(node.isOnetoOne()) {
			if len(node.outnodes) > 0 {
				for _, outnode := range node.outnodes {
					nonBranchingPath := make([]Node2, 0)
					nonBranchingPath = append(nonBranchingPath, node)
					fmt.Fprintln(logFile, "added ", node.id)
					nonBranchingPath = append(nonBranchingPath, outnode)
					fmt.Fprintln(logFile, "added ", outnode.id)
					for (len(outnode.outnodes) > 0) {
						fmt.Fprintln(logFile, "outnode has outedges")
						fmt.Println(outnode.id)
						outnode2 := outnode.outnodes[0]
						nonBranchingPath = append(nonBranchingPath, outnode2)
						outnode = outnode2
					}
					//for outnode.isOnetoOne() {
					//	fmt.Fprintln(logFile, "is not 1 to 1")
					//}
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

func (graph *Graph2) dfs(node Node2, connectedComp *connectedComponent, logFile *os.File) {
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
		//graph.addConnectedComp(connectedComp)
		return
		//graph.FindLongestPath(connectedComp, logFile)
	}
	for _, outnode := range node.outnodes {
		fmt.Fprintln(logFile, "edge found:", outnode.id)
		if !(graph.Nodes[outnode.id].isInConnectedComp(connectedComp)) {
			node := graph.Nodes[outnode.id]
			node.visited = true
			graph.Nodes[outnode.id] = node

			//fmt.Println("node that is about to be appended:", graph.Nodes[id].id)
			connectedComp.nodes = append(connectedComp.nodes, graph.Nodes[outnode.id])
			fmt.Fprintln(logFile, "node added")
			//nodes := graph.Nodes
			//delete(nodes, outnode.id)
			//graph.Nodes = nodes
			fmt.Fprintln(logFile, "accumulated comps:", len(connectedComp.nodes))

			graph.dfs(graph.Nodes[outnode.id], connectedComp, logFile)
		} else if !(isAConnectedComp2(graph, connectedComp)) {
			fmt.Fprintln(logFile, "Found a connected Component")
			//graph.addConnectedComp(connectedComp)
			return
			//graph.FindLongestPath(connectedComp, logFile)
		}
	}
}

func (node Node2) isInConnectedComp(connComp *connectedComponent) bool {
	for _, neigh := range (*connComp).nodes {
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

func isAConnectedComp2(graph *Graph2, connectedComp *connectedComponent) bool {
	for _, comps := range graph.connectedComponents {
		if comps.nodes[0].id == (*connectedComp).nodes[0].id {
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

			connectedComp := connectedComponent{make([]Node2, 0), make([][]Node2, 0)}
			connectedComp.nodes = append(connectedComp.nodes, graph.Nodes[i])

			temp := connectedComp
			graph.dfs(graph.Nodes[i], &temp, logFile)
			connectedComp = temp
			graph.addConnectedComp(connectedComp)
		}
	}
}

func (graph Graph2) getConnectedComponents() []connectedComponent {
	return graph.connectedComponents
}

func (graph *Graph2) addConnectedComp(cc connectedComponent) {
	//fmt.Println("length of connected comps before adding", len(graph.getConnectedComponents()))
	graph.connectedComponents = append(graph.connectedComponents, cc)
	graph.connectCompCount++
	//fmt.Println("length of connected comps after adding", len(graph.getConnectedComponents()))
}

func (graph *Graph2) removeConnectedComp(index int) {
	//fmt.Println("length of connected components before removal:", len(graph.connectedComponents))
	newSlice := make([]connectedComponent, 0)
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
			if len(graph.connectedComponents[i].nodes) == len(graph.connectedComponents[j].nodes) {
				id1 := graph.connectedComponents[j].nodes[0].id
				id2 := graph.connectedComponents[i].nodes[0].id
				if id1 == id2 {
					graph.removeConnectedComp(j)
					i--
				}
			} else if len(graph.connectedComponents[i].nodes) > len(graph.connectedComponents[j].nodes) {
				id1 := graph.connectedComponents[j].nodes[0].id
				//fmt.Println(id1)
				for _, node := range graph.connectedComponents[i].nodes {
					if node.id == id1 {
						//fmt.Println("removing connected comp and i > j")
						graph.removeConnectedComp(j)
						i--
					}
				}
			} else if len(graph.connectedComponents[j].nodes) > len(graph.connectedComponents[i].nodes) {
				id1 := graph.connectedComponents[i].nodes[0].id
				//fmt.Println(id1)
				for _, node := range graph.connectedComponents[j].nodes {
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
		for i := 0; i < len(con.nodes); i++ {
			numNodes++
		}
		totalNodes = append(totalNodes, numNodes)
	}
	fmt.Fprintln(summaryFile, "\t\ttotal connected comps:", len(graph.connectedComponents))
	fmt.Fprintln(summaryFile, "")
	tot := 0
	for _, lens := range totalNodes {
		tot += lens
	}
	avgLen := float64(tot) / float64(len(totalNodes))
	min := len(graph.connectedComponents[0].nodes)
	max := len(graph.connectedComponents[0].nodes)
	for _, cc := range graph.connectedComponents {
		if len(cc.nodes) < min {
			min = len(cc.nodes)
		}
		if len(cc.nodes) > max {
			max = len(cc.nodes)
		}
	}
	pathNodes := make([]int, 0)
	for _, path := range graph.LNBPs {
		numNodes := 0
		for i := 0; i < len(path); i++ {
			numNodes++
		}
		pathNodes = append(pathNodes, numNodes)
	}
	totpathNodes := 0
	for _, lens := range totalNodes {
		tot += lens
	}
	avgLenPath := 0
	minPath := 0
	maxPath := 0
	if (len(pathNodes) != 0) && (len(graph.LNBPs) != 0) {
		avgLenPath = totpathNodes / (len(pathNodes))
		minPath = len(graph.LNBPs[0])
		maxPath = len(graph.LNBPs[0])
	}
	for _, path := range graph.LNBPs {
		length := len(path)
		if length < minPath {
			min = length
		}
		if length > maxPath {
			maxPath = length
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
	fmt.Fprintln(summaryFile, "\t\taverage lnbp size:", avgLenPath)
	fmt.Fprintln(summaryFile, "")
	fmt.Fprintln(summaryFile, "\t\tmax lnbp size:", maxPath)
	fmt.Fprintln(summaryFile, "")
	fmt.Fprintln(summaryFile, "\t\tmin lnbp size:", minPath)
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

func (graph *Graph2) removeLastNodeNBP(nbp *[]Node2, logFile *os.File) *[]Node2 {
	//fmt.Println("length of connected components before removal:", len(graph.connectedComponents))
	newSlice := (*nbp)[:len(*nbp)-1]
	fmt.Fprintln(logFile, "****REMOVENODE METHOD**** length of newSlice:", len(newSlice))
	return &newSlice
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

func (connectedComp *connectedComponent) FindLongestPath(graph *Graph2, logFile *os.File) {
	longestPath := make([]Node2, 0)
	for _, node := range connectedComp.nodes {
		var tempPath *[]Node2
		tp := make([]Node2, 0)
		tempPath = &tp
		connectedComp.dfs1(node, &tempPath, logFile, graph)
	}
	longestPath = graph.maxNBP(connectedComp.nbps)
	fmt.Fprintln(logFile, "length of longest path", len(longestPath))
	graph.addLNBP(longestPath)
}

func (connectedComp *connectedComponent) dfs1(node Node2, nbp **[]Node2, logFile *os.File, graph *Graph2) {
	fmt.Fprintln(logFile, "entering dfs")
	fmt.Fprintln(logFile, "innodes:", len(node.innodes))
	fmt.Fprintln(logFile, "outnodes:", len(node.outnodes))
	fmt.Fprintln(logFile, "node of dfs:", node.id)

	// add this node to the path
	nbpTempPoint := *nbp
	nbpTemp := *nbpTempPoint
	nbpTemp = append(nbpTemp, node)
	nbpTempPoint = &nbpTemp
	nbp = &nbpTempPoint

	//checking each of the adjacent nodes
	for _, outnode := range node.outnodes {
		fmt.Fprint(logFile, "path at top of for: [ ")
		for _, node := range **nbp {
			fmt.Fprint(logFile, node.id, " ")
		}
		fmt.Fprintln(logFile, "]")
		fmt.Fprintln(logFile, "edge found:", outnode.id)
		inCC := graph.Nodes[outnode.id].isInConnectedCompPoint(*nbp)
		fmt.Fprintln(logFile, "is it in the cc?", graph.Nodes[outnode.id].isInConnectedCompPoint(*nbp))

		if !inCC {
			connectedComp.dfs1(graph.Nodes[outnode.id], nbp, logFile, graph)
			fmt.Fprint(logFile, "path at bottom of for: [ ")
			for _, node := range **nbp {
				fmt.Fprint(logFile, node.id, " ")
			}
			fmt.Fprintln(logFile, "]")
		}
	}

	//done visiting all of them add current path to list of paths and delete nyself
	fmt.Fprintln(logFile, "done with", node.id)
	fmt.Fprint(logFile, "path found: [ ")
	for _, node := range **nbp {
		fmt.Fprint(logFile, node.id, " ")
	}
	fmt.Fprintln(logFile, "]")
	//fmt.Fprintln(logFile, "length of nbp that's being added", len(*nbp))
	connectedComp.nbps = append(connectedComp.nbps, **nbp)
	*nbp = graph.removeLastNodeNBP(*nbp, logFile)
	fmt.Fprint(logFile, "path after deletion: [ ")
	for _, node := range **nbp {
		fmt.Fprint(logFile, node.id, " ")
	}
	fmt.Fprintln(logFile, "]")
	//fmt.Fprintln(logFile, "length of nbp after deletion", len(*nbp))
}

func (graph *Graph2) addLNBP(lnbp []Node2) {
	graph.LNBPs = append(graph.LNBPs, lnbp)
}

func (graph *Graph2) LNBPfinder(logFile *os.File) {
	for _, cc := range graph.connectedComponents {
		cc.FindLongestPath(graph, logFile)
	}
}
