package types

import (
	"gloomhaven-companion-service/internal/constants"
	"time"
)

type Stat struct {
	Normal    FigureCreateInput `dynamodbav:"normal,omitempty" json:"normal"`
	Elite     FigureCreateInput `dynamodbav:"elite,omitempty" json:"elite"`
	Character FigureCreateInput `dynamodbav:"character,omitempty" json:"character"`
	Summon    FigureCreateInput `dynamodbav:"summon,omitempty" json:"summon"`
}

type TemplateCreateInput struct {
	Class        *string      `dynamodbav:"class,omitempty" json:"class"`
	StandeeLimit *int         `dynamodbav:"standee_limit,omitempty" json:"standeeLimit"`
	Stats        map[int]Stat `dynamodbav:"stats,omitempty" json:"stats"`
	UpdatedAt    *string      `dynamodbav:"updated_at" json:"updatedAt"`
}

type TemplatePatchInput struct {
	TemplateCreateInput `dynamodbav:",inline"`
}

type TemplateItem struct {
	Item                `dynamodbav:",inline"`
	TemplateCreateInput `dynamodbav:",inline"`
}

func NewTemplateItem(input TemplateCreateInput, templateId string) TemplateItem {
	updatedAt := time.Now().UTC().Format(time.RFC3339)
	return TemplateItem{
		Item: Item{
			Parent: constants.TEMPLATE,
			Entity: constants.TEMPLATE + constants.SEPERATOR + templateId,
		},
		TemplateCreateInput: TemplateCreateInput{
			Class:        input.Class,
			Stats:        input.Stats,
			StandeeLimit: input.StandeeLimit,
			UpdatedAt:    &updatedAt,
		},
	}
}
