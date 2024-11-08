package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
)

const (
	EXIT = "exit"
	ECHO = "echo"
	TYPE = "type"
)

var BUILTIN = []string{EXIT, ECHO, TYPE}

var paths = strings.Split(os.Getenv("PATH"), ":")

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
		commandParams := delimitedInput[1:]
		// fmt.Printf("command: %v\n", command)
		if command == EXIT {
			break
		} else if command == ECHO {
			echoHandler(commandParams)
			continue
		} else if command == TYPE {
			typeHandler((commandParams))
			continue
		}

		pathCommandFound := false
		for _, path := range paths {
			fp := filepath.Join(path, command)
			_, err := os.Stat(fp)
			if err == nil {
				pathCommandFound = true
				break
			}
		}

		if pathCommandFound {
			cmd := exec.Command(command, commandParams...)
			stdout, err := cmd.Output()
			if err == nil {
				fmt.Printf("%v\n", string(stdout))
				continue
			}
			fmt.Printf("%v\n", err.Error())
		}

		fmt.Printf("%v: command not found\n", input)
	}
}

func echoHandler(commandParams []string) {
	toEcho := strings.Join(commandParams, " ")
	fmt.Printf("%v\n", toEcho)
}

func typeHandler(commandParams []string) {
	command := strings.Join(commandParams, " ")
	if slices.Contains(BUILTIN, command) {
		fmt.Printf("%v is a shell builtin\n", command)
		return
	}

	for _, path := range paths {
		fp := filepath.Join(path, command)
		_, err := os.Stat(fp)
		// if we're able to find the command from any of our paths then print and then return
		// since the command exists
		if err == nil {
			fmt.Printf("%v is %v\n", command, fp)
			return
		}
	}
	fmt.Printf("%v: not found\n", command)
}
