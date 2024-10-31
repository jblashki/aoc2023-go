package day14

import (
	"fmt"
	"os"

	filereader "github.com/jblashki/aoc-filereader-go"
)

const name = "Day 14"
const inputFile = "./days/day14/input"

const (
	NONE  = iota
	ROUND = iota
	CUBE  = iota
)

const (
	NORTH = iota
	EAST  = iota
	SOUTH = iota
	WEST  = iota
)

var cycles map[string]int

var states map[int][][]int

// RunDay runs Advent of Code Day 14 Puzzle
func RunDay(verbose bool) {
	var aResult int
	var bResult int
	var err error

	if verbose {
		fmt.Printf("\n%v Output:\n", name)
	}

	cycles = make(map[string]int)
	states = make(map[int][][]int)

	m, err := readInput(verbose)
	if err != nil {
		fmt.Printf("%v: **** Error: %q ****\n", name, err)
		os.Exit(1)
	}

	aResult, err = a(m, verbose)
	if err != nil {
		fmt.Printf("%va: **** Error: %q ****\n", name, err)
	} else {
		fmt.Printf("%va: Answer = %v\n", name, aResult)
	}
	m, err = readInput(verbose)
	bResult, err = b(m, verbose)
	if err != nil {
		fmt.Printf("%vb: **** Error: %q ****\n", name, err)
	} else {
		fmt.Printf("%vb: Answer = %v\n", name, bResult)
	}
}

func a(m [][]int, verbose bool) (int, error) {
	m = tip(m, NORTH, verbose)

	return calcLoad(m, verbose), nil

}

func b(m [][]int, verbose bool) (int, error) {
	key := getKey(m)
	cycles[key]++
	i := 0
	firstRepeat := -1
	loop := true
	for loop {
		m = tip(m, NORTH, verbose)
		m = tip(m, WEST, verbose)
		m = tip(m, SOUTH, verbose)
		m = tip(m, EAST, verbose)
		key = getKey(m)
		count, found := cycles[key]
		mcopy := copyMap(m)
		states[i] = mcopy
		if found {
			if firstRepeat == -1 {
				firstRepeat = i
			}
		}
		if count != 3 {
			cycles[key]++
			i++
		} else {
			loop = false
		}

	}

	num := 1000000000 - 1
	num -= firstRepeat
	num %= i - firstRepeat
	num += firstRepeat

	newMap := states[num]

	return calcLoad(newMap, verbose), nil
}

func readInput(verbose bool) ([][]int, error) {
	retValue := make([][]int, 0)

	lines, err := filereader.ReadLines(inputFile)
	if err != nil {
		return retValue, err
	}

	for _, line := range lines {
		newMapLine := make([]int, 0)
		for _, c := range line {
			value := NONE
			switch c {
			case '#':
				value = CUBE
			case 'O':
				value = ROUND
			default:
				value = NONE
			}
			newMapLine = append(newMapLine, value)
		}
		retValue = append(retValue, newMapLine)
	}

	return retValue, nil
}

func printMap(m [][]int) {
	for _, line := range m {
		for _, node := range line {
			switch node {
			case NONE:
				fmt.Printf(".")
			case ROUND:
				fmt.Printf("O")
			case CUBE:
				fmt.Printf("#")
			}
		}
		fmt.Println()
	}
}

func tip(m [][]int, direction int, verbose bool) [][]int {
	switch direction {
	case NORTH:
		return tipNorth(m, verbose)
	case EAST:
		return tipEast(m, verbose)
	case SOUTH:
		return tipSouth(m, verbose)
	case WEST:
		return tipWest(m, verbose)
	}

	return m
}

func tipNorth(m [][]int, verbose bool) [][]int {
	for row := 1; row < len(m); row++ {
		for col, node := range m[row] {
			if node == ROUND {
				newRow := row
				for newRow > 0 && m[newRow-1][col] == NONE {
					newRow--
				}
				if newRow != row {
					m[newRow][col] = ROUND
					m[row][col] = NONE
				}
			}
		}
	}
	return m
}

func tipEast(m [][]int, verbose bool) [][]int {
	for col := len(m[0]) - 1; col >= 0; col-- {
		for row := 0; row < len(m); row++ {
			node := m[row][col]
			if node == ROUND {
				newCol := col
				for newCol < len(m[0])-1 && m[row][newCol+1] == NONE {
					newCol++
				}

				if newCol != col {
					m[row][newCol] = ROUND
					m[row][col] = NONE
				}
			}
		}
	}
	return m
}

func tipSouth(m [][]int, verbose bool) [][]int {
	for row := len(m) - 1; row >= 0; row-- {
		for col, node := range m[row] {
			if node == ROUND {
				newRow := row
				for newRow < len(m)-1 && m[newRow+1][col] == NONE {
					newRow++
				}
				if newRow != row {
					m[newRow][col] = ROUND
					m[row][col] = NONE
				}
			}
		}
	}
	return m
}

func tipWest(m [][]int, verbose bool) [][]int {
	for col := 1; col < len(m[0]); col++ {
		for row := 0; row < len(m); row++ {
			node := m[row][col]
			if node == ROUND {
				newCol := col
				for newCol > 0 && m[row][newCol-1] == NONE {
					newCol--
				}

				if newCol != col {
					m[row][newCol] = ROUND
					m[row][col] = NONE
				}
			}
		}
	}
	return m
}

func calcLoad(m [][]int, verbose bool) int {
	retValue := 0
	totalRows := len(m)

	for i, line := range m {
		for _, rock := range line {
			if rock == ROUND {
				retValue += (totalRows - i)
			}
		}
	}

	return retValue
}

func getKey(m [][]int) string {
	return fmt.Sprintf("%v", m)
}

func copyMap(m [][]int) [][]int {
	returnMap := make([][]int, 0)
	for _, line := range m {
		newLine := make([]int, 0)
		for _, rock := range line {
			newLine = append(newLine, rock)
		}
		returnMap = append(returnMap, newLine)
	}

	return returnMap
}
