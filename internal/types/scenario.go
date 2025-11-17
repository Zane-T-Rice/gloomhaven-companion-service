package types

import "gloomhaven-companion-service/internal/constants"

type ScenarioCreateInput struct {
	Name          *string `dynamodbav:"name,omitempty" json:"name"`
	Groups        *string `dynamodbav:"groups,omitempty" json:"groups"`
	ScenarioLevel *int    `dynamodbav:"scenario_level,omitempty" json:"scenarioLevel"`
}

type ScenarioPatchInput struct {
	Name          *string `dynamodbav:"name,omitempty" json:"name"`
	Groups        *string `dynamodbav:"groups,omitempty" json:"groups"`
	ScenarioLevel *int    `dynamodbav:"scenario_level,omitempty" json:"scenarioLevel"`
}

type ScenarioItem struct {
	Item          `dynamodbav:",inline"`
	Name          *string `dynamodbav:"name,omitempty" json:"name"`
	Groups        *string `dynamodbav:"groups,omitempty" json:"groups"`
	ScenarioLevel *int    `dynamodbav:"scenario_level,omitempty" json:"scenarioLevel"`
}

func NewScenarioItem(input ScenarioCreateInput, campaignId string, scenarioId string) ScenarioItem {
	return ScenarioItem{
		Item: Item{
			Parent: constants.CAMPAIGN + constants.SEPERATOR + campaignId,
			Entity: constants.SCENARIO + constants.SEPERATOR + scenarioId,
		},
		Name:          input.Name,
		Groups:        input.Groups,
		ScenarioLevel: input.ScenarioLevel,
	}
}
