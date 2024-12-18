package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	HOME = '~'
)

var PATHS = strings.Split(os.Getenv("PATH"), ":")
var HOME_DIRECTORY = os.Getenv("HOME")

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
		delimitedInput := strings.Split(input, " ")
		command := delimitedInput[0]

		var tokens []string
		for i := 1; i < len(delimitedInput); i++ {
			normalizedToken := normalizeTokens(delimitedInput[i])
			/*
				Only add token if it is not a empty so we don't add extra spaces when we join
				the values together to format one string after we normalize by string quotes
			*/
			if len(normalizedToken) > 0 {
				tokens = append(tokens, normalizedToken)
			}
		}
		commandParams := strings.Join(tokens, " ")

		if handled := handleCommand(command, commandParams); !handled {
			fmt.Printf("%v: command not found\n", input)
		}
	}
}

func normalizeTokens(token string) string {
	// hack fix right now, just to remove all single quotes and double quotes
	token = strings.Trim(token, "'")
	token = strings.Trim(token, "\"")
	return token
}

func handleCommand(command, args string) bool {
	if _, ok := builtInHandlers[command]; ok {
		builtInHandlers[command](args)
		return true
	} else if ok := handleExternalCommands(command, args); ok {
		return true
	}
	return false
}

func handleExternalCommands(command, args string) bool {
	_, pathCommandFound := fileInPathVariables(command)

	if pathCommandFound {
		cmd := exec.Command(command, args)
		stdout, _ := cmd.Output()
		fmt.Print(string(stdout))
		return true
	}

	return false
}

func fileInPathVariables(command string) (string, bool) {
	for _, path := range PATHS {
		fp := filepath.Join(path, command)
		_, err := os.Stat(fp)
		// if we're able to find the command from any of our paths return the path
		if err == nil {
			return fp, true
		}
	}
	return "", false
}
