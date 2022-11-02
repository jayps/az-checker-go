package azure

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type AllOf struct {
	MetricName      string  `json:"metricName"`
	MetricNamespace string  `json:"metricNamespace"`
	Name            string  `json:"name"`
	Operator        string  `json:"operator"`
	Threshold       float32 `json:"threshold"`
	TimeAggregation string  `json:"timeAggregation"`
}

type AlertRuleCriteria struct {
	AllOf               []AllOf `json:"allOf"`
	Enabled             bool    `json:"enabled"`
	EvaluationFrequency bool    `json:"evaluationFrequency"`
	WindowSize          bool    `json:"windowSize"`
}

type AlertRule struct {
	Scopes   []string          `json:"scopes"`
	Name     string            `json:"name"`
	Id       string            `json:"id"`
	Criteria AlertRuleCriteria `json:"criteria"`
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
