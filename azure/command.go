package azure

import (
	"os/exec"
	"runtime"
)

func Execute(command string) ([]byte, error) {
	if runtime.GOOS == "windows" {
		return exec.Command("powershell", command).Output()
	}

	return exec.Command("bash", "-c", command).Output()
}

func RunCommand(command string) ([]byte, error) {
	return Execute(command)
}
