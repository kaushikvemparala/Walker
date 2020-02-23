package main

import (
	"fmt";
	"os"
)

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

func BuildPrefixIndex2(reads []string, prefixLength int) map[string]([]int) {
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

func BuildSuffixIndex2(reads []string, suffixLength int) map[string]([]int) {
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

func BuildPrefixIndex3(nodes *map[int]Node2, prefixLength, k int, errorRate float64, logFile *os.File) map[string]([]int) {

	index := make(map[string]([]int))

	//expectedShared := float64(ExpectedSharedkmers(prefixLength, errorRate, k))

	//populate our index
	for id := range *nodes {
		//fmt.Println("index", id, "th node")
		node := (*nodes)[id]
		if len((*nodes)[id].read) < prefixLength {
			panic("Error: reads too short to build prefix index.")
		}
		prefix := node.read[:prefixLength]

		if len(index) > 0 {
			for key := range index {
				//fmt.Println("length of key:", len(key))
				//fmt.Println("length of prefix:", len(prefix))
				countShared := float64(CountSharedKmers(key, prefix, 15))
				expectedShared := float64(ExpectedSharedkmers(prefixLength, errorRate, k))
				//fmt.Println("countShared:", countShared)
				//fmt.Println("expectedShared:", expectedShared)
				if countShared >= 0.5 * expectedShared {
					index[key] = append(index[key], id)
					//node.setprefkey(key)
					//(*nodes)[id] = node
					//fmt.Println(len(index[key]))
					//fmt.Println("found some matches")
					break
				} else { // we haven't seen this prefix before, create new list of occurrences
					index[prefix] = make([]int, 1)
					index[prefix][0] = id
					//fmt.Println("made a new entry")
				}
			}
		} else { // we haven't seen this prefix before, create new list of occurrences
			index[prefix] = make([]int, 1)
			index[prefix][0] = id
			//fmt.Println("made a new entry")
		}
		fmt.Fprintln(logFile, len(index))
	}
	fmt.Fprintln(logFile, "length of prefixIndex index:", len(index))
	return index
}

func BuildSuffixIndex3(nodes *map[int]Node2, suffixLength, k int, errorRate float64, logFile *os.File) map[string]([]int) {
	index := make(map[string]([]int))

	//expectedShared := float64(ExpectedSharedkmers(suffixLength, errorRate, k))

	//populate our index
	for id := range *nodes {
		//fmt.Println("index for", id, "th node")
		node := (*nodes)[id]
		if len((*nodes)[id].read) < suffixLength {
			panic("Error: reads too short to build prefix index.")
		}
		n := len(node.read)
		suffix := node.read[n-suffixLength:]

		if len(index) > 0 {
			for key := range index {
				if float64(CountSharedKmers(key, suffix, k)) >= 0.5 * float64(ExpectedSharedkmers(suffixLength, errorRate, k)) {
					index[key] = append(index[key], id)
					//node.setsuffkey(key)
					//(*nodes)[id] = node
					//fmt.Println(len(index[key]))
					break
				} else { // we haven't seen this prefix before, create new list of occurrences
					index[suffix] = make([]int, 1)
					index[suffix][0] = id
				}
			}
		} else { // we haven't seen this prefix before, create new list of occurrences
			index[suffix] = make([]int, 1)
			index[suffix][0] = id
		}
		fmt.Fprintln(logFile, len(index))
	}
	fmt.Fprintln(logFile, "length of suffix index:", len(index))
	return index
}
