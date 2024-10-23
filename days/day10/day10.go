package day10

import (
	"fmt"
	"os"

	filereader "github.com/jblashki/aoc-filereader-go"
)

const name = "Day 10 "
const inputFile = "./days/day10/input"

const (
	NORTH = iota
	EAST  = iota
	SOUTH = iota
	WEST  = iota
)

type coord struct {
	x int
	y int
}

type pipe struct {
	symbol     byte
	name       string
	directions [4]bool
	pipes      [4]*pipe
	position   coord
	distance   int
}

// RunDay runs Advent of Code Day 8 Puzzle
func RunDay(verbose bool) {
	var aResult int
	var bResult int
	var err error

	if verbose {
		fmt.Printf("\n%v Output:\n", name)
	}

	pipeMap, _, err := readInput()
	if err != nil {
		fmt.Printf("%v: **** Error: %q ****\n", name, err)
		os.Exit(0)
	}

	if verbose {
		printMap(pipeMap)
	}

	aResult, err = a(pipeMap, verbose)
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

func a(pipeMap [][]*pipe, verbose bool) (int, error) {
	retValue := 0

	for _, pipeLine := range pipeMap {
		for _, pipeNode := range pipeLine {
			if pipeNode.distance > retValue {
				retValue = pipeNode.distance
			}
		}
	}

	return retValue, nil
}

func b(verbose bool) (int, error) {
	retValue := 0

	return retValue, fmt.Errorf("not implemented yet")
}

func readInput() ([][]*pipe, *pipe, error) {
	retMap := make([][]*pipe, 0)
	var retStartPipe *pipe = nil
	lines, err := filereader.ReadLines(inputFile)
	if err != nil {
		return retMap, retStartPipe, err
	}

	for i := 0; i < len(lines); i++ {
		pipeLine := make([]*pipe, 0)
		for j := 0; j < len(lines[i]); j++ {
			newPipe, err := getNewPipe(lines[i][j])
			if err != nil {
				return retMap, retStartPipe, err
			}
			if lines[i][j] == 'S' {
				retStartPipe = newPipe
			}
			newPipe.position.x = j
			newPipe.position.y = i
			pipeLine = append(pipeLine, newPipe)
		}
		retMap = append(retMap, pipeLine)
	}

	for _, pipeLine := range retMap {
		for _, pipeNode := range pipeLine {
			for d := 0; d <= WEST; d++ {
				if pipeNode.directions[d] {
					otherPipeNode, err := getPipe(retMap, pipeNode.position, d)
					if err == nil {
						oppositeD := getOppositeDirection(d)
						if otherPipeNode.directions[oppositeD] {
							pipeNode.pipes[d] = otherPipeNode
						}
					}
				}
			}
		}
	}

	// Determine start pipe value
	if retStartPipe != nil {
		directionViable := [4]bool{false, false, false, false}

		for d := 0; d <= WEST; d++ {
			pipeNode, err := getPipe(retMap, retStartPipe.position, d)
			oppositeD := getOppositeDirection(d)
			if err == nil {
				if pipeNode.directions[oppositeD] {
					directionViable[d] = true
				}
			}
		}

		if directionViable[NORTH] && directionViable[SOUTH] {
			retStartPipe.symbol = '|'
			retStartPipe.name = "VERTICAL"
			retStartPipe.directions[NORTH] = true
			retStartPipe.directions[SOUTH] = true
		} else if directionViable[EAST] && directionViable[WEST] {
			retStartPipe.symbol = '-'
			retStartPipe.name = "HORIZONTAL"
			retStartPipe.directions[EAST] = true
			retStartPipe.directions[WEST] = true
		} else if directionViable[NORTH] && directionViable[EAST] {
			retStartPipe.symbol = 'L'
			retStartPipe.name = "NORTH-EAST BEND"
			retStartPipe.directions[NORTH] = true
			retStartPipe.directions[EAST] = true
		} else if directionViable[NORTH] && directionViable[WEST] {
			retStartPipe.symbol = 'J'
			retStartPipe.name = "NORTH-WEST BEND"
			retStartPipe.directions[NORTH] = true
			retStartPipe.directions[WEST] = true
		} else if directionViable[SOUTH] && directionViable[WEST] {
			retStartPipe.symbol = '7'
			retStartPipe.name = "SOUTH-WEST BEND"
			retStartPipe.directions[SOUTH] = true
			retStartPipe.directions[WEST] = true
		} else if directionViable[SOUTH] && directionViable[EAST] {
			retStartPipe.symbol = 'F'
			retStartPipe.name = "SOUTH-EAST BEND"
			retStartPipe.directions[SOUTH] = true
			retStartPipe.directions[EAST] = true
		}

		for d := 0; d <= WEST; d++ {
			if retStartPipe.directions[d] {
				otherPipeNode, err := getPipe(retMap, retStartPipe.position, d)
				if err == nil {
					retStartPipe.pipes[d] = otherPipeNode
					oppositeD := getOppositeDirection(d)
					otherPipeNode.pipes[oppositeD] = retStartPipe
				}
			}
		}
	}

	currentNode := retStartPipe
	origDirection := getOtherDirection(currentNode, -1)
	currentDirection := origDirection
	steps := 0
	for steps == 0 || currentNode != retStartPipe {
		currentNode = currentNode.pipes[currentDirection]
		currentDirection = getOppositeDirection(currentDirection)
		currentDirection = getOtherDirection(currentNode, currentDirection)
		steps++
		if currentNode.distance == -1 || steps < currentNode.distance {
			currentNode.distance = steps
		}
	}
	currentNode = retStartPipe
	currentDirection = getOtherDirection(currentNode, origDirection)
	steps = 0
	for steps == 0 || currentNode != retStartPipe {
		currentNode = currentNode.pipes[currentDirection]
		currentDirection = getOppositeDirection(currentDirection)
		currentDirection = getOtherDirection(currentNode, currentDirection)
		steps++
		if currentNode.distance == -1 || steps < currentNode.distance {
			currentNode.distance = steps
		}
	}
	retStartPipe.distance = 0

	return retMap, retStartPipe, nil
}

func getOppositeDirection(direction int) int {
	retValue := -1
	switch direction {
	case NORTH:
		retValue = SOUTH
	case SOUTH:
		retValue = NORTH
	case EAST:
		retValue = WEST
	case WEST:
		retValue = EAST
	}
	return retValue
}

func getOtherDirection(pipeNode *pipe, currentDirection int) int {
	for i, d := range pipeNode.directions {
		if i != currentDirection && d {
			return i
		}
	}

	return -1
}

func getNewPipe(symbol byte) (*pipe, error) {
	retValue := pipe{
		symbol:     '.',
		name:       "NONE",
		pipes:      [4]*pipe{nil, nil, nil, nil},
		directions: [4]bool{false, false, false, false},
		position: coord{
			x: -1,
			y: -1,
		},
		distance: -1,
	}
	switch symbol {
	case '|':
		retValue.symbol = '|'
		retValue.name = "VERTICAL"
		retValue.directions[NORTH] = true
		retValue.directions[SOUTH] = true
	case '-':
		retValue.symbol = '-'
		retValue.name = "HORIZONTAL"
		retValue.directions[EAST] = true
		retValue.directions[WEST] = true
	case 'L':
		retValue.symbol = 'L'
		retValue.name = "NORTH-EAST BEND"
		retValue.directions[NORTH] = true
		retValue.directions[EAST] = true
	case 'J':
		retValue.symbol = 'J'
		retValue.name = "NORTH-WEST BEND"
		retValue.directions[NORTH] = true
		retValue.directions[WEST] = true
	case '7':
		retValue.symbol = '7'
		retValue.name = "SOUTH-WEST BEND"
		retValue.directions[SOUTH] = true
		retValue.directions[WEST] = true
	case 'F':
		retValue.symbol = 'F'
		retValue.name = "SOUTH-EAST BEND"
		retValue.directions[SOUTH] = true
		retValue.directions[EAST] = true
	}

	return &retValue, nil
}

func printMap(pipeMap [][]*pipe) {
	for _, line := range pipeMap {
		for _, pipeNode := range line {
			fmt.Printf("%c", pipeNode.symbol)
			//fmt.Printf("(%d,%d) ", pipeNode.position.x, pipeNode.position.y)
			// if pipeNode.distance == -1 {
			// 	fmt.Printf(". ")
			// } else {
			// 	fmt.Printf("%v ", pipeNode.distance)
			// }
		}
		fmt.Printf("\n")
	}
}

func getPipe(pipeMap [][]*pipe, currentPos coord, direction int) (*pipe, error) {
	maxY := len(pipeMap)
	maxX := len(pipeMap[0])

	newCoord, err := getCoord(currentPos, direction, maxX, maxY)
	if err != nil {
		return nil, err
	}

	retPipe := pipeMap[newCoord.y][newCoord.x]

	return retPipe, nil
}

func getCoord(pos coord, direction int, maxX int, maxY int) (coord, error) {
	retCoord := pos
	switch direction {
	case NORTH:
		if pos.y == 0 {
			return retCoord, fmt.Errorf("invalid direction")
		}
		retCoord.y--
	case SOUTH:
		if pos.y == maxY-1 {
			return retCoord, fmt.Errorf("invalid direction")
		}
		retCoord.y++
	case EAST:
		if pos.x == maxX-1 {
			return retCoord, fmt.Errorf("invalid direction")
		}
		retCoord.x++
	case WEST:
		if pos.x == 0 {
			return retCoord, fmt.Errorf("invalid direction")
		}
		retCoord.x--
	}
	return retCoord, nil
}
