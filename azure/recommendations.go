package azure

import (
	"encoding/json"
	"fmt"
)

type ShortDescription struct {
	Problem string `json:"problem"`
}

type AdvisorRecommendation struct {
	Description      ShortDescription `json:"shortDescription"`
	Impact           string           `json:"impact"`
	ResourceType     string           `json:"impactedField"`
	AffectedResource string           `json:"impactedValue"`
	ResourceGroup    string           `json:"resourceGroup"`
	Category         string           `json:"category"`
}

func FetchAdvisorRecommendations() (map[string][]AdvisorRecommendation, error) {
	fmt.Println("Fetching advisor recommendations...")
	output, err := RunCommand("az advisor recommendation list")
	if err != nil {
		return nil, err
	}
	var recommendations []AdvisorRecommendation
	err = json.Unmarshal(output, &recommendations)

	if err != nil {
		return nil, err
	}

	result := make(map[string][]AdvisorRecommendation)

	for i := 0; i < len(recommendations); i++ {
		categoryName := recommendations[i].Category
		if _, ok := result[categoryName]; ok {
			result[categoryName] = append(result[categoryName], recommendations[i])
		} else {
			result[categoryName] = []AdvisorRecommendation{recommendations[i]}
		}
	}

	return result, nil
}
