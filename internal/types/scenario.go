package types

type ScenarioCreateInput struct {
	Name   *string `dynamodbav:"name,omitempty" json:"name"`
	Groups *string `dynamodbav:"groups,omitempty" json:"groups"`
}

type ScenarioPatchInput struct {
	Name   *string `dynamodbav:"name,omitempty" json:"name"`
	Groups *string `dynamodbav:"groups,omitempty" json:"groups"`
}

type ScenarioItem struct {
	Item   `dynamodbav:",inline"`
	Name   *string `dynamodbav:"name,omitempty" json:"name"`
	Groups *string `dynamodbav:"groups,omitempty" json:"groups"`
}
