package day13

import (
	"fmt"
	"os"

	filereader "github.com/jblashki/aoc-filereader-go"
)

const name = "Day 13"
const inputFile = "./days/day13/input"

type mirrorMap struct {
	lines [][]byte
}

const (
	HORIZ = iota
	VERT  = iota
)

type foundMirror struct {
	mirrorType int
	pos        int
}

type coord struct {
	row int
	col int
}

var foundMirrors []foundMirror = make([]foundMirror, 0)

// RunDay runs Advent of Code Day 13 Puzzle
func RunDay(verbose bool) {
	var aResult int
	var bResult int
	var err error

	if verbose {
		fmt.Printf("\n%v Output:\n", name)
	}

	mirrorMaps, err := readInput(verbose)
	if err != nil {
		fmt.Printf("%v: **** Error: %q ****\n", name, err)
		os.Exit(1)
	}

	if verbose {
		fmt.Printf("map count = %d\n", len(mirrorMaps))
		for i, mm := range mirrorMaps {
			fmt.Printf("=============\n")
			fmt.Printf("Mirror Map %d\n", i+1)
			fmt.Printf("=============\n")
			for _, line := range mm.lines {
				fmt.Printf("%s\n", line)
			}
			fmt.Println()
		}
	}

	aResult, err = a(mirrorMaps, verbose)
	if err != nil {
		fmt.Printf("%va: **** Error: %q ****\n", name, err)
	} else {
		fmt.Printf("%va: Answer = %v\n", name, aResult)
	}
	bResult, err = b(mirrorMaps, verbose)
	if err != nil {
		fmt.Printf("%vb: **** Error: %q ****\n", name, err)
	} else {
		fmt.Printf("%vb: Answer = %v\n", name, bResult)
	}
}

func getScore(mirrorMaps []mirrorMap, setFound bool, verbose bool) int {
	retValue := 0
	for i, mm := range mirrorMaps {
		if verbose {
			fmt.Printf("=======================\n")
			fmt.Printf("Checking Mirror Map %d\n", i+1)
			fmt.Printf("=======================\n")
		}

		found := foundMirror{
			mirrorType: -1,
			pos:        -1,
		}
		if i < len(foundMirrors) {
			found = foundMirrors[i]
		}
		for j := 0; j < len(mm.lines)-1; j++ {
			if isHorizontalMirror(mm, j) && (found.mirrorType != HORIZ || found.pos != j) {
				if setFound {
					found = foundMirror{
						mirrorType: HORIZ,
						pos:        j,
					}
					if i >= len(foundMirrors) {
						foundMirrors = append(foundMirrors, found)
					} else {
						foundMirrors[i] = found
					}
				}
				score := (j + 1) * 100
				retValue += score
				if verbose {
					fmt.Printf("Found mirror at row: %d->%d (%d)\n", j+1, j+2, score)
					printHorizontal(mm, j)
				}
			}
		}

		for j := 0; j < len(mm.lines[0])-1; j++ {
			if isVerticalMirror(mm, j) && (found.mirrorType != VERT || found.pos != j) {
				if setFound {
					found = foundMirror{
						mirrorType: VERT,
						pos:        j,
					}
					if i >= len(foundMirrors) {
						foundMirrors = append(foundMirrors, found)
					} else {
						foundMirrors[i] = found
					}
				}
				score := (j + 1)
				retValue += score
				if verbose {
					fmt.Printf("Found mirror at col: %d->%d (%d)\n", j+1, j+2, score)
					printVertical(mm, j)
				}
			}
		}
	}

	return retValue
}

func a(mirrorMaps []mirrorMap, verbose bool) (int, error) {
	retValue := getScore(mirrorMaps, true, verbose)

	return retValue, nil
}

func b(mirrorMaps []mirrorMap, verbose bool) (int, error) {
	mirrorMaps = fixSmudges(mirrorMaps, verbose)
	retValue := getScore(mirrorMaps, false, verbose)

	return retValue, nil
}

func fixSmudges(mirrorMap []mirrorMap, verbose bool) []mirrorMap {
	for i, mm := range mirrorMap {
		mm = fixSmudge(mm, verbose)
		mirrorMap[i] = mm
	}
	return mirrorMap
}

func fixSmudge(mm mirrorMap, verbose bool) mirrorMap {
	for j := 0; j < len(mm.lines)-1; j++ {
		smudge, err := findHorizSmudge(mm, j)
		if err == nil {
			if verbose {
				fmt.Printf("Found smudge @ (%d,%d)\n", smudge.col, smudge.row)
			}
			value := mm.lines[smudge.row][smudge.col]
			if value == '#' {
				value = '.'
			} else {
				value = '#'
			}
			mm.lines[smudge.row][smudge.col] = value
			return mm
		}
	}

	for j := 0; j < len(mm.lines[0])-1; j++ {
		smudge, err := findVertSmudge(mm, j)
		if err == nil {
			if verbose {
				fmt.Printf("Found smudge @ (%d,%d)\n", smudge.col, smudge.row)
			}
			value := mm.lines[smudge.row][smudge.col]
			if value == '#' {
				value = '.'
			} else {
				value = '#'
			}
			mm.lines[smudge.row][smudge.col] = value
			return mm
		}
	}

	return mm
}

func readInput(verbose bool) ([]mirrorMap, error) {
	retValue := make([]mirrorMap, 0)

	lines, err := filereader.ReadLines(inputFile)
	if err != nil {
		return retValue, err
	}

	newLines := make([][]byte, 0)
	for _, line := range lines {
		if line == "" {
			mm := mirrorMap{
				lines: newLines,
			}
			retValue = append(retValue, mm)
			newLines = make([][]byte, 0)
		} else {
			newLine := make([]byte, 0)
			for _, value := range line {
				newLine = append(newLine, (byte)(value))
			}
			newLines = append(newLines, newLine)
		}
	}
	mm := mirrorMap{
		lines: newLines,
	}
	retValue = append(retValue, mm)

	return retValue, nil
}

func isVerticalMirror(mm mirrorMap, col int) bool {
	for colA, colB := col, col+1; colB < len(mm.lines[0]) && colA >= 0; colA, colB = colA-1, colB+1 {
		for row := 0; row < len(mm.lines); row++ {
			if mm.lines[row][colA] != mm.lines[row][colB] {
				return false
			}
		}
	}

	return true
}

func countNonMatchCols(mm mirrorMap, col int) []int {
	retValue := make([]int, 0)
	for colA, colB := col, col+1; colB < len(mm.lines[0]) && colA >= 0; colA, colB = colA-1, colB+1 {
		nonMatchCount := 0
		for row := 0; row < len(mm.lines); row++ {
			if mm.lines[row][colA] != mm.lines[row][colB] {
				nonMatchCount++
			}
		}
		retValue = append(retValue, nonMatchCount)
	}

	return retValue
}

func findVertSmudge(mm mirrorMap, col int) (coord, error) {
	smudge := coord{
		row: -1,
		col: -1,
	}
	oneCount := 0
	for colA, colB := col, col+1; colB < len(mm.lines[0]) && colA >= 0; colA, colB = colA-1, colB+1 {
		nonMatchCount := 0
		for row := 0; row < len(mm.lines); row++ {
			if mm.lines[row][colA] != mm.lines[row][colB] {
				nonMatchCount++
				smudge.row = row
				smudge.col = colA
			}
		}
		if nonMatchCount > 1 {
			return coord{row: -1, col: -1}, fmt.Errorf("no smudge")
		}
		if nonMatchCount == 1 {
			oneCount++
		}
	}

	if oneCount != 1 {
		return coord{row: -1, col: -1}, fmt.Errorf("no smudge")
	}
	return smudge, nil
}

func isHorizontalMirror(mm mirrorMap, row int) bool {
	for col := 0; col < len(mm.lines[0]); col++ {
		for rowA, rowB := row, row+1; rowB < len(mm.lines) && rowA >= 0; rowA, rowB = rowA-1, rowB+1 {
			if mm.lines[rowA][col] != mm.lines[rowB][col] {
				return false
			}
		}
	}

	return true
}

func countNonMatchRows(mm mirrorMap, row int) []int {
	retValue := make([]int, 0)
	for col := 0; col < len(mm.lines[0]); col++ {
		nonMatchCount := 0
		for rowA, rowB := row, row+1; rowB < len(mm.lines) && rowA >= 0; rowA, rowB = rowA-1, rowB+1 {
			if mm.lines[rowA][col] != mm.lines[rowB][col] {
				nonMatchCount++
			}
		}

		retValue = append(retValue, nonMatchCount)
	}

	return retValue
}

func findHorizSmudge(mm mirrorMap, row int) (coord, error) {
	smudge := coord{
		row: -1,
		col: -1,
	}
	oneCount := 0
	for col := 0; col < len(mm.lines[0]); col++ {
		nonMatchCount := 0
		for rowA, rowB := row, row+1; rowB < len(mm.lines) && rowA >= 0; rowA, rowB = rowA-1, rowB+1 {
			if mm.lines[rowA][col] != mm.lines[rowB][col] {
				nonMatchCount++
				smudge.row = rowA
				smudge.col = col
			}
		}
		if nonMatchCount > 1 {
			return coord{row: -1, col: -1}, fmt.Errorf("no smudge")
		}

		if nonMatchCount == 1 {
			oneCount++
		}
	}

	if oneCount != 1 {
		return coord{row: -1, col: -1}, fmt.Errorf("no smudge")
	}
	return smudge, nil
}

func printHorizontal(mm mirrorMap, row int) {
	for i, line := range mm.lines {
		if i == row {
			fmt.Printf("v")
		} else if i == row+1 {
			fmt.Printf("^")
		} else {
			fmt.Printf(" ")
		}
		fmt.Printf("%s", line)
		if i == row {
			fmt.Printf("v")
		} else if i == row+1 {
			fmt.Printf("^")
		} else {
			fmt.Printf(" ")
		}
		fmt.Println()
	}
}

func printVertical(mm mirrorMap, col int) {
	for i := 0; i < len(mm.lines[0]); i++ {
		if i == col {
			fmt.Printf(">")
		} else if i == col+1 {
			fmt.Printf("<")
		} else {
			fmt.Printf(" ")
		}
	}
	fmt.Println()
	for _, line := range mm.lines {
		fmt.Printf("%s\n", line)
	}
	for i := 0; i < len(mm.lines[0]); i++ {
		if i == col {
			fmt.Printf(">")
		} else if i == col+1 {
			fmt.Printf("<")
		} else {
			fmt.Printf(" ")
		}
	}
	fmt.Println()
}
