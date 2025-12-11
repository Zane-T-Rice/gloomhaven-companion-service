package inputs

import (
	"encoding/json"
	"gloomhaven-companion-service/internal/dto"
	"gloomhaven-companion-service/internal/types"
)

// Apply incoming JSON input on top of existing Figure.
// This means fields which are not in the input JSON will be untouched,
// but fields which are null in the input JSON will be unset.
func NewFigurePatchInput(rawInput []byte, dto *dto.Figure) types.FigurePatchInput {
	figurePatchInput := types.FigurePatchInput{
		FigureCreateInput: types.FigureCreateInput{
			Name:           dto.Name,
			MaximumHP:      dto.MaximumHP,
			Damage:         dto.Damage,
			Class:          dto.Class,
			Shield:         dto.Shield,
			Retaliate:      dto.Retaliate,
			Rank:           dto.Rank,
			Number:         dto.Number,
			Move:           dto.Move,
			Attack:         dto.Attack,
			Pierce:         dto.Pierce,
			Range:          dto.Range,
			Target:         dto.Target,
			XP:             dto.XP,
			InnateDefenses: dto.InnateDefenses,
			InnateOffenses: dto.InnateOffenses,
			Statuses:       dto.Statuses,
			Special:        dto.Special,
			Alignment:      dto.Alignment,
			AttackPlusC:    dto.AttackPlusC,
			Flying:         dto.Flying,
			UpdatedAt:      dto.UpdatedAt,
		},
	}
	json.Unmarshal(rawInput, &figurePatchInput)
	return figurePatchInput
}
