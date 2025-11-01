package dto

type Scenario struct {
	Parent string `dynamodbav:"parent"`
	Entity string `dynamodbav:"entity"`
	Name   string `dynamodbav:"name"`
}
