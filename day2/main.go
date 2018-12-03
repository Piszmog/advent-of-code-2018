package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"os"
	"path/filepath"
)

const filename = "day2/ids.csv"

func main() {
	pathToFile, err := filepath.Abs(filename)
	if err != nil {
		panic(errors.Wrapf(err, "failed to get absolute path of %s", filename))
	}
	file, err := os.Open(pathToFile)
	if err != nil {
		panic(errors.Wrapf(err, "failed to open file %s", filename))
	}
	defer closeFile(file)
	reader := csv.NewReader(bufio.NewReader(file))
	line := 0
	numberOfTwoOccurrences := 0
	numberOfThreeOccurrences := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(errors.Wrapf(err, "failed to read line $d", line))
		}
		id := record[0]
		letterOccurrences := make(map[int32]int)
		for _, letter := range id {
			if _, ok := letterOccurrences[letter]; ok {
				letterOccurrences[letter] += 1
			} else {
				letterOccurrences[letter] = 1
			}
		}
		alreadyAddedTwo := false
		alreadyAddedThree := false
		for _, value := range letterOccurrences {
			if !alreadyAddedTwo && value == 2 {
				numberOfTwoOccurrences++
				alreadyAddedTwo = true
			} else if !alreadyAddedThree && value == 3 {
				numberOfThreeOccurrences++
				alreadyAddedThree = true
			} else if alreadyAddedTwo && alreadyAddedThree {
				break
			}
		}
	}
	fmt.Printf("The checksum value of %d x %d = %d", numberOfTwoOccurrences, numberOfThreeOccurrences, numberOfTwoOccurrences*numberOfThreeOccurrences)
}

func closeFile(file *os.File) {
	err := file.Close()
	if err != nil {
		panic(errors.Wrapf(err, "failed to close %s", filename))
	}
}
