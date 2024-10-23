package day6

import (
	"fmt"
	"strconv"
	"strings"

	filereader "github.com/jblashki/aoc-filereader-go"
)

const name = "Day 6"
const inputFile = "./days/day6/input"

type race struct {
	time     int
	distance int
}

// RunDay runs Advent of Code Day 6 Puzzle
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

	races, err := readInputA(verbose)
	if err != nil {
		return retValue, err
	}

	for i := 0; i < len(races); i++ {
		winCount, err := getWinCount(races[i].time, races[i].distance)
		if err != nil {
			return retValue, err
		}

		if retValue == 0 {
			retValue = winCount
		} else {
			retValue *= winCount
		}
	}

	return retValue, nil
}

func b(verbose bool) (int, error) {
	retValue := 0

	races, err := readInputB(verbose)
	if err != nil {
		return retValue, err
	}

	for i := 0; i < len(races); i++ {
		winCount, err := getWinCount(races[i].time, races[i].distance)
		if err != nil {
			return retValue, err
		}

		if retValue == 0 {
			retValue = winCount
		} else {
			retValue *= winCount
		}
	}

	return retValue, nil
}

func readInputA(verbose bool) ([]race, error) {
	retValue := make([]race, 0)

	lines, err := filereader.ReadLines(inputFile)
	if err != nil {
		return retValue, err
	}

	if len(lines) != 2 {
		return retValue, fmt.Errorf("invalid input file. more than 2 lines")
	}

	timeLine := strings.Split(lines[0], ":")
	if timeLine[0] != "Time" {
		return retValue, fmt.Errorf("invalid input file. line 0 header doesn't match '%v'", "Time")
	}

	timeData := strings.Split(timeLine[1], " ")

	for i := 0; i < len(timeData); i++ {
		if timeData[i] != "" {

			time, err := strconv.Atoi(timeData[i])
			if err != nil {
				return retValue, err
			}

			newRace := race{
				time:     time,
				distance: 0,
			}
			retValue = append(retValue, newRace)
		}
	}

	distanceLine := strings.Split(lines[1], ":")
	if distanceLine[0] != "Distance" {
		return retValue, fmt.Errorf("invalid input file. line 1 header doesn't match '%v'", "Distance")
	}

	distanceData := strings.Split(distanceLine[1], " ")

	for addCount, i := 0, 0; i < len(distanceData); i++ {
		if distanceData[i] != "" {
			if addCount >= len(retValue) {
				return retValue, fmt.Errorf("invalid input file. more distance values than time values")
			}

			distance, err := strconv.Atoi(distanceData[i])
			if err != nil {
				return retValue, err
			}
			retValue[addCount].distance = distance
			addCount++
		}
	}

	return retValue, nil
}

func readInputB(verbose bool) ([]race, error) {
	retValue := make([]race, 0)

	lines, err := filereader.ReadLines(inputFile)
	if err != nil {
		return retValue, err
	}

	if len(lines) != 2 {
		return retValue, fmt.Errorf("invalid input file. more than 2 lines")
	}

	timeLine := strings.Split(lines[0], ":")
	if timeLine[0] != "Time" {
		return retValue, fmt.Errorf("invalid input file. line 0 header doesn't match '%v'", "Time")
	}

	timeData := strings.Split(timeLine[1], " ")

	timeString := ""
	for i := 0; i < len(timeData); i++ {
		if timeData[i] != "" {
			timeString += timeData[i]
		}
	}

	distanceLine := strings.Split(lines[1], ":")
	if distanceLine[0] != "Distance" {
		return retValue, fmt.Errorf("invalid input file. line 1 header doesn't match '%v'", "Distance")
	}

	distanceData := strings.Split(distanceLine[1], " ")

	distanceString := ""
	for i := 0; i < len(distanceData); i++ {
		if distanceData[i] != "" {
			distanceString += distanceData[i]
		}
	}

	time, err := strconv.Atoi(timeString)
	if err != nil {
		return retValue, err
	}
	distance, err := strconv.Atoi(distanceString)
	if err != nil {
		return retValue, err
	}

	newRace := race{
		time:     time,
		distance: distance,
	}
	retValue = append(retValue, newRace)

	return retValue, nil
}

func calcDistance(secs int, raceLength int) (int, error) {
	retValue := 0

	retValue = secs * (raceLength - secs)

	if retValue < 0 {
		return 0, nil
	}

	return retValue, nil
}

func calcRaces(time int) ([]race, error) {
	retValue := make([]race, 0)

	for i := 1; i < time; i++ {
		distance, err := calcDistance(i, time)
		if err != nil {
			return retValue, err
		}
		race := race{
			time:     i,
			distance: distance,
		}
		retValue = append(retValue, race)
	}

	return retValue, nil
}

func getWinCount(time int, distance int) (int, error) {
	possibleRaces, err := calcRaces(time)
	if err != nil {
		fmt.Printf("%va: **** Error: %q ****\n", name, err)
	}
	winCount := 0
	for j := 0; j < len(possibleRaces); j++ {
		if possibleRaces[j].distance > distance {
			winCount++
		}
	}

	return winCount, nil
}
