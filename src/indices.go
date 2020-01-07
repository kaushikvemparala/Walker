package main

import "fmt"

//BuildPrefixIndex takes a collection of strings (of arbitrary length bigger than prefix length)
//and a prefix length.
//It returns a map of the prefixes of strings of length prefixLength to
//their occurrences in reads.
func BuildPrefixIndex(reads []string, prefixLength int) map[string]([]int) {
	index := make(map[string]([]int))

	//populate our index
	for i, read := range reads {
		if len(read) < prefixLength {
			panic("Error: reads too short to build prefix index.")
		}
		prefix := read[:prefixLength]
		// have we seen this prefix before?
		_, ok := index[prefix]
		if ok == true { // we have seen this prefix before :)
			index[prefix] = append(index[prefix], i)
		} else { // we haven't seen this prefix before, create new list of occurrences
			index[prefix] = make([]int, 1)
			index[prefix][0] = i
		}
		if i%100000 == 0 {
			fmt.Println("Update: We have indexed", i, "prefixes.")
		}
	}
	return index
}

func BuildSuffixIndex(reads []string, suffixLength int) map[string]([]int) {
	index := make(map[string]([]int))

	//populate our index
	for i, read := range reads {
		if len(read) < suffixLength {
			panic("Error: reads too short to build suffix index.")
		}
		n := len(read)
		suffix := read[n-suffixLength:] // we want suffix of length suffixLength
		// have we seen this suffix before?
		_, ok := index[suffix]
		if ok == true { // we have seen this before :)
			index[suffix] = append(index[suffix], i)
		} else { // we haven't seen this before, create new list of occurrences
			index[suffix] = make([]int, 1)
			index[suffix][0] = i
		}
		if i%100000 == 0 {
			fmt.Println("Update: We have indexed", i, "suffixes.")
		}
	}
	return index
}
