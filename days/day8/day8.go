package day8

import (
	"fmt"
	"os"
	"os/signal"
	"strings"

	filereader "github.com/jblashki/aoc-filereader-go"
)

const name = "Day 8	"
const inputFile = "./days/day8/input"

const (
	LEFT  = iota //0
	RIGHT = iota //1
)

type nodeStruct struct {
	name  string
	left  *nodeStruct
	right *nodeStruct
	line  int
}

type inputStruct struct {
	directions []int
	nodeMap    map[string]*nodeStruct
	firstNode  *nodeStruct
	aNodes     []*nodeStruct
}

type primeFactor struct {
	prime    int
	exponent int
}

var input inputStruct

// RunDay runs Advent of Code Day 8 Puzzle
func RunDay(verbose bool) {
	var aResult int
	var bResult int
	var err error

	if verbose {
		fmt.Printf("\n%v Output:\n", name)
	}

	err = readInput()
	if err != nil {
		fmt.Printf("%v: **** Error: %q ****\n", name, err)
		return
	}

	if verbose {
		printInput()
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
	currentNode := input.firstNode
	steps := 0
	stepsLen := len(input.directions)
	for currentNode != nil && currentNode.name != "ZZZ" {
		if verbose {
			fmt.Printf("%v: %v --> ", steps, currentNode.name)
		}
		switch input.directions[steps%stepsLen] {
		case LEFT:
			currentNode = currentNode.left
		case RIGHT:
			currentNode = currentNode.right
		default:
			return steps, fmt.Errorf("invalid direction: %v", input.directions[steps%stepsLen])
		}
		steps++
		if verbose {
			fmt.Printf("%v\n", currentNode.name)
		}
	}

	return steps, nil
}

func b(verbose bool) (int, error) {
	steps := 0
	stepsLen := len(input.directions)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			fmt.Printf("%v: Got to %v steps\n", sig, steps)

			os.Exit(0)
		}
	}()

	foundCount := 0
	stepsCount := make([]int, len(input.aNodes))
	// Find step count for all first nodes
	for checkNodes(input.aNodes) && foundCount != len(input.aNodes) {
		x := steps % stepsLen

		direction := input.directions[x]

		for i := 0; i < len(input.aNodes); i++ {
			if verbose && i == 0 {
				fmt.Printf("%v: %v --> ", steps, input.aNodes[i].name)
			}
			switch direction {
			case LEFT:
				input.aNodes[i] = input.aNodes[i].left
			case RIGHT:
				input.aNodes[i] = input.aNodes[i].right
			default:
				return steps, fmt.Errorf("invalid direction: %v", input.directions[steps%stepsLen])
			}

			if verbose && i == 0 {
				fmt.Printf("%v\n", input.aNodes[i].name)
			}
		}
		steps++

		foundCount = 0
		for i := 0; i < len(input.aNodes); i++ {
			if input.aNodes[i].name[2] == 'Z' && stepsCount[i] == 0 {
				stepsCount[i] = steps
			}
			if stepsCount[i] != 0 {
				foundCount++
			}
		}
	}

	retValue := 0
	// Get LCM of all stepcounts
	primeFactors := make(map[int]int)
	for i := 0; i < len(stepsCount); i++ {
		primes := getPrimeFactors(stepsCount[i])
		for j := 0; j < len(primes); j++ {
			primeExp := primeFactors[primes[j].prime]
			if primes[j].exponent > primeExp {
				primeFactors[primes[j].prime] = primes[j].exponent
			}
		}
	}
	for k, v := range primeFactors {
		value := k
		for j := v - 1; j > 0; j-- {
			value *= k
		}
		if retValue == 0 {
			retValue = value
		} else {
			retValue *= value
		}
	}

	return retValue, nil
}

func isPrime(n int) bool {
	if n <= 1 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func getNextPrime(num int) int {
	i := num + 1
	for !isPrime(i) {
		i++
	}
	return i
}

func getPrimeFactors(num int) []primeFactor {
	retValue := make([]primeFactor, 0)
	currentPrime := 2
	currentIdx := 0
	for num > 1 {
		if num%currentPrime == 0 {
			num /= currentPrime
			if currentIdx >= len(retValue) {
				factor := primeFactor{
					prime:    currentPrime,
					exponent: 1,
				}
				retValue = append(retValue, factor)
			} else {
				factor := retValue[currentIdx]
				factor.exponent++
				retValue[currentIdx] = factor
			}
		} else {
			currentPrime = getNextPrime(currentPrime)
			currentIdx++
		}
	}
	return retValue
}

func checkNodes(nodes []*nodeStruct) bool {
	total := len(nodes)
	count := 0
	for i := 0; i < total; i++ {
		if nodes[i].name[2] == 'Z' {
			count++
		}
	}

	if count > 2 {
		for i := 0; i < total; i++ {
			fmt.Printf("%v) %v\n", i, nodes[i].name)
		}
		fmt.Printf("Found %v/%v\n", count, len(nodes))
	}

	return count != total
}

func readInput() error {
	lines, err := filereader.ReadLines(inputFile)
	if err != nil {
		return err
	}

	input.directions = make([]int, 0)
	input.nodeMap = make(map[string]*nodeStruct)
	input.aNodes = make([]*nodeStruct, 0)

	for i := 0; i < len(lines[0]); i++ {
		switch lines[0][i] {
		case 'R':
			input.directions = append(input.directions, RIGHT)
		case 'L':
			input.directions = append(input.directions, LEFT)
		default:
			return fmt.Errorf("invalid input: unknown value '%v' on line 0", lines[0][i])
		}
	}

	for i := 2; i < len(lines); i++ {
		lineData := strings.Split(lines[i], "=")
		if len(lineData) != 2 {
			return fmt.Errorf("invalid input: line format on line %v invalid '%v'", i, lines[i])
		}
		header := strings.TrimSpace(lineData[0])
		data := strings.TrimSpace(lineData[1])

		node, err := getNode(header)
		if err != nil {
			return err
		}

		leftRight := strings.Split(data, ",")

		left := strings.Trim(leftRight[0], " ()")
		right := strings.Trim(leftRight[1], " ()")

		leftNode, err := getNode(left)
		if err != nil {
			return err
		}
		rightNode, err := getNode(right)
		if err != nil {
			return err
		}

		node.left = leftNode
		node.right = rightNode
		node.line = i

		if node.name == "AAA" && input.firstNode == nil {
			input.firstNode = node
		}

		if node.name[2] == 'A' {
			input.aNodes = append(input.aNodes, node)
		}
	}

	return nil
}

func getNode(name string) (*nodeStruct, error) {
	returnNode, ok := input.nodeMap[name]
	if !ok {
		newNode := nodeStruct{
			name:  name,
			left:  nil,
			right: nil,
		}
		input.nodeMap[name] = &newNode
		returnNode = &newNode
	}

	return returnNode, nil
}

func printInput() {
	fmt.Printf("Input\n")
	fmt.Printf("=====\n")
	fmt.Printf("Directions: \n")
	for i := 0; i < len(input.directions); i++ {
		switch input.directions[i] {
		case LEFT:
			fmt.Printf("%v) %v\n", i+1, "Left")
		case RIGHT:
			fmt.Printf("%v) %v\n", i+1, "Right")
		default:
			fmt.Printf("%v) %v\n", i+1, "[ERROR]")
		}
	}
	fmt.Println()
	fmt.Printf("Map: \n")
	for key, value := range input.nodeMap {
		fmt.Printf("%v(%v): %v <-> %v\n", key, value.name, value.left.name, value.right.name)
	}
}

// func printInputReplica() {
// 	for i := 0; i < len(input.directions); i++ {
// 		switch input.directions[i] {
// 		case LEFT:
// 			fmt.Printf("L")
// 		case RIGHT:
// 			fmt.Printf("R")
// 		default:
// 			fmt.Printf("%v) %v\n", i+1, "[ERROR]")
// 		}
// 	}
// 	fmt.Println()
// 	fmt.Println()

// 	order := make([]*nodeStruct, 0)
// 	for _, v := range input.nodeMap {
// 		order = append(order, v)
// 	}

// 	sort.SliceStable(order, func(i, j int) bool {
// 		return order[i].line < order[j].line
// 	})

// 	for _, v := range order {
// 		fmt.Printf("%v = (%v, %v)\n", v.name, v.left.name, v.right.name)
// 	}
// }
