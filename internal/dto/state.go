package dto

import "gloomhaven-companion-service/internal/types"

type State struct {
	Parent string  `json:"parent"`
	Entity string  `json:"entity"`
	Value  *string `json:"value"`
}

func NewState(item types.StateItem) State {
	return State{
		Parent: item.Parent,
		Entity: item.Entity,
		Value:  item.Value,
	}
}
