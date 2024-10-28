package day11

import (
	"fmt"
	"os"

	filereader "github.com/jblashki/aoc-filereader-go"
)

const name = "Day 11 "
const inputFile = "./days/day11/input"

// RunDay runs Advent of Code Day 11 Puzzle
func RunDay(verbose bool) {
	var aResult int
	var bResult int
	var err error

	if verbose {
		fmt.Printf("\n%v Output:\n", name)
	}

	galaxyMap, err := readInput()
	if err != nil {
		fmt.Printf("%v: **** Error: %q ****\n", name, err)
		os.Exit(1)
	}

	galaxyMap, _ = expandMap(galaxyMap)

	if verbose {
		printMap(galaxyMap)
	}

	aResult, err = a(galaxyMap, verbose)
	if err != nil {
		fmt.Printf("%va: **** Error: %q ****\n", name, err)
	} else {
		fmt.Printf("%va: Answer = %v\n", name, aResult)
	}
	bResult, err = b(galaxyMap, verbose)
	if err != nil {
		fmt.Printf("%vb: **** Error: %q ****\n", name, err)
	} else {
		fmt.Printf("%vb: Answer = %v\n", name, bResult)
	}
}

type coord struct {
	row int
	col int
}

func getGalaxyList(galaxyMap [][]byte, verbose bool) []coord {
	retValue := make([]coord, 0)

	for row, line := range galaxyMap {
		for col, c := range line {
			if c == '#' {
				pos := coord{
					row: row,
					col: col,
				}
				retValue = append(retValue, pos)
			}
		}
	}

	return retValue
}

func countDistance(a coord, b coord, verbose bool) int {
	retValue := 0

	distance := max(a.row, b.row) - min(a.row, b.row)
	distance += max(a.col, b.col) - min(a.col, b.col)
	retValue += distance
	return retValue
}

func a(galaxyMap [][]byte, verbose bool) (int, error) {
	retValue := 0

	galaxies := getGalaxyList(galaxyMap, verbose)

	for i, galaxy := range galaxies {
		for j := i + 1; j < len(galaxies); j++ {
			retValue += countDistance(galaxy, galaxies[j], verbose)
		}
	}

	return retValue, nil
}

func b(galaxyMap [][]byte, verbose bool) (int, error) {
	retValue := 0

	galaxies := getGalaxyList(galaxyMap, verbose)

	for i, galaxy := range galaxies {
		for j := i + 1; j < len(galaxies); j++ {
			expansion := countExpansion(galaxyMap, galaxy, galaxies[j])
			distance := countDistance(galaxy, galaxies[j], verbose)

			updatedDistance := distance
			updatedDistance -= expansion
			updatedDistance += expansion * (1000000 - 1)
			retValue += updatedDistance
		}
	}

	return retValue, nil
}

func readInput() ([][]byte, error) {
	retMap := make([][]byte, 0)
	lines, err := filereader.ReadLines(inputFile)
	if err != nil {
		return retMap, err
	}

	for _, line := range lines {
		galaxyLine := make([]byte, len(line))
		for j, c := range line {
			galaxyLine[j] = (byte)(c)
		}
		retMap = append(retMap, galaxyLine)
	}

	return retMap, nil
}

func printMap(galaxyMap [][]byte) {
	for _, line := range galaxyMap {
		for _, galaxy := range line {
			fmt.Printf("%c", galaxy)
		}
		fmt.Println()
	}
}

func expandMap(galaxyMap [][]byte) ([][]byte, error) {
	cols := getInsertColumns(galaxyMap)

	for _, c := range cols {
		galaxyMap = insertCol(galaxyMap, c)
	}

	rows := getInsertRows(galaxyMap)
	for _, r := range rows {
		galaxyMap = insertRow(galaxyMap, r)
	}

	return galaxyMap, nil
}

func getInsertColumns(galaxyMap [][]byte) []int {
	retCols := make([]int, 0)
	nonEmptyCols := make([]bool, len(galaxyMap[0]))

	for _, line := range galaxyMap {
		for j, c := range line {
			if c == '#' {
				nonEmptyCols[j] = true
			}
		}
	}

	for i, col := range nonEmptyCols {
		if !col {
			retCols = append(retCols, i+len(retCols))
		}
	}

	return retCols
}

func getInsertRows(galaxyMap [][]byte) []int {
	retRows := make([]int, 0)
	nonEmptyRows := make([]bool, len(galaxyMap))

	for i, line := range galaxyMap {
		for _, c := range line {
			if c == '#' {
				nonEmptyRows[i] = true
			}
		}
	}

	for i, row := range nonEmptyRows {
		if !row {
			retRows = append(retRows, i+len(retRows))
		}
	}

	return retRows
}

func insertCol(galaxyMap [][]byte, col int) [][]byte {
	for i, line := range galaxyMap {
		if col == len(line) {
			line = append(line, 'x')
		} else {
			line = append(line[:col+1], line[col:]...)
			line[col] = 'x'
		}
		galaxyMap[i] = line
	}
	return galaxyMap
}

func insertRow(galaxyMap [][]byte, row int) [][]byte {
	newRow := make([]byte, len(galaxyMap[0]))

	for i, _ := range newRow {
		newRow[i] = 'x'
	}
	if row == len(galaxyMap) {
		galaxyMap = append(galaxyMap, newRow)
	} else {
		galaxyMap = append(galaxyMap[:row+1], galaxyMap[row:]...)
		galaxyMap[row] = newRow
	}
	return galaxyMap
}

func countExpansion(galaxyMap [][]byte, from coord, to coord) int {
	retValue := 0

	rowStart := min(from.row, to.row)
	rowEnd := max(from.row, to.row)

	colStart := min(from.col, to.col)
	colEnd := max(from.col, to.col)

	col := colStart
	if rowStart < len(galaxyMap) {
		for row := rowStart + 1; row < rowEnd; row++ {
			if galaxyMap[row][col] == 'x' {
				retValue++
			}
		}
	}

	row := rowStart
	if colStart < len(galaxyMap[rowStart]) {
		for col := colStart + 1; col < colEnd; col++ {
			if galaxyMap[row][col] == 'x' {
				retValue++
			}
		}
	}

	return retValue
}
