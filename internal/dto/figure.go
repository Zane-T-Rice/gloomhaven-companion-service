package dto

import "gloomhaven-companion-service/internal/types"

type Figure struct {
	Parent         string  `dynamodbav:"parent" json:"parent"`
	Entity         string  `dynamodbav:"entity" json:"entity"`
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
	Range          *int    `dynamodbav:"range,omitempty" json:"range,omitempty"`
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
		Range:          item.Range,
		XP:             item.XP,
		InnateDefenses: item.InnateDefenses,
		InnateOffenses: item.InnateOffenses,
		Statuses:       item.Statuses,
		Special:        item.Special,
		Alignment:      item.Alignment,
		AttackPlusC:    item.AttackPlusC,
		Flying:         item.Flying,
		UpdatedAt:      item.UpdatedAt,
	}
}
