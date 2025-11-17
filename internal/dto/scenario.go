package dto

import "gloomhaven-companion-service/internal/types"

type Scenario struct {
	Parent        string  `dynamodbav:"parent" json:"parent"`
	Entity        string  `dynamodbav:"entity" json:"entity"`
	Name          *string `dynamodbav:"name" json:"name"`
	Groups        *string `dynamodbav:"groups,omitempty" json:"groups"`
	ScenarioLevel *int    `dynamodbav:"scenario_level,omitempty" json:"scenarioLevel"`
	UpdatedAt     *string `dynamodbav:"updated_at" json:"updatedAt"`
}

func NewScenario(item types.ScenarioItem) Scenario {
	return Scenario{
		Parent:        item.Parent,
		Entity:        item.Entity,
		Name:          item.Name,
		Groups:        item.Groups,
		ScenarioLevel: item.ScenarioLevel,
		UpdatedAt:     item.UpdatedAt,
	}
}
