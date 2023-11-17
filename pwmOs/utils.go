package pwmOs

import (
	"os"
	"os/exec"
	"runtime"
)

func ClearTerminalBuffer() {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = exec.Command("clear")
	case "darwin":
		cmd = exec.Command("clear")
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	}

	if cmd != nil {
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	}
}
