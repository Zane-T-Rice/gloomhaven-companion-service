package types

type Item struct {
	Parent string `dynamodbav:"parent"`
	Entity string `dynamodbav:"entity"`
}
