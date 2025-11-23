package types

import (
	"gloomhaven-companion-service/internal/constants"
	"time"
)

type FigureCreateInput struct {
	Name           *string `dynamodbav:"name,omitempty" json:"name,omitempty"`
	MaximumHP      *int    `dynamodbav:"maximum_hp,omitempty" json:"maximumHP,omitempty"`
	Damage         *int    `dynamodbav:"damage,omitempty" json:"damage,omitempty"`
	Class          *string `dynamodbav:"class,omitempty" json:"class,omitempty"`
	Shield         *int    `dynamodbav:"shield,omitempty" json:"shield,omitempty"`
	Retaliate      *int    `dynamodbav:"retaliate,omitempty" json:"retaliate,omitempty"`
	Rank           *string `dynamodbav:"rank,omitempty" json:"rank,omitempty"`
	Number         *int    `dynamodbav:"number,omitempty" json:"number,omitempty"`
	Move           *int    `dynamodbav:"move,omitempty" json:"move,omitempty"`
	Attack         *int    `dynamodbav:"attack,omitempty" json:"attack,omitempty"`
	Target         *int    `dynamodbav:"target,omitempty" json:"target,omitempty"`
	Pierce         *int    `dynamodbav:"pierce,omitempty" json:"pierce,omitempty"`
	XP             *int    `dynamodbav:"xp,omitempty" json:"xp,omitempty"`
	InnateDefenses *string `dynamodbav:"innate_defenses,omitempty" json:"innateDefenses,omitempty"`
	InnateOffenses *string `dynamodbav:"innate_offenses,omitempty" json:"innateOffenses,omitempty"`
	Statuses       *string `dynamodbav:"statuses,omitempty" json:"statuses,omitempty"`
	Special        *string `dynamodbav:"special,omitempty" json:"special,omitempty"`
	Alignment      *string `dynamodbav:"alignment,omitempty" json:"alignment,omitempty"`
	AttackPlusC    *bool   `dynamodbav:"attack_plus_c,omitempty" json:"attackPlusC,omitempty"`
	Flying         *bool   `dynamodbav:"flying,omitempty" json:"flying,omitempty"`
	UpdatedAt      *string `dynamodbav:"updated_at,omitempty" json:"updatedAt,omitempty"`
}

type FigurePatchInput struct {
	FigureCreateInput `dynamodbav:",inline"`
}

type FigureItem struct {
	Item              `dynamodbav:",inline"`
	FigureCreateInput `dynamodbav:",inline"`
}

func NewFigureItem(input FigureCreateInput, campaignId string, scenarioId string, figureId string) FigureItem {
	updatedAt := time.Now().UTC().Format(time.RFC3339)
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
			Retaliate:      input.Retaliate,
			Rank:           input.Rank,
			Number:         input.Number,
			Move:           input.Move,
			Attack:         input.Attack,
			Pierce:         input.Pierce,
			Target:         input.Target,
			XP:             input.XP,
			InnateDefenses: input.InnateDefenses,
			InnateOffenses: input.InnateOffenses,
			Statuses:       input.Statuses,
			Special:        input.Special,
			Alignment:      input.Alignment,
			AttackPlusC:    input.AttackPlusC,
			Flying:         input.Flying,
			UpdatedAt:      &updatedAt,
		}}
}
