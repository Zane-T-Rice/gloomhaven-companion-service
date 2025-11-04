package dto

import "gloomhaven-companion-service/internal/types"

type Figure struct {
	Parent    string  `dynamodbav:"parent" json:"parent"`
	Entity    string  `dynamodbav:"entity" json:"entity"`
	Name      *string `dynamodbav:"name" json:"name"`
	MaximumHP *int    `dynamodbav:"maximum_hp" json:"maximumHP"`
	Damage    *int    `dynamodbav:"damage" json:"damage"`
	Class     *string `dynamodbav:"class" json:"class"`
	Shield    *int    `dynamodbav:"shield" json:"shield"`
	Rank      *string `dynamodbav:"rank" json:"rank"`
	Number    *int    `dynamodbav:"number" json:"number"`
	Move      *int    `dynamodbav:"move" json:"move"`
	Attack    *int    `dynamodbav:"attack" json:"attack"`
}

func NewFigure(item types.FigureItem) Figure {
	return Figure{
		Parent:    item.Parent,
		Entity:    item.Entity,
		Name:      item.Name,
		MaximumHP: item.MaximumHP,
		Damage:    item.Damage,
		Class:     item.Class,
		Shield:    item.Shield,
		Rank:      item.Rank,
		Number:    item.Number,
		Move:      item.Move,
		Attack:    item.Attack,
	}
}
