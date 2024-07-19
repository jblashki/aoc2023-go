package day2

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	filereader "github.com/jblashki/aoc-filereader-go"
)

const name = "Day 2"
const inputFile = "./days/day2/input"

type count struct {
	redCount   int
	greenCount int
	blueCount  int
}

type game struct {
	number int
	rounds []count
	total  count
	max    count
	power  int
}

// RunDay runs Advent of Code Day 2 Puzzle
func RunDay(verbose bool) {
	var aResult int
	var bResult int
	var err error

	if verbose {
		fmt.Printf("\n%v Output:\n", name)
	}

	games, err := getGameArray(inputFile, verbose)

	if err != nil {
		fmt.Printf("%va: **** Error: %q ****\n", name, err)
		return
	}

	aResult, err = a(games, verbose)
	if err != nil {
		fmt.Printf("%va: **** Error: %q ****\n", name, err)
	} else {
		fmt.Printf("%va: Answer = %v\n", name, aResult)
	}
	bResult, err = b(games, verbose)
	if err != nil {
		fmt.Printf("%vb: **** Error: %q ****\n", name, err)
	} else {
		fmt.Printf("%vb: Answer = %v\n", name, bResult)
	}
}

func a(games []game, verbose bool) (int, error) {
	retValue := 0

	maxCount := count{redCount: 12, greenCount: 13, blueCount: 14}

	for i := 0; i < len(games); i++ {
		g := games[i]
		if g.max.redCount <= maxCount.redCount && g.max.greenCount <= maxCount.greenCount && g.max.blueCount <= maxCount.blueCount {
			if verbose {
				fmt.Printf("Game %d possible\n", g.number)
			}
			retValue += g.number
		}
	}

	return retValue, nil
}

func b(games []game, verbose bool) (int, error) {
	retValue := 0

	for i := 0; i < len(games); i++ {
		retValue += games[i].power
		if verbose {
			fmt.Printf("After Game %d (Power %d) Total = %d\n", games[i].number, games[i].power, retValue)
		}
	}

	return retValue, nil
}

func getGameArray(inputFile string, verbose bool) ([]game, error) {
	retValue := make([]game, 0)

	lines, err := filereader.ReadLines(inputFile)
	if err != nil {
		return retValue, err
	}

	for i := 0; i < len(lines); i++ {
		game, err := readGame(lines[i])
		if err != nil {
			return retValue, err
		}
		if verbose {
			fmt.Printf("Line: %v\n", lines[i])
			fmt.Printf("Game %d: POWER: %d --- TR: %d - TG: %d - TB: %d --- MR: %d - MG: %d - MB: %d\n", game.number, game.power, game.total.redCount, game.total.greenCount, game.total.blueCount, game.max.redCount, game.max.greenCount, game.max.blueCount)
		}
		retValue = append(retValue, game)
	}

	return retValue, nil
}

func readGame(line string) (game, error) {
	g := game{number: 0, rounds: make([]count, 0), total: count{redCount: 0, greenCount: 0, blueCount: 0}, max: count{redCount: 0, greenCount: 0, blueCount: 0}, power: 0}

	lineSplit := strings.Split(line, ":")
	if len(lineSplit) != 2 {
		return g, errors.New("unable to read game header")
	}

	num, err := getGameNum(lineSplit[0])

	if err != nil {
		return g, err
	}

	rounds := strings.Split(lineSplit[1], ";")

	for i := 0; i < len(rounds); i++ {
		r, err := readRound(rounds[i])

		g.rounds = append(g.rounds, r)

		if r.redCount > g.max.redCount {
			g.max.redCount = r.redCount
		}

		if r.greenCount > g.max.greenCount {
			g.max.greenCount = r.greenCount
		}

		if r.blueCount > g.max.blueCount {
			g.max.blueCount = r.blueCount
		}

		g.total.redCount += r.redCount
		g.total.greenCount += r.greenCount
		g.total.blueCount += r.blueCount

		if err != nil {
			return g, errors.New("unable to read round")
		}
	}

	g.power = g.max.redCount * g.max.greenCount * g.max.blueCount
	g.number = num

	return g, nil
}

func getGameNum(header string) (int, error) {
	headerSplit := strings.Split(header, " ")

	if len(headerSplit) != 2 || headerSplit[0] != "Game" {
		return 0, errors.New("unable to read game header")
	}

	numString := headerSplit[1]

	num, err := strconv.Atoi(strings.TrimSpace(numString))

	if err != nil {
		return 0, err
	}

	return num, nil
}

func readRound(roundString string) (count, error) {
	retCount := count{redCount: 0, greenCount: 0, blueCount: 0}

	roundSplit := strings.Split(roundString, ",")

	for i := 0; i < len(roundSplit); i++ {
		entrySplit := strings.Split(strings.TrimSpace(roundSplit[i]), " ")
		if len(entrySplit) != 2 {
			return retCount, errors.New("unable to read round")
		}
		countString := entrySplit[0]
		color := strings.ToUpper(entrySplit[1])
		count, err := strconv.Atoi(countString)

		if err != nil {
			return retCount, err
		}

		switch {
		case color == "RED":
			retCount.redCount += count
		case color == "GREEN":
			retCount.greenCount += count
		case color == "BLUE":
			retCount.blueCount += count
		default:
			return retCount, errors.New("unknown color")
		}
	}

	return retCount, nil
}
