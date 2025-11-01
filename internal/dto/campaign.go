package dto

type Campaign struct {
	Parent string  `dynamodbav:"parent" json:"parent"`
	Entity string  `dynamodbav:"entity" json:"entity"`
	Name   *string `dynamodbav:"name" json:"name"`
}
