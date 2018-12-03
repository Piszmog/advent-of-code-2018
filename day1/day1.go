package main

import (
	"fmt"
	"github.com/Piszmog/advent-of-code-2018/utils"
	"github.com/pkg/errors"
	"strconv"
)

const filename = "day1/frequencies.csv"

func main() {
	frequencyMap := make(map[int]bool)
	numberOfFileReads := 0
	startingFrequency := 0
	frequencyMap[startingFrequency] = false
	// read file once to get the result frequency for part 1
	frequency, duplicateFrequency, duplicateFrequencyFound := readFrequencyFile(startingFrequency, frequencyMap)
	resultFrequency := frequency
	// read file until we found the duplicate frequency
	for !duplicateFrequencyFound {
		frequency, duplicateFrequency, duplicateFrequencyFound = readFrequencyFile(frequency, frequencyMap)
		numberOfFileReads++
		// to prevent reading the file forever, stop after 500 times. May need to increase
		if numberOfFileReads >= 500 {
			panic("Read the file too many times to find duplicate frequency")
		}
	}
	// print results
	fmt.Printf("First duplicate is %d after reading file %d times\n", duplicateFrequency, numberOfFileReads)
	fmt.Printf("Final frequency is %d", resultFrequency)
}

func readFrequencyFile(frequency int, frequencyResultMap map[int]bool) (int, int, bool) {
	file := utils.OpenFile(filename)
	defer utils.CloseFile(file)
	duplicateFrequency := 0
	duplicateFrequencyFound := false
	utils.ReadFile(file, func(record []string, line int) {
		stringValue := record[0]
		newFrequency, err := strconv.Atoi(stringValue)
		if err != nil {
			panic(errors.Wrapf(err, "failed to convert %s to int on line $d", stringValue, line))
		}
		frequency += newFrequency
		if !duplicateFrequencyFound {
			if _, ok := frequencyResultMap[frequency]; ok {
				frequencyResultMap[frequency] = true
				duplicateFrequency = frequency
				duplicateFrequencyFound = true
			} else {
				frequencyResultMap[frequency] = false
			}
		}
	})
	return frequency, duplicateFrequency, duplicateFrequencyFound
}
