package main

import (
	"math/rand"
	//"fmt"
)

func ExpectedSharedkmers(stringLength int, errorRate float64, k int) int {
	// generate a random string
	str1 := GenerateRandomGenome(stringLength)

	//form second string by randomly mutating first string
	str2 := MutateDNAString(str1, errorRate)

	return CountSharedKmers(str1, str2, k)
}

func CountSharedKmers(str1, str2 string, k int) int {
	count := 0

	freqMap1 := FrequencyMap(str1, k)
	freqMap2 := FrequencyMap(str2, k)

	for pattern := range freqMap1 {
		// just take the minimum
		//countbefore := count
		count += Min2(freqMap1[pattern], freqMap2[pattern])
		//if count > countbefore {
		//			fmt.Println(pattern)
		//}
	}
	return count
}

func CountSharedKmersMod(freqMap1 map[string]int, str2 string, k int) int {
	count := 0

	freqMap2 := FrequencyMap(str2, k)

	for pattern := range freqMap1 {
		// just take the minimum
		count += Min2(freqMap1[pattern], freqMap2[pattern])
	}
	return count
}

func Min2(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func MutateDNAString(str string, errorRate float64) string {
	// string concatenation is slow, but generating arrays of bytes is fast
	symbols := make([]byte, len(str))

	// range over string, flip a coin, and mutate accordingly
	for i := range str {
		symbols[i] = MutateDNASymbol(str[i], errorRate)
	}

	return string(symbols)
}

//MutateDNASymbol mutates a given DNA symbol with probability equal to error rate given.
func MutateDNASymbol(symbol byte, errorRate float64) byte {
	x := rand.Float64()

	if x <= errorRate {
		//mutate!
		newSymbol := RandomDNASymbol()
		// if new == symbol, we don't want to return it
		for newSymbol == symbol {
			// generate another one
			newSymbol = RandomDNASymbol()
		}
		// we know we have a different symbol
		return newSymbol
	} else {
		return symbol
	}
}
func FrequentWords(text string, k int) []string {
	freqPatterns := make([]string, 0)

	freqMap := FrequencyMap(text, k)

	// find the max value of frequency map
	m := MaxMap(freqMap)

	// what achieves the max?
	for pattern, val := range freqMap {
		if val == m {
			// frequent pattern found! append it to our list
			freqPatterns = append(freqPatterns, pattern)
		}
	}

	return freqPatterns
}

func MaxMap(freqMap map[string]int) int {
	m := 0

	// range through map, and if something has higher value, update m!
	for pattern := range freqMap {
		if freqMap[pattern] > m {
			m = freqMap[pattern]
		}
	}
	// if all values in map were negative integers, this would return 0.
	// challenge: fix this bug so that it finds max value of any map of strings to ints.

	return m
}

//FrequencyMap takes a string text and an integer k. It produces a map
//of all k-mers in the string to their number of occurrences.
func FrequencyMap(text string, k int) map[string]int {
	// map declaration is analogous to slices
	// (we don't need to give an initial length)
	freq := make(map[string]int)
	n := len(text)
	for i := 0; i < n-k+1; i++ {
		pattern := text[i : i+k]
		// if freqMap[pattern] doesn't exist, create it.  How do we do this?
		/*
		   // approach #1
		   _, exists := freq[pattern]
		   if exists == false {
		     // create this value
		     freqMap[pattern] = 1
		   } else {
		     // we already have a value in the map
		     freqMap[pattern]++
		   }
		*/
		// approach #2
		// this says, if freqMap[pattern] exists, add one to it
		// if freqMap[pattern] doesn't exist, create it with a default value of 0, and add 1.
		freq[pattern]++
	}
	return freq
}

func PatternCount(pattern, text string) int {
	hits := StartingIndices(pattern, text)
	return len(hits)
}

func StartingIndices(pattern, text string) []int {
	hits := make([]int, 0)

	// append every starting position of pattern that we find in text

	n := len(text)
	k := len(pattern)

	for i := 0; i < n-k+1; i++ {
		if text[i:i+k] == pattern {
			// hit found!
			hits = append(hits, i)
		}
	}

	return hits
}

func SkewArray(genome string) []int {
	n := len(genome)
	array := make([]int, n+1)

	for i := range genome {
		/*
		   array[i+1] = array[i] + something
		   something = -1, 0, 1 depending genome[i]
		*/
		if genome[i] == 'A' || genome[i] == 'T' {
			array[i+1] = array[i]
		} else if genome[i] == 'C' {
			array[i+1] = array[i] - 1
		} else if genome[i] == 'G' {
			array[i+1] = array[i] + 1
		}
	}

	return array
}

func ReverseComplement(text string) string {
	return Reverse(Complement(text))
}

//Reverse takes a string and returns the reversed string.
func Reverse(text string) string {
	n := len(text)
	symbols := make([]byte, n)
	for i := range text {
		symbols[i] = text[n-i-1]
	}
	return string(symbols)
}

func Complement(text string) string {
	// as with arrays, we can use "range"

	n := len(text)
	symbols := make([]byte, n)

	for i := range text {
		if text[i] == 'A' {
			symbols[i] = 'T'
		} else if text[i] == 'T' {
			symbols[i] = 'A'
		} else if text[i] == 'C' {
			symbols[i] = 'G'
		} else if text[i] == 'G' {
			symbols[i] = 'C'
		}
	}

	return string(symbols)
}
