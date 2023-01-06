package azure

import (
	"encoding/json"
	"fmt"
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

func AssessPatches(vm *Resource) error {
	fmt.Println(fmt.Sprintf("Assessing patches for VM: %s... This might take a minute. Grab some coffee.", vm.Name))
	output, err := RunCommand(fmt.Sprintf("az vm assess-patches -n %s -g %s", vm.Name, vm.ResourceGroup))

	if err != nil {
		return err
	}

	var patchAssessmentResult PatchAssessmentResult
	err = json.Unmarshal(output, &patchAssessmentResult)
	vm.PatchAssessmentResult = patchAssessmentResult

	// Don't necessarily crash on this, just alert the user to it.
	if err != nil {
		return err
	}

	return nil
}
