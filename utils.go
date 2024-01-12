package main

import (
	"fmt"
	"os/exec"
	"runtime"
)

func open(filepath string) error {
	var (
		cmd  string
		args []string
	)

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, "file://"+filepath)

	if err := exec.Command(cmd, args...).Start(); err != nil {
		return fmt.Errorf("could not open file: %w", err)
	}

	return nil
}
