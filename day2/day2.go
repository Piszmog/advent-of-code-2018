package main

import (
	"fmt"
	"github.com/Piszmog/advent-of-code-2018/utils"
	"github.com/pkg/errors"
	"os"
	"strings"
)

const filename = "day2/ids.csv"

func main() {
	file := utils.OpenFile(filename)
	defer utils.CloseFile(file)
	ids, numberOfTwoOccurrences, numberOfThreeOccurrences := readFile(file)
	resultId := findId(ids)
	if len(resultId) == 0 {
		panic(errors.New("Failed to find ids with 1 letter difference"))
	}
	fmt.Printf("The checksum value of %d x %d = %d\n", numberOfTwoOccurrences, numberOfThreeOccurrences, numberOfTwoOccurrences*numberOfThreeOccurrences)
	fmt.Printf("The id with the fabic is %s", resultId)
}

func readFile(file *os.File) ([]string, int, int) {
	ids := make([]string, 250)
	numberOfTwoOccurrences := 0
	numberOfThreeOccurrences := 0
	utils.ReadFile(file, func(record []string, line int) {
		id := record[0]
		ids[line] = id
		letterOccurrences := getLetterOccurrences(id)
		numberOfTwoOccurrences, numberOfThreeOccurrences = checkOccurrences(letterOccurrences, numberOfTwoOccurrences, numberOfThreeOccurrences)
	})
	return ids, numberOfTwoOccurrences, numberOfThreeOccurrences
}

func getLetterOccurrences(id string) map[int32]int {
	letterOccurrences := make(map[int32]int)
	for _, letter := range id {
		if _, ok := letterOccurrences[letter]; ok {
			letterOccurrences[letter] += 1
		} else {
			letterOccurrences[letter] = 1
		}
	}
	return letterOccurrences
}

func checkOccurrences(letterOccurrences map[int32]int, numberOfTwoOccurrences int, numberOfThreeOccurrences int) (int, int) {
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
	return numberOfTwoOccurrences, numberOfThreeOccurrences
}

func findId(ids []string) string {
	for primaryIndex, primaryId := range ids {
		for secondaryIndex, secondaryId := range ids {
			if primaryIndex == secondaryIndex {
				continue
			} else {
				numberOfDifferences, matchingCharacters := getNumberOfDifferences(primaryId, secondaryId)
				if numberOfDifferences == 1 {
					return strings.Join(matchingCharacters, "")
				}
			}
		}
	}
	return ""
}

func getNumberOfDifferences(primaryId string, secondaryId string) (int, []string) {
	numberOfDifferences := 0
	splitPrimary := strings.Split(primaryId, "")
	splitSecondary := strings.Split(secondaryId, "")
	primaryLength := len(splitPrimary)
	secondaryLength := len(splitSecondary)
	if primaryLength != secondaryLength {
		panic(errors.Errorf("ids %s and %s have different lengths", primaryId, secondaryId))
	}
	matchingCharacters := make([]string, primaryLength)
	for i := 0; i < primaryLength; i++ {
		difference := strings.Compare(splitPrimary[i], splitSecondary[i])
		if difference != 0 {
			numberOfDifferences++
		} else {
			matchingCharacters[i] = splitPrimary[i]
		}
		if numberOfDifferences > 1 {
			break
		}
	}
	return numberOfDifferences, matchingCharacters
}
