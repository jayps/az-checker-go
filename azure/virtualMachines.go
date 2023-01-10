package azure

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"
)

type PatchAssessmentResult struct {
	AssessmentActivityId string `json:"assessmentActivityId"`
	AvailablePatches     []struct {
		ActivityId           string    `json:"activityId" default:"n/a"`
		AssessmentState      string    `json:"assessmentState" default:"n/a"`
		Classifications      []string  `json:"classifications" default:"n/a"`
		KbId                 string    `json:"kbId" default:"n/a"`
		LastModifiedDateTime time.Time `json:"lastModifiedDateTime" default:"n/a"`
		Name                 string    `json:"name" default:"n/a"`
		PatchId              string    `json:"patchId" default:"n/a"`
		PublishedDate        time.Time `json:"publishedDate" default:"n/a"`
		RebootBehavior       string    `json:"rebootBehavior" default:"n/a"`
		Version              string    `json:"version" default:"n/a"`
	} `json:"availablePatches"`
	CriticalAndSecurityPatchCount int `json:"criticalAndSecurityPatchCount"`
	Error                         struct {
		Code       string        `json:"code"`
		Details    []interface{} `json:"details"`
		Innererror interface{}   `json:"innererror"`
		Message    string        `json:"message"`
		Target     interface{}   `json:"target"`
	} `json:"error"`
	OtherPatchCount int       `json:"otherPatchCount"`
	RebootPending   bool      `json:"rebootPending"`
	StartDateTime   time.Time `json:"startDateTime"`
	Status          string    `json:"status"`
}

type PatchResult struct {
	VM  Resource
	Err error
}

func AssessPatches(vm Resource, patchResults chan<- PatchResult, wg *sync.WaitGroup) {
	// TODO: Make this nicelier. It's not great but I'm in a rush.
	complete := make(chan PatchResult, 1)
	go func() {
		fmt.Println(fmt.Sprintf("Assessing patches for VM: %s... This might take a minute. Grab some coffee.", vm.Name))
		output, err := RunCommand(fmt.Sprintf("az vm assess-patches -n %s -g %s", vm.Name, vm.ResourceGroup))

		if err != nil {
			complete <- PatchResult{vm, err}
		}

		var patchAssessmentResult PatchAssessmentResult
		err = json.Unmarshal(output, &patchAssessmentResult)
		vm.PatchAssessmentResult = patchAssessmentResult

		// Don't necessarily crash on this, just alert the user to it.
		if err != nil {
			complete <- PatchResult{vm, err}
		}

		complete <- PatchResult{vm, nil}
	}()
	select {
	case res := <-complete:
		patchResults <- res
	case <-time.After(5 * time.Minute):
		patchResults <- PatchResult{vm, errors.New(fmt.Sprintf("Timeout after 5 minutes while checking patches for VM: %s", vm.Name))}
	}

	wg.Done()
}
