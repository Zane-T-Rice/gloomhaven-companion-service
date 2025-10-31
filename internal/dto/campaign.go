package dto

type Campaign struct {
	Parent string `dynamodbav:"parent"`
	Entity string `dynamodbav:"entity"`
	Name   string `dynamodbav:"name"`
}
