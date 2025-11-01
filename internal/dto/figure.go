package dto

type Figure struct {
	Parent    string  `dynamodbav:"parent" json:"parent"`
	Entity    string  `dynamodbav:"entity" json:"entity"`
	Name      *string `dynamodbav:"name" json:"name"`
	MaximumHP *int    `dynamodbav:"maximum_hp" json:"maximumHP"`
	Damage    *int    `dynamodbav:"damage" json:"damage"`
}
