package day1

import (
	"errors"
	"fmt"
	"strings"

	filereader "github.com/jblashki/aoc-filereader-go"
)

const name = "Day 1"
const inputFile = "./days/day1/input"

var nums = []string{
	"ONE",
	"TWO",
	"THREE",
	"FOUR",
	"FIVE",
	"SIX",
	"SEVEN",
	"EIGHT",
	"NINE",
}

var numsRev = []string{
	"ENO",
	"OWT",
	"EERHT",
	"RUOF",
	"EVIF",
	"XIS",
	"NEVES",
	"THGIE",
	"ENIN",
}

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
	lines, err := filereader.ReadLines(inputFile)

	if err != nil {
		return 0, err
	}

	for i := 0; i < len(lines); i++ {
		num, err := getNum(strings.ToUpper(lines[i]), false)
		if err != nil {
			return 0, err
		}
		retValue += num

		if verbose {
			fmt.Printf("Line: %v = %v\n", lines[i], num)
		}
	}

	return retValue, nil
}

func b(verbose bool) (int, error) {
	retValue := 0
	lines, err := filereader.ReadLines(inputFile)

	if err != nil {
		return 0, err
	}

	for i := 0; i < len(lines); i++ {
		num, err := getNum(strings.ToUpper(lines[i]), true)
		if err != nil {
			return 0, err
		}
		retValue += num

		if verbose {
			fmt.Printf("Line: %v = %v\n", lines[i], num)
		}
	}

	return retValue, nil
}

func getNum(s string, adv bool) (int, error) {
	retNum := 0

	digit, err := getDigit(s, false, adv)
	if err != nil {
		return 0, err
	}
	retNum += digit * 10
	revS := Reverse(s)
	digit, err = getDigit(revS, true, adv)
	if err != nil {
		return 0, err
	}
	retNum += digit

	return retNum, nil
}

func getDigit(s string, rev bool, adv bool) (int, error) {
	for i := 0; i < len(s); i++ {
		if s[i] >= '0' && s[i] <= '9' {
			return (int)(s[i] - '0'), nil
		}
		if adv {
			var arr = []string{}
			if rev {
				arr = numsRev
			} else {
				arr = nums
			}

			for j := 0; j < len(arr); j++ {
				if checkNumWord(s, arr[j], i) {
					return j + 1, nil
				}
			}

			// check numerics
		}
	}

	return 0, errors.New("no digit found")
}

func Reverse(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func checkNumWord(s string, number string, startPos int) bool {
	for i := 0; i < len(number); i++ {
		if s[i+startPos] == 0 {
			return false
		}
		if s[i+startPos] != number[i] {
			return false
		}
	}

	return true
}
