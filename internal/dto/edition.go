package dto

import (
	"gloomhaven-companion-service/internal/constants"
	"gloomhaven-companion-service/internal/types"
)

type Edition struct {
	Parent    string  `json:"parent"`
	Entity    string  `json:"entity"`
	Name      *string `json:"name"`
	Scope     string  `json:"scope"`
	UpdatedAt *string `json:"updatedAt,omitempty"`
}

func NewEdition(item types.EditionItem, playerId string) Edition {
	scope := ScopeOf(item.Parent, playerId)
	return Edition{
		Parent:    maskOwner(item.Parent, scope),
		Entity:    item.Entity,
		Name:      item.Name,
		Scope:     scope,
		UpdatedAt: item.UpdatedAt,
	}
}

// The edition that every template without one is in.
func NewDefaultEdition() Edition {
	name := constants.DEFAULT_EDITION_NAME
	return Edition{
		Parent: constants.EDITION,
		Entity: constants.EDITION + constants.SEPERATOR + constants.DEFAULT_EDITION_ID,
		Name:   &name,
		Scope:  constants.SCOPE_NAME_SYSTEM,
	}
}
