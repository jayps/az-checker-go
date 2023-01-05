package azure

import (
	"encoding/json"
	"fmt"
	"time"
)

type PatchAssessmentResult struct {
	AssessmentActivityId string `json:"assessmentActivityId"`
	AvailablePatches     []struct {
		ActivityId           string      `json:"activityId"`
		AssessmentState      string      `json:"assessmentState"`
		Classifications      []string    `json:"classifications"`
		KbId                 interface{} `json:"kbId"`
		LastModifiedDateTime time.Time   `json:"lastModifiedDateTime"`
		Name                 string      `json:"name"`
		PatchId              string      `json:"patchId"`
		PublishedDate        time.Time   `json:"publishedDate"`
		RebootBehavior       string      `json:"rebootBehavior"`
		Version              string      `json:"version"`
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

func AssessPatches(vm *Resource) {
	fmt.Println(fmt.Sprintf("Assessing patches for VM: %s...", vm.Name))
	output := RunCommand(fmt.Sprintf("az vm assess-patches -n %s -g %s", vm.Name, vm.ResourceGroup))
	var patchAssessmentResult PatchAssessmentResult
	err := json.Unmarshal(output, &patchAssessmentResult)
	vm.PatchAssessmentResult = patchAssessmentResult

	if err != nil {
		fmt.Println("Could not load patch assessment", err.Error())
	}
}
