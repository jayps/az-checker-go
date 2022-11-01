package azure

import (
	"fmt"
	"os/exec"
)

func SetSubscription(subscriptionId string) {
	exec.Command(fmt.Sprintf("az account set --subscription %s", subscriptionId))
}
