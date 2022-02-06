package main

import (
	"errors"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type SumAndCount struct {
	sum int
	count int
}

func main() {
	// Play bingo
	// A grid is a structure of 5 lines, each containing 5 numbers
	// A grid is the winning grid if one of the lines had all its numbers drawn

	// Part 1: we'll create 5 structures:
	// 1. A list (a queue) of the drawn numbers, in the order they were drawn
	// 2. An array of integers whose:
	//    * indexes represent the row number in the grids
	//    * values are the sum of all numbers on that row and the count of numbers drawn on that row
	//    Rows in grid 1 are at indexes 0 to 4, rows in grid 2 are at indexes 5 to 9, etc
	// 3. A map of <integer,integer> whose
	//    * keys are numbers presents in the grids' rows
	//    * value for a key is a list of the rows where that number is present (indexes in the previous structure)
	// 4. An array of integers whose:
	//    * indexes represent the grid and column number in the grids
	//    * values are the sum of all numbers on that column for that grid and the count of numbers drawn on that column
	//    Columns in grid 1 are at indexes 0 to 4, columns in grid 2 are at indexes 5 to 9, etc
	// 5. A map of <integer,integer> whose
	//    * keys are numbers presents in the grids' columns
	//    * value for a key is a list of the columns where that number is present (indexes in the previous structure)
	drawnNumbers, rowSums, rowReverseIndex, colSums, colReverseIndex, err := getStructFromInput("day4/input.txt")
	if err != nil {
		log.Fatal(err)
	}

	result, winningDrawIndex, err := processPart1(drawnNumbers, rowSums, rowReverseIndex, colSums, colReverseIndex)
	if err != nil {
		log.Printf("Part 1 - %s", err)
	} else {
		log.Printf("Part 1 - %d", result)
	}

	// Part 2: we play until our last grid wins. For that we need to keep the number of winning grids
	// We'll actually keep a count of grids that did not win
	remainingGrids := countRemainingNonWinningGrids(rowSums, colSums)
	result, err = processPart2(drawnNumbers, rowSums, rowReverseIndex, colSums, colReverseIndex, remainingGrids, winningDrawIndex)
	if err != nil {
		log.Printf("Part 2 - %s", err)
	} else {
		log.Printf("Part 2 - %d", result)
	}
}

// See in main the 5 structures that this function will return
// It returns the 5 structures mentioned in main in the order mentioned:
// * the list of drawn numbers
// * the sums for each row in every grid
// * the reverse index for each number in the grids and which row they can be found in
// * the sums for each column in every grid
// * the reverse index for each number in the grids and which column they can be found in
func getStructFromInput(path string) (
	[]int,
	[]SumAndCount,
	map[int][]int,
	[]SumAndCount,
	map[int][]int,
	error) {
	// Read the input file. It contains
	// * A first line with a series of number in order in which they were drawn. separated by `,`
	// * A series of 5 consecutive lines with each 5 numbers seperated by spaces
	// * the series of 5 lines are separated by an empty line

	file, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, nil, nil, nil, err
	}
	lines:=strings.Split(string(file), "\n") // lines in the file

	drawnNumbers := make([]int,0)
	rowReverseIndex := make(map[int][]int)
	rowSums := make([]SumAndCount, 0)
	colReverseIndex := make(map[int][]int)
	colSums := make([]SumAndCount, 0)

	// Process the first line that contains the drawn numbers
	firstLine := lines[0] // line in the file
	lines = lines[2:] // discard that first line and the next empty line
	for _ , number := range strings.Split(firstLine, ",") {
		value, err := strconv.Atoi(number)
		if err != nil {
			return nil, nil, nil, nil, nil, err
		}
		drawnNumbers = append(drawnNumbers, value)
	}

	// Process the other lines that contains the bingo grid lines
	lineIndex := 0 // Can't use the range index as it would include empty lines in between grids
	for _ , line := range lines {
		// each line contains numbers split by a whitespace
		// except empty lines in between bingo grids
		if len(line) != 0 {
			re := regexp.MustCompile("\\s+")
			line = strings.TrimSpace(line) // remove leading and trailing white space, because re.Split does not
			numbers:=re.Split(line, -1)
			rowSums = append(rowSums,SumAndCount{0, 0})
			for colIndex , number := range numbers {
				if lineIndex%5 == 0 { // for the 1st time we encounter a new column in this grid
					colSums = append(colSums, SumAndCount{0, 0})
				}
				value, err := strconv.Atoi(number)
				if err != nil {
					return nil, nil, nil, nil, nil, err
				}

				// Process row
				rowSumAndCount := rowSums[lineIndex]
				rowSumAndCount.sum += value
				rowSumAndCount.count += 1
				rowSums[lineIndex] = rowSumAndCount
				if rowReverseIndex[value] == nil {
					rowReverseIndex[value] = make([]int, 0)
				}
				rowReverseIndex[value] = append(rowReverseIndex[value], lineIndex)

				// Process column
				colSumAndCount := colSums[((lineIndex/5)*5)+colIndex]
				colSumAndCount.sum += value
				colSumAndCount.count++
				colSums[((lineIndex/5)*5)+colIndex] = colSumAndCount
				if colReverseIndex[value] == nil {
					colReverseIndex[value] = make([]int, 0)
				}
				colReverseIndex[value] = append(colReverseIndex[value], ((lineIndex/5)*5)+colIndex)
			}
			lineIndex++
		}
	}

	return drawnNumbers, rowSums, rowReverseIndex, colSums, colReverseIndex, nil
}

// Part 1: Find the winning grid. Then sum the rest of the numbers that were
// not drawn in that grid. Then multiply that sum by the winning number.
// To do that, we pop a number from the queue of drawn numbers (structure 1)
// we find that number in the map (structure 3) and for each line represented
// by the values at that key, we decrease the value at the corresponding index
// in the sum array (structure 2) by the value of that drawn number. If the sum
// reaches 0, we have a winning line. We then need to sum the remaining sums in
// the lines from the same grid. To find those, we take lines from the winning
// line index, modulo 5, to find the 1st line in the grid and sum the remaining
// value from that and the next 4 lines.
// It also returns the index of the winning number
func processPart1(
	drawnNumbers []int,
	rowSums []SumAndCount,
	rowReverseIndex map[int][]int,
	colSums []SumAndCount,
	colReverseIndex map[int][]int) (int, int, error) {
	// Processing the drawn number 1 by 1
	for index , draw := range drawnNumbers {
		// For each drawn number, we look in the rowReverseIndex map in which row we'll find them
		for _ , gridLine := range rowReverseIndex[draw] {
			sumAndCount := rowSums[gridLine]
			sumAndCount.sum -= draw
			sumAndCount.count--
			rowSums[gridLine] = sumAndCount
			// Check if we have a winner
			if sumAndCount.count == 0 {
				// We process the columns as well to keep it consistent, before we return
				for _ , gridCol := range colReverseIndex[draw] {
					sumAndCount := colSums[gridCol]
					sumAndCount.sum -= draw
					sumAndCount.count--
					colSums[gridCol] = sumAndCount
				}
				return weHaveAWinner(draw, gridLine, rowSums), index, nil
			}
		}

		// For each drawn number, we look in the colReverseIndex map in which column we'll find them
		for _ , gridCol := range colReverseIndex[draw] {
			sumAndCount := colSums[gridCol]
			sumAndCount.sum -= draw
			sumAndCount.count--
			colSums[gridCol] = sumAndCount
			// Check if we have a winner
			if sumAndCount.count == 0 {
				return weHaveAWinner(draw, gridCol, colSums), index, nil
			}
		}
	}

	// No winner --> return an error
	return 0, 0, errors.New("no winner")
}

// Return the multiplication of the winning number by the sum of the remaining values in the same grid
func weHaveAWinner(winningDrawNumber int, lineIndex int, lineSums []SumAndCount) int {
	sumRemainingInGrid := 0
	for i := (lineIndex/5)*5; i < ((lineIndex/5)*5) + 5; i++ {
		sumRemainingInGrid += lineSums[i].sum
	}
	return sumRemainingInGrid*winningDrawNumber
}

// Part 2: we play until the last winning grid. Then we need to return a similar output
// as in part 1: the multiplication of the winning drawn number by the sum of the remaining
// values in the winning grid
func processPart2(
	drawnNumbers []int,
	rowSums []SumAndCount,
	rowReverseIndex map[int][]int,
	colSums []SumAndCount,
	colReverseIndex map[int][]int,
	remainingGrids int,
	startIndexDrawnNumber int) (int, error) {
	// We keep playing
	for i := startIndexDrawnNumber+1; i < len(drawnNumbers); i++ {
		draw := drawnNumbers[i]
		// For each drawn number, we look in the rowReverseIndex map in which row we'll find them
		for _ , gridLine := range rowReverseIndex[draw] {
			sumAndCount := rowSums[gridLine]
			sumAndCount.sum -= draw
			sumAndCount.count--
			rowSums[gridLine] = sumAndCount
			// We have another winner. We should check whether that particular grid has already won or not
			// and whether this is the last grid to win before we decide to keep going
			if sumAndCount.count == 0 {
				if !gridHasAlreadyWon(gridLine, rowSums, colSums) {
					remainingGrids--
					// Now check whether that's the last winning grid
					if remainingGrids == 0 {
						return weHaveAWinner(draw, gridLine, rowSums), nil
					}
				}
			}
		}

		// For each drawn number, we look in the colReverseIndex map in which column we'll find them
		for _ , gridCol := range colReverseIndex[draw] {
			sumAndCount := colSums[gridCol]
			sumAndCount.sum -= draw
			sumAndCount.count--
			colSums[gridCol] = sumAndCount
			// We have another winner. We should check whether that particular grid has already won or not
			// and whether this is the last grid to win before we decide to keep going
			if sumAndCount.count == 0 {
				if !gridHasAlreadyWon(gridCol, rowSums, colSums) {
					remainingGrids--
					// Now check whether that's the last winning grid
					if remainingGrids == 0 {
						return weHaveAWinner(draw, gridCol, colSums), nil
					}
				}
			}
		}
	}

	// all numbers have been drawn and we have multiple remaining grids
	return 0, errors.New("multiple remaining grids")
}

// Returns whether the grid for the given lineIndex has already won
// return: true if the grid had already won previously
func gridHasAlreadyWon(index int, rowSums []SumAndCount, colSums []SumAndCount) bool {
	// To know if a grid has already won, we count the nb of rows and columns that have already won in that
	// grid. If there is only one, that grid has not won before
	nbOfWinningRowOrCol := 0
	for i := (index / 5) * 5; i < ((index/5)*5)+5; i++ {
		if rowSums[i].count == 0 {
			nbOfWinningRowOrCol++
		}
		if colSums[i].count == 0 {
			nbOfWinningRowOrCol++
		}
	}
	if nbOfWinningRowOrCol == 1 {
		return false
	} else {
		return true
	}
}

func countRemainingNonWinningGrids(rowSums []SumAndCount, colSums []SumAndCount) int {
	countRemainingGridsByRow := countRemainingNonWinningGridsByX(rowSums)
	countRemainingGridsByCol := countRemainingNonWinningGridsByX(colSums)
	if countRemainingGridsByRow < countRemainingGridsByCol {
		return countRemainingGridsByRow
	} else {
		return countRemainingGridsByCol
	}
}

func countRemainingNonWinningGridsByX(xSums []SumAndCount) int {
	countRemainingGridsByX := 0
	tmpCount := 0

	for index, sumAndCount := range xSums {
		if index%5 == 0 {
			tmpCount = 0
		}
		if sumAndCount.count > 0 {
			tmpCount++
		}
		if (index%5 == 4) && (tmpCount == 5) {
			countRemainingGridsByX++
		}
	}

	return countRemainingGridsByX
}