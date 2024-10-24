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

const (
	NOT_SET = iota
	INSIDE  = iota
	OUTSIDE = iota
)

type coord struct {
	x int
	y int
}

type pipe struct {
	symbol      byte
	name        string
	directions  [4]bool
	pipes       [4]*pipe
	position    coord
	distance    int
	insideOut   int
	virtualNode bool
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
	bResult, err = b(pipeMap, verbose)
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

func b(pipeMap [][]*pipe, verbose bool) (int, error) {
	retValue := 0

	//calcOutside(pipeMap, verbose)

	column := findVerticalPosAdjust(pipeMap)
	for column != -1 {
		fmt.Printf("Found column %v\n", column)
		insertColumn(pipeMap, column)
		column = findVerticalPosAdjust(pipeMap)
	}

	// row := findRowPosAdjust(pipeMap)
	// for row != -1 {
	// 	fmt.Printf("Found row %v\n", row)
	// 	insertRow(pipeMap, row)
	// 	row = findRowPosAdjust(pipeMap)
	// }

	printMap(pipeMap)

	return retValue, fmt.Errorf("not implemented yet")
}

func insertColumn(pipeMap [][]*pipe, column int) {
	fmt.Printf("Column: %v\n", column)
	for y, pipeLine := range pipeMap {
		westPipe, err := getPipe(pipeMap, pipeLine[column].position, WEST)
		if err != nil {
			westPipe = nil
		}
		eastPipe := pipeLine[column]
		newPipe, _ := getNewPipe('.')
		if westPipe.directions[EAST] && eastPipe.directions[WEST] {
			newPipe, _ = getNewPipe('-')
			westPipe.pipes[EAST] = newPipe
			eastPipe.pipes[WEST] = newPipe
			newPipe.pipes[EAST] = eastPipe
			newPipe.pipes[WEST] = westPipe
		}
		newPipe.position.y = y
		newPipe.position.x = column
		newPipe.virtualNode = true

		pipeLine = append(pipeLine[:column+1], pipeLine[column:]...)
		pipeLine[column] = newPipe

		for i, pipeNode := range pipeLine {
			pipeNode.position.x = i
		}

		pipeMap[y] = pipeLine

	}
}

func insertRow(pipeMap [][]*pipe, row int) {
	fmt.Printf("Column: %v\n", row)

	// Needs to be written

	// for y, pipeLine := range pipeMap {
	// 	westPipe, err := getPipe(pipeMap, pipeLine[column].position, WEST)
	// 	if err != nil {
	// 		westPipe = nil
	// 	}
	// 	eastPipe := pipeLine[column]
	// 	newPipe, _ := getNewPipe('.')
	// 	if westPipe.directions[EAST] && eastPipe.directions[WEST] {
	// 		newPipe, _ = getNewPipe('-')
	// 		westPipe.pipes[EAST] = newPipe
	// 		eastPipe.pipes[WEST] = newPipe
	// 		newPipe.pipes[EAST] = eastPipe
	// 		newPipe.pipes[WEST] = westPipe
	// 	}
	// 	newPipe.position.y = y
	// 	newPipe.position.x = column
	// 	newPipe.virtualNode = true

	// 	pipeLine = append(pipeLine[:column+1], pipeLine[column:]...)
	// 	pipeLine[column] = newPipe

	// 	for i, pipeNode := range pipeLine {
	// 		pipeNode.position.x = i
	// 	}

	// 	pipeMap[y] = pipeLine

	// }
}

func findVerticalPosAdjust(pipeMap [][]*pipe) int {
	for _, pipeLine := range pipeMap {
		for _, pipeNode := range pipeLine {
			if pipeNode.symbol != '.' {
				if pipeNode.position.x == 0 && !pipeNode.directions[WEST] {
					fmt.Printf("HERE1\n")
					return pipeNode.position.x
				} else if pipeNode.position.x == len(pipeLine)-1 && !pipeNode.directions[EAST] {
					fmt.Printf("HERE2\n%v - %v\n", len(pipeLine), pipeNode)
					return pipeNode.position.x + 1
				} else {
					eastPipe, err := getPipe(pipeMap, pipeNode.position, EAST)
					if err != nil || (eastPipe.symbol != '.' && (!pipeNode.directions[EAST] || !eastPipe.directions[WEST])) {
						fmt.Printf("HERE3\n")
						return pipeNode.position.x + 1
					}
				}
			}
		}
	}

	return -1
}

func findRowPosAdjust(pipeMap [][]*pipe) int {
	for _, pipeLine := range pipeMap {
		for _, pipeNode := range pipeLine {
			if pipeNode.symbol != '.' {
				if pipeNode.position.y == 0 && !pipeNode.directions[NORTH] {
					return pipeNode.position.y
				} else if pipeNode.position.y == len(pipeMap)-1 && !pipeNode.directions[SOUTH] {
					return pipeNode.position.y + 1
				} else {
					southPipe, err := getPipe(pipeMap, pipeNode.position, SOUTH)
					if err != nil || (southPipe.symbol != '.' && (!pipeNode.directions[SOUTH] || !southPipe.directions[NORTH])) {
						return pipeNode.position.y + 1
					}
				}
			}
		}
	}

	return -1
}

func calcOutside(pipeMap [][]*pipe, verbose bool) {
	firstLine := pipeMap[0]
	for _, pipeNode := range firstLine {
		setOutside(pipeMap, pipeNode.position, verbose)
	}

	for i := 1; i < len(pipeMap)-1; i++ {
		pipeLine := pipeMap[i]
		// first value
		pipeNode := pipeLine[0]
		setOutside(pipeMap, pipeNode.position, verbose)
		// last value
		pipeNode = pipeLine[len(pipeLine)-1]
		setOutside(pipeMap, pipeNode.position, verbose)
	}

	lastLine := pipeMap[len(pipeMap)-1]
	for _, pipeNode := range lastLine {
		setOutside(pipeMap, pipeNode.position, verbose)
	}
}

func setOutside(pipeMap [][]*pipe, pos coord, verbose bool) {
	pipeNode := pipeMap[pos.y][pos.x]

	/* Nothing to do */
	if pipeNode.symbol != '.' || pipeNode.insideOut != NOT_SET {
		return
	}

	pipeNode.insideOut = OUTSIDE

	if verbose {
		printMap(pipeMap)
	}

	nextCoord, err := getCoordAdv(pipeMap, pos, NORTH)
	if err == nil {
		setOutside(pipeMap, nextCoord, verbose)
	}
	nextCoord, err = getCoordAdv(pipeMap, pos, EAST)
	if err == nil {
		setOutside(pipeMap, nextCoord, verbose)
	}
	nextCoord, err = getCoordAdv(pipeMap, pos, SOUTH)
	if err == nil {
		setOutside(pipeMap, nextCoord, verbose)
	}
	nextCoord, err = getCoordAdv(pipeMap, pos, WEST)
	if err == nil {
		setOutside(pipeMap, nextCoord, verbose)
	}
}

func getCoordAdv(pipeMap [][]*pipe, pos coord, direction int) (coord, error) {
	maxY := len(pipeMap)
	maxX := len(pipeMap[0])

	loop := true
	retCoord := pos
	var err error
	for loop {
		retCoord, err = getCoord(retCoord, direction, maxX, maxY)
		if err != nil {
			return retCoord, err
		}
		loop = false
	}

	return retCoord, nil
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
		distance:    -1,
		insideOut:   NOT_SET,
		virtualNode: false,
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
			//if !pipeNode.virtualNode {
			if pipeNode.symbol == '.' {
				switch pipeNode.insideOut {
				case INSIDE:
					fmt.Printf("%c", 'I')
				case OUTSIDE:
					fmt.Printf("%c", 'O')
				default:
					fmt.Printf("%c", '?')
				}
			} else {
				fmt.Printf("%c", pipeNode.symbol)
			}
			//fmt.Printf("(%d,%d) ", pipeNode.position.x, pipeNode.position.y)
			// if pipeNode.distance == -1 {
			// 	fmt.Printf(". ")
			// } else {
			// 	fmt.Printf("%v ", pipeNode.distance)
			// }
			//}
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
