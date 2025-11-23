package types

import (
	"gloomhaven-companion-service/internal/constants"
	"time"
)

type Stat struct {
	Normal    *FigureCreateInput `dynamodbav:"normal,omitempty" json:"normal,omitempty"`
	Elite     *FigureCreateInput `dynamodbav:"elite,omitempty" json:"elite,omitempty"`
	Boss      *FigureCreateInput `dynamodbav:"boss,omitempty" json:"boss,omitempty"`
	Character *FigureCreateInput `dynamodbav:"character,omitempty" json:"character,omitempty"`
	Summon    *FigureCreateInput `dynamodbav:"summon,omitempty" json:"summon,omitempty"`
}

type TemplateCreateInput struct {
	Type         *string       `dynamodbav:"type,omitempty" json:"type,omitempty"`
	StandeeLimit *int          `dynamodbav:"standee_limit,omitempty" json:"standeeLimit,omitempty"`
	Stats        *map[int]Stat `dynamodbav:"stats,omitempty" json:"stats,omitempty"`
	UpdatedAt    *string       `dynamodbav:"updated_at,omitempty" json:"updatedAt,omitempty"`
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
			Type:         input.Type,
			Stats:        input.Stats,
			StandeeLimit: input.StandeeLimit,
			UpdatedAt:    &updatedAt,
		},
	}
}
