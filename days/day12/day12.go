package day12

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	filereader "github.com/jblashki/aoc-filereader-go"
)

const name = "Day 12 "
const inputFile = "./days/day12/input"

const (
	OPERATIONAL = '.'
	DAMAGED     = '#'
	UNKNOWN     = '?'
)

type springLine struct {
	springs       []byte
	missing       []int
	damagedGroups []int
}

// RunDay runs Advent of Code Day 12 Puzzle
func RunDay(verbose bool) {
	var aResult int
	var bResult int
	var err error

	if verbose {
		fmt.Printf("\n%v Output:\n", name)
	}

	springLines, err := readInput()
	if err != nil {
		fmt.Printf("%v: **** Error: %q ****\n", name, err)
		os.Exit(1)
	}

	aResult, err = a(springLines, verbose)
	if err != nil {
		fmt.Printf("%va: **** Error: %q ****\n", name, err)
	} else {
		fmt.Printf("%va: Answer = %v\n", name, aResult)
	}
	bResult, err = b(springLines, verbose)
	if err != nil {
		fmt.Printf("%vb: **** Error: %q ****\n", name, err)
	} else {
		fmt.Printf("%vb: Answer = %v\n", name, bResult)
	}
}

func getCount(springLines []springLine, verbose bool) (int, error) {
	retValue := 0
	for i, springLineDetail := range springLines {
		newValue := getValidComb(springLineDetail.springs, springLineDetail.damagedGroups)
		fmt.Printf("Line %d: %s - %d\n", i+1, springLineDetail.springs, newValue)
		retValue += newValue

	}

	return retValue, nil
}

func a(springLines []springLine, verbose bool) (int, error) {
	retValue, err := getCount(springLines, verbose)
	if err != nil {
		return retValue, err
	}

	return retValue, nil
}

func b(springLines []springLine, verbose bool) (int, error) {
	retValue := 0
	//printSpringLines(springLines)
	for i, sl := range springLines {
		//fmt.Printf("Line %d\n", i)
		newSprings := make([]byte, 0)
		for j := 0; j < 5; j++ {
			if j >= 1 {
				newSprings = append(newSprings, '?')
			}
			newSprings = append(newSprings, sl.springs...)
		}
		sl.springs = newSprings

		newGroups := make([]int, 0)
		for j := 0; j < 5; j++ {
			newGroups = append(newGroups, sl.damagedGroups...)
		}
		sl.damagedGroups = newGroups

		newMissing := make([]int, 0)
		for j, spring := range sl.springs {
			//sl.springs = append(sl.springs, byte(spring))
			if byte(spring) == UNKNOWN {
				newMissing = append(newMissing, j)
			}
		}
		sl.missing = newMissing
		springLines[i] = sl
	}

	//printSpringLines(springLines)

	retValue, err := getCount(springLines, verbose)
	if err != nil {
		return retValue, err
	}

	return retValue, nil
}

func readInput() ([]springLine, error) {
	retValue := make([]springLine, 0)
	lines, err := filereader.ReadLines(inputFile)
	if err != nil {
		return retValue, err
	}

	for i, line := range lines {

		springLineDetail := springLine{
			springs:       make([]byte, 0),
			missing:       make([]int, 0),
			damagedGroups: make([]int, 0),
		}
		data := strings.Split(line, " ")
		if len(data) != 2 {
			return retValue, fmt.Errorf("invalid input line %d: %v", i, line)
		}

		for j, spring := range data[0] {
			springLineDetail.springs = append(springLineDetail.springs, byte(spring))
			if byte(spring) == UNKNOWN {
				springLineDetail.missing = append(springLineDetail.missing, j)
			}
		}

		groups := strings.Split(data[1], ",")
		for j, group := range groups {
			groupInt, err := strconv.Atoi(group)
			if err != nil {
				return retValue, fmt.Errorf("invalid input line %d. Group %v invalid: %v", i, j, line)
			}
			springLineDetail.damagedGroups = append(springLineDetail.damagedGroups, groupInt)
		}
		retValue = append(retValue, springLineDetail)
	}

	return retValue, nil
}

func printSprintLine(line int, sl springLine) {
	fmt.Printf("%d) %v - %v - %v\n", line, (string)(sl.springs), sl.missing, sl.damagedGroups)
}

func printSpringLines(springLines []springLine) {
	for i, sl := range springLines {
		printSprintLine(i, sl)

	}
}

// const (
// 	BITMASK = 0x01
// )

// func getBitArray(value int, arrayLen int) []bool {
// 	retArray := make([]bool, arrayLen)

// 	for i := 0; i < arrayLen; i++ {
// 		bit := (value >> i) & BITMASK
// 		if bit == BITMASK {
// 			retArray[i] = true
// 		}
// 	}

// 	return retArray
// }

// func getSprings(springLineDetail springLine, value int) []byte {
// 	retArray := append([]byte(nil), springLineDetail.springs...)
// 	bitArrayLen := len(springLineDetail.missing)
// 	bitArray := getBitArray(value, bitArrayLen)
// 	for i, bit := range bitArray {
// 		pos := springLineDetail.missing[i]
// 		if bit {
// 			retArray[pos] = OPERATIONAL
// 		} else {
// 			retArray[pos] = DAMAGED
// 		}
// 	}
// 	return retArray
// }

// func getGroups(springs []byte) []int {
// 	retArray := make([]int, 0)

// 	inGroup := false
// 	curGroup := -1
// 	curCount := 0
// 	for i := 0; i < len(springs); i++ {
// 		if springs[i] == DAMAGED {
// 			if !inGroup {
// 				curGroup++
// 				inGroup = true
// 			}
// 			curCount++
// 		} else {
// 			if inGroup {
// 				retArray = append(retArray, curCount)
// 			}
// 			inGroup = false
// 			curCount = 0
// 		}
// 	}
// 	if inGroup {
// 		retArray = append(retArray, curCount)
// 	}
// 	return retArray
// }

// func IntArrayEquals(a []int, b []int) bool {
// 	if len(a) != len(b) {
// 		return false
// 	}
// 	for i, v := range a {
// 		if v != b[i] {
// 			return false
// 		}
// 	}
// 	return true
// }

func getValidComb(springs []byte, damagedGroups []int) int {
	//fmt.Printf("%s -> %v\n", springs, damagedGroups)
	if len(damagedGroups) == 0 {
		if contains(springs, DAMAGED) {
			return 0
		} else {
			return 1
		}
	}

	if len(springs) == 0 {
		if len(damagedGroups) == 0 {
			return 1
		} else {
			return 0
		}
	}

	retValue := 0

	if springs[0] == OPERATIONAL || springs[0] == UNKNOWN {
		retValue += getValidComb(springs[1:], damagedGroups)
	}

	if springs[0] == DAMAGED || springs[0] == UNKNOWN {
		if isValidComb(springs, damagedGroups) {
			if damagedGroups[0] >= len(springs) {
				retValue += getValidComb(nil, damagedGroups[1:])
			} else {
				retValue += getValidComb(springs[damagedGroups[0]+1:], damagedGroups[1:])
			}
		}
	}

	return retValue
}

func isValidComb(springs []byte, damagedGroups []int) bool {
	return (damagedGroups[0] <= len(springs) &&
		!contains(springs[:damagedGroups[0]], OPERATIONAL) &&
		(damagedGroups[0] == len(springs) ||
			springs[damagedGroups[0]] != DAMAGED))
}

func contains(springArray []byte, b byte) bool {
	for _, spring := range springArray {
		if spring == b {
			return true
		}
	}
	return false
}
