package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	exit = "exit"
)

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		// Wait for user input
		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			log.Panicf("error reading in user input: %v", err)
		}

		// remove the trailing "\n" when we read the user input in.
		input = strings.Trim(input, "\n")
		// fmt.Printf("input line: %v\n", input)

		delimitedInput := strings.Split(input, " ")
		// fmt.Printf("%v\n", delimitedInput)

		command := delimitedInput[0]
		// fmt.Printf("command: %v\n", command)
		if command == exit {
			break
		}
		fmt.Printf("%v: command not found\n", input)
	}
}
