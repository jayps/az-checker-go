package azure

import (
	"encoding/json"
	"errors"
	"fmt"
)

func FetchBackupsForVM(vmId string) (*Resource, error) {
	// The vault ID will be returned
	fmt.Println(fmt.Sprintf("Checking backups for VM %s...", vmId))
	backupVaultId := string(RunCommand(fmt.Sprintf("az backup protection check-vm --vm %s", vmId)))
	if backupVaultId != "" {
		output := FetchResourceDetails(backupVaultId)
		var vault Resource
		err := json.Unmarshal(output, &vault)

		return &vault, err
	}

	return nil, errors.New("No backup fault found")
}

func FetchVMBackups(vms map[string]Resource) {
	for id, vm := range vms {
		backupVault, err := FetchBackupsForVM(id)
		if err != nil {
			vm.BackupVault = nil
		} else {
			vm.BackupVault = backupVault
		}
		vms[id] = vm
	}
}
