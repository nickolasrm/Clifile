package util

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// TryThrow checks if an error is not nil and then prints it to stderr and
// finished the program execution
func TryThrow(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Shell executes a script depending on what os you are
func Shell(script string) error {
	var shellCmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		shellCmd = exec.Command("powershell", "-Command", script)
	case "linux", "darwin":
		shellCmd = exec.Command(os.Getenv("SHELL"), "-c", script)
	default:
		return fmt.Errorf("unsupported operating system '%s'", runtime.GOOS)
	}
	shellCmd.Stdin = os.Stdin
	shellCmd.Stdout = os.Stdout
	shellCmd.Stderr = os.Stderr
	return shellCmd.Run()
}
