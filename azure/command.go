package azure

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
)

func Execute(command string) ([]byte, error) {
	if runtime.GOOS == "windows" {
		return exec.Command("powershell", command).Output()
	}

	return exec.Command("bash", "-c", command).Output()
}

func RunCommand(command string) []byte {
	output, err := Execute(command)
	if err != nil {
		log.Fatal(fmt.Sprintf("Command '%s' failed with error %s", command, err.Error()))
	}

	return output
}
