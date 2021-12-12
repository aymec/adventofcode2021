package main

import (
	"errors"
	"log"
	"os"
	"strings"
)

type Rate struct {
	value []bool
}

func main() {
	// Read the input file. It contains a list of binary numbers
	structuredInput, err := getStructFromInput("day3/input.txt")
	if err != nil {
		log.Fatal(err)
	}

	// Part 1: multiply the gamma rate by the epsilon rate
	// sumsOfOnes will contain the count of '1' at each index over the whole input
	sumsOfOnes := getCountsOfOnes(structuredInput)
	// Now, to find Gamma and Epsilon rates, we need to verify whether each value
	//  in sumsOfOnes is more or less than half the number of inputs
	gammaRate := 0
	epsilonRate := 0
	for index, count := range sumsOfOnes {
		if count > (len(structuredInput)/2) {
			gammaRate += 1<<(len(sumsOfOnes)-index-1)
		} else {
			epsilonRate += 1<<(len(sumsOfOnes)-index-1)
		}
	}
	log.Printf("Part 1 - power consumptoin: %d\n", gammaRate*epsilonRate)

	// Part 2: multiply th oxygen generator rating by the CO2 scrubber rating = life support rating

	// Calculate oxygen rate
	oxygenRate := 0
	oRate, err := getRating(structuredInput, true, 0)
	if err != nil {
		log.Fatal(err)
	}
	for index, value := range oRate.value {
		if value {
			oxygenRate += 1 << (len(oRate.value) - index - 1)
		}
	}

	// Calculate CO2 rate
	co2Rate := 0
	co2RateStruct, err := getRating(structuredInput, false, 0)
	if err != nil {
		log.Fatal(err)
	}
	for index, value := range co2RateStruct.value {
		if value {
			co2Rate += 1 << (len(co2RateStruct.value) - index - 1)
		}
	}
	
	log.Printf("Part 2 - life support rating: %d\n", oxygenRate * co2Rate)
}

func getStructFromInput(path string) ([]Rate, error) {
	// Read the input file. It contains a list of instruction composed of
	// a string and an integer
	// up 3, down 5, forward 7, etc
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lines:=strings.Split(string(file), "\n")

	// Well, I know the size of the input, so let's just use that information
	// When knowing the size, it's better to allocate the right size immediately
	// as append() has a cost
	// https://medium.com/vendasta/golang-the-time-complexity-of-append-2177dcfb6bad
	// /!\ Do not use make([]Elements, 1000) as it will give it can AND size 1000, and
	// appending to it will just append after element 1000, so the first 1000 elements will be 0
	structuredInput := make([]Rate, 0, 1000)

	for _ , line := range lines {
		boolArr := make([]bool, 0, len(line)) // From the input, all elements are 12 bits long
		// Get the integer value from the line
		for _, c := range line {
			if c == '0' {
				boolArr = append(boolArr, false)
			} else {
				boolArr = append(boolArr, true)
			}
		}
		structuredInput = append(structuredInput, Rate{boolArr})
	}

	return structuredInput, nil
}

func getCountsOfOnes(structuredInput []Rate) []int {
	if len(structuredInput) == 0 {
		return nil
	}
	sumsOfOnes := make([]int, len(structuredInput[0].value))
	for _, rate := range structuredInput {
		for index, zeroOrOne := range rate.value {
			if zeroOrOne {
				sumsOfOnes[index] += 1
			}
		}
	}
	return sumsOfOnes
}

// input: a list of rates (the input from the exercise, successively filtered)
// defaultKeep: is used to define which rates should be kept in case the
// counts of 1 and 0 at the given index for the given input are equal
// To find the oxygen rate, use defaultKeep = 1, to find the CO2 rate, use defaultKeep = 0
// index: the bit index to check in the given rates
func getRating(input []Rate, defaultKeep bool, index int) (Rate, error) {
	// No input
	if len(input) == 0 {
		return Rate{nil}, errors.New("no input provided to getRating")
	}

	if index > len(input[0].value) {
		return Rate{nil}, errors.New("index out of bounds in getRating")
	}

	//Get the counts of '1' at each position for the given input
	sumsOfOnes := getCountsOfOnes(input)

	// Should we keep numbers in 0 or 1?
	keep := defaultKeep
	if float32(sumsOfOnes[index]) > (float32(len(input))/2) {
		// Most common value is 1
		keep = defaultKeep
	} else if float32(sumsOfOnes[index]) < (float32(len(input))/2) {
		// Most common value is 0
		keep = !defaultKeep
	}

	newInput := make([]Rate, 0, len(input))
	for _, rate := range input {
		if rate.value[index] == keep {
			newInput = append(newInput, rate)
		}
	}
	if len(newInput) == 1 {
		return newInput[0], nil
	} else {
		return getRating(newInput, defaultKeep, index+1)
	}
}
