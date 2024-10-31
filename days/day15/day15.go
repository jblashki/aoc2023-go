package day15

import (
	"fmt"
	"os"
	"strconv"

	filereader "github.com/jblashki/aoc-filereader-go"
)

const (
	UNKNOWN = iota
	ADD     = iota
	REMOVE  = iota
)

type lensStruct struct {
	label    string
	focalLen int
}

type operation struct {
	lens      lensStruct
	operation int
}

type lensBoxNode struct {
	lens     lensStruct
	nextNode *lensBoxNode
	prevNode *lensBoxNode
}

type lensBox struct {
	id            int
	firstLensNode *lensBoxNode
	lastLensNode  *lensBoxNode
	lensPos       map[string]int
	lensCount     int
}

const name = "Day 15"
const inputFile = "./days/day15/input"

// RunDay runs Advent of Code Day 15 Puzzle
func RunDay(verbose bool) {
	var aResult int
	var bResult int
	var err error

	if verbose {
		fmt.Printf("\n%v Output:\n", name)
	}

	input, err := filereader.ReadCSVStrings(inputFile)
	if err != nil {
		fmt.Printf("%v: **** Error: %q ****\n", name, err)
		os.Exit(1)
	}

	aResult, err = a(input, verbose)
	if err != nil {
		fmt.Printf("%va: **** Error: %q ****\n", name, err)
	} else {
		fmt.Printf("%va: Answer = %v\n", name, aResult)
	}
	bResult, err = b(input, verbose)
	if err != nil {
		fmt.Printf("%vb: **** Error: %q ****\n", name, err)
	} else {
		fmt.Printf("%vb: Answer = %v\n", name, bResult)
	}
}

func a(input []string, verbose bool) (int, error) {
	retValue := 0
	for _, s := range input {
		retValue += hash(s, verbose)
	}

	return retValue, nil
}

func b(input []string, verbose bool) (int, error) {
	retValue := 0

	lensBoxes := make(map[int]lensBox)

	for _, opInput := range input {
		op, err := parseOp(opInput, verbose)
		if err != nil {
			return retValue, err
		}

		lensBoxes = applyOp(lensBoxes, op, verbose)
		if verbose {
			if op.operation == ADD {
				fmt.Printf("After \"%s=%d\"\n", op.lens.label, op.lens.focalLen)
			} else {
				fmt.Printf("After \"%s-\"\n", op.lens.label)
			}
			printBoxes(lensBoxes)
			fmt.Println()
		}
	}

	retValue = calcValue(lensBoxes)

	return retValue, nil
}

func hash(s string, verbose bool) int {
	retValue := 0
	for _, c := range s {
		asciiValue := (byte)(c)
		if verbose {
			fmt.Printf("%c) Ascii Value = %d\n", c, (int)(asciiValue))
		}
		retValue += int(asciiValue)
		if verbose {
			fmt.Printf("%c) Step 1 Value = %d\n", c, retValue)
		}
		retValue *= 17
		if verbose {
			fmt.Printf("%c) Step 2 Value = %d\n", c, retValue)
		}
		retValue %= 256
		if verbose {
			fmt.Printf("%c) Step 3 Value = %d\n", c, retValue)
		}
	}

	return retValue
}

func parseOp(input string, verbose bool) (operation, error) {
	retOp := operation{
		lens: lensStruct{
			label:    "",
			focalLen: -1,
		},
		operation: UNKNOWN,
	}
	var err error

	readingLens := false
	lensString := ""
	for _, c := range input {
		if !readingLens {
			if c == '=' {
				retOp.operation = ADD
				readingLens = true
			} else if c == '-' {
				retOp.operation = REMOVE
			} else {
				retOp.lens.label = fmt.Sprintf("%s%c", retOp.lens.label, c)
			}
		} else {
			lensString = fmt.Sprintf("%s%c", lensString, c)
		}
	}

	if len(lensString) > 0 {
		retOp.lens.focalLen, err = strconv.Atoi(lensString)
		if err != nil {
			return retOp, err
		}
	}

	return retOp, nil
}

func applyOp(lensBoxes map[int]lensBox, op operation, verbose bool) map[int]lensBox {
	hashInt := hash(op.lens.label, false)
	box, found := lensBoxes[hashInt]
	if !found {
		box = lensBox{
			id:            hashInt,
			firstLensNode: nil,
			lastLensNode:  nil,
			lensPos:       make(map[string]int),
			lensCount:     0,
		}
	}

	switch op.operation {
	case ADD:
		box = addLens(box, op.lens, verbose)
	case REMOVE:
		box = removeLens(box, op.lens.label, verbose)
		if box.lensCount <= 0 {
			delete(lensBoxes, hashInt)
		}
	}

	if box.lensCount > 0 {
		lensBoxes[hashInt] = box
	}

	return lensBoxes
}

func addLens(box lensBox, lens lensStruct, verbose bool) lensBox {
	if verbose {
		fmt.Printf("Adding Lens '%s' with focal length %d to box %d\n", lens.label, lens.focalLen, box.id)
	}

	lensPos, found := box.lensPos[lens.label]
	var node *lensBoxNode = nil
	if !found {
		newNode := lensBoxNode{
			lens:     lens,
			nextNode: nil,
			prevNode: nil,
		}
		node = &newNode
		if box.lastLensNode == nil {
			box.firstLensNode = node
			box.lastLensNode = node
		} else {
			box.lastLensNode.nextNode = node
			node.prevNode = box.lastLensNode
			box.lastLensNode = node
		}
		box.lensPos[newNode.lens.label] = box.lensCount
		box.lensCount++
	} else {
		node = getLensNode(box, lensPos)
		if node == nil {
			return box
		}
		node.lens = lens
	}

	return box
}

func removeLens(box lensBox, label string, verbose bool) lensBox {
	if verbose {
		fmt.Printf("Removing Lens '%s' from box %d\n", label, box.id)
	}

	lensPos, found := box.lensPos[label]

	if !found {
		return box
	}

	node := getLensNode(box, lensPos)

	if node == nil {
		return box
	}

	prevNode := node.prevNode
	nextNode := node.nextNode

	if prevNode == nil {
		box.firstLensNode = nextNode
	} else {
		prevNode.nextNode = nextNode
	}

	if nextNode == nil {
		box.lastLensNode = prevNode
	} else {
		nextNode.prevNode = prevNode
	}

	for k, pos := range box.lensPos {
		if pos > lensPos {
			box.lensPos[k] = pos - 1
		}
	}

	delete(box.lensPos, label)

	box.lensCount--

	return box
}

func getLensNode(box lensBox, lensPos int) *lensBoxNode {
	lensBoxNode := box.firstLensNode

	for i := 0; lensBoxNode != nil && i < lensPos; i++ {
		lensBoxNode = lensBoxNode.nextNode
	}

	return lensBoxNode
}

func printBoxes(boxes map[int]lensBox) {
	for k, box := range boxes {
		fmt.Printf("Box %d: ", k)
		for boxNode := box.firstLensNode; boxNode != nil; boxNode = boxNode.nextNode {
			fmt.Printf("[ %s %d ] ", boxNode.lens.label, boxNode.lens.focalLen)
		}
		fmt.Println()
	}
}

func calcValue(lensBoxes map[int]lensBox) int {
	retValue := 0
	for k, box := range lensBoxes {

		for i, boxNode := 0, box.firstLensNode; boxNode != nil; i, boxNode = i+1, boxNode.nextNode {
			retValue += (k + 1) * (i + 1) * boxNode.lens.focalLen
		}
	}

	return retValue
}
