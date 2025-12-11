package dto

import "gloomhaven-companion-service/internal/types"

type Figure struct {
	Parent         string  `json:"parent"`
	Entity         string  `json:"entity"`
	Name           *string `json:"name,omitempty"`
	MaximumHP      *int    `json:"maximumHP,omitempty"`
	Damage         *int    `json:"damage,omitempty"`
	Class          *string `json:"class,omitempty"`
	Shield         *int    `json:"shield,omitempty"`
	Retaliate      *int    `json:"retaliate,omitempty"`
	Rank           *string `json:"rank,omitempty"`
	Number         *int    `json:"number,omitempty"`
	Move           *int    `json:"move,omitempty"`
	Attack         *int    `json:"attack,omitempty"`
	Target         *int    `json:"target,omitempty"`
	Pierce         *int    `json:"pierce,omitempty"`
	Range          *int    `json:"range,omitempty"`
	XP             *int    `json:"xp,omitempty"`
	InnateDefenses *string `json:"innateDefenses,omitempty"`
	InnateOffenses *string `json:"innateOffenses,omitempty"`
	Statuses       *string `json:"statuses,omitempty"`
	Special        *string `json:"special,omitempty"`
	Alignment      *string `json:"alignment,omitempty"`
	AttackPlusC    *bool   `json:"attackPlusC,omitempty"`
	Flying         *bool   `json:"flying,omitempty"`
	UpdatedAt      *string `json:"updatedAt,omitempty"`
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

func NewFigureItem(dto *Figure) types.FigureItem {
	return types.FigureItem{Item: types.Item{
		Parent: dto.Parent,
		Entity: dto.Entity,
	},
		FigureCreateInput: types.FigureCreateInput{
			Name:           dto.Name,
			MaximumHP:      dto.MaximumHP,
			Damage:         dto.Damage,
			Class:          dto.Class,
			Shield:         dto.Shield,
			Retaliate:      dto.Retaliate,
			Rank:           dto.Rank,
			Number:         dto.Number,
			Move:           dto.Move,
			Attack:         dto.Attack,
			Pierce:         dto.Pierce,
			Range:          dto.Range,
			Target:         dto.Target,
			XP:             dto.XP,
			InnateDefenses: dto.InnateDefenses,
			InnateOffenses: dto.InnateOffenses,
			Statuses:       dto.Statuses,
			Special:        dto.Special,
			Alignment:      dto.Alignment,
			AttackPlusC:    dto.AttackPlusC,
			Flying:         dto.Flying,
			UpdatedAt:      dto.UpdatedAt,
		}}
}
