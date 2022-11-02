package azure

import (
	"fmt"
	"log"
	"os/exec"
)

func SetSubscription(subscriptionId string) {
	_, err := exec.Command("powershell", fmt.Sprintf("az account set --subscription %s", subscriptionId)).Output()
	if err != nil {
		log.Fatal(fmt.Sprintf("Could not set subscription ID: %s", err.Error()))
	}
}
