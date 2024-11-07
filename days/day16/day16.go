package day16

import (
	"fmt"
	"os"
)

const name = "Day 16"
const inputFile = "./days/day16/input"

const (
	EMPTY    = iota
	MIRROR   = iota
	SPLITTER = iota
)

const (
	NORTH = 0x01
	EAST  = 0x02
	SOUTH = 0x04
	WEST  = 0x08
)

const (
	NORTH_IDX = iota
	EAST_IDX  = iota
	SOUTH_IDX = iota
	WEST_IDX  = iota
)

type mapCoord struct {
	row int
	col int
}

type mapNodeStruct struct {
	piece      int // Piece Type
	pieceDir   int // Direction of Piece
	lightHist  int // History of light entry
	pos        mapCoord
	neighbours [](*mapNodeStruct)
}

type lightBeamStuct struct {
	pos       mapCoord
	direction int
}

type lightMap struct {
	m          [][]*mapNodeStruct
	lightBeams []lightBeamStuct
}

// RunDay runs Advent of Code Day 16 Puzzle
func RunDay(verbose bool) {
	var aResult int
	var bResult int
	var err error

	if verbose {
		fmt.Printf("\n%v Output:\n", name)
	}

	lmap, err := readMap(verbose)
	if err != nil {
		fmt.Printf("%v: **** Error: %q ****\n", name, err)
		os.Exit(1)
	}

	printMap(lmap)

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

	return retValue, fmt.Errorf("not implemented yet")
}

func b(verbose bool) (int, error) {
	retValue := 0

	return retValue, fmt.Errorf("not implemented yet")
}

func readMap(verbose bool) (lightMap, error) {
	retValue := lightMap{
		m:          make([][]*mapNodeStruct, 0),
		lightBeams: make([]lightBeamStuct, 0),
	}

	return retValue, fmt.Errorf("not implemented yet")
}

func printMap(m lightMap) {
	fmt.Printf("Light Map:\n")
}

func readNode(c byte, pos mapCoord) *mapNodeStruct {
	retNode := mapNodeStruct{
		piece:      EMPTY,
		pieceDir:   0,
		lightHist:  0,
		pos:        pos,
		neighbours: make([]*mapNodeStruct, 4),
	}

	switch c {
	case '/':
		retNode.piece = MIRROR
		retNode.pieceDir = SOUTH + WEST
	case '\\':
		retNode.piece = MIRROR
		retNode.pieceDir = NORTH + WEST
	case '-':
		retNode.piece = SPLITTER
		retNode.pieceDir = EAST + WEST
	case '|':
		retNode.piece = SPLITTER
		retNode.pieceDir = NORTH + SOUTH
	case '.':
		fallthrough
	default:
		retNode.piece = EMPTY
	}

	return &retNode
}
