package azure

import (
	"fmt"
)

func SetSubscription(subscriptionId string) error {
	_, err := Execute(fmt.Sprintf("az account set --subscription %s", subscriptionId))
	if err != nil {
		return err
	}

	return nil
}
