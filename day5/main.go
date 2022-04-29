package main

import (
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type SumAndCount struct {
	sum   int
	count int
}

func main() {

	rawCoordinates, err := getStructFromInput("day5/input.txt")
	if err != nil {
		log.Fatal(err)
	}

	// Part 1:
	result, err := processPart1(rawCoordinates)
	if err != nil {
		log.Printf("Part 1 - %s", err)
	} else {
		log.Printf("Part 1 - %d", result)
	}

	// Part 2: TODO
	result, err = processPart2(rawCoordinates)
	if err != nil {
		log.Printf("Part 2 - %s", err)
	} else {
		log.Printf("Part 2 - %d", result)
	}
}

// Read the input file. It contains numbers representing 2 set of x,y coordinates
// These lines are written as `x1,y1 -> x2,y2`
// Returns an array in which each row contains the set of 4 coordinates
func getStructFromInput(path string) ([][]int, error) {
	// Read the input file.
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(file), "\n") // lines in the file

	rawCoords := make([][]int, 0)

	for _, line := range lines {
		if len(line) != 0 {
			re := regexp.MustCompile("\\d+")
			coordStr := re.FindAllString(line, 4)
			coord := make([]int, 4)
			for idx, valStr := range coordStr {
				// Get the integer value from the line
				coord[idx], err = strconv.Atoi(valStr)
				if err != nil {
					return nil, err
				}
			}

			// At this point, coord[] contains x1,y1,x2,y2
			rawCoords = append(rawCoords, coord)
		}
	}

	return rawCoords, nil
}

// Part 1: Based on the coordinates, trace each vertical and horizontal lines
// and count the number of points where more than one line pass over it
func processPart1(rawCoordinates [][]int) (int, error) {
	countPtsOver1 := 0
	// The following map will contain points where
	// * the key is the coordinates of the point
	// * the value is the number of lines that pass by that point
	ptMap := make(map[string]int)

	for _, coord := range rawCoordinates {
		// We care only about vertical and horizontal lines
		if coord[0] == coord[2] {
			// Vertical lines
			start := 0
			end := 0
			if coord[1] <= coord[3] {
				start = coord[1]
				end = coord[3]
			} else {
				start = coord[3]
				end = coord[1]
			}
			for i := start; i <= end; i++ {
				ptMap[string(coord[0])+"-"+string(i)] += 1
				if ptMap[string(coord[0])+"-"+string(i)] == 2 {
					countPtsOver1++
				}
			}
		} else if coord[1] == coord[3] {
			// Horizontal lines
			start := 0
			end := 0
			if coord[0] <= coord[2] {
				start = coord[0]
				end = coord[2]
			} else {
				start = coord[2]
				end = coord[0]
			}
			for i := start; i <= end; i++ {
				ptMap[string(i)+"-"+string(coord[1])] += 1
				if ptMap[string(i)+"-"+string(coord[1])] == 2 {
					countPtsOver1++
				}
			}
		}
	}

	return countPtsOver1, nil
}

// Part 2: Based on the coordinates, trace each vertical, horizontal and diagonal lines
// and count the number of points where more than one line pass over it
func processPart2(rawCoordinates [][]int) (int, error) {
	countPtsOver1 := 0
	// The following map will contain points where
	// * the key is the coordinates of the point
	// * the value is the number of lines that pass by that point
	ptMap := make(map[string]int)

	for _, coord := range rawCoordinates {
		// Let's do vertical, horizontal and finally diagonal lines
		if coord[0] == coord[2] {
			// Vertical lines
			start := 0
			end := 0
			if coord[1] <= coord[3] {
				start = coord[1]
				end = coord[3]
			} else {
				start = coord[3]
				end = coord[1]
			}
			for i := start; i <= end; i++ {
				ptMap[string(coord[0])+"-"+string(i)] += 1
				if ptMap[string(coord[0])+"-"+string(i)] == 2 {
					countPtsOver1++
				}
			}
		} else if coord[1] == coord[3] {
			// Horizontal lines
			start := 0
			end := 0
			if coord[0] <= coord[2] {
				start = coord[0]
				end = coord[2]
			} else {
				start = coord[2]
				end = coord[0]
			}
			for i := start; i <= end; i++ {
				ptMap[string(i)+"-"+string(coord[1])] += 1
				if ptMap[string(i)+"-"+string(coord[1])] == 2 {
					countPtsOver1++
				}
			}
		} else {
			// Time for diagonal lines
			xFactor, yFactor := 0, 0
			if coord[0] <= coord[2] {
				xFactor = 1
			} else {
				xFactor = -1
			}
			if coord[1] <= coord[3] {
				yFactor = 1
			} else {
				yFactor = -1
			}

			i := coord[0]
			j := coord[1]
			for cpt := 0; cpt <= int(math.Abs(float64(coord[2]-coord[0]))); cpt++ {
				ptMap[string(i)+"-"+string(j)] += 1
				if ptMap[string(i)+"-"+string(j)] == 2 {
					countPtsOver1++
				}
				i += xFactor
				j += yFactor
			}
		}
	}

	return countPtsOver1, nil
}
