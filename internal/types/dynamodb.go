package types

type Item struct {
	Parent string `dynamodbav:"parent,omitempty" json:"parent"`
	Entity string `dynamodbav:"entity,omitempty" json:"entity"`
}
