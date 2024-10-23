package day5

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	filereader "github.com/jblashki/aoc-filereader-go"
)

const name = "Day 5"
const inputFile = "./days/day5/input"

type rangeDesc struct {
	start    int
	offset   int
	rangeLen int
}

type almanac struct {
	seeds        []int
	seedsToSoil  []rangeDesc
	soilToFert   []rangeDesc
	fertToWater  []rangeDesc
	waterToLight []rangeDesc
	lightToTemp  []rangeDesc
	tempToHumid  []rangeDesc
	humidToLoc   []rangeDesc
}

const (
	SEEDS  = iota
	HEADER = iota
	MAP    = iota
)

// RunDay runs Advent of Code Day 5 Puzzle
func RunDay(verbose bool) {
	var aResult int
	var bResult int
	var err error

	if verbose {
		fmt.Printf("\n%v Output:\n", name)
		fmt.Printf("Reading Input...\n")
	}
	al, err := readInput(verbose)
	if verbose {
		fmt.Printf("Reading Input Complete\n")
	}

	if err != nil {
		fmt.Printf("%va: **** Error: %q ****\n", name, err)
		return
	}

	aResult, err = a(verbose, al)
	if err != nil {
		fmt.Printf("%va: **** Error: %q ****\n", name, err)
	} else {
		fmt.Printf("%va: Answer = %v\n", name, aResult)
	}
	bResult, err = b(verbose, al)
	if err != nil {
		fmt.Printf("%vb: **** Error: %q ****\n", name, err)
	} else {
		fmt.Printf("%vb: Answer = %v\n", name, bResult)
	}
}

func a(verbose bool, al almanac) (int, error) {
	retValue := -1

	for i := 0; i < len(al.seeds); i++ {
		location, err := calcLocation(verbose, al.seeds[i], al)
		if err != nil {
			return retValue, err
		}

		if retValue == -1 || retValue > location {
			retValue = location
		}
	}

	return retValue, nil
}

func b(verbose bool, al almanac) (int, error) {
	retValue := -1

	for i := 0; i < len(al.seeds); i++ {
		if i+1 >= len(al.seeds) {
			return retValue, fmt.Errorf("invalid seed format. missing range number")
		}

		line := (i / 2) + 1
		seedStart := al.seeds[i]
		seedRange := al.seeds[i+1]
		i++

		for j := 0; j < seedRange; j++ {
			seed := seedStart + j

			if verbose {
				if j%1000000 == 0 {
					if j >= 1000000 {
						fmt.Printf("Line %v: %.00fM / %.02fM\n", line, (float64)(j)/1000000, (float64)(seedRange)/1000000)
					} else if j >= 1000 {
						fmt.Printf("Line %v: %.00fK / %.02fM\n", line, (float64)(j)/1000, (float64)(seedRange)/1000000)
					} else {
						fmt.Printf("Line %v: %d / %.02fM\n", line, j, (float64)(seedRange)/1000000)
					}
				}
			}

			location, err := calcLocation(verbose, seed, al)
			if err != nil {
				return retValue, err
			}

			if retValue == -1 || retValue > location {
				if verbose {
					fmt.Printf("Changing value from %v to %v\n", retValue, location)
				}
				retValue = location
			}
		}
	}

	return retValue, nil
}

func calcLocation(verbose bool, seed int, al almanac) (int, error) {
	retValue := -1

	soil := convert(seed, al.seedsToSoil)
	fert := convert(soil, al.soilToFert)
	water := convert(fert, al.fertToWater)
	light := convert(water, al.waterToLight)
	temp := convert(light, al.lightToTemp)
	humid := convert(temp, al.tempToHumid)
	location := convert(humid, al.humidToLoc)
	if verbose {
		fmt.Printf("%v->%v->%v->%v->%v->%v->%v->%v\n", seed, soil, fert, water, light, temp, humid, location)
	}

	if retValue == -1 || retValue > location {
		retValue = location
	}

	return retValue, nil
}

func readInput(verbose bool) (almanac, error) {
	var retValue almanac

	retValue.seeds = make([]int, 0)
	retValue.seedsToSoil = make([]rangeDesc, 0)
	retValue.soilToFert = make([]rangeDesc, 0)
	retValue.fertToWater = make([]rangeDesc, 0)
	retValue.waterToLight = make([]rangeDesc, 0)
	retValue.lightToTemp = make([]rangeDesc, 0)
	retValue.tempToHumid = make([]rangeDesc, 0)
	retValue.humidToLoc = make([]rangeDesc, 0)

	var currentRangeDesc *[]rangeDesc

	lines, err := filereader.ReadLines(inputFile)
	if err != nil {
		return retValue, err
	}

	state := SEEDS
	headerValue := ""
	currentRangeDesc = nil
	for i := 0; i < len(lines); i++ {
		switch state {
		case SEEDS:
			line := strings.Split(lines[i], ":")
			if len(line) != 2 {
				return retValue, fmt.Errorf("invalid seed line[%d]ne: %v", i, lines[i])
			}
			header := line[0]
			if header != "seeds" {
				return retValue, fmt.Errorf("invalid seed line[%d]: %v", i, lines[i])
			}
			data := strings.Split(line[1], " ")
			if data[0] == "" {
				data = data[1:]
			}
			for i := 0; i < len(data); i++ {
				value, err := strconv.Atoi(data[i])
				if err != nil {
					return retValue, err
				}
				retValue.seeds = append(retValue.seeds, value)
			}
			state = HEADER
		case HEADER:
			if lines[i] != "" {
				headerValue = strings.Split(lines[i], ":")[0]
				switch headerValue {
				case "seed-to-soil map":
					currentRangeDesc = &(retValue.seedsToSoil)
					state = MAP
				case "soil-to-fertilizer map":
					currentRangeDesc = &(retValue.soilToFert)
					state = MAP
				case "fertilizer-to-water map":
					currentRangeDesc = &(retValue.fertToWater)
					state = MAP
				case "water-to-light map":
					currentRangeDesc = &(retValue.waterToLight)
					state = MAP
				case "light-to-temperature map":
					currentRangeDesc = &(retValue.lightToTemp)
					state = MAP
				case "temperature-to-humidity map":
					currentRangeDesc = &(retValue.tempToHumid)
					state = MAP
				case "humidity-to-location map":
					currentRangeDesc = &(retValue.humidToLoc)
					state = MAP
				}

			}
		case MAP:
			if lines[i] == "" {
				// Sort slice here
				if currentRangeDesc != nil {
					sort.Slice((*currentRangeDesc), func(i, j int) bool {
						return (*currentRangeDesc)[i].start < (*currentRangeDesc)[j].start
					})
				}
				state = HEADER
			} else {
				data := strings.Split(lines[i], " ")
				if len(data) != 3 {
					return retValue, fmt.Errorf("invalid map(%v) line[%d]: %v", headerValue, i, lines[i])
				}

				destStart, err := strconv.Atoi(data[0])
				if err != nil {
					return retValue, err
				}
				sourceStart, err := strconv.Atoi(data[1])
				if err != nil {
					return retValue, err
				}
				rangeLen, err := strconv.Atoi(data[2])
				if err != nil {
					return retValue, err
				}

				if verbose {
					fmt.Printf("Map \"%v\" Line: %d %d %d\n", headerValue, destStart, sourceStart, rangeLen)
				}

				newRange := rangeDesc{
					start:    sourceStart,
					offset:   destStart - sourceStart,
					rangeLen: rangeLen,
				}

				(*currentRangeDesc) = append((*currentRangeDesc), newRange)
			}
		default:
			fmt.Printf("UNKNOWN\n")
		}
	}

	return retValue, nil
}

func convert(source int, rng []rangeDesc) int {
	for i := 0; i < len(rng); i++ {
		if source >= rng[i].start && source < rng[i].start+rng[i].rangeLen {
			return source + rng[i].offset
		}
	}

	return source
}
