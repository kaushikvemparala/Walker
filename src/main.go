package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	fmt.Println("Assembling genomes (or trying at least)!")

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
	reads := CollectReadsFromFASTA(filename)
	fmt.Println("We have", len(reads), "total reads.")
	PrintStatistics(reads)

	minReadLength := 1000
	fmt.Println("Let's throw out short reads of length <", minReadLength)
	reads = DiscardShortReads(reads, minReadLength)
	fmt.Println("Updated read stats.")
	PrintStatistics(reads)

	fmt.Println("Calling assembler.")
	minMatchLength := 800
	indexLength := 15
	k := 7
	errorRate := 0.11
	contigs := GenomeAssembler4(reads, minMatchLength, indexLength, errorRate, k)
	PrintStatistics(contigs)
	fmt.Println("Finally, we write contigs to file.")
	outFilename := "assembly_contigs.fasta"
	WriteContigsToFile(contigs, outFilename)
}
