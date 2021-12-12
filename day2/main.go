package main

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type Elements struct {
	word string
	value int
}

type Position struct {
	aim int
	depth int
	horizontal int
}

func main() {
	// Read the input file. It contains a list of integers representing depth measures
	structuredInput, err := getStructFromInput("day2/input.txt")
	if err != nil {
		log.Fatal(err)
	}

	// Part 1: multiply depth by horizontal distance
	m := make(map[string]int, 3)
	for _, element := range structuredInput {
		m[element.word] += element.value
	}
	result := m["forward"] * (m["down"] - m["up"])

	log.Printf("Part 1 - depth * horizontal distance: %d\n", result)

	// Part 2: Different instructions, run new depth * horizontal distance
	position := Position{0, 0, 0}
	for index, element := range structuredInput {
		switch element.word {
		case "down":
			position.aim += element.value
		case "up":
			position.aim -= element.value
		case "forward":
			position.horizontal += element.value
			position.depth += position.aim * element.value
		default:
			log.Fatalf("Unexpected instruction at line %d.\nExpected: up, down or forward\nFound: %s", index, element.word)
		}
	}

	log.Printf("Part 2 - depth * horizontal distance: %d\n", position.depth * position.horizontal)
}

func getStructFromInput(path string) ([]Elements, error) {
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
	structuredInput := make([]Elements, 0, 1000)

	for _ , line := range lines {
		// A line is supposed to be composed of a single word, a white space and an integer
		parts:=strings.Split(line, " ")
		if len(parts) != 2 {
			log.Fatalf("Unexpected content in input file.\nExpected: \"string int\"\nFound: %s", line)
		}
		// Get the integer value from the line
		value, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}
		// get the value for the corresponding word from the map
		structuredInput = append(structuredInput, Elements {parts[0], value})
	}

	return structuredInput, nil
}


