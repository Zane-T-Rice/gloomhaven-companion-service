package dto

import "gloomhaven-companion-service/internal/types"

type Scenario struct {
	Parent string  `dynamodbav:"parent" json:"parent"`
	Entity string  `dynamodbav:"entity" json:"entity"`
	Name   *string `dynamodbav:"name" json:"name"`
}

func NewScenario(item types.ScenarioItem) Scenario {
	return Scenario{
		Parent: item.Parent,
		Entity: item.Entity,
		Name:   item.Name,
	}
}
