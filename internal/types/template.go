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
	// The edition this template is listed under. Templates without one are in the default edition.
	EditionId *string `dynamodbav:"edition_id,omitempty" json:"editionId,omitempty"`
	UpdatedAt *string `dynamodbav:"updated_at,omitempty" json:"updatedAt,omitempty"`
}

type TemplatePatchInput struct {
	TemplateCreateInput `dynamodbav:",inline"`
}

type TemplateItem struct {
	Item                `dynamodbav:",inline"`
	TemplateCreateInput `dynamodbav:",inline"`
}

// NewTemplateItem makes a template for the given owner, which is either constants.TEMPLATE for a
// system template or the parent of a player's items.
func NewTemplateItem(input TemplateCreateInput, templateId string, owner string) TemplateItem {
	updatedAt := time.Now().UTC().Format(time.RFC3339)
	return TemplateItem{
		Item: Item{
			Parent: owner,
			Entity: constants.TEMPLATE + constants.SEPERATOR + templateId,
		},
		TemplateCreateInput: TemplateCreateInput{
			Type:         input.Type,
			Stats:        input.Stats,
			StandeeLimit: input.StandeeLimit,
			EditionId:    input.EditionId,
			UpdatedAt:    &updatedAt,
		},
	}
}
