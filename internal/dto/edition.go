package dto

import (
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
