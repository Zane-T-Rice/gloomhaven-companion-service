package types

import "gloomhaven-companion-service/internal/constants"

type FigureCreateInput struct {
	Name           *string `dynamodbav:"name" json:"name"`
	MaximumHP      *int    `dynamodbav:"maximum_hp" json:"maximumHP"`
	Damage         *int    `dynamodbav:"damage" json:"damage"`
	Class          *string `dynamodbav:"class" json:"class"`
	Rank           *string `dynamodbav:"rank" json:"rank"`
	Shield         *int    `dynamodbav:"shield" json:"shield"`
	Number         *int    `dynamodbav:"number" json:"number"`
	Move           *int    `dynamodbav:"move" json:"move"`
	Attack         *int    `dynamodbav:"attack" json:"attack"`
	XP             *int    `dynamodbav:"xp" json:"xp"`
	InnateDefenses *string `dynamodbav:"innate_defenses" json:"innateDefenses"`
	InnateOffenses *string `dynamodbav:"innate_offenses" json:"innateOffenses"`
	Statuses       *string `dynamodbav:"statuses" json:"statuses"`
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
			Name:           input.Name,
			MaximumHP:      input.MaximumHP,
			Damage:         input.Damage,
			Class:          input.Class,
			Shield:         input.Shield,
			Rank:           input.Rank,
			Number:         input.Number,
			Move:           input.Move,
			Attack:         input.Attack,
			XP:             input.XP,
			InnateDefenses: input.InnateDefenses,
			InnateOffenses: input.InnateOffenses,
			Statuses:       input.Statuses,
		}}
}
