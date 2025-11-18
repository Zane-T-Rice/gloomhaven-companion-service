package dto

import "gloomhaven-companion-service/internal/types"

type Figure struct {
	Parent         string  `dynamodbav:"parent" json:"parent"`
	Entity         string  `dynamodbav:"entity" json:"entity"`
	Name           *string `dynamodbav:"name" json:"name"`
	MaximumHP      *int    `dynamodbav:"maximum_hp" json:"maximumHP"`
	Damage         *int    `dynamodbav:"damage" json:"damage"`
	Class          *string `dynamodbav:"class" json:"class"`
	Shield         *int    `dynamodbav:"shield" json:"shield"`
	Retaliate      *int    `dynamodbav:"retaliate" json:"retaliate"`
	Rank           *string `dynamodbav:"rank" json:"rank"`
	Number         *int    `dynamodbav:"number" json:"number"`
	Move           *int    `dynamodbav:"move" json:"move"`
	Attack         *int    `dynamodbav:"attack" json:"attack"`
	Target         *int    `dynamodbav:"target" json:"target"`
	Pierce         *int    `dynamodbav:"pierce" json:"pierce"`
	XP             *int    `dynamodbav:"xp" json:"xp"`
	InnateDefenses *string `dynamodbav:"innate_defenses" json:"innateDefenses"`
	InnateOffenses *string `dynamodbav:"innate_offenses" json:"innateOffenses"`
	Statuses       *string `dynamodbav:"statuses" json:"statuses"`
	Special        *string `dynamodbav:"special" json:"special"`
	UpdatedAt      *string `dynamodbav:"updated_at" json:"updatedAt"`
}

func NewFigure(item types.FigureItem) Figure {
	return Figure{
		Parent:         item.Parent,
		Entity:         item.Entity,
		Name:           item.Name,
		MaximumHP:      item.MaximumHP,
		Damage:         item.Damage,
		Class:          item.Class,
		Shield:         item.Shield,
		Retaliate:      item.Retaliate,
		Rank:           item.Rank,
		Number:         item.Number,
		Move:           item.Move,
		Attack:         item.Attack,
		Target:         item.Target,
		Pierce:         item.Pierce,
		XP:             item.XP,
		InnateDefenses: item.InnateDefenses,
		InnateOffenses: item.InnateOffenses,
		Statuses:       item.Statuses,
		Special:        item.Special,
		UpdatedAt:      item.UpdatedAt,
	}
}
