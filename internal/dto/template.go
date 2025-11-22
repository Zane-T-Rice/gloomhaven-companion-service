package dto

import "gloomhaven-companion-service/internal/types"

type Template struct {
	Parent string `dynamodbav:"parent" json:"parent"`
	Entity string `dynamodbav:"entity" json:"entity"`

	Type         *string            `dynamodbav:"type,omitempty" json:"type"`
	StandeeLimit *int               `dynamodbav:"standee_limit,omitempty" json:"standeeLimit"`
	Stats        map[int]types.Stat `dynamodbav:"stats,omitempty" json:"stats"`
	UpdatedAt    *string            `dynamodbav:"updated_at" json:"updatedAt"`
}

func NewTemplate(item types.TemplateItem) Template {
	return Template{
		Parent:       item.Parent,
		Entity:       item.Entity,
		Type:         item.Type,
		Stats:        item.Stats,
		StandeeLimit: item.StandeeLimit,
		UpdatedAt:    item.UpdatedAt,
	}
}
