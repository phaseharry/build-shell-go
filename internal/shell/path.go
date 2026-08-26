package shell

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var errHomeNotSet = errors.New("HOME not set")

func lookupExecutable(command string) (string, bool) {
	path, err := exec.LookPath(command)
	return path, err == nil
}

func expandHomePath(path string) (string, error) {
	if path != "~" && !strings.HasPrefix(path, "~/") {
		return path, nil
	}

	home := os.Getenv("HOME")
	if home == "" {
		return "", errHomeNotSet
	}
	if path == "~" {
		return home, nil
	}

	return filepath.Join(home, strings.TrimPrefix(path, "~/")), nil
}
