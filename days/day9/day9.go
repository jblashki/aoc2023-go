package day9

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	filereader "github.com/jblashki/aoc-filereader-go"
)

const name = "Day 9 "
const inputFile = "./days/day9/input"

// RunDay runs Advent of Code Day 8 Puzzle
func RunDay(verbose bool) {
	var aResult int
	var bResult int
	var err error

	if verbose {
		fmt.Printf("\n%v Output:\n", name)
	}

	series, err := readInput()
	if err != nil {
		fmt.Printf("%v: **** Error: %q ****\n", name, err)
		os.Exit(0)
	}

	aResult, err = a(series, verbose)
	if err != nil {
		fmt.Printf("%va: **** Error: %q ****\n", name, err)
	} else {
		fmt.Printf("%va: Answer = %v\n", name, aResult)
	}
	bResult, err = b(series, verbose)
	if err != nil {
		fmt.Printf("%vb: **** Error: %q ****\n", name, err)
	} else {
		fmt.Printf("%vb: Answer = %v\n", name, bResult)
	}
}

func a(series [][]int, verbose bool) (int, error) {
	retValue := 0

	for _, v := range series {
		nextValue, err := getNextValue(v)
		if err != nil {
			fmt.Printf("%v: **** Error: %q ****\n", name, err)
			return retValue, nil
		}
		retValue += nextValue
	}

	return retValue, nil
}

func b(series [][]int, verbose bool) (int, error) {
	retValue := 0

	for _, v := range series {
		nextValue, err := getPrevValue(v)
		if err != nil {
			fmt.Printf("%v: **** Error: %q ****\n", name, err)
			return retValue, nil
		}
		retValue += nextValue
	}

	return retValue, nil
}

func readInput() ([][]int, error) {
	retValue := make([][]int, 0)
	lines, err := filereader.ReadLines(inputFile)
	if err != nil {
		return nil, err
	}

	for i := 0; i < len(lines); i++ {
		numbersAsStrings := strings.Split(lines[i], " ")

		numbers := make([]int, 0)
		for j := 0; j < len(numbersAsStrings); j++ {
			num, err := strconv.Atoi(numbersAsStrings[j])
			if err != nil {
				return retValue, err
			}
			numbers = append(numbers, num)
		}
		retValue = append(retValue, numbers)
	}

	return retValue, nil
}

func getNextValue(numbers []int) (int, error) {
	seriesPyramid := make([][]int, 0)

	newSeries := make([]int, 0)
	for i := 0; i < len(numbers); i++ {
		newSeries = append(newSeries, numbers[i])
	}
	newSeries = append(newSeries, 0)
	seriesPyramid = append(seriesPyramid, newSeries)

	loop := true
	for i := 0; loop; i++ {
		newSeries := make([]int, 0)
		loop = false
		for j := 1; j < len(seriesPyramid[i])-1; j++ {

			newValue := seriesPyramid[i][j] - seriesPyramid[i][j-1]
			newSeries = append(newSeries, newValue)
			if newValue != 0 {
				loop = true
			}
		}
		newSeries = append(newSeries, 0)
		seriesPyramid = append(seriesPyramid, newSeries)
	}

	prevValue := 0
	for i := len(seriesPyramid) - 1; i >= 0; i-- {
		series := seriesPyramid[i]

		prevValue = prevValue + series[len(series)-2]
		series[len(series)-1] = prevValue
	}

	return prevValue, nil
}

func getPrevValue(numbers []int) (int, error) {
	seriesPyramid := make([][]int, 0)

	newSeries := make([]int, 0)
	for i := 0; i < len(numbers); i++ {
		newSeries = append(newSeries, numbers[i])
	}
	newSeries = append([]int{0}, newSeries...)
	seriesPyramid = append(seriesPyramid, newSeries)

	loop := true
	for i := 0; loop; i++ {
		newSeries := make([]int, 0)
		loop = false
		for j := 2; j < len(seriesPyramid[i]); j++ {
			newValue := seriesPyramid[i][j-1] - seriesPyramid[i][j]
			newSeries = append(newSeries, newValue)
			if newValue != 0 {
				loop = true
			}
		}
		newSeries = append([]int{0}, newSeries...)
		seriesPyramid = append(seriesPyramid, newSeries)
	}

	prevValue := 0
	for i := len(seriesPyramid) - 1; i >= 0; i-- {
		series := seriesPyramid[i]
		prevValue = series[1] + prevValue
		series[0] = prevValue
	}

	return prevValue, nil
}
