package main

import (
	"bufio"
	"fmt"
	"os"
)

func CollectReadsFromFASTA(filename string) []string {
	file, err := os.Open(filename)

	if err != nil {
		// error in opening file
		panic("Error: something went wrong with file open (probably you gave wrong filename).")
	}

	scanner := bufio.NewScanner(file) // think of this as a "reader bot"
	reads := make([]string, 0)

	// let's use same trick of using map and not read in duplicate reads.
	readCount := make(map[string]int)
	currentRead := ""
	counter := 0 // for updating user

	// go for as long as the reader bot can still see text
	for scanner.Scan() {
		if counter == 100000 {
			break
		}
		currentLine := scanner.Text() // grabs one line of text and returns a strings
		if currentLine[0] != '>' {
			// append the current line to our growing read
			currentRead += currentLine
		} else { // we are at a header
			// the current read is complete! :) append it
			if currentRead != "" && ValidDNAString(currentRead) {
				readCount[currentRead]++
				counter++
				currentRead = ""
				if counter%20000 == 0 {
					fmt.Println("Update: we have processed", counter, "reads.")
				}
			}
		}
	}

	// we have read everything in
	if scanner.Err() != nil {
		panic("Error: issue in scanning process.")
	}

	file.Close()

	// our reads are living as the keys of the readCount map. Grab them.
	for read := range readCount {
		reads = append(reads, read)
	}

	return reads
}

func ValidDNAString(dna string) bool {
	// check if dna is composed of A's, C's, G's, T's
	for _, symbol := range dna {
		if symbol != 'A' && symbol != 'C' && symbol != 'G' && symbol != 'T' {
			return false
		}
	}
	// we made it!
	return true
}

func WriteContigsToFile(contigs []string, outFilename string) {
	outFile, err := os.Create(outFilename)
	if err != nil {
		panic("Sorry, couldn't create file!")
	}
	for _, str := range contigs {
		fmt.Fprintln(outFile, str)
	}
	outFile.Close()
}
