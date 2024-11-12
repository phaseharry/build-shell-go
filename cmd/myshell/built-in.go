package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

const (
	CD   = "cd"
	ECHO = "echo"
	EXIT = "exit"
	PWD  = "pwd"
	TYPE = "type"
)

var builtInHandlers map[string]func(string)

// init gets called when the program is started
func init() {
	builtInHandlers = map[string]func(string){
		CD:   cdCommand,
		ECHO: echoCommand,
		EXIT: exitCommand,
		PWD:  pwdCommand,
		TYPE: typeCommand,
	}
}

func exitCommand(args string) {
	exitCode, _ := strconv.Atoi(args)
	os.Exit(exitCode)
}

func echoCommand(args string) {
	fmt.Println(args)
}

func pwdCommand(_ string) {
	workingDir, _ := os.Getwd()
	fmt.Println(workingDir)
}

func cdCommand(args string) {
	targetPath := args
	nextDir := ""
	currentDir, _ := os.Getwd()

	if targetPath == "" || targetPath == string(HOME) {
		nextDir = HOME_DIRECTORY
	} else if targetPath[0] == HOME { // using home directory as a base to change directory
		nextDir = filepath.Join(HOME_DIRECTORY, targetPath[1:])
	} else if targetPath[0] == '/' { // absolute path, so just overwrite
		nextDir = targetPath
	} else { // relative path so will use current working dir to generate next path
		nextDir = filepath.Join(currentDir, targetPath)
	}

	if err := os.Chdir(nextDir); err != nil {
		fmt.Printf("cd: %v: No such file or directory\n", nextDir)
	}
}

func typeCommand(args string) {
	command := args

	// checking if command is within builtInHandlers map
	if _, ok := builtInHandlers[command]; ok {
		fmt.Printf("%v is a shell builtin\n", command)
	} else if path, ok := fileInPathVariables(command); ok {
		fmt.Printf("%v is %v\n", command, path)
	} else {
		fmt.Printf("%v: not found\n", command)
	}

	for _, path := range PATHS {
		fp := filepath.Join(path, command)
		_, err := os.Stat(fp)
		// if we're able to find the command from any of our paths then print and then return
		// since the command exists
		if err == nil {
			fmt.Printf("%v is %v\n", command, fp)
			return
		}
	}
}
