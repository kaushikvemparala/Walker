package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	summaryFileName := "summary1000TEST.txt"
	summaryFile, err1 := os.Create(summaryFileName)
	if err1 != nil {
		panic("Sorry, couldn't create file!")
	}
	logFileName := "log1000TEST.txt"
	logFile, err := os.Create(logFileName)
	if err != nil {
		panic("Sorry, couldn't create file!")
	}

	fmt.Fprintln(summaryFile, "\t************************************")
	fmt.Fprintln(summaryFile, "\t******** Assembling genomes ********")
	fmt.Fprintln(summaryFile, "\t************************************")
	fmt.Fprintln(summaryFile, "")

	/*
		// part 1: initial assembler with perfect conditions
		//generate random genome
		length := 3000000
		k := 200
		randomGenome := GenerateRandomGenome(length)
		kmers := KmerComposition(randomGenome, k)
		// assemble genome
		//constructedGenome := GenomeAssembler1(kmers)
		constructedGenome := GenomeAssembler2(kmers)
		if constructedGenome == randomGenome {
			fmt.Println("Yay")
		}
	*/

	rand.Seed(time.Now().UnixNano()) // start the pseudo random number generation in a seemingly random place
	/*
		newGenome := GenerateRandomGenome(10)
		fmt.Println(newGenome)
	*/

	/*
		// part 2: reads have different lengths, imperfect coverage
		length := 100000
		genome := GenerateRandomGenome(length)
		minReadLength := 500
		maxReadLength := 1000
		coverage := 300
		reads := SimulateReads(genome, minReadLength, maxReadLength, coverage)
		fmt.Println("We have:", len(reads), "total reads.")
		minMatchLength := 300
		indexLength := 150
		contigs := GenomeAssembler3(reads, minMatchLength, indexLength)
		if contigs[0] == genome {
			fmt.Println("Good")
		}
		fmt.Println(len(contigs[0]), len(contigs[1]), len(contigs[2]))
		fmt.Println(len(contigs), "total contigs.")
	*/

	/*
		// part 3: applying imperfect coverage, varying length assembler to real reads.
		// the sequencing errors strike back
		filename := "data/BS_2GG.fasta.txt"
		reads := CollectReadsFromFASTA(filename)
		fmt.Println("We have", len(reads), "total reads.")
		PrintStatistics(reads)
		minReadLength := 1000
		fmt.Println("Let's throw out short reads of length <", minReadLength)
		reads = DiscardShortReads(reads, minReadLength)
		fmt.Println("Updated read stats.")
		PrintStatistics(reads)
		fmt.Println("Calling assembler.")
		minMatchLength := 300
		indexLength := 150
		contigs := GenomeAssembler3(reads, minMatchLength, indexLength)
		PrintStatistics(contigs)
	*/

	/*
		for i := range contigs {
			fmt.Println(len(contigs[i]))
		}
	*/

	/*
		// how similar are two mutated string with respect to k-mers?
		stringLength := 1000
		errorRate := 0.11
		k := 7
		numKmers := ExpectedSharedkmers(stringLength, errorRate, k)
		fmt.Println(numKmers)
		random1 := GenerateRandomGenome(stringLength)
		random2 := GenerateRandomGenome(stringLength)
		fmt.Println(CountSharedKmers(random1, random2, k))
	*/

	// part 4: saving our assembler OR coder's revenge


	filename := "data/BS_2GG.fasta.txt"
	fmt.Fprintln(summaryFile, "\tLoading reads...")
	fmt.Fprintln(summaryFile, "")
	reads := CollectReadsFromFASTA(filename)
	fmt.Fprintln(summaryFile, "\t\tLoaded", len(reads), "total reads.")
	fmt.Fprintln(summaryFile, "")
	PrintStatistics(reads, summaryFile)

	minReadLength := 1000
	fmt.Fprintln(summaryFile, "\t\tLet's throw out short reads of length <", minReadLength)
	fmt.Fprintln(summaryFile, "")
	reads = DiscardShortReads(reads, minReadLength)
	fmt.Fprintln(summaryFile, "\t\tUpdated read stats.")
	fmt.Fprintln(summaryFile, "")
	PrintStatistics(reads, summaryFile)

	fmt.Fprintln(summaryFile, "\tCalling assembler...")
	fmt.Fprintln(summaryFile, "")
	minMatchLength := 300
	indexLength := 15
	k := 7
	errorRate := 0.21
	//contigs := GenomeAssembler4(reads, minMatchLength, indexLength, errorRate, k)
	//fmt.Println(len(contigs))
	//PrintStatistics(contigs)
	//fmt.Println("Finally, we write contigs to file.")
	//outFilename := "assembly_contigs.fasta"
	//WriteContigsToFile(contigs, outFilename)
	//var graph Graph
	//start1 := time.Now()
	numOfReads := len(reads)
	fmt.Fprintln(summaryFile, "\t\tRunning assembly on", numOfReads, "reads.")
	fmt.Fprintln(summaryFile, "")
	var graph Graph2
	graph = CreateReadNetwork5Index(reads[:numOfReads], minMatchLength, k, indexLength, errorRate, logFile, summaryFile) //CreateReadNetwork3Index(reads[:numOfReads], minMatchLength, k, indexLength, errorRate, logFile, summaryFile)
	//graph = GetTestGraph52()
	//fmt.Println("network created")
	//graph.PrintGraph()
	pointerToGraph := &graph
	//start := time.Now()
	pointerToGraph.FindConnectedComponents(logFile, summaryFile)
	pointerToGraph.PruneConnectedComps(logFile, summaryFile)
	//pointerToGraph.LNBPfinder(logFile)
	fmt.Fprintln(summaryFile, "\tAssembly summary")
	fmt.Fprintln(summaryFile, "")
	fmt.Fprintln(summaryFile, "\t\tminMatchLength:", minMatchLength)
	fmt.Fprintln(summaryFile, "")
	fmt.Fprintln(summaryFile, "\t\tkmer length:", k)
	fmt.Fprintln(summaryFile, "")
	fmt.Fprintln(summaryFile, "\t\tindexLength:", indexLength)
	fmt.Fprintln(summaryFile, "")
	fmt.Fprintln(summaryFile, "\t\terrorRate:", errorRate)
	fmt.Fprintln(summaryFile, "")
	Graph2Statistics(graph, summaryFile)
	fmt.Fprintln(logFile, "number of edges built: ", len(graph.Edges))


	//fmt.Println("read:", reads[0])
	//fmt.Println()
	//fmt.Println()
	//fmt.Println()
	//fmt.Println()
	//fmt.Println("read:", reads[1])
	//elapsed := time.Since(start)
	//elapsed1 := time.Since(start1)
	//graph.PrintGraph()
	//fmt.Println("DFS took:", elapsed)
	//fmt.Println("Edges took:", elapsed1)
	//graph.PrintGraph2()

	//fmt.Println("Making read network")
	//CreateReadNetwork(reads, minMatchLength, k, indexLength, errorRate)
	//var graph Graph
	//graph.PrintGraph()
	//graph = GenomeAssemblerKaushik4(reads, 200, 7, .11, 15) // GetTestGraph5()
	//graph.PrintGraph()
	//pointerToGraph := &graph
	//pointerToGraph.FindConnectedComponents()
	//pointerToGraph.PruneConnectedComps()
	//graph.PrintGraph()

	//CreateReadNetwork(reads, 300)

	//stringone := GenerateRandomGenome(5)
	//fmt.Println("string one:", stringone)
	//stringtwo := GenerateRandomGenome(5)
	//fmt.Println("string two:", stringtwo)
	//sharedKmers := CountSharedKmers(stringone, stringtwo, 3)
	//fmt.Println("shared kmers:", sharedKmers)

	/*
		testGraph := GetTestGraph112()
		pointerToGraph := &testGraph
		pointerToGraph.FindConnectedComponents(logFile, summaryFile)
		pointerToGraph.PruneConnectedComps(logFile, summaryFile)
		 testGraph.LNBPs = testGraph.getNonBranchingPaths(logFile)
		Graph2Statistics(testGraph, summaryFile)
		for i, conn := range testGraph.connectedComponents {
			fmt.Print("cc ", i+1, ":	[")
			for _, node := range conn.nodes {
				fmt.Print(" ", node.id, " ")
			}
			fmt.Println("]")
		}
		*/

}

func GetTestGraph() Graph {
	testGraph := MakeGraph()
	testGraph.addNode(0, "A")
	testGraph.addNode(1, "B")
	testGraph.addNode(2, "C")
	testGraph.addNode(3, "D")
	testGraph.addNode(4, "E")
	testGraph.addNode(5, "F")
	testGraph.addNode(6, "G")
	testGraph.addEdge(0, 1)
	return testGraph
}

func GetTestGraph1() Graph2 {
	testGraph := MakeGraph2()
	testGraph.addNode2(0, "A")
	return testGraph
}

func GetTestGraph22() Graph2 {
	testGraph := MakeGraph2()
	testGraph.addNode2(0, "A")
	testGraph.addNode2(1, "B")
	testGraph.addEdge2(0, 1, "Over")
	return testGraph
}

func GetTestGraph2() Graph {
	testGraph := MakeGraph()
	testGraph.addNode(0, "A")
	testGraph.addNode(1, "B")
	testGraph.addNode(2, "C")
	testGraph.addNode(3, "D")
	testGraph.addNode(4, "E")
	testGraph.addNode(5, "F")
	testGraph.addNode(6, "G")
	testGraph.addEdge(0, 1)
	testGraph.addEdge(2, 5)
	testGraph.addEdge(1, 2)
	return testGraph
}

func GetTestGraph3() Graph {
	testGraph := MakeGraph()
	testGraph.addNode(0, "A")
	testGraph.addNode(1, "B")
	testGraph.addNode(2, "C")
	testGraph.addNode(3, "D")
	testGraph.addNode(4, "E")
	testGraph.addNode(5, "F")
	testGraph.addNode(6, "G")
	testGraph.addEdge(0, 1)
	testGraph.addEdge(2, 1)
	//testGraph.addEdge(1, 2)
	return testGraph
}

func GetTestGraph4() Graph {
	testGraph := MakeGraph()
	testGraph.addNode(0, "A")
	testGraph.addNode(1, "B")
	testGraph.addNode(2, "C")
	testGraph.addNode(3, "D")
	testGraph.addNode(4, "E")
	testGraph.addNode(5, "F")
	testGraph.addNode(6, "G")
	testGraph.addEdge(0, 1)
	testGraph.addEdge(2, 0)
	testGraph.addEdge(1, 2)
	return testGraph
}

func GetTestGraph5() Graph {
	testGraph := MakeGraph()
	testGraph.addNode(0, "A")
	testGraph.addNode(1, "B")
	testGraph.addNode(2, "C")
	testGraph.addNode(3, "D")
	testGraph.addNode(4, "E")
	testGraph.addNode(5, "F")
	testGraph.addNode(6, "G")
	testGraph.addEdge(0, 1)
	testGraph.addEdge(2, 0)
	testGraph.addEdge(1, 2)
	testGraph.addEdge(5, 6)
	testGraph.addEdge(4, 5)
	return testGraph
}

func GetTestGraph32() Graph2 {
	testGraph := MakeGraph2()
	for i := 0; i < 4999; i++ {
		testGraph.addNode2(i, "A")
	}
	for i := 0; i < 4998; i++ {
		testGraph.addEdge2(i, i+1, "Over")
	}
	testGraph.addEdge2(4999, 0, "Over")
	return testGraph
}

func GetTestGraph42() Graph2 {
	testGraph := MakeGraph2()
	for i := 0; i < 5000; i++ {
		testGraph.addNode2(i, "A")
	}
	return testGraph
}

func GetTestGraph52() Graph2 {
	testGraph := MakeGraph2()
	testGraph.addNode2(0, "A")
	testGraph.addNode2(1, "B")
	testGraph.addNode2(2, "C")
	testGraph.addNode2(3, "D")
	testGraph.addNode2(4, "E")
	testGraph.addEdge2(0, 1, "Over")
	testGraph.addEdge2(2, 1, "Over")
	testGraph.addEdge2(1, 3, "Over")
	testGraph.addEdge2(3, 2, "Over")
	return testGraph
}

func GetTestGraph62() Graph2 {
	testGraph := MakeGraph2()
	testGraph.addNode2(0, "A")
	testGraph.addNode2(1, "B")
	testGraph.addNode2(2, "C")
	testGraph.addNode2(3, "D")
	testGraph.addNode2(4, "E")
	testGraph.addNode2(5, "E")
	testGraph.addNode2(6, "E")
	testGraph.addEdge2(0, 1, "Over")
	testGraph.addEdge2(0, 2, "Over")
	testGraph.addEdge2(0, 3, "Over")
	testGraph.addEdge2(2, 5, "Over")
	testGraph.addEdge2(1, 4, "Over")
	testGraph.addEdge2(3, 6, "Over")
	return testGraph
}

func GetTestGraph72() Graph2 {
	testGraph := MakeGraph2()
	testGraph.addNode2(0, "A")
	testGraph.addNode2(1, "B")
	testGraph.addNode2(2, "C")
	testGraph.addNode2(3, "D")
	testGraph.addNode2(4, "E")
	testGraph.addNode2(5, "E")
	testGraph.addNode2(6, "E")
	testGraph.addNode2(7, "E")
	testGraph.addNode2(8, "E")
	testGraph.addEdge2(0, 1, "Over")
	testGraph.addEdge2(0, 2, "Over")
	testGraph.addEdge2(0, 3, "Over")
	testGraph.addEdge2(2, 5, "Over")
	testGraph.addEdge2(1, 4, "Over")
	testGraph.addEdge2(3, 6, "Over")
	testGraph.addEdge2(6, 7, "Over")
	testGraph.addEdge2(3, 5, "Over")
	testGraph.addEdge2(5, 8, "Over")
	testGraph.addEdge2(8, 7, "Over")
	return testGraph
}

func GetTestGraph82() Graph2 {
	testGraph := MakeGraph2()
	testGraph.addNode2(0, "A")
	testGraph.addNode2(1, "B")
	testGraph.addNode2(2, "C")
	return testGraph
}

func GetTestGraph92() Graph2 {
	testGraph := MakeGraph2()
	testGraph.addNode2(0, "A")
	testGraph.addNode2(1, "B")
	testGraph.addNode2(2, "C")
	testGraph.addNode2(3, "A")
	testGraph.addNode2(4, "B")
	testGraph.addNode2(5, "C")
	testGraph.addEdge2(0, 1, "Over")
	testGraph.addEdge2(1, 2, "Over")
	testGraph.addEdge2(2, 3, "Over")
	testGraph.addEdge2(3, 4, "Over")
	testGraph.addEdge2(4, 5, "Over")
	testGraph.addEdge2(5, 0, "Over")
	return testGraph
}

func GetTestGraph102() Graph2 {
	testGraph := MakeGraph2()
	testGraph.addNode2(0, "A")
	testGraph.addNode2(1, "B")
	testGraph.addNode2(2, "C")
	testGraph.addEdge2(0, 1, "Over")
	testGraph.addEdge2(1, 0, "Over")
	testGraph.addEdge2(1, 2, "Over")
	testGraph.addEdge2(2, 1, "Over")
	return testGraph
}

func GetTestGraph112() Graph2 {
	testGraph := MakeGraph2()
	testGraph.addNode2(0, "A")
	testGraph.addNode2(1, "B")
	testGraph.addNode2(2, "C")
	testGraph.addNode2(3, "C")
	testGraph.addNode2(4, "C")
	testGraph.addNode2(5, "C")
	testGraph.addNode2(6, "C")
	testGraph.addEdge2(0, 1, "Over")
	testGraph.addEdge2(1, 2, "Over")
	testGraph.addEdge2(3, 4, "Over")
	testGraph.addEdge2(3, 5, "Over")
	testGraph.addEdge2(5, 6, "Over")
	testGraph.addEdge2(6, 5, "Over")
	return testGraph
}
