package azure

import (
	"fmt"
	"log"
	"os/exec"
)

func RunCommand(command string) []byte {
	output, err := exec.Command("powershell", command).Output()
	if err != nil {
		log.Fatal(fmt.Sprintf("Command '%s' failed with error %s", command, err.Error()))
	}

	return output
}
