package day16

import (
	"fmt"
	"os"
	"time"

	filereader "github.com/jblashki/aoc-filereader-go"
)

const name = "Day 16"
const inputFile = "./days/day16/input"

const (
	EMPTY    = iota
	MIRROR   = iota
	SPLITTER = iota
)

const (
	NORTH       = 0x01
	EAST        = 0x02
	SOUTH       = 0x04
	WEST        = 0x08
	NORTH_SOUTH = NORTH + SOUTH
	EAST_WEST   = EAST + WEST
	NORTH_EAST  = NORTH + EAST
	SOUTH_EAST  = SOUTH + EAST
	SOUTH_WEST  = SOUTH + WEST
	NORTH_WEST  = NORTH + WEST
)

const (
	NORTH_IDX = iota
	EAST_IDX  = iota
	SOUTH_IDX = iota
	WEST_IDX  = iota
)

const (
	PIECE_MAP      = iota
	COORD_MAP      = iota
	LIGHT_MAP      = iota
	JUST_LIGHT_MAP = iota
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
	coord     mapCoord
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

	aResult, err = a(lmap, verbose)
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

func getPieceDirection(piece int, pieceDir int, dir int) int {
	switch piece {
	case SPLITTER:
		if pieceDir&dir == dir {
			return dir
		} else if pieceDir&NORTH == NORTH {
			return NORTH
		} else { // EAST or WEST
			return EAST
		}
	case MIRROR:
		if pieceDir&NORTH == NORTH {
			switch dir {
			case NORTH:
				return WEST
			case SOUTH:
				return EAST
			case EAST:
				return SOUTH
			case WEST:
				return NORTH
			}
		} else {
			switch dir {
			case NORTH:
				return EAST
			case SOUTH:
				return WEST
			case EAST:
				return NORTH
			case WEST:
				return SOUTH
			}
		}
	}
	return dir
}

func getOppositeDirection(dir int) int {
	switch dir {
	case NORTH:
		return SOUTH
	case EAST:
		return WEST
	case SOUTH:
		return NORTH
	case WEST:
		return EAST
	}

	return dir
}

func getNeihbourIndex(dir int) int {
	switch dir {
	case NORTH:
		return NORTH_IDX
	case EAST:
		return EAST_IDX
	case SOUTH:
		return SOUTH_IDX
	case WEST:
		return WEST_IDX
	}

	return -1
}

func moveLight(lmap *lightMap) bool {

	if len(lmap.lightBeams) <= 0 {
		return false
	}

	removeIdxs := make([]int, 0)
	addBeam := make([]lightBeamStuct, 0)

	for idx, lightBeam := range lmap.lightBeams {
		//fmt.Printf("%d) %v\n", idx, lightBeam)
		node := lmap.m[lightBeam.coord.row][lightBeam.coord.col]
		// if node.pos.row == 14 && node.pos.col == 52 {
		// 	fmt.Printf("HELLO\n")
		// }
		node.lightHist |= lightBeam.direction
		if node.piece == SPLITTER {
			if node.pieceDir&lightBeam.direction != lightBeam.direction {
				// Split
				dir1 := getPieceDirection(node.piece, node.pieceDir, lightBeam.direction)
				dir2 := getOppositeDirection(dir1)

				neighBourIdx := getNeihbourIndex(dir1)
				nextNode := node.neighbours[neighBourIdx]
				if nextNode == nil || nextNode.lightHist&dir1 == dir1 {
					removeIdxs = append(removeIdxs, idx)
				} else {
					lightBeam.coord = nextNode.pos
					lightBeam.direction = dir1
					lmap.lightBeams[idx] = lightBeam
				}

				neighBourIdx = getNeihbourIndex(dir2)
				nextNode = node.neighbours[neighBourIdx]
				if nextNode != nil && nextNode.lightHist&dir2 != dir2 {
					newLightBeam := lightBeamStuct{coord: nextNode.pos, direction: dir2}
					addBeam = append(addBeam, newLightBeam)
				}
			} else {
				// Straight through
				neighBourIdx := getNeihbourIndex(lightBeam.direction)
				nextNode := node.neighbours[neighBourIdx]
				if nextNode == nil || nextNode.lightHist&lightBeam.direction == lightBeam.direction {
					removeIdxs = append(removeIdxs, idx)
				} else {
					lightBeam.coord = nextNode.pos
					lmap.lightBeams[idx] = lightBeam
				}
			}
		} else {
			dir := getPieceDirection(node.piece, node.pieceDir, lightBeam.direction)
			neighBourIdx := getNeihbourIndex(dir)
			nextNode := node.neighbours[neighBourIdx]
			if nextNode == nil || nextNode.lightHist&dir == dir {
				removeIdxs = append(removeIdxs, idx)
			} else {
				lightBeam.coord = nextNode.pos
				lightBeam.direction = dir
				lmap.lightBeams[idx] = lightBeam
			}
		}
		//node.lightHist |= lightBeam.direction
	}

	for i, idx := range removeIdxs {
		lmap.lightBeams = remove(lmap.lightBeams, idx-i)
	}

	lmap.lightBeams = append(lmap.lightBeams, addBeam...)

	return true
}

func remove(slice []lightBeamStuct, idx int) []lightBeamStuct {
	return append(slice[:idx], slice[idx+1:]...)
}

func a(lmap lightMap, verbose bool) (int, error) {
	retValue := 0

	lmap.lightBeams = append(lmap.lightBeams, lightBeamStuct{coord: mapCoord{row: 0, col: 0}, direction: EAST})

	for moveLight(&lmap) {
		if verbose {
			fmt.Print("\033[H\033[2J")
			printMap(lmap, LIGHT_MAP)
			fmt.Println()
			time.Sleep(10000000)
		}
	}
	if verbose {
		printMap(lmap, LIGHT_MAP)
		fmt.Println()
		printMap(lmap, JUST_LIGHT_MAP)
		fmt.Println()
	}

	retValue = countLight(lmap)

	return retValue, nil
}

func b(verbose bool) (int, error) {
	retValue := 0

	return retValue, fmt.Errorf("not implemented yet")
}

func countLight(lmap lightMap) int {
	retValue := 0

	for _, line := range lmap.m {
		for _, node := range line {
			if node.lightHist != 0 {
				retValue++
			}
		}
	}

	return retValue
}

func readMap(verbose bool) (lightMap, error) {
	retValue := lightMap{
		m:          make([][]*mapNodeStruct, 0),
		lightBeams: make([]lightBeamStuct, 0),
	}

	lines, err := filereader.ReadLines(inputFile)
	if err != nil {
		return retValue, err
	}

	for row, line := range lines {
		newNodeLine := make([]*mapNodeStruct, 0)
		for col, nodeChar := range line {
			newNode := readNode(byte(nodeChar), mapCoord{row: row, col: col})
			newNodeLine = append(newNodeLine, newNode)
		}
		retValue.m = append(retValue.m, newNodeLine)
	}

	for row, mapLine := range retValue.m {
		for col, node := range mapLine {
			node.neighbours[NORTH_IDX] = getNode(retValue.m, row-1, col)
			node.neighbours[EAST_IDX] = getNode(retValue.m, row, col+1)
			node.neighbours[SOUTH_IDX] = getNode(retValue.m, row+1, col)
			node.neighbours[WEST_IDX] = getNode(retValue.m, row, col-1)
		}
	}

	return retValue, nil
}

func getNode(lmap [][]*mapNodeStruct, row int, col int) *mapNodeStruct {
	if row < 0 || col < 0 || row >= len(lmap) || col >= len(lmap[0]) {
		return nil
	}

	return lmap[row][col]
}

func getPieceChar(node *mapNodeStruct) rune {
	retValue := '.'

	switch node.piece {
	case MIRROR:
		if node.pieceDir&NORTH == NORTH {
			retValue = '\\'
		} else {
			retValue = '/'
		}
	case SPLITTER:
		if node.pieceDir&NORTH == NORTH {
			retValue = '|'
		} else {
			retValue = '-'
		}
	case EMPTY:
		fallthrough
	default:
		retValue = '.'
	}

	return retValue
}

func getLightBeamId(node *mapNodeStruct, lightBeams []lightBeamStuct) int {
	retId := -1

	for i, light := range lightBeams {
		if light.coord == node.pos {
			if retId != -1 {
				return -2
			}
			retId = i + 1
		}
	}
	return retId
}

func getLightOnlyChar(node *mapNodeStruct) rune {
	if node.lightHist != 0 {
		return '#'
	}

	return '.'
}

func getLightChar(node *mapNodeStruct, lightBeams []lightBeamStuct) rune {
	retValue := '.'

	lightBeamID := getLightBeamId(node, lightBeams)
	if lightBeamID >= 0 {
		if lightBeamID > 9 {
			retValue = (rune)('@') + (rune)(lightBeamID-9)
		} else {
			retValue = (rune)('0') + (rune)(lightBeamID)
		}
	} else if node.piece == EMPTY {
		if lightBeamID == -2 {
			retValue = '*'
		} else if lightBeamID == -1 {
			switch node.lightHist {
			case NORTH:
				retValue = '^'
			case NORTH_SOUTH:
				retValue = 'I'
			case SOUTH:
				retValue = 'v'
			case EAST:
				retValue = '>'
			case WEST:
				retValue = '<'
			case EAST_WEST:
				retValue = '='
			case 0:
				retValue = '.'
			default:
				retValue = '+'
			}
		}
	} else {
		retValue = getPieceChar(node)
	}

	return retValue
}

func printMapNode(node *mapNodeStruct, lightBeams []lightBeamStuct, mapType int) {

	switch mapType {
	case COORD_MAP:
		fmt.Printf("(%d,%d) ", node.pos.row, node.pos.col)
	case PIECE_MAP:
		outchar := getPieceChar(node)
		fmt.Printf("%c", outchar)
	case LIGHT_MAP:
		outchar := getLightChar(node, lightBeams)
		fmt.Printf("%c", outchar)
	case JUST_LIGHT_MAP:
		outchar := getLightOnlyChar(node)
		fmt.Printf("%c", outchar)
	}
}

func printMap(m lightMap, mapType int) {
	fmt.Printf("Light Map:\n")

	for _, line := range m.m {
		for _, node := range line {
			printMapNode(node, m.lightBeams, mapType)
		}
		fmt.Println()
	}
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
