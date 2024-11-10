package day16

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
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

const (
	LIGHT       = iota
	MIRR        = iota
	LIGHT_TRAIL = iota
	OTHER       = iota
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

func Run(screen tcell.Screen, defStyle tcell.Style, verbose bool) {
	var aResult int
	var bResult int
	var erra error
	var errb error

	lmap, err := readMap(verbose)
	if err != nil {
		fmt.Printf("%v: **** Error: %q ****\n", name, err)
		os.Exit(1)
	}

	aResult, erra = a(lmap, screen, defStyle, verbose)
	bResult, errb = b(verbose)
	screen.Clear()
	if erra != nil {
		output := fmt.Sprintf("%va: **** Error: %q ****", name, erra)
		screen.SetContent(0, 0, ' ', []rune(output), defStyle)
	} else {
		output := fmt.Sprintf("%va: Answer = %v", name, aResult)
		screen.SetContent(0, 0, ' ', []rune(output), defStyle)
	}
	if errb != nil {
		output := fmt.Sprintf("%vb: **** Error: %q ****", name, errb)
		screen.SetContent(0, 1, ' ', []rune(output), defStyle)
	} else {
		output := fmt.Sprintf("%vb: Answer = %v", name, bResult)
		screen.SetContent(0, 1, ' ', []rune(output), defStyle)
	}
	screen.SetContent(0, 2, ' ', []rune("[Press ESC to Continue]"), defStyle)
	screen.Show()
}

// RunDay runs Advent of Code Day 16 Puzzle
func RunDay(verbose bool) {

	var err error

	if verbose {
		fmt.Printf("\n%v Output:\n", name)
	}

	screen, err := tcell.NewScreen()

	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	defStyle := tcell.StyleDefault.Background(tcell.ColorBlack).
		Foreground(tcell.ColorWhite)
	screen.SetStyle(defStyle)

	go Run(screen, defStyle, verbose)

	for {

		switch event := screen.PollEvent().(type) {
		case *tcell.EventResize:
			screen.Sync()
		case *tcell.EventKey:
			if event.Key() == tcell.KeyEscape || event.Key() == tcell.KeyCtrlC {
				screen.Fini()
				os.Exit(0)
			}
		}
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
		node := lmap.m[lightBeam.coord.row][lightBeam.coord.col]

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

func a(lmap lightMap, screen tcell.Screen, defStyle tcell.Style, verbose bool) (int, error) {
	retValue := 0

	lmap.lightBeams = append(lmap.lightBeams, lightBeamStuct{coord: mapCoord{row: 0, col: 0}, direction: EAST})

	for moveLight(&lmap) {
		if verbose {
			screen.Clear()
			printMap(lmap, screen, defStyle, LIGHT_MAP)
			time.Sleep(1000)
		}
	}
	if verbose {
		printMap(lmap, screen, defStyle, LIGHT_MAP)
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

func getLightChar(node *mapNodeStruct, lightBeams []lightBeamStuct) (rune, int) {
	retValue := '.'
	retType := OTHER

	lightBeamID := getLightBeamId(node, lightBeams)
	if lightBeamID >= 0 {
		// if lightBeamID > 9 {
		// 	retValue = (rune)('@') + (rune)(lightBeamID-9)
		// } else {
		// 	retValue = (rune)('0') + (rune)(lightBeamID)
		// }
		retValue = '*'
		retType = LIGHT
	} else if node.piece == EMPTY {
		if lightBeamID == -2 {
			retValue = '*'
			retType = LIGHT
		} else if lightBeamID == -1 {
			retType = LIGHT_TRAIL
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
				retType = OTHER
			default:
				retValue = '+'
			}
		}
	} else {
		retValue = getPieceChar(node)
		if retValue == '.' {
			retType = OTHER
		} else {
			retType = MIRR
		}
	}

	return retValue, retType
}

func printMapNode(node *mapNodeStruct, lightBeams []lightBeamStuct, screen tcell.Screen, defStyle tcell.Style, mapType int) {
	lightStyle := tcell.StyleDefault.Background(tcell.ColorBlack).
		Foreground(tcell.ColorYellow)
	lightTrailStyle := tcell.StyleDefault.Background(tcell.ColorBlack).
		Foreground(tcell.ColorLightYellow)
	mirrortStyle := tcell.StyleDefault.Background(tcell.ColorBlack).
		Foreground(tcell.ColorBlue)
	otherStyle := tcell.StyleDefault.Background(tcell.ColorBlack).
		Foreground(tcell.ColorBlack)
	switch mapType {
	case COORD_MAP:
		fmt.Printf("(%d,%d) ", node.pos.row, node.pos.col)
	case PIECE_MAP:
		outchar := getPieceChar(node)
		screen.SetContent(node.pos.col, node.pos.row, outchar, nil, defStyle)
	case LIGHT_MAP:
		outchar, t := getLightChar(node, lightBeams)
		switch t {
		case LIGHT:
			screen.SetContent(node.pos.col, node.pos.row, outchar, nil, lightStyle)
		case LIGHT_TRAIL:
			screen.SetContent(node.pos.col, node.pos.row, outchar, nil, lightTrailStyle)
		case MIRR:
			screen.SetContent(node.pos.col, node.pos.row, outchar, nil, mirrortStyle)
		default:
			screen.SetContent(node.pos.col, node.pos.row, outchar, nil, otherStyle)
		}

	case JUST_LIGHT_MAP:
		outchar := getLightOnlyChar(node)
		screen.SetContent(node.pos.col, node.pos.row, outchar, nil, defStyle)
	}
}

func printMap(m lightMap, screen tcell.Screen, defStyle tcell.Style, mapType int) {
	//fmt.Printf("Light Map:\n")

	for _, line := range m.m {
		for _, node := range line {
			printMapNode(node, m.lightBeams, screen, defStyle, mapType)
		}
		//fmt.Println()
	}
	screen.Show()
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
