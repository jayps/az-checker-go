package azure

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type Resource struct {
	Id            string `json:"id"`
	Type          string `json:"type"`
	Name          string `json:"name"`
	ResourceGroup string `json:"resourceGroup"`
	AlertRules    []AlertRule
	BackupVault   *Resource // For VMs only, I'll separate this later.
}

func getResourceList(command string) ([]Resource, error) {
	output := RunCommand(command)
	var vms []Resource
	err := json.Unmarshal(output, &vms)
	return vms, err
}

func getResourceMap(command string, name string) map[string]Resource {
	fmt.Println(fmt.Sprintf("Fetching %s...", name))
	vms, err := getResourceList(command)

	if err != nil {
		log.Fatalln(fmt.Sprintf("Could not fetch %s: ", name), err.Error())
	}

	result := make(map[string]Resource)
	for i := 0; i < len(vms); i++ {
		result[strings.ToLower(vms[i].Id)] = vms[i]
	}

	return result
}

func FetchVMs() map[string]Resource {
	command := "az vm list"
	resourceName := "virtual machines"

	return getResourceMap(command, resourceName)
}

func FetchAKSClusters() map[string]Resource {
	command := "az aks list"
	resourceName := "AKS clusters"

	return getResourceMap(command, resourceName)
}

func FetchMySQLServers() map[string]Resource {
	command := "az mysql server list"
	resourceName := "mysql servers"

	return getResourceMap(command, resourceName)
}

func FetchFlexibleMySQLServers() map[string]Resource {
	command := "az mysql flexible-server list"
	resourceName := "flexible mysql servers"

	return getResourceMap(command, resourceName)
}

func FetchSQLServers() map[string]Resource {
	command := "az sql server list"
	resourceName := "sql servers"

	return getResourceMap(command, resourceName)
}

func FetchStorageAccounts() map[string]Resource {
	command := "az storage account list"
	resourceName := "storage accounts"

	return getResourceMap(command, resourceName)
}

func FetchWebApps() map[string]Resource {
	command := "az webapp list"
	resourceName := "web apps"

	return getResourceMap(command, resourceName)
}

func FetchResourceDetails(resourceId string) []byte {
	command := fmt.Sprintf("az resource show --ids %s", resourceId)

	return RunCommand(command)
}
