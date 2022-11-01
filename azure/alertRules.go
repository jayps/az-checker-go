package azure

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type AlertRule struct {
	Scopes []string `json:"scopes"`
	Id     string   `json:"id"`
}

func FetchAlertRules() []AlertRule {
	fmt.Println("Fetching alert rules...")
	output := RunCommand("az monitor metrics alert list")
	var alertRules []AlertRule
	err := json.Unmarshal(output, &alertRules)

	if err != nil {
		log.Fatalln("Could not load alert rules", err.Error())
	}

	return alertRules
}

func AssignAlertRulesToResources(rules []AlertRule, resources map[string]Resource) {
	for i := 0; i < len(rules); i++ {
		for scopeIndex := 0; scopeIndex < len(rules[i].Scopes); scopeIndex++ {
			scope := strings.ToLower(rules[i].Scopes[scopeIndex])
			if resource, ok := resources[scope]; ok {
				resource.AlertRules = append(resource.AlertRules, rules[i])
				resources[scope] = resource
			}
		}
	}
}
