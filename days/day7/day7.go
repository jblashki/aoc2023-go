package day7

import (
	"cmp"
	"fmt"
	"slices"
	"sort"
	"strconv"
	"strings"

	filereader "github.com/jblashki/aoc-filereader-go"
)

const name = "Day 7"
const inputFile = "./days/day7/input"

const (
	HIGHCARD     = iota
	ONEPAIR      = iota
	TWOPAIR      = iota
	THREEOFAKIND = iota
	FULLHOUSE    = iota
	FOUROFAKIND  = iota
	FIVEOFAKIND  = iota
)

type hand struct {
	cards string
	rank  int
	bid   int
}

// RunDay runs Advent of Code Day 7 Puzzle
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
	retValue, err := run(verbose, "a")

	return retValue, err
}

func b(verbose bool) (int, error) {
	retValue, err := run(verbose, "b")

	return retValue, err
}

func run(verbose bool, part string) (int, error) {
	retValue := 0

	hands, err := readInput(verbose, part)
	if err != nil {
		return retValue, err
	}

	if verbose {
		fmt.Printf("Hand Rankings:\n")
	}

	for i := 0; i < len(hands); i++ {
		score := hands[i].bid * (i + 1)
		retValue += score
		if verbose {
			fmt.Printf("%v) %v (R:%v) (B:%v) (S:%v)\n", i+1, hands[i].cards, hands[i].rank, hands[i].bid, score)
		}
	}

	return retValue, nil
}

func readInput(verbose bool, part string) ([]hand, error) {
	retValue := make([]hand, 0)

	lines, err := filereader.ReadLines(inputFile)
	if err != nil {
		return retValue, err
	}

	for i := 0; i < len(lines); i++ {
		data := strings.Split(lines[i], " ")
		if len(data) != 2 {
			return retValue, fmt.Errorf("invalid line: %v", lines[i])
		}
		cards := data[0]
		if len(cards) != 5 {
			return retValue, fmt.Errorf("invalid line: %v", lines[i])
		}
		bid, err := strconv.Atoi(data[1])
		if err != nil {
			return retValue, fmt.Errorf("invalid line: %v", lines[i])
		}
		rank, err := getRank(verbose, cards, part)
		if err != nil {
			return retValue, fmt.Errorf("invalid line: %v", lines[i])
		}

		newHand := hand{
			cards: data[0],
			rank:  rank,
			bid:   bid,
		}
		retValue = append(retValue, newHand)
	}

	retValue, err = sortHands(retValue, part)
	if err != nil {
		return retValue, err
	}

	return retValue, nil
}

func sortHands(hands []hand, part string) ([]hand, error) {
	var retErr error = nil
	slices.SortFunc(hands, func(a, b hand) int {
		if a.rank != b.rank {
			return cmp.Compare(a.rank, b.rank)
		}

		for k := 0; k < len(a.cards); k++ {
			if a.cards[k] != b.cards[k] {
				aValue, err := getCardValue(a.cards[k], part)
				if err != nil {
					retErr = err
					return -1
				}
				bValue, err := getCardValue(b.cards[k], part)
				if err != nil {
					retErr = err
					return -1
				}

				return cmp.Compare(aValue, bValue)
			}
		}
		return 0
	})

	return hands, retErr
}

func getCardValue(card byte, part string) (int, error) {
	retValue := (int)(card - '0')
	if retValue < 2 || retValue > 9 {
		switch card {
		case 'A':
			retValue = 14
		case 'K':
			retValue = 13
		case 'Q':
			retValue = 12
		case 'J':
			if part == "b" {
				retValue = 1
			} else {
				retValue = 11
			}
		case 'T':
			retValue = 10
		default:
			return retValue, fmt.Errorf("invalid card value %v", card)
		}
	}

	return retValue, nil
}

func SortString(w string) string {
	s := strings.Split(w, "")
	sort.Strings(s)
	return strings.Join(s, "")
}

func getRank(verbose bool, cards string, part string) (int, error) {
	retValue := HIGHCARD

	if len(cards) != 5 {
		return -1, fmt.Errorf("invalid cards: %v", cards)
	}

	cardCount := make([]int, 13)
	jokerCount := 0
	for i := 0; i < len(cards); i++ {
		value, err := getCardValue(cards[i], part)
		if err != nil {
			return retValue, err
		}
		if value == 1 {
			jokerCount++
		} else {
			(cardCount[value-2])++
		}
	}

	sort.Ints(cardCount)

	if verbose {
		for i := len(cardCount) - 1; i >= 0; i-- {
			fmt.Printf("%v: %v\n", i, cardCount[i])
		}
	}

	if cardCount[len(cardCount)-1] == 5 ||
		(cardCount[len(cardCount)-1] == 4 && jokerCount >= 1) ||
		(cardCount[len(cardCount)-1] == 3 && jokerCount >= 2) ||
		(cardCount[len(cardCount)-1] == 2 && jokerCount >= 3) ||
		jokerCount >= 4 {
		retValue = FIVEOFAKIND
	} else if cardCount[len(cardCount)-1] == 4 ||
		(cardCount[len(cardCount)-1] == 3 && jokerCount >= 1) ||
		(cardCount[len(cardCount)-1] == 2 && jokerCount >= 2) ||
		(cardCount[len(cardCount)-1] == 1 && jokerCount >= 3) {
		retValue = FOUROFAKIND
	} else if (cardCount[len(cardCount)-1] == 3 && cardCount[len(cardCount)-2] == 2) ||
		(cardCount[len(cardCount)-1] == 2 && cardCount[len(cardCount)-2] == 2 && jokerCount >= 1) {
		retValue = FULLHOUSE
	} else if cardCount[len(cardCount)-1] == 3 ||
		(cardCount[len(cardCount)-1] == 2 && jokerCount >= 1) ||
		(cardCount[len(cardCount)-1] == 1 && jokerCount >= 2) {
		retValue = THREEOFAKIND
	} else if cardCount[len(cardCount)-1] == 2 && cardCount[len(cardCount)-2] == 2 {
		retValue = TWOPAIR
	} else if cardCount[len(cardCount)-1] == 2 ||
		jokerCount == 2 {
		retValue = ONEPAIR
	}

	if verbose {
		fmt.Printf("Cards: %v (Rank: %v) (%v)\n", cards, retValue, cardCount[len(cardCount)-1])
	}

	return retValue, nil
}
