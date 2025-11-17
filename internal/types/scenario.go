package types

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
