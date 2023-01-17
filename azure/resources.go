package azure

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Resource struct {
	Id                    string `json:"id"`
	Type                  string `json:"type"`
	Name                  string `json:"name"`
	ResourceGroup         string `json:"resourceGroup"`
	AlertRules            []AlertRule
	BackupVault           *Resource             // For VMs only, I'll separate this later.
	PatchAssessmentResult PatchAssessmentResult // For VMs only
}

func getResourceList(command string) ([]Resource, error) {
	output, err := RunCommand(command)

	if err != nil {
		return nil, err
	}

	var vms []Resource
	err = json.Unmarshal(output, &vms)

	return vms, err
}

func getResourceMap(command string, name string) (map[string]Resource, error) {
	fmt.Println(fmt.Sprintf("Fetching %s...", name))
	vms, err := getResourceList(command)

	if err != nil {
		return nil, err
	}

	result := make(map[string]Resource)
	for i := 0; i < len(vms); i++ {
		result[strings.ToLower(vms[i].Id)] = vms[i]
	}

	return result, nil
}

func FetchVMs() (map[string]Resource, error) {
	command := "az vm list -d --query \"[?powerState=='VM running']\""
	resourceName := "virtual machines"

	return getResourceMap(command, resourceName)
}

func FetchDeallocatedVMs() (map[string]Resource, error) {
	command := "az vm list -d --query \"[?powerState!='VM running']\""
	resourceName := "deallocated virtual machines"

	return getResourceMap(command, resourceName)
}

func FetchAKSClusters() (map[string]Resource, error) {
	command := "az aks list"
	resourceName := "AKS clusters"

	return getResourceMap(command, resourceName)
}

func FetchMySQLServers() (map[string]Resource, error) {
	command := "az mysql server list"
	resourceName := "mysql servers"

	return getResourceMap(command, resourceName)
}

func FetchFlexibleMySQLServers() (map[string]Resource, error) {
	command := "az mysql flexible-server list"
	resourceName := "flexible mysql servers"

	return getResourceMap(command, resourceName)
}

func FetchSQLServers() (map[string]Resource, error) {
	command := "az sql server list"
	resourceName := "sql servers"

	return getResourceMap(command, resourceName)
}

func FetchStorageAccounts() (map[string]Resource, error) {
	command := "az storage account list"
	resourceName := "storage accounts"

	return getResourceMap(command, resourceName)
}

func FetchWebApps() (map[string]Resource, error) {
	command := "az webapp list"
	resourceName := "web apps"

	return getResourceMap(command, resourceName)
}

func FetchResourceDetails(resourceId string) ([]byte, error) {
	command := fmt.Sprintf("az resource show --ids %s", resourceId)

	return RunCommand(command)
}
