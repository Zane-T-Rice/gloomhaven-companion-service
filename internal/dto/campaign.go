package dto

type Campaign struct {
	Id     string `dynamodbav:"id"`
	Entity string `dynamodbav:"entity"`
	name   string `dynamodbav:"name"`
}
