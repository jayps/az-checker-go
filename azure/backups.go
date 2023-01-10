package azure

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
)

type VMBackupResult struct {
	Vault            *Resource
	VirtualMachineId string
	err              error
}

func FetchBackupsForVM(vmId string, backups chan<- VMBackupResult, wg *sync.WaitGroup) {
	fmt.Println(fmt.Sprintf("Checking backups for VM %s...", vmId))
	output, err := RunCommand(fmt.Sprintf("az backup protection check-vm --vm %s", vmId))
	backupVaultId := string(output)

	if err != nil {
		fmt.Println(err.Error())
		backups <- VMBackupResult{nil, vmId, err}
		wg.Done()
		return
	}

	if backupVaultId != "" {
		output, err := FetchResourceDetails(backupVaultId)

		if err != nil {
			backups <- VMBackupResult{nil, vmId, err}
		}

		var vault Resource
		err = json.Unmarshal(output, &vault)
		if err == nil {
			fmt.Println(fmt.Sprintf("Found backup vault for %s.", vmId))
		}
		backups <- VMBackupResult{&vault, vmId, err}
	} else {
		fmt.Println(fmt.Sprintf("No backup vault for %s.", vmId))
		backups <- VMBackupResult{nil, vmId, errors.New(fmt.Sprintf("No backup vault provided for VM %s", vmId))}
	}
	wg.Done()
}

func FetchVMBackups(vms map[string]Resource) error {
	backups := make(chan VMBackupResult, len(vms))
	var wg sync.WaitGroup

	for id := range vms {
		wg.Add(1)
		go FetchBackupsForVM(id, backups, &wg)
	}

	go func() {
		defer close(backups)
		wg.Wait()
	}()

	for b := range backups {
		if b.err == nil {
			vm := vms[b.VirtualMachineId]
			vm.BackupVault = b.Vault
			vms[b.VirtualMachineId] = vm
		}
	}

	return nil
}
