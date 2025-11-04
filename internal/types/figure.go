package types

import "gloomhaven-companion-service/internal/constants"

type FigureCreateInput struct {
	Name      *string `dynamodbav:"name,omitempty" json:"name"`
	MaximumHP *int    `dynamodbav:"maximum_hp,omitempty" json:"maximumHP"`
	Damage    *int    `dynamodbav:"damage,omitempty" json:"damage"`
	Class     *string `dynamodbav:"class,omitempty" json:"class"`
	Shield    *int    `dynamodbav:"shield,omitempty" json:"shield"`
	Rank      *string `dynamodbav:"rank,omitempty" json:"rank"`
	Number    *int    `dynamodbav:"number,omitempty" json:"number"`
	Move      *int    `dynamodbav:"move" json:"move"`
	Attack    *int    `dynamodbav:"attack" json:"attack"`
}

type FigurePatchInput struct {
	FigureCreateInput `dynamodbav:",inline"`
}

type FigureItem struct {
	Item              `dynamodbav:",inline"`
	FigureCreateInput `dynamodbav:",inline"`
}

func NewFigureItem(input FigureCreateInput, campaignId string, scenarioId string, figureId string) FigureItem {
	return FigureItem{Item: Item{
		Parent: constants.CAMPAIGN + constants.SEPERATOR + campaignId + constants.SCENARIO + constants.SEPERATOR + scenarioId,
		Entity: constants.FIGURE + constants.SEPERATOR + figureId,
	},
		FigureCreateInput: FigureCreateInput{
			Name:      input.Name,
			MaximumHP: input.MaximumHP,
			Damage:    input.Damage,
			Class:     input.Class,
			Shield:    input.Shield,
			Rank:      input.Rank,
			Number:    input.Number,
			Move:      input.Move,
			Attack:    input.Attack,
		}}
}
