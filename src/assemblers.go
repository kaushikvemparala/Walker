package main

import "fmt"

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
