package main

import "math/rand"

//GenerateRandomGenome takes a parameter length and returns
//a random DNA string of this length where each nucleotide has equal probability.
func GenerateRandomGenome(length int) string {
	// generate an array of random symbols
	symbols := make([]byte, length)
	for i := 0; i < length; i++ {
		symbols[i] = RandomDNASymbol()
	}
	// combine your array into a string
	return string(symbols)
}

//RandomDNASymbol takes no inputs and produces a symbol from alphabet
//{A, C, G, T} with equal probability.
func RandomDNASymbol() byte {
	number := rand.Intn(4)
	/*
	  if number == 0 {
	    return 'A'
	  } else if number == 1 {
	    return 'C'
	  } else if number == 2 {
	    return 'G'
	  } else if number == 3 {
	    return 'T'
	  }
	*/
	switch number {
	case 0:
		return 'A'
	case 1:
		return 'C'
	case 2:
		return 'G'
	case 3:
		return 'T'
	}
	panic("Error: something really weird is happening in RandomDNASymbol()")
}

//KmerComposition returns the k-mer composition (all k-mer substrings) of a given genome.
func KmerComposition(genome string, k int) []string {
	n := len(genome)
	kmers := make([]string, n-k+1)
	// range through and grab all substrings
	for i := 0; i < n-k+1; i++ {
		kmers[i] = genome[i : i+k]
	}

	return kmers
}

// now let's produce a function that simulates reads sampled from a genome
// let's assume we have generated a genome already

func SimulateReads(genome string, minReadLength, maxReadLength, coverage int) []string {
	n := len(genome)
	reads := make([]string, 0)

	// map will be useful to detect if a sampled read is already present
	patternMap := make(map[string]int)

	averageReadLength := (minReadLength + maxReadLength) / 2

	//what should numTrials be?
	numTrials := int(float64(coverage) * float64(n) / float64(averageReadLength))

	for i := 0; i < numTrials; i++ {
		// we need a random # between minReadLength and maxReadLength
		randNum := rand.Intn(maxReadLength - minReadLength + 1) // 0 to maxReadLength - minReadLength
		readLength := randNum + minReadLength                   // :)

		// grab a read from the genome at random starting position
		startingPos := rand.Intn(n - readLength + 1)
		read := genome[startingPos : startingPos+readLength]
		// reads = append(reads, read) // I don't want to repeatedly sample the same spot
		// idea 1: if read already appears in our reads, don't add it.
		// idea 2: every time, reduce size of genome, redeclare n.
		// idea 3: make a map whose keys are starting positions, and now it's quick
		// to figure out if something exists
		patternMap[read]++ // if key = read doesn't exist, set value to 1. otherwise, add 1 to it.
	}

	// I just need to set reads = KEYS of the map we created
	for read := range patternMap {
		reads = append(reads, read)
	}

	return reads
}

// challenge: write this equals function
/*
func Equals(kmers1, kmers2 []string) bool {
	for i, str := range kmers1 {
		// look for a match
		if Contains(str, kmers2) {

		}
	}
}
*/
