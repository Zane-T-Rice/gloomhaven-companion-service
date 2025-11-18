package dto

import "gloomhaven-companion-service/internal/types"

type Template struct {
	Parent string `dynamodbav:"parent" json:"parent"`
	Entity string `dynamodbav:"entity" json:"entity"`

	Class        *string            `dynamodbav:"class,omitempty" json:"class"`
	StandeeLimit *int               `dynamodbav:"standee_limit,omitempty" json:"standeeLimit"`
	Stats        map[int]types.Stat `dynamodbav:"stats,omitempty" json:"stats"`
	UpdatedAt    *string            `dynamodbav:"updated_at" json:"updatedAt"`
}

func NewTemplate(item types.TemplateItem) Template {
	return Template{
		Parent:       item.Parent,
		Entity:       item.Entity,
		Class:        item.Class,
		Stats:        item.Stats,
		StandeeLimit: item.StandeeLimit,
		UpdatedAt:    item.UpdatedAt,
	}
}
