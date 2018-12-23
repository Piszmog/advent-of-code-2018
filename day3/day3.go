package main

import (
	"fmt"
	"github.com/Piszmog/advent-of-code-2018/utils"
	"github.com/pkg/errors"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const filename = "day3/fabricClaims.csv"

type Coordinate struct {
	x int
	y int
}

type FabricClaim struct {
	id            string
	topLeftCorner Coordinate
	height        int
	width         int
	area          int
}

func (fabricClaim FabricClaim) getOverlapClaim(secondaryFabricClaim FabricClaim) *FabricClaim {
	area := 0
	dxMin := int(math.Min(float64(fabricClaim.topLeftCorner.x+fabricClaim.width), float64(secondaryFabricClaim.topLeftCorner.x+secondaryFabricClaim.width)))
	dxMax := int(math.Max(float64(fabricClaim.topLeftCorner.x), float64(secondaryFabricClaim.topLeftCorner.x)))
	overlapWidth := dxMin - dxMax
	dyMin := int(math.Min(float64(fabricClaim.topLeftCorner.y+fabricClaim.height), float64(secondaryFabricClaim.topLeftCorner.y+secondaryFabricClaim.height)))
	dyMax := int(math.Max(float64(fabricClaim.topLeftCorner.y), float64(secondaryFabricClaim.topLeftCorner.y)))
	overlapHeight := dyMin - dyMax
	if overlapWidth > 0 && overlapHeight > 0 {
		area = overlapWidth * overlapHeight
		ids := []string{fabricClaim.id, secondaryFabricClaim.id}
		sort.Strings(ids)
		return &FabricClaim{
			id: ids[0] + ids[1],
			topLeftCorner: Coordinate{
				x: dxMin,
				y: dyMin,
			},
			width:  overlapWidth,
			height: overlapHeight,
			area:   area,
		}
	} else {
		return nil
	}
}

func main() {
	defer utils.RunTime(time.Now())
	file := utils.OpenFile(filename)
	defer utils.CloseFile(file)
	fabricClaims := make([]FabricClaim, 1397)
	lines := make(chan []string, 200)
	done := make(chan bool)
	go mapFabricClaim(lines, fabricClaims, done)
	readFile(file, lines)
	<-done
	totalOverlapArea := 0
	// todo complete
	//overlapMap := make(map[string][]string)
	//for primaryIndex, primaryFabricClaim := range fabricClaims {
	//	var finalOverlapClaim *FabricClaim
	//	for secondaryIndex, secondaryFabricClaim := range fabricClaims {
	//		if finalOverlapClaim == nil {
	//			if primaryIndex == secondaryIndex {
	//				continue
	//			}
	//			overlapFabricClaim := primaryFabricClaim.getOverlapClaim(secondaryFabricClaim)
	//			if overlapFabricClaim == nil {
	//				continue
	//			} else if containsId(secondaryFabricClaim.id, overlapMap[primaryFabricClaim.id]) && containsId(primaryFabricClaim.id, overlapMap[secondaryFabricClaim.id]) {
	//				continue
	//			}
	//			overlapMap[primaryFabricClaim.id] = append(overlapMap[primaryFabricClaim.id], secondaryFabricClaim.id)
	//			overlapMap[secondaryFabricClaim.id] = append(overlapMap[secondaryFabricClaim.id], primaryFabricClaim.id)
	//			finalOverlapClaim = overlapFabricClaim
	//		} else {
	//			overlapClaim := finalOverlapClaim.getOverlapClaim(secondaryFabricClaim)
	//			if overlapClaim != nil {
	//				finalOverlapClaim = overlapClaim
	//			}
	//		}
	//	}
	//}
	fmt.Printf("Total overlap area is %d\n", totalOverlapArea)
}

func readFile(file *os.File, lines chan []string) {
	defer close(lines)
	utils.ReadFileWithDelimiter(file, ' ', func(record []string, line int) {
		lines <- record
	})
}

func mapFabricClaim(lines chan []string, fabricClaims []FabricClaim, done chan bool) {
	i := 0
	for record := range lines {
		edgeDistances := record[2]
		distances := strings.Split(edgeDistances, ",")
		xLeftCorner, err := strconv.Atoi(distances[0])
		if err != nil {
			panic(errors.Wrapf(err, "failed to convert %s to int"))
		}
		yLeftCorner, err := strconv.Atoi(strings.Trim(distances[1], ":"))
		if err != nil {
			panic(errors.Wrapf(err, "failed to convert %s to int"))
		}
		fabricDimensions := record[3]
		dimensions := strings.Split(fabricDimensions, "x")
		width, err := strconv.Atoi(dimensions[0])
		if err != nil {
			panic(errors.Wrapf(err, "failed to convert %s to int"))
		}
		height, err := strconv.Atoi(dimensions[1])
		if err != nil {
			panic(errors.Wrapf(err, "failed to convert %s to int"))
		}
		fabricClaim := FabricClaim{
			id:            record[0],
			topLeftCorner: Coordinate{x: xLeftCorner, y: yLeftCorner},
			height:        height,
			width:         width,
			area:          width * height,
		}
		fabricClaims[i] = fabricClaim
		i++
	}
	close(done)
}

func containsId(id string, listOfIds []string) bool {
	for _, existingId := range listOfIds {
		if id == existingId {
			return true
		}
	}
	return false
}
