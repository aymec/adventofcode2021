package main

import (
	"log"
	"os"
	"strconv"
	"strings"
)


func main() {
	// Read the input file. It contains a list of integers representing depth measures
	// in the order they are made
	data, err := os.ReadFile("day1/input.txt")
	if err != nil {
		log.Fatal(err)
		return
	}
	input:=strings.Split(string(data), "\n")

	// Read the first value, we need to start the comparison somewhere
	previous, err := strconv.Atoi(input[0])
	if err != nil {
		log.Fatal(err)
		return
	}
	count := 0

	// Part 1: Count the number of times a depth measurement increases
	for _, line := range input {
		// The first comparison is useless, at least I can use `range`
		value, err := strconv.Atoi(line)
		if err != nil {
			log.Fatal(err)
		} else {
			if value > previous {
				count++
			}
			previous = value
		}
	}

	log.Printf("Part 1 - Total number of times a depth measurement increases: %d\n", count)

	// Part 2: Make triplets of measures, as a sliding window, and count the number of times
	// the sun of measurements increases over the previous one
	// Reset necessary variables
	previous, err = strconv.Atoi(input[0])
	if err != nil {
		log.Fatal(err)
		return
	}
	count = 0
	// Value of the current window
	window := 0
	// Keep the value at n-3 so it can be removed and the new one add
	// to reduce the number of additions
	toRemove := previous

	for index , line := range input {
		value, err := strconv.Atoi(line)
		if err != nil {
			log.Fatal(err)
			return
		}
		previous = window
		window += value

		if index >= 3 {
			// Remove the value at index-3 from the window
			window -= toRemove
			// Keep the value to remove next
			toRemove, err = strconv.Atoi(input[index-2])
			if err != nil {
				log.Fatal(err)
				return
			}
			if window > previous {
				count++
			}
		}
	}
	log.Printf("Part 2 - Total number of times a depth window measurement increases: %d\n", count)


}

func sumOfSlice(slice []int) int {
	count := 0
	for _, value := range slice {
		count += value
	}
	return count
}

func stringArrayToIntArray(arr []string) ([]int, error) {
	ret := make([]int, len(arr))
	return ret, nil
}

