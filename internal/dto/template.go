package dto

import (
	"gloomhaven-companion-service/internal/constants"
	"gloomhaven-companion-service/internal/types"
)

type Template struct {
	Parent string `dynamodbav:"parent" json:"parent"`
	Entity string `dynamodbav:"entity" json:"entity"`

	Type         *string             `dynamodbav:"type,omitempty" json:"type,omitempty"`
	StandeeLimit *int                `dynamodbav:"standee_limit,omitempty" json:"standeeLimit,omitempty"`
	Stats        *map[int]types.Stat `dynamodbav:"stats,omitempty" json:"stats,omitempty"`
	// "system", "mine" or "shared". See ScopeOf.
	Scope     string  `json:"scope"`
	EditionId string  `json:"editionId"`
	UpdatedAt *string `dynamodbav:"updated_at" json:"updatedAt,omitempty"`
}

// ScopeOf says whose an item is, from the point of view of the player asking.
func ScopeOf(owner string, playerId string) string {
	if owner == constants.TEMPLATE || owner == constants.EDITION {
		return constants.SCOPE_NAME_SYSTEM
	}
	if owner == constants.PLAYER+constants.SEPERATOR+playerId {
		return constants.SCOPE_NAME_MINE
	}
	return constants.SCOPE_NAME_SHARED
}

// The parent is only given out for system items. For a player's items it would be their login id.
func maskOwner(owner string, scope string) string {
	if scope == constants.SCOPE_NAME_SYSTEM {
		return owner
	}
	return constants.PLAYER
}

func NewTemplate(item types.TemplateItem, playerId string) Template {
	scope := ScopeOf(item.Parent, playerId)
	editionId := constants.DEFAULT_EDITION_ID
	if item.EditionId != nil && *item.EditionId != "" {
		editionId = *item.EditionId
	}
	return Template{
		Parent:       maskOwner(item.Parent, scope),
		Entity:       item.Entity,
		Type:         item.Type,
		Stats:        item.Stats,
		StandeeLimit: item.StandeeLimit,
		Scope:        scope,
		EditionId:    editionId,
		UpdatedAt:    item.UpdatedAt,
	}
}
