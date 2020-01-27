package main

import (
	"fmt"
	"os"
)

func AverageStringLength(patterns []string) float64 {
	numStrings := len(patterns)
	return float64(TotalStringLength(patterns)) / float64(numStrings)
}

func TotalStringLength(patterns []string) int {
	s := 0
	for _, pattern := range patterns {
		s += len(pattern)
	}
	return s
}

func MinimumStringLength(patterns []string) int {
	m := len(patterns[0])
	for i := 1; i < len(patterns); i++ {
		if len(patterns[i]) < m {
			m = len(patterns[i])
		}
	}
	return m
}

func MaximumStringLength(patterns []string) int {
	m := len(patterns[0])
	for i := 1; i < len(patterns); i++ {
		if len(patterns[i]) > m {
			m = len(patterns[i])
		}
	}
	return m
}

func PrintStatistics(patterns []string, summaryFile *os.File) {
	fmt.Fprintln(summaryFile, "\t\tShortest read:", MinimumStringLength(patterns))
	fmt.Fprintln(summaryFile, "")
	fmt.Fprintln(summaryFile, "\t\tLongest Read:", MaximumStringLength(patterns))
	fmt.Fprintln(summaryFile, "")
	fmt.Fprintln(summaryFile, "\t\tTotal length:", TotalStringLength(patterns))
	fmt.Fprintln(summaryFile, "")
	fmt.Fprintln(summaryFile, "\t\tAverage read length:", AverageStringLength(patterns))
	fmt.Fprintln(summaryFile, "")
}

func DiscardShortReads(patterns []string, minReadLength int) []string {
	//challenge: why do I go from end of reads backward instead of just ranging?
	for j := len(patterns) - 1; j >= 0; j-- {
		if len(patterns[j]) < minReadLength {
			patterns = Remove(patterns, j)
		}
	}
	return patterns
}
