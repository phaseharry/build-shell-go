package main

import (
	"bufio"
	"fmt"
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
			fmt.Fprintln(os.Stderr, err)
			return
		}

		// remove the trailing "\n" when we read the user input in.
		input = strings.Trim(input, "\r\n")

		var tokens []string
		for {
			/*
				get the first index of '. If it exists then it is the first single
				quote and we need to find the closing single quote of it. We append
				everything in between the two single quotes. We do this until we do not have any
				single quotes left.
			*/
			startIdx := strings.Index(input, "'")
			/*
				if there is no single quotes then we can just append the rest of the input
				to our tokens
			*/
			if startIdx == -1 {
				/*
					using string.Fields to return a slice of strings with no spaces.
					the elements are delimited by spaces. This is only done for tokens before the first
					single quote. Since we want to keep any spaces between 2 single quotes.
					ex) echo 'test     example'. We don't want to filter out the spaces between the first
					single quote and the last, but we do want to filter out spaces for
					ex) echo    'test example'
				*/
				tokens = append(tokens, strings.Fields(input)...)
				break
			}
			// appending every token before the first single quote, but stripping out the spaces before the first single quote
			tokens = append(tokens, strings.Fields(input[:startIdx])...)
			// updating the existing input string to remove all tokens already appended to tokens slice
			input = input[startIdx+1:]
			/*
				getting index of the ending single quote and append all tokens before it
				to tokens
			*/
			endIdx := strings.Index(input, "'")
			tokenBetweenStartIdxAndEndIdx := input[:endIdx] // getting all values include spaces between first single quote and end single quote
			tokens = append(tokens, tokenBetweenStartIdxAndEndIdx)
			input = input[endIdx+1:]
		}

		command := tokens[0]
		commandParams := tokens[1:]
		// log.Printf("command: %v", command)
		// log.Printf("commandParams: %v", commandParams)
		// for _, val := range commandParams {
		// 	log.Println(val)
		// }
		if handled := handleCommand(command, commandParams); !handled {
			fmt.Printf("%v: command not found\n", input)
		}
	}
}

func handleCommand(command string, args []string) bool {
	if _, ok := builtInHandlers[command]; ok {
		builtInHandlers[command](args)
		return true
	} else if ok := handleExternalCommands(command, args); ok {
		return true
	}
	return false
}

func handleExternalCommands(command string, args []string) bool {
	_, pathCommandFound := fileInPathVariables(command)

	if pathCommandFound {
		cmd := exec.Command(command, args...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			fmt.Printf("error executing command '%v'. err: %v", command, err)
			return false
		}
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
