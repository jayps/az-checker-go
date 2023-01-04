package azure

import (
	"fmt"
	"log"
)

func SetSubscription(subscriptionId string) {
	_, err := Execute(fmt.Sprintf("az account set --subscription %s", subscriptionId))
	if err != nil {
		log.Fatal(fmt.Sprintf("Could not set subscription ID: %s", err.Error()))
	}
}
