package day4

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"slices"
	"strconv"
	"strings"

	filereader "github.com/jblashki/aoc-filereader-go"
)

const name = "Day 4"
const inputFile = "./days/day4/input"

type card struct {
	numbers        []int
	winningNumbers []int
	matching       []int
}

// RunDay runs Advent of Code Day 4 Puzzle
func RunDay(verbose bool) {
	var aResult int
	var bResult int
	var err error

	if verbose {
		fmt.Printf("\n%v Output:\n", name)
	}

	cards, err := readInput(verbose, inputFile)
	if err != nil {
		fmt.Printf("%va: **** Error: %q ****\n", name, err)
		return
	}

	aResult, err = a(verbose, cards)
	if err != nil {
		fmt.Printf("%va: **** Error: %q ****\n", name, err)
	} else {
		fmt.Printf("%va: Answer = %v\n", name, aResult)
	}
	bResult, err = b(verbose, cards)
	if err != nil {
		fmt.Printf("%vb: **** Error: %q ****\n", name, err)
	} else {
		fmt.Printf("%vb: Answer = %v\n", name, bResult)
	}
}

func a(verbose bool, cards []card) (int, error) {
	retValue := 0

	for i := 0; i < len(cards); i++ {

		matchCount := len(cards[i].matching)
		points := 0

		if matchCount > 0 {
			matchCount -= 1
			points = int(math.Pow(2, float64(matchCount)))
		}

		if verbose {
			fmt.Printf("Card: %d\n", i+1)
			fmt.Printf("\t\tWinning:  %v\n", cards[i].winningNumbers)
			fmt.Printf("\t\tNum:      %v\n", cards[i].numbers)
			fmt.Printf("\t\tMatching: %v\n", cards[i].matching)
			fmt.Printf("\t\tPoints:    %d\n", points)
		}

		retValue += points
	}

	return retValue, nil
}

func b(verbose bool, cards []card) (int, error) {
	retValue := 0

	copies := make([]int, len(cards))
	for i := 0; i < len(copies); i++ {
		copies[i] = 1
	}

	for i := 0; i < len(cards); i++ {
		if verbose {
			fmt.Printf("Processing card %d. Matches: %d. Has %d copies\n", i+1, len(cards[i].matching), copies[i])
		}
		matchCount := len(cards[i].matching)
		for j := 0; j < copies[i]; j++ {
			for k := i + 1; k <= i+matchCount && k < len(cards); k++ {
				copies[k]++
			}
			retValue++
		}
	}

	return retValue, nil
}

func readInput(verbose bool, inputFile string) ([]card, error) {
	retCards := make([]card, 0)

	lines, err := filereader.ReadLines(inputFile)
	if err != nil {
		return retCards, err
	}

	for i := 0; i < len(lines); i++ {
		if verbose {
			fmt.Printf("Line %d: %v\n", i, lines[i])
		}

		data := strings.Split(lines[i], ":")[1]
		if verbose {
			fmt.Printf("Data: %v\n", data)
		}

		numberLists := strings.Split(data, "|")

		re := regexp.MustCompile("  ")

		//winningNumsString := numberLists[0]
		winningNumsString := re.ReplaceAllString(strings.TrimSpace(numberLists[0]), " ")
		if verbose {
			fmt.Printf("Winning Nums: %v\n", winningNumsString)
		}
		winningNumsStringList := strings.Split(winningNumsString, " ")

		numsString := re.ReplaceAllString(strings.TrimSpace(numberLists[1]), " ")
		if verbose {
			fmt.Printf("Nums: %v\n", numsString)
		}
		numStringList := strings.Split(numsString, " ")

		winningNums := make([]int, len(winningNumsStringList))
		for i := 0; i < len(winningNums); i++ {
			winningNums[i], err = strconv.Atoi(winningNumsStringList[i])
			if err != nil {
				return retCards, err
			}
		}
		slices.Sort(winningNums)

		nums := make([]int, len(numStringList))
		for i := 0; i < len(nums); i++ {
			nums[i], err = strconv.Atoi(numStringList[i])
			if err != nil {
				return retCards, err
			}
		}
		slices.Sort(nums)

		matching := make([]int, 0)

		for i := 0; i < len(nums); i++ {
			_, found := slices.BinarySearch(winningNums, nums[i])
			if found {
				matching = append(matching, nums[i])
			}
		}
		slices.Sort(matching)

		newCard := card{numbers: nums, winningNumbers: winningNums, matching: matching}

		retCards = append(retCards, newCard)
	}

	return retCards, nil
}

func readCard(line string) (card, error) {
	retCard := card{numbers: make([]int, 0), winningNumbers: make([]int, 0), matching: make([]int, 0)}
	return retCard, errors.New("not implemented")
}
