package types

type ScenarioCreateInput struct {
	Name *string `dynamodbav:"name,omitempty" json:"name"`
}

type ScenarioPatchInput struct {
	Name *string `dynamodbav:"name,omitempty" json:"name"`
}

type ScenarioItem struct {
	Item `dynamodbav:",inline"`
	Name *string `dynamodbav:"name,omitempty" json:"name"`
}
