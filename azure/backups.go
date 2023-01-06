package azure

import (
	"encoding/json"
	"fmt"
)

func FetchBackupsForVM(vmId string) (*Resource, error) {
	// The vault ID will be returned
	fmt.Println(fmt.Sprintf("Checking backups for VM %s...", vmId))
	output, err := RunCommand(fmt.Sprintf("az backup protection check-vm --vm %s", vmId))
	backupVaultId := string(output)

	if err != nil {
		return nil, err
	}

	if backupVaultId != "" {
		output, err := FetchResourceDetails(backupVaultId)

		if err != nil {
			return nil, err
		}

		var vault Resource
		err = json.Unmarshal(output, &vault)

		return &vault, err
	}

	return nil, nil
}

func FetchVMBackups(vms map[string]Resource) error {
	for id, vm := range vms {
		backupVault, err := FetchBackupsForVM(id)
		if err != nil {
			return err
		}
		vm.BackupVault = backupVault
		vms[id] = vm
	}

	return nil
}
