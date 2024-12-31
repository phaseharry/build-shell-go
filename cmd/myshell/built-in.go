package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	CD   = "cd"
	ECHO = "echo"
	EXIT = "exit"
	PWD  = "pwd"
	TYPE = "type"
)

var builtInHandlers map[string]func([]string)

// init gets called when the program is started
func init() {
	builtInHandlers = map[string]func([]string){
		CD:   cdCommand,
		ECHO: echoCommand,
		EXIT: exitCommand,
		PWD:  pwdCommand,
		TYPE: typeCommand,
	}
}

func exitCommand(args []string) {
	if len(args) == 0 {
		os.Exit(0)
	}
	exitCode, _ := strconv.Atoi(args[0])
	os.Exit(exitCode)
}

func echoCommand(args []string) {
	fmt.Println(strings.Join(args, " "))
}

func pwdCommand(_ []string) {
	workingDir, _ := os.Getwd()
	fmt.Println(workingDir)
}

func cdCommand(args []string) {
	var targetPath string
	if len(args) == 0 {
		targetPath = string(HOME)
	} else {
		targetPath = args[0]
	}
	nextDir := ""
	currentDir, _ := os.Getwd()

	if targetPath == string(HOME) {
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

func typeCommand(args []string) {
	command := args[0]

	// checking if command is within builtInHandlers map
	if _, ok := builtInHandlers[command]; ok {
		fmt.Printf("%v is a shell builtin\n", command)
	} else if path, ok := fileInPathVariables(command); ok {
		fmt.Printf("%v is %v\n", command, path)
	} else {
		fmt.Printf("%v: not found\n", command)
	}
}
