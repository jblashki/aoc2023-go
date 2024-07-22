package day4

import (
	"fmt"
)

const name = "Day 4"
const inputFile = "./days/day4/input"

// RunDay runs Advent of Code Day 4 Puzzle
func RunDay(verbose bool) {
	var aResult int
	var bResult int
	var err error

	if verbose {
		fmt.Printf("\n%v Output:\n", name)
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

	if verbose {
		fmt.Printf("func a not implemented yet\n")
	}

	return retValue, nil
}

func b(verbose bool) (int, error) {
	retValue := 0

	if verbose {
		fmt.Printf("func b not implemented yet\n")
	}

	return retValue, nil
}