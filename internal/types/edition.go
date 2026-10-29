package types

import (
	"gloomhaven-companion-service/internal/constants"
	"time"
)

// An edition is a named group of templates, i.e. "Gloomhaven 2E" or a player's "Homebrew".
type EditionCreateInput struct {
	Name      *string `dynamodbav:"name,omitempty" json:"name,omitempty"`
	UpdatedAt *string `dynamodbav:"updated_at,omitempty" json:"updatedAt,omitempty"`
}

type EditionPatchInput struct {
	EditionCreateInput `dynamodbav:",inline"`
}

type EditionItem struct {
	Item               `dynamodbav:",inline"`
	EditionCreateInput `dynamodbav:",inline"`
}

func NewEditionItem(input EditionCreateInput, editionId string, owner string) EditionItem {
	updatedAt := time.Now().UTC().Format(time.RFC3339)
	return EditionItem{
		Item: Item{
			Parent: owner,
			Entity: constants.EDITION + constants.SEPERATOR + editionId,
		},
		EditionCreateInput: EditionCreateInput{
			Name:      input.Name,
			UpdatedAt: &updatedAt,
		},
	}
}
