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
	num      int
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

	symbols, numberList, err2 := readNumbersAndSymbols(inputFile)
	if err2 != nil {
		fmt.Printf("%va: **** Error: %q ****\n", name, err)
		return
	}

	aResult, err = a(verbose, symbols, numberList)
	if err != nil {
		fmt.Printf("%va: **** Error: %q ****\n", name, err)
	} else {
		fmt.Printf("%va: Answer = %v\n", name, aResult)
	}
	bResult, err = b(verbose, symbols, numberList)
	if err != nil {
		fmt.Printf("%vb: **** Error: %q ****\n", name, err)
	} else {
		fmt.Printf("%vb: Answer = %v\n", name, bResult)
	}
}

func a(verbose bool, symbols [][]byte, numberList []number) (int, error) {
	retValue := 0

	for i := 0; i < len(numberList); i++ {
		if checkNumber(symbols, numberList[i]) {
			if verbose {
				fmt.Printf("%d matched\n", numberList[i].num)
			}
			retValue += numberList[i].num
		}
	}

	return retValue, nil
}

func b(verbose bool, symbols [][]byte, numberList []number) (int, error) {
	retValue := 0

	for i := 0; i < len(symbols); i++ {
		for j := 0; j < len(symbols[i]); j++ {
			if symbols[i][j] != 0 {
				if verbose {
					fmt.Printf("Checking symbol: %c (%d, %d)\n", symbols[i][j], i, j)
				}

				adjacentNumbers := getAdjacentNumbers(i, j, numberList)
				if len(adjacentNumbers) == 2 {
					ratio := adjacentNumbers[0].num * adjacentNumbers[1].num
					if verbose {
						fmt.Printf("Found gear @ (%d,%d). Ratio %dx%d = %d\n", i, j, adjacentNumbers[0].num, adjacentNumbers[1].num, ratio)
					}
					retValue += ratio
				}
			}
		}
		// if checkNumber(symbols, numberList[i]) {
		// 	if verbose {
		// 		fmt.Printf("%d matched\n", numberList[i].num)
		// 	}
		// 	retValue += numberList[i].num
		// }
	}

	return retValue, nil

	return retValue, nil
}

func getAdjacentNumbers(row int, column int, numberList []number) []number {
	retNumbers := make([]number, 0)

	for i := 0; i < len(numberList); i++ {
		for j := 0; j < len(numberList[i].adjacent); j++ {
			if row == numberList[i].adjacent[j].row && column == numberList[i].adjacent[j].column {
				retNumbers = append(retNumbers, numberList[i])
				break
			}
		}
	}
	return retNumbers
}

func checkNumber(symbols [][]byte, num number) bool {
	for i := 0; i < len(num.adjacent); i++ {
		if num.adjacent[i].row < len(symbols) && num.adjacent[i].column < len(symbols[num.adjacent[i].row]) {
			if symbols[num.adjacent[i].row][num.adjacent[i].column] != 0 {
				return true
			}
		}
	}

	return false
}

func readNumbersAndSymbols(inputFile string) ([][]byte, []number, error) {

	retNumbers := make([]number, 0)

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

		lineNumbers, err := readNumbers(lines[i], i)
		retNumbers = append(retNumbers, lineNumbers...)
		if err != nil {
			return nil, nil, err
		}
	}

	return retSymbols, retNumbers, nil
}

func readNumbers(line string, row int) ([]number, error) {
	retNumbers := make([]number, 0)

	currNum := 0
	currStart := -1

	for i := 0; i < len(line); i++ {
		if line[i] >= '0' && line[i] <= '9' {
			if currStart == -1 {
				currStart = i
			}
			currNum *= 10
			currNum += (int)(line[i] - '0')
		} else {
			if currStart != -1 {
				newNumber, err := makeNumber(currNum, currStart, row)
				if err != nil {
					return nil, err
				}
				retNumbers = append(retNumbers, newNumber)
				currStart = -1
				currNum = 0
			}
		}
	}

	if currStart != -1 {
		newNumber, err := makeNumber(currNum, currStart, row)
		if err != nil {
			return nil, err
		}
		retNumbers = append(retNumbers, newNumber)
	}

	return retNumbers, nil
}

func makeNumber(num int, col int, row int) (number, error) {
	adjacentList := make([]idx, 0)
	startCol := col - 1
	numLen := getDigitCount(num)
	endCol := col + numLen

	if row > 0 {
		// Row before
		for i := startCol; i <= endCol; i++ {
			if i >= 0 {
				adjacentList = append(adjacentList, idx{column: i, row: row - 1})
			}
		}
	}
	if startCol >= 0 {
		adjacentList = append(adjacentList, idx{column: startCol, row: row})
	}
	adjacentList = append(adjacentList, idx{column: endCol, row: row})

	// Row after
	for i := startCol; i <= endCol; i++ {
		if i >= 0 {
			adjacentList = append(adjacentList, idx{column: i, row: row + 1})
		}
	}

	retNumber := number{num: num, adjacent: adjacentList}
	return retNumber, nil
}

func getDigitCount(num int) int {
	count := 0
	for num > 0 {
		num /= 10
		count++
	}

	return count
}
