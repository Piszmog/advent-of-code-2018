package main

import (
	"fmt"
	"github.com/Piszmog/advent-of-code-2018/utils"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

const filename = "day1/frequencies.csv"

func main() {
	defer utils.RunTime(time.Now())
	frequencyMap := make(map[int]bool)
	startingFrequency := 0
	frequencyMap[startingFrequency] = false
	frequencyChannel := make(chan int, 100)
	done := make(chan bool)
	frequencyList := make([]int, 959)
	// read file once to get the result frequency for part 1
	go func() {
		i := 0
		for frequency := range frequencyChannel {
			frequencyList[i] = frequency
			i++
		}
		done <- true
		close(done)
	}()
	readFrequencyFile(frequencyChannel)
	<-done
	resultFrequency, duplicateFrequencyFound, duplicateFrequency := processFrequencies(startingFrequency, frequencyMap, frequencyList)
	reprocessAttempts := 0
	if !duplicateFrequencyFound {
		startingFrequency = resultFrequency
		maxReprocesses := 500
		for !duplicateFrequencyFound && reprocessAttempts < maxReprocesses {
			startingFrequency, duplicateFrequencyFound, duplicateFrequency = processFrequencies(startingFrequency, frequencyMap, frequencyList)
			reprocessAttempts++
		}
	}
	if !duplicateFrequencyFound {
		panic("failed to find duplicate frequency after 500 attempts")
	}
	// print results
	fmt.Printf("First duplicate is %d after reprocessing %d times\n", duplicateFrequency, reprocessAttempts)
	fmt.Printf("Final frequency is %d", resultFrequency)
}

func readFrequencyFile(frequencyChannel chan int) {
	file := utils.OpenFile(filename)
	defer utils.CloseFile(file)
	defer close(frequencyChannel)
	utils.ReadFile(file, func(record []string, line int) {
		stringValue := record[0]
		newFrequency, err := strconv.Atoi(stringValue)
		if err != nil {
			panic(errors.Wrapf(err, "failed to convert %s to int on line $d", stringValue, line))
		}
		frequencyChannel <- newFrequency
	})
}

func processFrequencies(startingFrequency int, frequencyMap map[int]bool, frequencyList []int) (int, bool, int) {
	duplicateFrequencyFound := false
	currentFrequency := startingFrequency
	resultFrequency := currentFrequency
	duplicateFrequency := currentFrequency
	for _, frequency := range frequencyList {
		resultFrequency += frequency
		if !duplicateFrequencyFound {
			if _, ok := frequencyMap[resultFrequency]; ok {
				frequencyMap[resultFrequency] = true
				duplicateFrequency = resultFrequency
				duplicateFrequencyFound = true
			} else {
				frequencyMap[resultFrequency] = false
			}
		}
		currentFrequency = resultFrequency
	}
	return resultFrequency, duplicateFrequencyFound, duplicateFrequency
}
