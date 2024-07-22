package day3

import (
	"errors"
	"fmt"

	filereader "github.com/jblashki/aoc-filereader-go"
)

const name = "Day 3"
const inputFile = "./days/day3/input"

type idx struct {
	row    int
	column int
}

type number struct {
	number   int
	adjacent []idx
}

// RunDay runs Advent of Code Day 3 Puzzle
func RunDay(verbose bool) {
	var aResult int
	var bResult int
	var err error

	if verbose {
		fmt.Printf("\n%v Output:\n", name)
	}

	symbols, _, err2 := readNumbersAndSymbols(inputFile)
	if err2 != nil {
		fmt.Printf("%va: **** Error: %q ****\n", name, err)
		return
	}

	for i := 0; i < len(symbols); i++ {
		for j := 0; j < len(symbols[i]); j++ {
			if symbols[i][j] == 0 {
				fmt.Printf(" ")
			} else {
				fmt.Printf("%c", symbols[i][j])
			}
		}
		fmt.Printf("\n")
	}

	aResult, err = a(verbose)
	if err != nil {
		fmt.Printf("%va: **** Error: %q ****\n", name, err)
	} else {
		fmt.Printf("%va: Answer = %v\n", name, aResult)
	}
	bResult, err = b(verbose)
	if err != nil {
		fmt.Printf("%vb: **** Error: %q ****\n", name, err)
	} else {
		fmt.Printf("%vb: Answer = %v\n", name, bResult)
	}
}

func a(verbose bool) (int, error) {
	retValue := 0

	return retValue, nil
}

func b(verbose bool) (int, error) {
	retValue := 0

	return retValue, nil
}

func readNumbersAndSymbols(inputFile string) ([][]byte, []number, error) {

	lines, err := filereader.ReadLines(inputFile)
	if err != nil {
		return nil, nil, err
	}

	if len(lines) <= 0 || len(lines[0]) <= 0 {
		return nil, nil, errors.New("invalid input")
	}

	lineLen := len(lines[0])

	retSymbols := make([][]byte, len(lines))
	for i := 0; i < len(lines); i++ {
		retSymbols[i] = make([]byte, lineLen)
	}

	fmt.Printf(("HERE\n"))

	for i := 0; i < len(lines); i++ {
		if len(lines[i]) != lineLen {
			return nil, nil, fmt.Errorf("invalid line length on line %d", i)
		}

		for j := 0; j < lineLen; j++ {
			// Load Symbols array
			if (lines[i][j] < '0' || lines[i][j] > '9') &&
				lines[i][j] != '.' {
				retSymbols[i][j] = lines[i][j]
			}
		}
	}

	return retSymbols, nil, nil
}
