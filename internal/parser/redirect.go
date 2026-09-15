package parser

import (
	"strconv"
	"strings"
)

type Redirect struct {
	FileDescriptor int    // 1 for stdout, 2 for stderr
	Target         string // filename
	Append         bool   // >> is append to file while > is overwrite
}

func getFileDescriptor(operator string) int {
	digits := strings.TrimRight(operator, ">")

	// if there is no digits (> and >>) has no explicit descriptors and means stdout
	if digits == "" {
		return 1
	}

	// if there was an error, just use stdout
	fileDescriptor, err := strconv.Atoi(digits)
	if err != nil {
		return 1
	}

	return fileDescriptor
}
