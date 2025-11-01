package types

type FigureCreateInput struct {
	Name      *string `dynamodbav:"name,omitempty" json:"name"`
	MaximumHP *int    `dynamodbav:"maximum_hp,omitempty" json:"maximumHP"`
	Damage    *int    `dynamodbav:"damage,omitempty" json:"damage"`
}

type FigurePatchInput struct {
	Name      *string `dynamodbav:"name,omitempty" json:"name"`
	MaximumHP *int    `dynamodbav:"maximum_hp,omitempty" json:"maximumHP"`
	Damage    *int    `dynamodbav:"damage,omitempty" json:"damage"`
}

type FigureItem struct {
	Item      `dynamodbav:",inline"`
	Name      *string `dynamodbav:"name,omitempty" json:"name"`
	MaximumHP *int    `dynamodbav:"maximum_hp,omitempty" json:"maximumHP"`
	Damage    *int    `dynamodbav:"damage,omitempty" json:"damage"`
}
