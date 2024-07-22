package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"aoc2023/days/day1"
	"aoc2023/days/day2"
	"aoc2023/days/day3"
	"aoc2023/days/day4"
)

var argCalls = map[string]func(bool){
	"DAY1": day1.RunDay,
	"DAY2": day2.RunDay,
	"DAY3": day3.RunDay,
	"DAY4": day4.RunDay,
}

var functionPointers = []func(bool){
	day1.RunDay,
	day2.RunDay,
	day3.RunDay,
	day4.RunDay,
}

var programName = ""

func usage() {
	fmt.Printf("\nUsage: %v [-v] [Day?]...\n", path.Base(os.Args[0]))
	flag.PrintDefaults()
}

func main() {
	programName = path.Base(os.Args[0])
	var verboseFlag = flag.Bool("v", false, "verbose mode")

	flag.Usage = func() {
		usage()
	}

	flag.Parse()

	args := flag.Args()

	if len(args) > 0 {
		// Validate args first
		idx := -1
		var err error = nil
		for i := 0; i < len(args); i++ {
			if len(args[i]) < len("day") {
				idx, err = strconv.Atoi(args[i])
				if err != nil {
					fmt.Printf("Invalid argument %q\n", args[i])
					usage()
					return
				}
				idx--
			} else {
				if strings.ToUpper(args[i][:3]) != "DAY" {
					idx, err = strconv.Atoi(args[i])
					if err != nil {
						fmt.Printf("Invalid argument %q\n", args[i])
						usage()
						return
					}
					idx--
				} else {
					idx, err = strconv.Atoi(args[i][3:])
					if err != nil {
						fmt.Printf("Invalid argument %q\n", args[i])
						usage()
						return
					}
					idx--
				}
			}

			if idx < 0 || idx >= len(functionPointers) {
				fmt.Printf("Invalid argument %q\n", args[i])
				usage()
				return
			}
		}
	}

	fmt.Println("Advent of Code (Go) 2023")
	fmt.Println("========================")

	if len(args) > 0 {
		for i := 0; i < len(args); i++ {
			idx := -1
			if len(args[i]) < len("day") {
				idx, _ = strconv.Atoi(args[i])
				idx--
			} else {
				if strings.ToUpper(args[i][:3]) != "DAY" {
					idx, _ = strconv.Atoi(args[i])
					idx--
				} else {
					idx, _ = strconv.Atoi(args[i][3:])
					idx--
				}
			}
			call := functionPointers[idx]

			call(*verboseFlag)
		}
	} else {
		// Call all options
		for i := 0; i < len(functionPointers); i++ {
			call := functionPointers[i]
			call(*verboseFlag)
		}
	}
}
