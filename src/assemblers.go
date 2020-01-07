package main

import (
	"fmt"
	"time"
)

//GenomeAssembler1 takes a collection of strings and returns a genome whose
//k-mer composition is these strings. It makes the following assumptions.
//1. "Perfect coverage" -- every k-mer is detected
//2. No errors in reads
//3. Every read has equal length (k)
//4. DNA is single-stranded
//5. (No k-mer repeats)
func GenomeAssembler1(kmers []string) string {
	// greedy algorithm: look for whatever helps me the most (overlap of k-1 symbols).
	if len(kmers) == 0 {
		panic("Error: No kmers given to GenomeAssembler!")
	}
	// start with arbitrary kmer
	// first, what is k? length of first read
	k := len(kmers[0])

	genome := kmers[len(kmers)/2] // midpoint k-mer

	// let's throw out everything we have used
	kmers = Remove(kmers, len(kmers)/2)

	// while we still have reads, try to extend current read
	for len(kmers) > 0 {
		// note: we need to remember to delete any kmer we use or else hit an infinite loop
		for i, kmer := range kmers {
			// try to extend genome to left and right
			// a hit means that we match k-1 nucleotides to end of genome
			if genome[0:k-1] == kmer[1:] { // extending left
				// update genome by adding first symbol of kmer to left
				genome = kmer[0:1] + genome
				// throw out read
				kmers = Remove(kmers, i)
				// stop the for loop so we don't have an index out of bounds error
				break // breaks innermost loop you are in
			} else if genome[len(genome)-k+1:len(genome)] == kmer[:k-1] { // extending right
				genome = genome + kmer[k-1:]
				kmers = Remove(kmers, i)
				break
			}
		}
	}

	return genome
}

//Remove takes a collection of strings and an index.
//It removes the string at the given index and returns the updated array.
func Remove(patterns []string, index int) []string {
	// remember our trick for deleting an element
	patterns = append(patterns[:index], patterns[index+1:]...)
	return patterns
}

// issue 1: current assembler is slowwwwwww (won't scale to 3M bp)
// reason why is because most of the time, when it is looking for a match, it can't find one.
// this gets worse the bigger the genome gets.
// solution: build prefix and suffix indices before looking for matches.

func GenomeAssembler2(kmers []string) string {
	if len(kmers) == 0 {
		panic("Error: no kmers given to assembler")
	}
	k := len(kmers[0])
	genome := kmers[len(kmers)/2]

	//build a prefix and suffix index
	indexLength := k - 1

	fmt.Println("Building indices.")
	prefixIndex := BuildPrefixIndex(kmers, indexLength)
	fmt.Println("Prefix index built.")
	suffixIndex := BuildSuffixIndex(kmers, indexLength)
	fmt.Println("Suffix index built. Ready to assemble!")

	// while we continue to find things, keep going
	keepLooping := true
	counter := 0
	for keepLooping == true {
		keepLooping = false
		// update keepLooping to true when we find an overlap
		// first, check the right side of genome
		prefix := genome[len(genome)-k+1:]
		// is prefix in the prefix index?
		matches1, exists1 := prefixIndex[prefix] // ok1 is true if this exists in map
		if exists1 == true {                     // match found :)
			// make sure we keep going!
			keepLooping = true
			// where do I need to look in my kmers?
			nextRead := kmers[matches1[0]] // always take the first match we see
			genome = genome + nextRead[k-1:]
			counter++
			if counter%100000 == 0 {
				fmt.Println("Update: We have overlapped", counter, "reads.")
			}
		}
		// now, try to extend to the left too
		suffix := genome[:k-1]
		matches2, exists2 := suffixIndex[suffix]
		if exists2 { // we found a match
			keepLooping = true
			prevRead := kmers[matches2[0]]
			// extend genome left
			genome = prevRead[:len(prevRead)-(k-1)] + genome
			counter++
			if counter%100000 == 0 {
				fmt.Println("Update: We have overlapped", counter, "reads.")
			}
		}
	}

	return genome
}

// part 3: relaxing the assumptions of "perfect coverage" and equal read lengths
// how well does a collection of reads "cover" a genome?
// first, we should be able to simulate a dataset to test the algorithm we develop
// minMatchLength is the minimum perfect match we will allow between overlapping reads.
// minMatchLength must be bigger than the index length
// now we will produce contigs too.
// these reads now have variable length (bigger than indexLength)

func GenomeAssembler3(reads []string, minMatchLength, indexLength int) []string {
	if len(reads) == 0 {
		panic("Error: No reads given to GenomeAssembler.")
	}

	if minMatchLength <= indexLength {
		panic("Error: minMatchLength must be bigger than indexLength.")
	}

	contigs := make([]string, 0)

	fmt.Println("Building a prefix and suffix index for reads.")
	prefixIndex := BuildPrefixIndex(reads, indexLength)
	fmt.Println("Prefix index built!")
	suffixIndex := BuildSuffixIndex(reads, indexLength)
	fmt.Println("Suffix index built!")

	currentReadIndex := 0                  // or whatever
	currentRead := reads[currentReadIndex] // get corresponding read

	// idea: whenever we use a read, let's delete it from the prefix index (and suffix index).
	// continue for as long as we have elements still in the prefix index.
	for len(prefixIndex) > 0 {
		// let's throw out the elements of the indices corresponding to current read.
		prefix := currentRead[:indexLength]
		suffix := currentRead[len(currentRead)-indexLength:]
		delete(prefixIndex, prefix)
		delete(suffixIndex, suffix)

		//extend currentRead to right and extend to left as far as I can.
		contig1 := ExtendContigRight(currentRead, prefixIndex, suffixIndex, reads, minMatchLength, indexLength)
		contig2 := ExtendContigLeft(currentRead, prefixIndex, suffixIndex, reads, minMatchLength, indexLength)

		// join into one contig and append to our set
		contig := contig2 + contig1[len(currentRead):]

		//previously, we appended every contig we found, even if it wasn't good (i.e., short).
		//because coverage is high, let's just keep longer contigs.
		if len(contig) > 100000 {
			contigs = append(contigs, contig)
			fmt.Println("We have generated", len(contigs), "contigs.")
			fmt.Println("Prefix index is down to", len(prefixIndex), "elements.")
		}

		// we need a new starting point (currentRead) if still stuff in prefix index
		if len(prefixIndex) > 0 {
			// note: we know which reads haven't been used!
			// they're the elements still in the prefix index.
			// so just range over the prefix index, grab the first thing we see, and break
			for prefix := range prefixIndex {
				currentReadIndex = (prefixIndex[prefix])[0]
				currentRead = reads[currentReadIndex]
				break // stop as soon as we grab a value
			}
		}
	}

	return contigs
}

//ExtendContigRight takes an initial string (currentRead) along with everything we need for assembly. It iteratively extends our initial string to the right by looking for exact matches in the prefix index. As it goes, it deletes elements from the indices. It returns a string corresponding to a contig.
func ExtendContigRight(currentRead string, prefixIndex, suffixIndex map[string][]int, reads []string, minMatchLength, indexLength int) string {
	contig := currentRead

	keepLooping := true
	// while we can keep going right
	for keepLooping == true {
		keepLooping = false
		// if we find anything, we will update it to true

		// range over all possible overlap lengths and pick the first place that we find a long, matching string.
		n := len(currentRead)
		for j := 1; j <= n-minMatchLength; j++ {
			// let's try overlapping this.
			prefix := currentRead[j : j+indexLength]
			// is this prefix present in the index?
			matchList, exists := prefixIndex[prefix]
			if exists {
				// grab first element as matching read
				matchedRead := reads[matchList[0]]
				// does this string match completely? AND is it long enough?
				if len(matchedRead) > n-j && currentRead[j:] == matchedRead[:n-j] {
					// success!
					keepLooping = true
					contig += matchedRead[n-j:]
					//update currentRead and its length
					currentRead = matchedRead
					n = len(currentRead)
					// clean up the indices too by throwing out its prefix and suffix.
					delete(prefixIndex, prefix)
					suffix := currentRead[n-indexLength:] // what we overlapped
					delete(suffixIndex, suffix)
					break // stop the outer looping process since we found a match.
				}
			}
		}
	}

	return contig
}

//ExtendContigLeft takes an initial string (currentRead) along with everything we need for assembly. It iteratively extends our initial string to the left by looking for exact matches in the suffix index. As it goes, it deletes elements from the indices. It returns a string corresponding to a contig.
func ExtendContigLeft(currentRead string, prefixIndex, suffixIndex map[string][]int, reads []string, minMatchLength, indexLength int) string {
	contig := currentRead

	keepLooping := true
	// while we can keep going right
	for keepLooping == true {
		keepLooping = false
		// if we find anything, we will update it to true
		// range over all possible overlap lengths and pick the first place that we find a long, matching string.
		n := len(currentRead)
		for j := 1; j <= n-minMatchLength; j++ { // j represents a count from right end of string
			// let's try overlapping this.
			suffix := currentRead[n-j-indexLength : n-j]
			// is this prefix present in the index?
			matchList, exists := suffixIndex[suffix]
			if exists {
				// grab first element as matching read
				matchedRead := reads[matchList[0]]
				// does this string match completely? AND is it long enough?
				if len(matchedRead) > n-j && currentRead[:n-j] == matchedRead[len(matchedRead)-(n-j):] {
					// success!
					keepLooping = true
					contig = matchedRead[:len(matchedRead)-(n-j)] + contig
					//update currentRead and its length
					currentRead = matchedRead
					n = len(currentRead)
					// clean up the indices too by throwing out its prefix and suffix.
					prefix := currentRead[:indexLength]
					delete(prefixIndex, prefix)
					delete(suffixIndex, suffix)
					break // stop the outer looping process since we found a match.
				}
			}
		}
	}
	return contig
}

func GenomeAssembler4(reads []string, minMatchLength, indexLength int, errorRate float64, k int) []string {
	if len(reads) == 0 {
		panic("Error: No reads given to GenomeAssembler.")
	}

	if minMatchLength <= indexLength {
		panic("Error: minMatchLength must be bigger than indexLength.")
	}

	contigs := make([]string, 0)

	fmt.Println("Building a prefix and suffix index for reads.")
	prefixIndex := BuildPrefixIndex(reads, indexLength)
	fmt.Println("Prefix index built!")
	suffixIndex := BuildSuffixIndex(reads, indexLength)
	fmt.Println("Suffix index built!")

	currentReadIndex := 0                  // or whatever
	currentRead := reads[currentReadIndex] // get corresponding read

	// idea: whenever we use a read, let's delete it from the prefix index (and suffix index).
	// continue for as long as we have elements still in the prefix index.
	for len(prefixIndex) > 0 {
		// let's throw out the elements of the indices corresponding to current read.
		prefix := currentRead[:indexLength]
		suffix := currentRead[len(currentRead)-indexLength:]
		delete(prefixIndex, prefix)
		delete(suffixIndex, suffix)

		//extend currentRead to right and extend to left as far as I can.
		contig1 := ExtendContigRightInexact(currentRead, prefixIndex, suffixIndex, reads, minMatchLength, indexLength, errorRate, k)
		contig2 := ExtendContigLeftInexact(currentRead, prefixIndex, suffixIndex, reads, minMatchLength, indexLength, errorRate, k)

		// join into one contig and append to our set
		contig := contig2 + contig1[len(currentRead):]

		//previously, we appended every contig we found, even if it wasn't good (i.e., short).
		//because coverage is high, let's just keep longer contigs.
		if len(contig) > 100000 {
			contigs = append(contigs, contig)
			fmt.Println("We have generated", len(contigs), "contigs.")
			fmt.Println("Prefix index is down to", len(prefixIndex), "elements.")
		}

		// we need a new starting point (currentRead) if still stuff in prefix index
		if len(prefixIndex) > 0 {
			// note: we know which reads haven't been used!
			// they're the elements still in the prefix index.
			// so just range over the prefix index, grab the first thing we see, and break
			for prefix := range prefixIndex {
				currentReadIndex = (prefixIndex[prefix])[0]
				currentRead = reads[currentReadIndex]
				break // stop as soon as we grab a value
			}
		}
	}

	return contigs
}

func ExtendContigRightInexact(currentRead string, prefixIndex, suffixIndex map[string][]int, reads []string, minMatchLength, indexLength int, errorRate float64, k int) string {
	contig := currentRead

	keepLooping := true
	// while we can keep going right
	for keepLooping == true {
		keepLooping = false
		// if we find anything, we will update it to true

		// range over all possible overlap lengths and pick the first place that we find a long, matching string.
		n := len(currentRead)
		for j := 1; j <= n-minMatchLength; j++ {
			// let's try overlapping this.
			prefix := currentRead[j : j+indexLength]
			// is this prefix present in the index?
			matchList, exists := prefixIndex[prefix]
			if exists {
				// grab first element as matching read
				matchedRead := reads[matchList[0]]
				// does this string match completely? AND is it long enough?
				if len(matchedRead) > n-j && float64(CountSharedKmers(currentRead[j:], matchedRead[:n-j], k)) >= 0.9*float64(ExpectedSharedkmers(len(currentRead[j:]), errorRate, k)) {
					// success!
					keepLooping = true
					contig += matchedRead[n-j:]
					//update currentRead and its length
					currentRead = matchedRead
					n = len(currentRead)
					// clean up the indices too by throwing out its prefix and suffix.
					delete(prefixIndex, prefix)
					suffix := currentRead[n-indexLength:] // what we overlapped
					delete(suffixIndex, suffix)
					break // stop the outer looping process since we found a match.
				}
			}
		}
	}

	return contig
}

func ExtendContigLeftInexact(currentRead string, prefixIndex, suffixIndex map[string][]int, reads []string, minMatchLength, indexLength int, errorRate float64, k int) string {
	contig := currentRead

	keepLooping := true
	// while we can keep going right
	for keepLooping == true {
		keepLooping = false
		// if we find anything, we will update it to true
		// range over all possible overlap lengths and pick the first place that we find a long, matching string.
		n := len(currentRead)
		for j := 1; j <= n-minMatchLength; j++ { // j represents a count from right end of string
			// let's try overlapping this.
			suffix := currentRead[n-j-indexLength : n-j]
			// is this prefix present in the index?
			matchList, exists := suffixIndex[suffix]
			if exists {
				// grab first element as matching read
				matchedRead := reads[matchList[0]]
				// does this string match completely? AND is it long enough?
				if len(matchedRead) > n-j && float64(CountSharedKmers(currentRead[:n-j], matchedRead[len(matchedRead)-(n-j):], k)) >= 0.9*float64(ExpectedSharedkmers(len(currentRead[:n-j]), errorRate, k)) {
					// success!
					keepLooping = true
					contig = matchedRead[:len(matchedRead)-(n-j)] + contig
					//update currentRead and its length
					currentRead = matchedRead
					n = len(currentRead)
					// clean up the indices too by throwing out its prefix and suffix.
					prefix := currentRead[:indexLength]
					delete(prefixIndex, prefix)
					delete(suffixIndex, suffix)
					break // stop the outer looping process since we found a match.
				}
			}
		}
	}
	return contig
}

// CreateReadNetwork takes the array of reads and the minMatchLength and creates the network of reads.
func CreateReadNetwork(reads []string, minMatchLength, k, indexLength int, errorRate float64) Graph {
	network := MakeGraph()
	pointerToNetwork := &network

	for i, read := range reads {
		fmt.Println("Adding read number", i+1)
		pointerToNetwork.addNode(i, read)
	}
	count := 0
	for id, node := range network.Nodes {
		fmt.Println("node:", id)
		n := len(node.read)
		expectedKmers := 0.5 * float64(ExpectedSharedkmers(len(node.read[:minMatchLength]), errorRate, k))
		expectedKmers2 := 0.5 * float64(ExpectedSharedkmers(len(node.read[n-minMatchLength:]), errorRate, k))
		BuildEdges(node, minMatchLength, k, errorRate, pointerToNetwork, reads, expectedKmers, expectedKmers2)
		count++
		fmt.Println("done with", count, "of", len(network.Nodes))
	}
	//pathMap := make(map[int][]Node)
	//pointerToPathMap := &pathMap

	//fmt.Println("Building a prefix and suffix index for reads.")
	//prefixIndex := BuildPrefixIndex(reads, indexLength)
	//fmt.Println("Prefix index built! And the length is:", len(prefixIndex))
	//suffixIndex := BuildSuffixIndex(reads, indexLength)
	//fmt.Println("Suffix index built! And the length is:", len(suffixIndex))

	//currentReadIndex := 0
	//currentRead := reads[currentReadIndex]
	//currentNode := network.Nodes[currentReadIndex]
	//for len(prefixIndex) > 0 {
	// let's throw out the elements of the indices corresponding to current read.
	//prefix := currentRead[:indexLength]
	//suffix := currentRead[len(currentRead)-indexLength:]
	//delete(prefixIndex, prefix)
	//delete(suffixIndex, suffix)

	//extend currentRead to right and extend to left as far as I can.
	//rightNodes := BuildEdgesRight(currentNode, minMatchLength, k, indexLength, errorRate, prefixIndex, suffixIndex, pointerToNetwork, reads)
	//leftNodes := BuildEdgesLeft(currentNode, minMatchLength, k, indexLength, errorRate, prefixIndex, suffixIndex, pointerToNetwork, reads)

	//path := leftNodes
	//for i := range rightNodes {
	//	if i != 0{
	//		path = append(path, rightNodes[i])
	//	}
	//}

	//if len(path) != 0 {
	//	pathMap[path[0].getID()] = path
	//} else {
	//	pathMap[currentNode.getID()] = make([]Node, 0)
	//}

	//previously, we appended every contig we found, even if it wasn't good (i.e., short).
	//because coverage is high, let's just keep longer contigs.
	//if len(contig) > 100000 {
	//	contigs = append(contigs, contig)
	//	fmt.Println("We have generated", len(contigs), "contigs.")
	//	fmt.Println("Prefix index is down to", len(prefixIndex), "elements.")
	//}

	// we need a new starting point (currentRead) if still stuff in prefix index
	//if len(prefixIndex) > 0 {
	// note: we know which reads haven't been used!
	// they're the elements still in the prefix index.
	// so just range over the prefix index, grab the first thing we see, and break
	//	for prefix := range prefixIndex {
	//		currentReadIndex = (prefixIndex[prefix])[0]
	//		currentNode = network.Nodes[currentReadIndex]
	//		currentRead = reads[currentReadIndex]
	//  	break // stop as soon as we grab a value
	//	}
	//}
	//}

	//fmt.Println("we have found", len(pathMap), "paths.")

	//fmt.Println("Read network created!")

	return network
}

func BuildEdges(node Node, minMatchLength, k int, errorRate float64, graph *Graph, reads []string, expectedKmers, expectedKmers2 float64) {
	n := len(node.read)
	for i := range graph.Nodes {
		matchedNode := graph.Nodes[i]
		matchedRead := matchedNode.read
		//fmt.Println("testing:", node.getID(), "and", matchedNode.getID())
		if len(matchedRead) > minMatchLength && float64(CountSharedKmers(node.read[:minMatchLength], matchedRead[len(matchedRead)-minMatchLength:], k)) >= expectedKmers {
			// success!
			fmt.Println("found a match")
			if !(node.isANeighbor(matchedNode)) {
				fmt.Println("building an edge", node.getID(), matchedNode.getID())
				graph.addEdge(matchedNode.getID(), node.getID())
			}
		} else if len(matchedRead) > minMatchLength && float64(CountSharedKmers(node.read[n-minMatchLength:], matchedRead[:minMatchLength], k)) >= expectedKmers2 {
			fmt.Println("found a match")
			if !(matchedNode.isANeighbor(node)) {
				fmt.Println("building an edge", node.getID(), matchedNode.getID())
				graph.addEdge(node.getID(), matchedNode.getID())
			}
		}

	}
}

func BuildEdgesLeft(node Node, minMatchLength, k, indexLength int, errorRate float64, prefixIndex, suffixIndex map[string][]int, graph *Graph, reads []string) []Node {

	currentNode := node
	currentRead := currentNode.read

	leftPath := make([]Node, 0)

	keepLooping := true
	// while we can keep going right
	for keepLooping == true {
		keepLooping = false
		// if we find anything, we will update it to true
		// range over all possible overlap lengths and pick the first place that we find a long, matching string.
		n := len(currentRead)
		for j := 1; j <= n-minMatchLength; j++ { // j represents a count from right end of string
			// let's try overlapping this.
			suffix := currentRead[n-j-indexLength : n-j]
			// is this prefix present in the index?
			matchList, exists := suffixIndex[suffix]
			if exists {
				fmt.Println("found some matches")
				// grab first element as matching read
				matchedRead := reads[matchList[0]]
				matchedNode := graph.Nodes[matchList[0]]
				// does this string match completely? AND is it long enough?
				if len(matchedRead) > n-j && float64(CountSharedKmers(currentRead[:n-j], matchedRead[len(matchedRead)-(n-j):], k)) >= 0.5*float64(ExpectedSharedkmers(len(currentRead[:n-j]), errorRate, k)) {
					// success!
					keepLooping = true
					fmt.Println("found a match")
					fmt.Println("building an edge", node.getID(), matchedNode.getID())
					graph.addEdge(matchedNode.getID(), node.getID())
					leftPath = append(leftPath, matchedNode)
					//update currentRead and its length
					currentNode = matchedNode
					currentRead = matchedRead
					n = len(currentRead)
					// clean up the indices too by throwing out its prefix and suffix.
					prefix := currentRead[:indexLength]
					delete(prefixIndex, prefix)
					delete(suffixIndex, suffix)
					break // stop the outer looping process since we found a match.
				}
			}
		}
	}

	/*
		currentRead := node.read
		n := len(currentRead)
		for j := 1; j <= n-minMatchLength; j++ { // j represents a count from right end of string
			// let's try overlapping this.
			suffix := currentRead[n-j-indexLength : n-j]
			// is this prefix present in the index?
			matchList, exists := suffixIndex[suffix]
			if exists {
				matchedNode := graph.Nodes[matchList[0]]
				matchedRead := matchedNode.read
				// does this string match completely? AND is it long enough?
				if !(node.isANeighbor(matchedNode)) && node.getID() != matchedNode.getID() {
					if len(matchedRead) > n-j && float64(CountSharedKmers(currentRead[:n-j], matchedRead[len(matchedRead)-(n-j):], k)) >= 0.9*float64(ExpectedSharedkmers(len(currentRead[:n-j]), errorRate, k)) {
						fmt.Println("building an edge", node.getID(), matchedNode.getID())
						graph.addEdge(matchedNode.getID(), node.getID())
						prefix := currentRead[:indexLength]
						delete(prefixIndex, prefix)
						delete(suffixIndex, suffix)
						//fmt.Println("length of pathMap array before adding:", len(pathMap[node.getID()]))
						//addToPath(node, matchedNode, pathMap)
						//fmt.Println("length of pathMap array after adding:", len(pathMap[node.getID()]))
						//if !(matchedNode.visited) {
						//	noode := graph.Nodes[node.getID()]
						//	node.visited = true
						//	graph.Nodes[node.getID()] = noode
						//	BuildEdgesLeft(matchedNode, minMatchLength, k, indexLength, errorRate, prefixIndex, suffixIndex, graph, pathMap)
						//}
					}
					//break
				}
			}
		}
	*/
	/*
		n := len(node.read)
		for j := 1; j <= n-minMatchLength; j++ {
			freqMap := FrequencyMap(node.read[j:j+minMatchLength], k)
			expectedKmers := .7 * float64(ExpectedSharedkmers(len(node.read[j:]), errorRate, k))
			// let's try overlapping this.
			for prefixQ, matchList := range prefixIndex {
				for _, i := range matchList {
					matchedNode := graph.Nodes[i]
					if !(node.isANeighbor(matchedNode)) {
						if len(prefixQ) > n-j && float64(CountSharedKmersMod(freqMap, prefixQ[:n-j], k)) >= expectedKmers {
							fmt.Println("building an edge", node.getID(), matchedNode.getID())
							graph.addEdge(matchedNode.getID(), node.getID())
						}
					}
				}
			}
		}
	*/

	leftPathReverse := make([]Node, 0)
	for i := len(leftPath) - 1; i >= 0; i-- {
		leftPathReverse = append(leftPathReverse, leftPath[i])
	}
	return leftPathReverse
}

func BuildEdgesRight(node Node, minMatchLength, k, indexLength int, errorRate float64, prefixIndex, suffixIndex map[string][]int, graph *Graph, reads []string) []Node {
	currentNode := node
	currentRead := currentNode.read

	rightPath := make([]Node, 0)

	keepLooping := true
	// while we can keep going right
	for keepLooping == true {
		keepLooping = false
		// if we find anything, we will update it to true

		// range over all possible overlap lengths and pick the first place that we find a long, matching string.
		n := len(currentRead)
		for j := 1; j <= n-minMatchLength; j++ {
			// let's try overlapping this.
			prefix := currentRead[j : j+indexLength]
			// is this prefix present in the index?
			matchList, exists := prefixIndex[prefix]
			if exists {
				// grab first element as matching read
				fmt.Println("found some matches")
				matchedRead := reads[matchList[0]]
				matchedNode := graph.Nodes[matchList[0]]
				// does this string match completely? AND is it long enough?
				if len(matchedRead) > n-j && float64(CountSharedKmers(currentRead[j:], matchedRead[:n-j], k)) >= 0.5*float64(ExpectedSharedkmers(len(currentRead[j:]), errorRate, k)) {
					// success!
					keepLooping = true
					fmt.Println("found a match")
					fmt.Println("building an edge", node.getID(), matchedNode.getID())
					graph.addEdge(node.getID(), matchedNode.getID())
					rightPath = append(rightPath, matchedNode)
					//update currentRead and its length
					currentRead = matchedRead
					currentNode = matchedNode
					n = len(currentRead)
					// clean up the indices too by throwing out its prefix and suffix.
					delete(prefixIndex, prefix)
					suffix := currentRead[n-indexLength:] // what we overlapped
					delete(suffixIndex, suffix)
					break // stop the outer looping process since we found a match.
				}
			}
		}
	}
	/*
		currentRead := node.read
		n := len(currentRead)
		for j := 1; j <= n-minMatchLength; j++ { // j represents a count from right end of string
			// let's try overlapping this.
			prefix := currentRead[j : j+indexLength]
			// is this prefix present in the index?
			matchList, exists := prefixIndex[prefix]
			if exists {
				// grab first element as matching read
				matchedNode := graph.Nodes[matchList[0]]
				matchedRead := matchedNode.read
				// does this string match completely? AND is it long enough?
				if !(matchedNode.isANeighbor(node)) && node.getID() != matchedNode.getID() {
					if len(matchedRead) > n-j && float64(CountSharedKmers(currentRead[j:], matchedRead[:n-j], k)) >= 0.9*float64(ExpectedSharedkmers(len(currentRead[j:]), errorRate, k)) {
						fmt.Println("building an edge", node.getID(), matchedNode.getID())
						graph.addEdge(node.getID(), matchedNode.getID())
						delete(prefixIndex, prefix)
						suffix := currentRead[n-indexLength:] // what we overlapped
						delete(suffixIndex, suffix)
						//if matchedNode.visited {
						//	noode := graph.Nodes[node.getID()]
						//	node.visited = true
						//	graph.Nodes[node.getID()] = noode
						//	BuildEdgesRight(matchedNode, minMatchLength, k, indexLength, errorRate, prefixIndex, suffixIndex, graph, pathMap)
						//}
					}
					//break
				}
			}
		}
	*/

	/*
		n := len(node.read)
		for j := 1; j <= n-minMatchLength; j++ {
			freqMap := FrequencyMap(node.read[n-j-minMatchLength:n-j], k)
			expectedKmers := .7 * float64(ExpectedSharedkmers(len(node.read[:n-j]), errorRate, k))
			// let's try overlapping this.
			for suffixQ, matchList := range suffixIndex {
				for _, i := range matchList {
					matchedNode := graph.Nodes[i]
					if !(node.isANeighbor(matchedNode)) {
						if len(suffixQ) > n-j && float64(CountSharedKmersMod(freqMap, suffixQ[:n-j], k)) >= expectedKmers {
							fmt.Println("building an edge", node.getID(), matchedNode.getID())
							graph.addEdge(node.getID(), matchedNode.getID())
						}
					}
				}
			}
		}
	*/
	return rightPath
}

func GenomeAssemblerKaushik(reads []string, minMatchLength, indexLength int, errorRate float64, k int) []string {
	if len(reads) == 0 {
		panic("Error: No reads given to GenomeAssembler.")
	}

	if minMatchLength <= indexLength {
		panic("Error: minMatchLength must be bigger than indexLength.")
	}

	contigs := make([]string, 0)

	//pathMap := make(map[int][]pathIncrement, 0)
	pathMap := make(map[int][]int, 0)

	fmt.Println("Building a prefix and suffix index for reads.")
	prefixIndex := BuildPrefixIndex(reads, indexLength)
	fmt.Println("Prefix index built!")
	suffixIndex := BuildSuffixIndex(reads, indexLength)
	fmt.Println("Suffix index built!")

	currentReadIndex := 0                  // or whatever
	currentRead := reads[currentReadIndex] // get corresponding read

	// idea: whenever we use a read, let's delete it from the prefix index (and suffix index).
	// continue for as long as we have elements still in the prefix index.
	i := 0
	for i < len(reads) {
		i++
		// let's throw out the elements of the indices corresponding to current read.
		//prefix := currentRead[:indexLength]
		//suffix := currentRead[len(currentRead)-indexLength:]
		//delete(prefixIndex, prefix)
		//delete(suffixIndex, suffix)

		fmt.Println(i, "of", len(reads))

		//extend currentRead to right and extend to left as far as I can.
		rightReads := make([]int, 0)
		rightReads = append(rightReads, currentReadIndex)
		leftReads := make([]int, 0)
		contig1 := ExtendContigRightInexactKaushik(currentRead, prefixIndex, suffixIndex, reads, minMatchLength, indexLength, errorRate, k, &rightReads)
		contig2 := ExtendContigLeftInexactKaushik(currentRead, prefixIndex, suffixIndex, reads, minMatchLength, indexLength, errorRate, k, &leftReads)

		// join into one contig and append to our set
		contig := contig2 + contig1[len(currentRead):]
		nodeArray := leftReads
		for _, i := range rightReads {
			nodeArray = append(nodeArray, i)
		}

		//previously, we appended every contig we found, even if it wasn't good (i.e., short).
		//because coverage is high, let's just keep longer contigs.
		if len(contig) > 100000 {

			contigs = append(contigs, contig)
			pathMap[nodeArray[0]] = nodeArray
			fmt.Println("We have generated", len(contigs), "contigs.")
			fmt.Println("Prefix index is down to", len(prefixIndex), "elements.")
			//for _, nodeID := range nodeArray {
			//	fmt.Print(nodeID, "")
			//}
			fmt.Println()
		}

		// we need a new starting point (currentRead) if still stuff in prefix index
		if len(prefixIndex) > 0 {
			// note: we know which reads haven't been used!
			// they're the elements still in the prefix index.
			// so just range over the prefix index, grab the first thing we see, and break
			for prefix := range prefixIndex {
				if len(prefixIndex[prefix]) != 0 {
					currentReadIndex = (prefixIndex[prefix])[0]
					currentRead = reads[currentReadIndex]
					break // stop as soon as we grab a value
				}
			}
		}
	}

	//for _, arr := range pathMap {
	//	for _, nodeID := range arr {
	//		fmt.Print(nodeID, " ")
	//	}
	//	fmt.Println()
	//}

	PrintMapStats(pathMap)

	return contigs
}

func ExtendContigRightInexactKaushik(currentRead string, prefixIndex, suffixIndex map[string][]int, reads []string, minMatchLength, indexLength int, errorRate float64, k int, readsInContig *[]int) string {
	contig := currentRead

	keepLooping := true
	// while we can keep going right
	for keepLooping == true {
		keepLooping = false
		// if we find anything, we will update it to true

		// range over all possible overlap lengths and pick the first place that we find a long, matching string.
		n := len(currentRead)
		for j := 1; j <= n-minMatchLength; j++ {
			// let's try overlapping this.
			prefix := currentRead[j : j+indexLength]
			// is this prefix present in the index?
			matchList, exists := prefixIndex[prefix]
			if exists {
				matchedRead := ""
				indiciesToBeDeleted := make([]int, 0)
				for i := range matchList {
					// grab first element as matching read
					//fmt.Println("	matches:", i, "of", len(matchList))
					matchedRead = reads[matchList[i]]
					//fmt.Println("matched read found")
					// does this string match completely? AND is it long enough?
					if len(matchedRead) > n-j && float64(CountSharedKmers(currentRead[j:], matchedRead[:n-j], k)) >= 0.9*float64(ExpectedSharedkmers(len(currentRead[j:]), errorRate, k)) {
						// success!
						keepLooping = true
						contig += matchedRead[n-j:]
						//fmt.Println("new read in contig", matchList[i])
						*readsInContig = append(*readsInContig, matchList[i])

						//add the index to be deleted
						indiciesToBeDeleted = append(indiciesToBeDeleted, i)

						if i == len(matchList)-1 {
							currentRead = matchedRead
							n = len(currentRead)
						}
						//readsInContig = append(readsInContig, matchList[0])
						// clean up the indices too by throwing out its prefix and suffix.
						//delete(prefixIndex, prefix)
						//suffix := currentRead[n-indexLength:] // what we overlapped
						//delete(suffixIndex, suffix)
					}
					//break // stop the outer looping process since we found a match.
				}

				//delete all the indicies from the matchList
				newSlice := make([]int, 0)
				for _, a := range prefixIndex[prefix] {
					if !(isInSlice(a, matchList)) {
						newSlice = append(newSlice, a)
					}
				}
				prefixIndex[prefix] = newSlice
			}
		}
	}

	return contig
}

func isInSlice(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func ExtendContigLeftInexactKaushik(currentRead string, prefixIndex, suffixIndex map[string][]int, reads []string, minMatchLength, indexLength int, errorRate float64, k int, readsInContig *[]int) string {
	contig := currentRead

	leftPathBack := make([]int, 0)

	keepLooping := true
	// while we can keep going right
	for keepLooping == true {
		keepLooping = false
		// if we find anything, we will update it to true
		// range over all possible overlap lengths and pick the first place that we find a long, matching string.
		n := len(currentRead)
		for j := 1; j <= n-minMatchLength; j++ { // j represents a count from right end of string
			// let's try overlapping this.
			suffix := currentRead[n-j-indexLength : n-j]
			// is this prefix present in the index?
			matchList, exists := suffixIndex[suffix]
			if exists {
				matchedRead := ""
				indiciesToBeDeleted := make([]int, 0)
				for i := range matchList {
					// grab first element as matching read
					//fmt.Println("	matches:", i, "of", len(matchList))
					matchedRead = reads[matchList[i]]
					// does this string match completely? AND is it long enough?
					if len(matchedRead) > n-j && float64(CountSharedKmers(currentRead[:n-j], matchedRead[len(matchedRead)-(n-j):], k)) >= 0.9*float64(ExpectedSharedkmers(len(currentRead[:n-j]), errorRate, k)) {
						// success!
						keepLooping = true
						contig = matchedRead[:len(matchedRead)-(n-j)] + contig
						//fmt.Println("new read in contig", matchList[i])
						leftPathBack = append(leftPathBack, matchList[i])
						indiciesToBeDeleted = append(indiciesToBeDeleted, i)
						//fmt.Println("read added")
						if i == len(matchList)-1 {
							currentRead = matchedRead
							n = len(currentRead)
						}
						// clean up the indices too by throwing out its prefix and suffix.
						prefix := currentRead[:indexLength]
						delete(prefixIndex, prefix)
						delete(suffixIndex, suffix)
					}
					//break // stop the outer looping process since we found a match.
				}
				//delete all the indicies from the matchList
				newSlice := make([]int, 0)
				for _, a := range suffixIndex[suffix] {
					if !(isInSlice(a, matchList)) {
						newSlice = append(newSlice, a)
					}
				}
				suffixIndex[suffix] = newSlice
			}
		}
	}
	for i := len(leftPathBack) - 1; i >= 0; i-- {
		*readsInContig = append(*readsInContig, leftPathBack[i])
	}
	return contig
}

func ExtendContigRightInexactKaushik2(currentRead string, prefixIndex, suffixIndex map[string][]int, reads []string, minMatchLength, indexLength int, errorRate float64, k int, readsInContig *[]int) string {
	contig := currentRead

	// range over all possible overlap lengths and pick the first place that we find a long, matching string.
	n := len(currentRead)
	for j := 1; j <= n-minMatchLength; j++ {
		// let's try overlapping this.
		prefix := currentRead[j : j+indexLength]
		// is this prefix present in the index?
		matchList, exists := prefixIndex[prefix]
		if exists {
			for i := 0; i < 2; i++ {
				// grab first element as matching read
				matchedRead := reads[matchList[i]]
				// does this string match completely? AND is it long enough?
				if len(matchedRead) > n-j && float64(CountSharedKmers(currentRead[j:], matchedRead[:n-j], k)) >= 0.9*float64(ExpectedSharedkmers(len(currentRead[j:]), errorRate, k)) {
					// success!
					fmt.Println("new read in contig", matchList[i])
					*readsInContig = append(*readsInContig, matchList[i])
					contig += matchedRead[n-j:] + ExtendContigRightInexactKaushik2(matchedRead, prefixIndex, suffixIndex, reads, minMatchLength, indexLength, errorRate, k, readsInContig)
					fmt.Println("length of contig:", len(contig))
					currentRead = matchedRead
					n = len(currentRead)
					//readsInContig = append(readsInContig, matchList[0])
					//update currentRead and its length
					// clean up the indices too by throwing out its prefix and suffix.

				}
				break // stop the outer looping process since we found a match.
			}
			delete(prefixIndex, prefix)
			suffix := currentRead[n-indexLength:] // what we overlapped
			delete(suffixIndex, suffix)
		}
	}

	return contig
}

func ExtendContigLeftInexactKaushik2(currentRead string, prefixIndex, suffixIndex map[string][]int, reads []string, minMatchLength, indexLength int, errorRate float64, k int, readsInContig *[]int) string {
	contig := currentRead

	leftPathBack := make([]int, 0)
	// range over all possible overlap lengths and pick the first place that we find a long, matching string.
	n := len(currentRead)
	for j := 1; j <= n-minMatchLength; j++ { // j represents a count from right end of string
		// let's try overlapping this.
		suffix := currentRead[n-j-indexLength : n-j]
		// is this prefix present in the index?
		matchList, exists := suffixIndex[suffix]
		if exists {
			for i := 0; i < 2; i++ {
				// grab first element as matching read
				matchedRead := reads[matchList[i]]
				// does this string match completely? AND is it long enough?
				if len(matchedRead) > n-j && float64(CountSharedKmers(currentRead[:n-j], matchedRead[len(matchedRead)-(n-j):], k)) >= 0.9*float64(ExpectedSharedkmers(len(currentRead[:n-j]), errorRate, k)) {
					// success!
					fmt.Println("new read in contig", matchList[i])
					leftPathBack = append(leftPathBack, matchList[i])
					contig = matchedRead[:len(matchedRead)-(n-j)] + contig + ExtendContigLeftInexactKaushik2(matchedRead, prefixIndex, suffixIndex, reads, minMatchLength, indexLength, errorRate, k, readsInContig)
					currentRead = matchedRead
					n = len(currentRead)
					//update currentRead and its length
					// clean up the indices too by throwing out its prefix and suffix.

				}
				break // stop the outer looping process since we found a match.
			}
			prefix := currentRead[:indexLength]
			delete(prefixIndex, prefix)
			delete(suffixIndex, suffix)
		}
	}

	for i := len(leftPathBack) - 1; i >= 0; i-- {
		*readsInContig = append(*readsInContig, leftPathBack[i])
	}
	return contig
}

func PrintMapStats(pathMap map[int][]int) {
	for i := range pathMap {
		fmt.Println("Root:", i)
		fmt.Println("Length of arr:", len(pathMap[i]))
		fmt.Println()
	}
}

type pathIncrement struct {
	node1        int
	node2        int
	overlapIndex int
}

//func GenerateContigs(pathMap map[int][]int) []strings {

//}

func GenomeAssemblerKaushik2(reads []string, minMatchLength, indexLength int, errorRate float64, k int) []string {
	if len(reads) == 0 {
		panic("Error: No reads given to GenomeAssembler.")
	}

	if minMatchLength <= indexLength {
		panic("Error: minMatchLength must be bigger than indexLength.")
	}

	contigs := make([]string, 0)

	//pathMap := make(map[int][]pathIncrement, 0)
	pathMap := make(map[int][]int, 0)

	fmt.Println("Building a prefix and suffix index for reads.")
	prefixIndex := BuildPrefixIndex(reads, indexLength)
	fmt.Println("Prefix index built!")
	suffixIndex := BuildSuffixIndex(reads, indexLength)
	fmt.Println("Suffix index built!")

	currentReadIndex := 0                  // or whatever
	currentRead := reads[currentReadIndex] // get corresponding read

	// idea: whenever we use a read, let's delete it from the prefix index (and suffix index).
	// continue for as long as we have elements still in the prefix index.
	i := 0
	for i < len(reads) {
		i++
		// let's throw out the elements of the indices corresponding to current read.
		//prefix := currentRead[:indexLength]
		//suffix := currentRead[len(currentRead)-indexLength:]
		//delete(prefixIndex, prefix)
		//delete(suffixIndex, suffix)

		//extend currentRead to right and extend to left as far as I can.
		rightReads := make([]int, 0)
		rightReads = append(rightReads, currentReadIndex)
		leftReads := make([]int, 0)
		contig1 := ExtendContigRightInexactKaushikk(currentRead, prefixIndex, suffixIndex, reads, minMatchLength, indexLength, errorRate, k, &rightReads)
		contig2 := ExtendContigLeftInexactKaushikk(currentRead, prefixIndex, suffixIndex, reads, minMatchLength, indexLength, errorRate, k, &leftReads)

		// join into one contig and append to our set
		contig := contig2 + contig1[len(currentRead):]
		nodeArray := leftReads
		for _, i := range rightReads {
			nodeArray = append(nodeArray, i)
		}

		//previously, we appended every contig we found, even if it wasn't good (i.e., short).
		//because coverage is high, let's just keep longer contigs.
		if len(contig) > 100000 {
			contigs = append(contigs, contig)
			pathMap[nodeArray[0]] = nodeArray
			fmt.Println("We have generated", len(contigs), "contigs.")
			fmt.Println("Prefix index is down to", len(prefixIndex), "elements.")
			//for _, nodeID := range nodeArray {
			//	fmt.Print(nodeID, "")
			//}
			fmt.Println()
		}

		// we need a new starting point (currentRead) if still stuff in prefix index
		if len(prefixIndex) > 0 {
			// note: we know which reads haven't been used!
			// they're the elements still in the prefix index.
			// so just range over the prefix index, grab the first thing we see, and break
			for prefix := range prefixIndex {
				currentReadIndex = (prefixIndex[prefix])[0]
				currentRead = reads[currentReadIndex]
				break // stop as soon as we grab a value
			}
		}
	}

	//for _, arr := range pathMap {
	//	for _, nodeID := range arr {
	//		fmt.Print(nodeID, " ")
	//	}
	//	fmt.Println()
	//}

	PrintMapStats(pathMap)

	return contigs
}

func ExtendContigRightInexactKaushikk(currentRead string, prefixIndex, suffixIndex map[string][]int, reads []string, minMatchLength, indexLength int, errorRate float64, k int, readsInContig *[]int) string {
	contig := currentRead

	keepLooping := true
	// while we can keep going right
	for keepLooping == true {
		keepLooping = false
		// if we find anything, we will update it to true

		// range over all possible overlap lengths and pick the first place that we find a long, matching string.
		n := len(currentRead)
		for j := 1; j <= n-minMatchLength; j++ {
			// let's try overlapping this.
			prefix := currentRead[j : j+indexLength]
			// is this prefix present in the index?
			matchList, exists := prefixIndex[prefix]
			if exists {
				for i := range matchList {
					// grab first element as matching read
					matchedRead := reads[matchList[i]]
					// does this string match completely? AND is it long enough?
					if len(matchedRead) > n-j && float64(CountSharedKmers(currentRead[j:], matchedRead[:n-j], k)) >= 0.9*float64(ExpectedSharedkmers(len(currentRead[j:]), errorRate, k)) {
						// success!
						keepLooping = true
						contig += matchedRead[n-j:]
						//fmt.Println("new read in contig", matchList[i])
						*readsInContig = append(*readsInContig, matchList[i])
						//readsInContig = append(readsInContig, matchList[0])
						//update currentRead and its length
						currentRead = matchedRead
						n = len(currentRead)
						// clean up the indices too by throwing out its prefix and suffix.
						//delete(prefixIndex, prefix)
						//suffix := currentRead[n-indexLength:] // what we overlapped
						//delete(suffixIndex, suffix)
					}
					break // stop the outer looping process since we found a match.
				}
			}
		}
	}

	return contig
}

func ExtendContigLeftInexactKaushikk(currentRead string, prefixIndex, suffixIndex map[string][]int, reads []string, minMatchLength, indexLength int, errorRate float64, k int, readsInContig *[]int) string {
	contig := currentRead

	leftPathBack := make([]int, 0)

	keepLooping := true
	// while we can keep going right
	for keepLooping == true {
		keepLooping = false
		// if we find anything, we will update it to true
		// range over all possible overlap lengths and pick the first place that we find a long, matching string.
		n := len(currentRead)
		for j := 1; j <= n-minMatchLength; j++ { // j represents a count from right end of string
			// let's try overlapping this.
			suffix := currentRead[n-j-indexLength : n-j]
			// is this prefix present in the index?
			matchList, exists := suffixIndex[suffix]
			if exists {
				for i := range matchList {
					// grab first element as matching read
					matchedRead := reads[matchList[i]]
					// does this string match completely? AND is it long enough?
					if len(matchedRead) > n-j && float64(CountSharedKmers(currentRead[:n-j], matchedRead[len(matchedRead)-(n-j):], k)) >= 0.9*float64(ExpectedSharedkmers(len(currentRead[:n-j]), errorRate, k)) {
						// success!
						keepLooping = true
						contig = matchedRead[:len(matchedRead)-(n-j)] + contig
						//fmt.Println("new read in contig", matchList[i])
						leftPathBack = append(leftPathBack, matchList[i])
						//update currentRead and its length
						currentRead = matchedRead
						n = len(currentRead)
						// clean up the indices too by throwing out its prefix and suffix.
						//prefix := currentRead[:indexLength]
						//delete(prefixIndex, prefix)
						//delete(suffixIndex, suffix)
					}
					break // stop the outer looping process since we found a match.
				}
			}
		}
	}
	for i := len(leftPathBack) - 1; i >= 0; i-- {
		*readsInContig = append(*readsInContig, leftPathBack[i])
	}
	return contig
}

func GenomeAssemblerKaushik3(reads []string, minMatchLength, indexLength int, errorRate float64, k int) []string {
	if len(reads) == 0 {
		panic("Error: No reads given to GenomeAssembler.")
	}

	if minMatchLength <= indexLength {
		panic("Error: minMatchLength must be bigger than indexLength.")
	}

	contigs := make([]string, 0)

	//pathMap := make(map[int][]pathIncrement, 0)
	pathMap := make(map[int][]int, 0)

	fmt.Println("Building a prefix and suffix index for reads.")
	prefixIndex := BuildPrefixIndex(reads, indexLength)
	fmt.Println("Prefix index built!")
	suffixIndex := BuildSuffixIndex(reads, indexLength)
	fmt.Println("Suffix index built!")

	currentReadIndex := 0                  // or whatever
	currentRead := reads[currentReadIndex] // get corresponding read

	// idea: whenever we use a read, let's delete it from the prefix index (and suffix index).
	// continue for as long as we have elements still in the prefix index.
	i := 0
	for i < len(reads) {
		i++
		// let's throw out the elements of the indices corresponding to current read.
		//prefix := currentRead[:indexLength]
		//suffix := currentRead[len(currentRead)-indexLength:]
		//delete(prefixIndex, prefix)
		//delete(suffixIndex, suffix)

		fmt.Println(i, "of", len(reads))

		//extend currentRead to right and extend to left as far as I can.
		rightReads := make([]int, 0)
		rightReads = append(rightReads, currentReadIndex)
		leftReads := make([]int, 0)
		contig1 := ExtendContigRightInexactKaushikRec(currentRead, prefixIndex, suffixIndex, reads, minMatchLength, indexLength, errorRate, k, &rightReads)
		contig2 := ExtendContigLeftInexactKaushikRec(currentRead, prefixIndex, suffixIndex, reads, minMatchLength, indexLength, errorRate, k, &leftReads)

		// join into one contig and append to our set
		contig := contig2 + contig1[len(currentRead):]
		nodeArray := leftReads
		for _, i := range rightReads {
			nodeArray = append(nodeArray, i)
		}

		//previously, we appended every contig we found, even if it wasn't good (i.e., short).
		//because coverage is high, let's just keep longer contigs.
		if len(contig) > 300000 {

			contigs = append(contigs, contig)
			pathMap[nodeArray[0]] = nodeArray
			fmt.Println("We have generated", len(contigs), "contigs.")
			fmt.Println("Prefix index is down to", len(prefixIndex), "elements.")
			//for _, nodeID := range nodeArray {
			//	fmt.Print(nodeID, "")
			//}
			fmt.Println()
		}

		// we need a new starting point (currentRead) if still stuff in prefix index
		if len(prefixIndex) > 0 {
			// note: we know which reads haven't been used!
			// they're the elements still in the prefix index.
			// so just range over the prefix index, grab the first thing we see, and break
			for prefix := range prefixIndex {
				if len(prefixIndex[prefix]) != 0 {
					currentReadIndex = (prefixIndex[prefix])[0]
					currentRead = reads[currentReadIndex]
					break // stop as soon as we grab a value
				}
			}
		}
	}

	//for _, arr := range pathMap {
	//	for _, nodeID := range arr {
	//		fmt.Print(nodeID, " ")
	//	}
	//	fmt.Println()
	//}

	PrintMapStats(pathMap)

	return contigs
}

func ExtendContigRightInexactKaushik3(currentRead string, prefixIndex, suffixIndex map[string][]int, reads []string, minMatchLength, indexLength int, errorRate float64, k int, readsInContig *[]int) string {
	contig := currentRead

	keepLooping := true
	// while we can keep going right
	for keepLooping == true {
		keepLooping = false
		// if we find anything, we will update it to true

		// range over all possible overlap lengths and pick the first place that we find a long, matching string.
		n := len(currentRead)
		for j := 1; j <= n-minMatchLength; j++ {
			// let's try overlapping this.
			prefix := currentRead[j : j+indexLength]
			// is this prefix present in the index?
			matchList, exists := prefixIndex[prefix]
			if exists {
				matchedRead := ""
				indiciesToBeDeleted := make([]int, 0)
				for i := range matchList {
					// grab first element as matching read
					//fmt.Println("	matches:", i, "of", len(matchList))
					if !(isInSlice(matchList[i], *readsInContig)) {
						//fmt.Println("	matched read found:", matchList[i])
						matchedRead = reads[matchList[i]]
						// does this string match completely? AND is it long enough?
						if len(matchedRead) > n-j && float64(CountSharedKmers(currentRead[j:], matchedRead[:n-j], k)) >= 0.9*float64(ExpectedSharedkmers(len(currentRead[j:]), errorRate, k)) {
							// success!
							keepLooping = true
							contig += matchedRead[n-j:]
							fmt.Println("	new read in contig", matchList[i])
							*readsInContig = append(*readsInContig, matchList[i])

							//add the index to be deleted
							indiciesToBeDeleted = append(indiciesToBeDeleted, i)

							if i == len(matchList)-1 {
								currentRead = matchedRead
								n = len(currentRead)
							}
							//readsInContig = append(readsInContig, matchList[0])
							// clean up the indices too by throwing out its prefix and suffix.
							//delete(prefixIndex, prefix)
							//suffix := currentRead[n-indexLength:] // what we overlapped
							//delete(suffixIndex, suffix)
						}
						//break // stop the outer looping process since we found a match.
					}
				}

				//delete all the indicies from the matchList
				newSlice := make([]int, 0)
				for _, a := range prefixIndex[prefix] {
					if !(isInSlice(a, indiciesToBeDeleted)) {
						newSlice = append(newSlice, a)
					}
				}
				prefixIndex[prefix] = newSlice
			}
		}
	}

	return contig
}

func ExtendContigLeftInexactKaushik3(currentRead string, prefixIndex, suffixIndex map[string][]int, reads []string, minMatchLength, indexLength int, errorRate float64, k int, readsInContig *[]int) string {
	contig := currentRead

	leftPathBack := make([]int, 0)

	keepLooping := true
	// while we can keep going right
	for keepLooping == true {
		keepLooping = false
		// if we find anything, we will update it to true
		// range over all possible overlap lengths and pick the first place that we find a long, matching string.
		n := len(currentRead)
		for j := 1; j <= n-minMatchLength; j++ { // j represents a count from right end of string
			// let's try overlapping this.
			suffix := currentRead[n-j-indexLength : n-j]
			// is this prefix present in the index?
			matchList, exists := suffixIndex[suffix]
			if exists {
				matchedRead := ""
				indiciesToBeDeleted := make([]int, 0)
				for i := range matchList {
					// grab first element as matching read
					//fmt.Println("	matches:", i, "of", len(matchList))
					if !(isInSlice(matchList[i], *readsInContig)) {
						//fmt.Println("matched read found:", matchList[i])
						matchedRead = reads[matchList[i]]
						// does this string match completely? AND is it long enough?
						if len(matchedRead) > n-j && float64(CountSharedKmers(currentRead[:n-j], matchedRead[len(matchedRead)-(n-j):], k)) >= 0.9*float64(ExpectedSharedkmers(len(currentRead[:n-j]), errorRate, k)) {
							// success!
							keepLooping = true
							fmt.Println("	new read in contig", matchList[i])
							leftPathBack = append(leftPathBack, matchList[i])
							indiciesToBeDeleted = append(indiciesToBeDeleted, i)
							contig = matchedRead[:len(matchedRead)-(n-j)] + contig
							//fmt.Println("read added")
							if i == len(matchList)-1 {
								currentRead = matchedRead
								n = len(currentRead)
							}
							// clean up the indices too by throwing out its prefix and suffix.
							//prefix := currentRead[:indexLength]
							//delete(prefixIndex, prefix)
							//delete(suffixIndex, suffix)
						}
						//break // stop the outer looping process since we found a match.
					}
				}
				//delete all the indicies from the matchList
				newSlice := make([]int, 0)
				for _, a := range suffixIndex[suffix] {
					if !(isInSlice(a, indiciesToBeDeleted)) {
						newSlice = append(newSlice, a)
					}
				}
				suffixIndex[suffix] = newSlice
			}
		}
	}
	for i := len(leftPathBack) - 1; i >= 0; i-- {
		*readsInContig = append(*readsInContig, leftPathBack[i])
	}
	return contig
}

func CreateReadNetwork5(reads []string, minMatchLength, indexLength int, errorRate float64, k int) Graph {
	if len(reads) == 0 {
		panic("Error: No reads given to GenomeAssembler.")
	}

	if minMatchLength <= indexLength {
		panic("Error: minMatchLength must be bigger than indexLength.")
	}

	network := MakeGraph()
	pointerToNetwork := &network

	for i, read := range reads {
		fmt.Println("Adding read number", i+1)
		pointerToNetwork.addNode(i, read)
	}

	contigs := make([]string, 0)

	//pathMap := make(map[int][]pathIncrement, 0)
	pathMap := make(map[int][]int, 0)

	fmt.Println("Building a prefix and suffix index for reads.")
	prefixIndex := BuildPrefixIndex(reads, indexLength)
	fmt.Println("Prefix index built!")
	suffixIndex := BuildSuffixIndex(reads, indexLength)
	fmt.Println("Suffix index built!")

	currentReadIndex := 0                  // or whatever
	currentRead := reads[currentReadIndex] // get corresponding read

	// idea: whenever we use a read, let's delete it from the prefix index (and suffix index).
	// continue for as long as we have elements still in the prefix index.
	i := 0
	for i < len(reads) {
		i++
		// let's throw out the elements of the indices corresponding to current read.
		//prefix := currentRead[:indexLength]
		//suffix := currentRead[len(currentRead)-indexLength:]
		//delete(prefixIndex, prefix)
		//delete(suffixIndex, suffix)

		fmt.Println(i, "of", len(reads))

		//extend currentRead to right and extend to left as far as I can.
		rightReads := make([]int, 0)
		rightReads = append(rightReads, currentReadIndex)
		leftReads := make([]int, 0)
		contig1 := ExtendContigRightInexactKaushikRec(currentRead, prefixIndex, suffixIndex, reads, minMatchLength, indexLength, errorRate, k, &rightReads /*, pointerToNetwork*/)
		contig2 := ExtendContigLeftInexactKaushikRec(currentRead, prefixIndex, suffixIndex, reads, minMatchLength, indexLength, errorRate, k, &leftReads /*, pointerToNetwork*/)

		// join into one contig and append to our set
		contig := contig2 + contig1[len(currentRead):]
		nodeArray := leftReads
		for _, i := range rightReads {
			nodeArray = append(nodeArray, i)
		}

		//previously, we appended every contig we found, even if it wasn't good (i.e., short).
		//because coverage is high, let's just keep longer contigs.
		if len(contig) > 100000 {

			contigs = append(contigs, contig)
			pathMap[nodeArray[0]] = nodeArray
			fmt.Println("We have generated", len(contigs), "contigs.")
			fmt.Println("Prefix index is down to", len(prefixIndex), "elements.")
			//for _, nodeID := range nodeArray {
			//	fmt.Print(nodeID, "")
			//}
			fmt.Println()
		}

		// we need a new starting point (currentRead) if still stuff in prefix index
		if len(prefixIndex) > 0 {
			// note: we know which reads haven't been used!
			// they're the elements still in the prefix index.
			// so just range over the prefix index, grab the first thing we see, and break
			for prefix := range prefixIndex {
				if len(prefixIndex[prefix]) != 0 {
					currentReadIndex = (prefixIndex[prefix])[0]
					currentRead = reads[currentReadIndex]
					break // stop as soon as we grab a value
				}
			}
		}
	}

	//for _, arr := range pathMap {
	//	for _, nodeID := range arr {
	//		fmt.Print(nodeID, " ")
	//	}
	//	fmt.Println()
	//}

	PrintMapStats(pathMap)

	return network
}

func ExtendContigRightInexactKaushik4(currentRead string, currentReadIndex int, prefixIndex, suffixIndex map[string][]int, reads []string, minMatchLength, indexLength int, errorRate float64, k int, readsInContig *[]int, network *Graph) string {
	contig := currentRead

	keepLooping := true
	// while we can keep going right
	for keepLooping == true {
		keepLooping = false
		// if we find anything, we will update it to true

		// range over all possible overlap lengths and pick the first place that we find a long, matching string.
		n := len(currentRead)
		for j := 1; j <= n-minMatchLength; j++ {
			// let's try overlapping this.
			prefix := currentRead[j : j+indexLength]
			// is this prefix present in the index?
			matchList, exists := prefixIndex[prefix]
			if exists {
				matchedRead := ""
				//indiciesToBeDeleted := make([]int, 0)
				for i := range matchList {
					// grab first element as matching read
					//fmt.Println("	matches:", i, "of", len(matchList))
					if !(isInSlice(matchList[i], *readsInContig)) {
						matchedRead = reads[matchList[i]]
						//fmt.Println("matched read found")
						// does this string match completely? AND is it long enough?
						if len(matchedRead) > n-j && float64(CountSharedKmers(currentRead[j:], matchedRead[:n-j], k)) >= 0.9*float64(ExpectedSharedkmers(len(currentRead[j:]), errorRate, k)) {
							// success!
							keepLooping = true
							contig += matchedRead[n-j:]
							//fmt.Println("new read in contig", matchList[i])
							*readsInContig = append(*readsInContig, matchList[i])
							network.addEdge(currentReadIndex, matchList[i])
							//add the index to be deleted
							//indiciesToBeDeleted = append(indiciesToBeDeleted, i)

							if i == len(matchList)-1 {
								currentRead = matchedRead
								n = len(currentRead)
							}
							//readsInContig = append(readsInContig, matchList[0])
							// clean up the indices too by throwing out its prefix and suffix.
							//delete(prefixIndex, prefix)
							//suffix := currentRead[n-indexLength:] // what we overlapped
							//delete(suffixIndex, suffix)
						}
						//break // stop the outer looping process since we found a match.
					}
				}

				//delete all the indicies from the matchList
				//newSlice := make([]int, 0)
				//for _, a := range prefixIndex[prefix] {
				//	if !(isInSlice(a, matchList)) {
				//		newSlice = append(newSlice, a)
				//	}
				//}
				//prefixIndex[prefix] = newSlice
			}
		}
	}

	return contig
}

func ExtendContigLeftInexactKaushik4(currentRead string, currentReadIndex int, prefixIndex, suffixIndex map[string][]int, reads []string, minMatchLength, indexLength int, errorRate float64, k int, readsInContig *[]int, network *Graph) string {
	contig := currentRead

	leftPathBack := make([]int, 0)

	keepLooping := true
	// while we can keep going right
	for keepLooping == true {
		keepLooping = false
		// if we find anything, we will update it to true
		// range over all possible overlap lengths and pick the first place that we find a long, matching string.
		n := len(currentRead)
		for j := 1; j <= n-minMatchLength; j++ { // j represents a count from right end of string
			// let's try overlapping this.
			suffix := currentRead[n-j-indexLength : n-j]
			// is this prefix present in the index?
			matchList, exists := suffixIndex[suffix]
			if exists {
				matchedRead := ""
				indiciesToBeDeleted := make([]int, 0)
				for i := range matchList {
					// grab first element as matching read
					//fmt.Println("	matches:", i, "of", len(matchList))
					if !(isInSlice(matchList[i], *readsInContig)) {
						matchedRead = reads[matchList[i]]
						// does this string match completely? AND is it long enough?
						if len(matchedRead) > n-j && float64(CountSharedKmers(currentRead[:n-j], matchedRead[len(matchedRead)-(n-j):], k)) >= 0.9*float64(ExpectedSharedkmers(len(currentRead[:n-j]), errorRate, k)) {
							// success!
							keepLooping = true
							contig = matchedRead[:len(matchedRead)-(n-j)] + contig
							//fmt.Println("new read in contig", matchList[i])
							leftPathBack = append(leftPathBack, matchList[i])
							indiciesToBeDeleted = append(indiciesToBeDeleted, i)
							//fmt.Println("read added")
							if i == len(matchList)-1 {
								currentRead = matchedRead
								n = len(currentRead)
							}
							// clean up the indices too by throwing out its prefix and suffix.
							//prefix := currentRead[:indexLength]
							//delete(prefixIndex, prefix)
							//delete(suffixIndex, suffix)
						}
						//break // stop the outer looping process since we found a match.
					}
				}
				//delete all the indicies from the matchList
				//newSlice := make([]int, 0)
				//for _, a := range suffixIndex[suffix] {
				//	if !(isInSlice(a, matchList)) {
				//		newSlice = append(newSlice, a)
				//	}
				//}
				//suffixIndex[suffix] = newSlice
			}
		}
	}
	for i := len(leftPathBack) - 1; i >= 0; i-- {
		*readsInContig = append(*readsInContig, leftPathBack[i])
	}
	return contig
}

func ExtendContigRightInexactKaushikRec(currentRead string, prefixIndex, suffixIndex map[string][]int, reads []string, minMatchLength, indexLength int, errorRate float64, k int, readsInContig *[]int) string {
	contig := currentRead

	keepLooping := true
	// while we can keep going right
	for keepLooping == true {
		keepLooping = false
		// if we find anything, we will update it to true

		// range over all possible overlap lengths and pick the first place that we find a long, matching string.
		n := len(currentRead)
		for j := 1; j <= n-minMatchLength; j++ {
			// let's try overlapping this.
			prefix := currentRead[j : j+indexLength]
			// is this prefix present in the index?
			matchList, exists := prefixIndex[prefix]
			if exists {
				matchedRead := ""
				indiciesToBeDeleted := make([]int, 0)
				for i := range matchList {
					// grab first element as matching read
					fmt.Println("	matches:", i, "of", len(matchList))
					if !(isInSlice(matchList[i], *readsInContig)) {
						fmt.Println("	matched read found:", matchList[i])
						matchedRead = reads[matchList[i]]
						// does this string match completely? AND is it long enough?
						if len(matchedRead) > n-j && float64(CountSharedKmers(currentRead[j:], matchedRead[:n-j], k)) >= 0.9*float64(ExpectedSharedkmers(len(currentRead[j:]), errorRate, k)) {
							// success!
							keepLooping = true
							fmt.Println("	new read in contig", matchList[i])
							indiciesToBeDeleted = append(indiciesToBeDeleted, i)
							*readsInContig = append(*readsInContig, matchList[i])
							contig += matchedRead[n-j:]
							contig += FindRestOfPathRight(contig, prefixIndex, suffixIndex, reads, minMatchLength, indexLength, errorRate, k, readsInContig)

							if i == len(matchList)-1 {
								currentRead = matchedRead
								n = len(currentRead)
							}
							//readsInContig = append(readsInContig, matchList[0])
							// clean up the indices too by throwing out its prefix and suffix.
							//delete(prefixIndex, prefix)
							//suffix := currentRead[n-indexLength:] // what we overlapped
							//delete(suffixIndex, suffix)
						}
						//break // stop the outer looping process since we found a match.
					}
				}

				//delete all the indicies from the matchList
				newSlice := make([]int, 0)
				for _, a := range prefixIndex[prefix] {
					if !(isInSlice(a, indiciesToBeDeleted)) {
						newSlice = append(newSlice, a)
					}
				}
				prefixIndex[prefix] = newSlice
			}
		}
	}

	return contig
}

func ExtendContigLeftInexactKaushikRec(currentRead string, prefixIndex, suffixIndex map[string][]int, reads []string, minMatchLength, indexLength int, errorRate float64, k int, readsInContig *[]int) string {
	contig := currentRead

	leftPathBack := make([]int, 0)

	keepLooping := true
	// while we can keep going right
	for keepLooping == true {
		keepLooping = false
		// if we find anything, we will update it to true
		// range over all possible overlap lengths and pick the first place that we find a long, matching string.
		n := len(currentRead)
		for j := 1; j <= n-minMatchLength; j++ { // j represents a count from right end of string
			// let's try overlapping this.
			suffix := currentRead[n-j-indexLength : n-j]
			// is this prefix present in the index?
			matchList, exists := suffixIndex[suffix]
			if exists {
				matchedRead := ""
				indiciesToBeDeleted := make([]int, 0)
				for i := range matchList {
					// grab first element as matching read
					fmt.Println("	matches:", i, "of", len(matchList))
					if !(isInSlice(matchList[i], *readsInContig)) && !(isInSlice(matchList[i], leftPathBack)) {
						fmt.Println("matched read found:", matchList[i])
						matchedRead = reads[matchList[i]]
						// does this string match completely? AND is it long enough?
						if len(matchedRead) > n-j && float64(CountSharedKmers(currentRead[:n-j], matchedRead[len(matchedRead)-(n-j):], k)) >= 0.9*float64(ExpectedSharedkmers(len(currentRead[:n-j]), errorRate, k)) {
							// success!
							keepLooping = true
							fmt.Println("	new read in contig", matchList[i])
							leftPathBack = append(leftPathBack, matchList[i])
							indiciesToBeDeleted = append(indiciesToBeDeleted, i)
							contig = matchedRead[:len(matchedRead)-(n-j)] + contig
							contig = FindRestOfPathLeft(contig, prefixIndex, suffixIndex, reads, minMatchLength, indexLength, errorRate, k, readsInContig) + contig

							//fmt.Println("read added")
							if i == len(matchList)-1 {
								currentRead = matchedRead
								n = len(currentRead)
							}
							// clean up the indices too by throwing out its prefix and suffix.
							//prefix := currentRead[:indexLength]
							//delete(prefixIndex, prefix)
							//delete(suffixIndex, suffix)
						}
						//break // stop the outer looping process since we found a match.
					}
				}
				//delete all the indicies from the matchList
				newSlice := make([]int, 0)
				for _, a := range suffixIndex[suffix] {
					if !(isInSlice(a, indiciesToBeDeleted)) {
						newSlice = append(newSlice, a)
					}
				}
				suffixIndex[suffix] = newSlice
			}
		}
	}
	for i := len(leftPathBack) - 1; i >= 0; i-- {
		*readsInContig = append(*readsInContig, leftPathBack[i])
	}
	return contig
}

func FindRestOfPathLeft(currentRead string, prefixIndex, suffixIndex map[string][]int, reads []string, minMatchLength, indexLength int, errorRate float64, k int, readsInContig *[]int) string {
	contig := currentRead
	n := len(currentRead)
	leftPathBack := make([]int, 0)

	for j := 1; j <= n-minMatchLength; j++ { // j represents a count from right end of string
		// let's try overlapping this.
		suffix := currentRead[n-j-indexLength : n-j]
		// is this prefix present in the index?
		matchList, exists := suffixIndex[suffix]
		if exists {
			matchedRead := ""
			indiciesToBeDeleted := make([]int, 0)
			for i := range matchList {
				// grab first element as matching read
				fmt.Println("	matches:", i, "of", len(matchList))
				if !(isInSlice(matchList[i], *readsInContig)) && !(isInSlice(matchList[i], leftPathBack)) {
					fmt.Println("matched read found:", matchList[i])
					matchedRead = reads[matchList[i]]
					// does this string match completely? AND is it long enough?
					if len(matchedRead) > n-j && float64(CountSharedKmers(currentRead[:n-j], matchedRead[len(matchedRead)-(n-j):], k)) >= 0.9*float64(ExpectedSharedkmers(len(currentRead[:n-j]), errorRate, k)) {
						// success!
						fmt.Println("	new read in contig", matchList[i])
						leftPathBack = append(leftPathBack, matchList[i])
						indiciesToBeDeleted = append(indiciesToBeDeleted, i)
						contig = matchedRead[:len(matchedRead)-(n-j)] + contig
						contig = FindRestOfPathLeft(contig, prefixIndex, suffixIndex, reads, minMatchLength, indexLength, errorRate, k, readsInContig) + contig

						//fmt.Println("read added")
					} else {
						newSlice := make([]int, 0)
						for _, a := range suffixIndex[suffix] {
							if !(isInSlice(a, indiciesToBeDeleted)) {
								newSlice = append(newSlice, a)
							}
						}
						suffixIndex[suffix] = newSlice
						for i := len(leftPathBack) - 1; i >= 0; i-- {
							*readsInContig = append(*readsInContig, leftPathBack[i])
						}
						return contig
					}
					//break // stop the outer looping process since we found a match.
				}
			}
			//delete all the indicies from the matchList
			newSlice := make([]int, 0)
			for _, a := range suffixIndex[suffix] {
				if !(isInSlice(a, indiciesToBeDeleted)) {
					newSlice = append(newSlice, a)
				}
			}
			suffixIndex[suffix] = newSlice
		}
	}
	for i := len(leftPathBack) - 1; i >= 0; i-- {
		*readsInContig = append(*readsInContig, leftPathBack[i])
	}
	return contig
}

func FindRestOfPathRight(currentRead string, prefixIndex, suffixIndex map[string][]int, reads []string, minMatchLength, indexLength int, errorRate float64, k int, readsInContig *[]int) string {
	contig := currentRead
	n := len(currentRead)

	for j := 1; j <= n-minMatchLength; j++ { // j represents a count from right end of string
		// let's try overlapping this.
		prefix := currentRead[j : j+indexLength]
		// is this prefix present in the index?
		matchList, exists := prefixIndex[prefix]
		if exists {
			matchedRead := ""
			indiciesToBeDeleted := make([]int, 0)
			for i := range matchList {
				// grab first element as matching read
				fmt.Println("	matches:", i, "of", len(matchList))
				if !(isInSlice(matchList[i], *readsInContig)) {
					fmt.Println("matched read found:", matchList[i])
					matchedRead = reads[matchList[i]]
					// does this string match completely? AND is it long enough?
					if len(matchedRead) > n-j && float64(CountSharedKmers(currentRead[j:], matchedRead[:n-j], k)) >= 0.9*float64(ExpectedSharedkmers(len(currentRead[j:]), errorRate, k)) {
						// success!
						fmt.Println("	new read in contig", matchList[i])
						*readsInContig = append(*readsInContig, matchList[i])
						indiciesToBeDeleted = append(indiciesToBeDeleted, i)
						contig += matchedRead[n-j:]
						contig += FindRestOfPathLeft(contig, prefixIndex, suffixIndex, reads, minMatchLength, indexLength, errorRate, k, readsInContig)
						//fmt.Println("read added")
					} else {
						newSlice := make([]int, 0)
						for _, a := range prefixIndex[prefix] {
							if !(isInSlice(a, matchList)) {
								newSlice = append(newSlice, a)
							}
						}
						prefixIndex[prefix] = newSlice
						return contig
					}
					//break // stop the outer looping process since we found a match.
				}
			}
			//delete all the indicies from the matchList
			newSlice := make([]int, 0)
			for _, a := range prefixIndex[prefix] {
				if !(isInSlice(a, indiciesToBeDeleted)) {
					newSlice = append(newSlice, a)
				}
			}
			prefixIndex[prefix] = newSlice
		}
	}
	return contig
}

func GenomeAssemblerKaushik4(reads []string, minMatchLength, indexLength int, errorRate float64, k int) Graph {
	if len(reads) == 0 {
		panic("Error: No reads given to GenomeAssembler.")
	}

	if minMatchLength <= indexLength {
		panic("Error: minMatchLength must be bigger than indexLength.")
	}

	network := MakeGraph()
	pointerToNetwork := &network

	for i, read := range reads {
		fmt.Println("Adding node number", i+1)
		pointerToNetwork.addNode(i, read)
	}

	contigs := make([]string, 0)

	fmt.Println("Building a prefix and suffix index for reads.")
	prefixIndex := BuildPrefixIndex(reads, indexLength)
	fmt.Println("Prefix index built!")
	suffixIndex := BuildSuffixIndex(reads, indexLength)
	fmt.Println("Suffix index built!")

	// idea: whenever we use a read, let's delete it from the prefix index (and suffix index).
	// continue for as long as we have elements still in the prefix index.
	for i := 0; i < len(reads); i++ {
		// let's throw out the elements of the indices corresponding to current read.
		//prefix := currentRead[:indexLength]
		//suffix := currentRead[len(currentRead)-indexLength:]
		//delete(prefixIndex, prefix)
		//delete(suffixIndex, suffix)
		currentReadIndex := i                  // or whatever
		currentRead := reads[currentReadIndex] // get corresponding read
		currentNode := network.Nodes[currentReadIndex]

		//extend currentRead to right and extend to left as far as I can.
		contig1 := ExtendContigRightInexactKaushik5(currentNode, currentRead, prefixIndex, suffixIndex, reads, minMatchLength, indexLength, errorRate, k, pointerToNetwork)
		contig2 := ExtendContigLeftInexactKaushik5(currentNode, currentRead, prefixIndex, suffixIndex, reads, minMatchLength, indexLength, errorRate, k, pointerToNetwork)

		// join into one contig and append to our set
		contig := contig2 + contig1[len(currentRead):]

		//previously, we appended every contig we found, even if it wasn't good (i.e., short).
		//because coverage is high, let's just keep longer contigs.
		//if len(contig) > 100000 {
		contigs = append(contigs, contig)
		//	fmt.Println("We have generated", len(contigs), "contigs.")
		//	fmt.Println("Prefix index is down to", len(prefixIndex), "elements.")
		//}

		// we need a new starting point (currentRead) if still stuff in prefix index
		//if len(prefixIndex) > 0 {
		// note: we know which reads haven't been used!
		// they're the elements still in the prefix index.
		// so just range over the prefix index, grab the first thing we see, and break
		//	for prefix := range prefixIndex {
		//		currentReadIndex = (prefixIndex[prefix])[0]
		//		currentRead = reads[currentReadIndex]
		//		break // stop as soon as we grab a value
		//	}
		//}
		fmt.Println("done with", i, "nodes")
	}

	return network
}

func ExtendContigRightInexactKaushik5(node Node, currentRead string, prefixIndex, suffixIndex map[string][]int, reads []string, minMatchLength, indexLength int, errorRate float64, k int, graph *Graph) string {
	contig := currentRead

	//keepLooping := true
	// while we can keep going right
	//for keepLooping == true {
	//keepLooping = false
	// if we find anything, we will update it to true

	// range over all possible overlap lengths and pick the first place that we find a long, matching string.
	n := len(currentRead)
	for j := 1; j <= n-minMatchLength; j++ {
		// let's try overlapping this.
		prefix := currentRead[j : j+indexLength]
		// is this prefix present in the index?
		matchList, exists := prefixIndex[prefix]
		if exists {
			fmt.Println("found some matches")
			for _, i := range matchList {
				// grab first element as matching read
				matchedRead := reads[i]
				matchedNode := graph.Nodes[i]
				// does this string match completely? AND is it long enough?
				if len(matchedRead) > n-j && float64(CountSharedKmers(currentRead[j:], matchedRead[:n-j], k)) >= 0.9*float64(ExpectedSharedkmers(len(currentRead[j:]), errorRate, k)) {
					fmt.Println("adding edge")
					graph.addEdge(node.getID(), matchedNode.getID())
					// success!
					//keepLooping = true
					//contig += matchedRead[n-j:]
					//update currentRead and its length
					//currentRead = matchedRead
					//n = len(currentRead)
					// clean up the indices too by throwing out its prefix and suffix.
					//delete(prefixIndex, prefix)
					//suffix := currentRead[n-indexLength:] // what we overlapped
					//delete(suffixIndex, suffix)
					//break // stop the outer looping process since we found a match.
				}
			}

		}
	}
	//}

	return contig
}

func ExtendContigLeftInexactKaushik5(node Node, currentRead string, prefixIndex, suffixIndex map[string][]int, reads []string, minMatchLength, indexLength int, errorRate float64, k int, graph *Graph) string {
	contig := currentRead

	//keepLooping := true
	// while we can keep going right
	//for keepLooping == true {
	//	keepLooping = false
	// if we find anything, we will update it to true
	// range over all possible overlap lengths and pick the first place that we find a long, matching string.
	n := len(currentRead)
	for j := 1; j <= n-minMatchLength; j++ { // j represents a count from right end of string
		// let's try overlapping this.
		suffix := currentRead[n-j-indexLength : n-j]
		// is this prefix present in the index?
		matchList, exists := suffixIndex[suffix]
		if exists {
			fmt.Println("found some matches")
			for _, i := range matchList {
				// grab first element as matching read
				matchedRead := reads[i]
				matchedNode := graph.Nodes[i]
				// does this string match completely? AND is it long enough?
				if len(matchedRead) > n-j && float64(CountSharedKmers(currentRead[:n-j], matchedRead[len(matchedRead)-(n-j):], k)) >= 0.9*float64(ExpectedSharedkmers(len(currentRead[:n-j]), errorRate, k)) {
					// success!
					fmt.Println("adding edge")
					graph.addEdge(matchedNode.getID(), node.getID())
					//keepLooping = true
					//contig = matchedRead[:len(matchedRead)-(n-j)] + contig
					//update currentRead and its length
					//currentRead = matchedRead
					//n = len(currentRead)
					// clean up the indices too by throwing out its prefix and suffix.
					//prefix := currentRead[:indexLength]
					//delete(prefixIndex, prefix)
					//delete(suffixIndex, suffix)
					//break // stop the outer looping process since we found a match.
				}
			}
		}
	}
	//}
	return contig
}

// -- Now we are using the Graph2 structure

/*
func CreateReadNetwork2(reads []string, minMatchLength, k, indexLength int, errorRate float64) Graph2 {
	network := MakeGraph2()
	pointerToNetwork := &network
	for i, read := range reads {
		fmt.Println("Adding read number", i+1)
		pointerToNetwork.addNode2(i, read)
	}
	count := 0
	for id, node := range network.Nodes {
		fmt.Println("node:", id)
		n := len(node.read)
		expectedKmers := 0.9 * float64(ExpectedSharedkmers(len(node.read[:minMatchLength]), errorRate, k))
		expectedKmers2 := 0.9 * float64(ExpectedSharedkmers(len(node.read[n-minMatchLength:]), errorRate, k))
		BuildEdges2(node, minMatchLength, k, errorRate, pointerToNetwork, reads, expectedKmers, expectedKmers2)
		count++
		fmt.Println("done with", count, "of", len(network.Nodes))
	}
	return network
}
func BuildEdges2(node Node, minMatchLength, k int, errorRate float64, graph *Graph2, reads []string, expectedKmers, expectedKmers2 float64) {
	n := len(node.read)
	for i := range graph.Nodes {
		matchedNode := graph.Nodes[i]
		matchedRead := matchedNode.read
		//fmt.Println("testing:", node.getID(), "and", matchedNode.getID())
		if len(matchedRead) > minMatchLength && float64(CountSharedKmers(node.read[:minMatchLength], matchedRead[len(matchedRead)-minMatchLength:], k)) >= expectedKmers {
			// success!
			fmt.Println("found a match")
			if !(node.isANeighbor(matchedNode)) {
				fmt.Println("building an edge", node.getID(), matchedNode.getID())
				graph.addEdge2(matchedNode.getID(), node.getID(), node.read[:minMatchLength])
			}
		} else if len(matchedRead) > minMatchLength && float64(CountSharedKmers(node.read[n-minMatchLength:], matchedRead[:minMatchLength], k)) >= expectedKmers2 {
			fmt.Println("found a match")
			if !(matchedNode.isANeighbor(node)) {
				fmt.Println("building an edge", node.getID(), matchedNode.getID())
				graph.addEdge2(node.getID(), matchedNode.getID(), node.read[n-minMatchLength:])
			}
		}
	}
	graph.removeNode(node.getID())
}
*/

func CreateReadNetwork2Index(reads []string, minMatchLength, k, indexLength int, errorRate float64) Graph2 {

	fmt.Println("Building prefixIndex")
	prefixIndex := BuildPrefixIndex(reads, minMatchLength)
	fmt.Println("prefixIndex built!")
	fmt.Println("Building suffixIndex")
	suffixIndex := BuildSuffixIndex(reads, minMatchLength)
	fmt.Println("suffixIndex built!")

	network := MakeGraph2()
	pointerToNetwork := &network

	for i, read := range reads {
		fmt.Println("Adding read number", i+1)
		pointerToNetwork.addNode2(i, read)
	}

	for _, node := range network.Nodes {
		suffix := node.read[:minMatchLength]
		suffList, suffExists := suffixIndex[suffix]
		if suffExists {
			fmt.Println("found a bunch of matches!")
			for _, matchedID := range suffList {
				pointerToNetwork.addEdge2(node.getID2(), matchedID, suffix)
			}
		}
		prefix := node.read[len(node.read)-minMatchLength:]
		prefList, preExists := prefixIndex[prefix]
		if preExists {
			fmt.Println("found a bunch of matches!")
			for _, matchedID := range prefList {
				pointerToNetwork.addEdge2(node.getID2(), matchedID, suffix)
			}
		}
	}
	return network
}

func CreateReadNetwork3Index(reads []string, minMatchLength, k, indexLength int, errorRate float64) Graph2 {

	nodeTimes := make([]float64, 0)

	network := MakeGraph2()
	pointerToNetwork := &network

	for i, read := range reads {
		fmt.Println("Adding read number", i+1)
		pointerToNetwork.addNode2(i, read)
	}

	pointertoNodes := &(network.Nodes)

	fmt.Println("Building prefixIndex")
	prefixIndex := BuildPrefixIndex3(pointertoNodes, minMatchLength, k, errorRate)
	fmt.Println("prefixIndex built!")
	fmt.Println("Building suffixIndex")
	suffixIndex := BuildSuffixIndex3(pointertoNodes, minMatchLength, k, errorRate)
	fmt.Println("suffixIndex built!")

	fmt.Println("length of prefixIndex:", len(prefixIndex))
	fmt.Println("length of suffixIndex:", len(suffixIndex))

	//expectedShared := float64(ExpectedSharedkmers(minMatchLength, errorRate, k))

	for i := 0; i < len(network.Nodes); i++ {
		node := network.Nodes[i]
		start := time.Now()
		fmt.Println("entering the for loop of the", node.getID2(), "th node")
		//fmt.Println("prefKey:", node.prefkey)
		//list := prefixIndex[node.getprefkey()]
		//fmt.Println(len(list))
		n := len(node.read)
		prefix := node.read[:minMatchLength]
		suffix := node.read[n-minMatchLength:]
		for key := range prefixIndex {
			if float64(CountSharedKmers(key, suffix, k)) >= 0.7*float64(ExpectedSharedkmers(len(key), errorRate, k)) {
				for _, id2 := range prefixIndex[key] {
					//build an edge from node to the id2
					//pointerToNode := network.Nodes[id2]
					network.addEdge2(id2, node.getID2(), node.getsuffkey())
					//fmt.Println("outnode before:", len(node.outnodes))
					//(&node).addOutNode(network.Nodes[id2])
					//fmt.Println("outnode after:", len(node.outnodes))
					//fmt.Println("innode before:", len(pointerToNode.innodes))
					//(&pointerToNode).addInNode(node)
					//fmt.Println("innode after:", len(pointerToNode.innodes))
					fmt.Println("Edge built!!!!")
				}
			}
			break
		}

		for key := range suffixIndex {
			if float64(CountSharedKmers(key, prefix, k)) >= 0.7*float64(ExpectedSharedkmers(len(key), errorRate, k)) {
				for _, id2 := range suffixIndex[node.prefkey] {
					//pointerToNode := network.Nodes[id2]
					network.addEdge2(id2, node.getID2(), node.getsuffkey())
					//fmt.Println("innode before:", len(node.innodes))
					//(&node).addInNode(network.Nodes[id2])
					//fmt.Println("innode after:", len(node.innodes))
					//fmt.Println("outnode before:", len(pointerToNode.outnodes))
					//(&pointerToNode).addOutNode(node)
					//fmt.Println("outnode after:", len(pointerToNode.outnodes))
					fmt.Println("Edge built!!!!")
				}
			}
			break
		}

		fmt.Println("outnodes:", len(node.outnodes))
		fmt.Println("innodes:", len(node.innodes))

		elapsed := time.Since(start).Minutes()
		nodeTimes = append(nodeTimes, elapsed)
		network.Nodes[i] = node
	}
	totTime := 0.0
	for _, times := range nodeTimes {
		totTime += times
	}
	fmt.Println("average time per node:", totTime/(float64(len(nodeTimes))))
	//network.PrintGraph2()
	return network
}
