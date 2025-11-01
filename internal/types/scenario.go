package types

type ScenarioCreateInput struct {
	Name string `dynamodbav:"name"`
}

type ScenarioPatchInput struct {
	Name string `dynamodbav:"name"`
}

type ScenarioItem struct {
	Item `dynamodbav:",inline"`
	Name string `dynamodbav:"name"`
}
