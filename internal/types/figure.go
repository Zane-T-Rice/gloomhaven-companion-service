package types

import (
	"encoding/json"
	"gloomhaven-companion-service/internal/constants"
	"strconv"
	"time"
)

type FigureCreateInput struct {
	Name           *string `dynamodbav:"name" json:"name"`
	MaximumHP      *int    `dynamodbav:"maximum_hp" json:"maximumHP"`
	Damage         *int    `dynamodbav:"damage" json:"damage"`
	Class          *string `dynamodbav:"class" json:"class"`
	Shield         *int    `dynamodbav:"shield" json:"shield"`
	Retaliate      *int    `dynamodbav:"retaliate" json:"retaliate"`
	Rank           *string `dynamodbav:"rank" json:"rank"`
	Number         *int    `dynamodbav:"number" json:"number"`
	Move           *int    `dynamodbav:"move" json:"move"`
	Attack         *int    `dynamodbav:"attack" json:"attack"`
	Target         *int    `dynamodbav:"target" json:"target"`
	Pierce         *int    `dynamodbav:"pierce" json:"pierce"`
	Range          *int    `dynamodbav:"range" json:"range"`
	XP             *int    `dynamodbav:"xp" json:"xp"`
	InnateDefenses *string `dynamodbav:"innate_defenses" json:"innateDefenses"`
	InnateOffenses *string `dynamodbav:"innate_offenses" json:"innateOffenses"`
	Statuses       *string `dynamodbav:"statuses" json:"statuses"`
	Special        *string `dynamodbav:"special" json:"special"`
	Alignment      *string `dynamodbav:"alignment" json:"alignment"`
	AttackPlusC    *bool   `dynamodbav:"attack_plus_c" json:"attackPlusC"`
	Flying         *bool   `dynamodbav:"flying" json:"flying"`
	UpdatedAt      *string `dynamodbav:"updated_at" json:"updatedAt"`
}

type FigurePatchInput struct {
	FigureCreateInput `dynamodbav:",inline"`
}

type FigureItem struct {
	Item              `dynamodbav:",inline"`
	FigureCreateInput `dynamodbav:",inline"`
}

func NewFigureItem(input FigureCreateInput, campaignId string, scenarioId string, figureId string) FigureItem {
	updatedAt := time.Now().UTC().Format(time.RFC3339)
	return FigureItem{Item: Item{
		Parent: constants.CAMPAIGN + constants.SEPERATOR + campaignId + constants.SCENARIO + constants.SEPERATOR + scenarioId,
		Entity: constants.FIGURE + constants.SEPERATOR + figureId,
	},
		FigureCreateInput: FigureCreateInput{
			Name:           input.Name,
			MaximumHP:      input.MaximumHP,
			Damage:         input.Damage,
			Class:          input.Class,
			Shield:         input.Shield,
			Retaliate:      input.Retaliate,
			Rank:           input.Rank,
			Number:         input.Number,
			Move:           input.Move,
			Attack:         input.Attack,
			Pierce:         input.Pierce,
			Range:          input.Range,
			Target:         input.Target,
			XP:             input.XP,
			InnateDefenses: input.InnateDefenses,
			InnateOffenses: input.InnateOffenses,
			Statuses:       input.Statuses,
			Special:        input.Special,
			Alignment:      input.Alignment,
			AttackPlusC:    input.AttackPlusC,
			Flying:         input.Flying,
			UpdatedAt:      &updatedAt,
		}}
}

// The nightmarish FigureNumber, FigureString, and FigureBool as well as the custom
// unmarshalers are to enable unsetting fields during partial update operations.
// If a field is null or "" in JSON, it is unset. If a field is missing in the
// JSON it is not updated. If a field is set in the JSON, then it is set to the
// new value.
type FigureNumber struct {
	IsNull bool
	Value  *int
}

func (t *FigureNumber) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		// The value "null" for a number input means to unset the field.
		*t = FigureNumber{IsNull: true}
	} else {
		value, err := strconv.Atoi(string(data))
		if err != nil {
			// Ignore this value if it cannot be parsed to a number.
			return err
		}
		*t = FigureNumber{IsNull: false, Value: &value}
	}
	return nil
}

type FigureString struct {
	IsNull bool
	Value  *string
}

func (t *FigureString) UnmarshalJSON(data []byte) error {
	if string(data) == "true" || string(data) == "false" {
		// true and false are not valid values for a string field, so ignore them.
	} else if string(data) == "null" || string(data) == "" {
		// The values "null" and "" in the input JSON means to unset the field.
		*t = FigureString{IsNull: true}
	} else {
		// Trim off the ""
		value := string(data)[1 : len(data)-1]
		*t = FigureString{IsNull: false, Value: &value}
	}
	return nil
}

type FigureBool struct {
	IsNull bool
	Value  *bool
}

func (t *FigureBool) UnmarshalJSON(data []byte) error {
	if string(data) == "true" || string(data) == "false" {
		value := string(data) == "true"
		*t = FigureBool{IsNull: false, Value: &value}
	} else if string(data) == "null" {
		// The value "null" means to unset the field.
		*t = FigureBool{IsNull: true}
	}
	return nil
}

type FigurePatchUnmarshalInput struct {
	Name           FigureString
	MaximumHP      FigureNumber
	Damage         FigureNumber
	Class          FigureString
	Shield         FigureNumber
	Retaliate      FigureNumber
	Rank           FigureString
	Number         FigureNumber
	Move           FigureNumber
	Attack         FigureNumber
	Target         FigureNumber
	Pierce         FigureNumber
	Range          FigureNumber
	XP             FigureNumber
	InnateDefenses FigureString
	InnateOffenses FigureString
	Statuses       FigureString
	Special        FigureString
	Alignment      FigureString
	AttackPlusC    FigureBool
	Flying         FigureBool
	UpdatedAt      FigureString
}

func setNumberField(value *int, figurePatchInputValue **int, isNull bool) {
	if isNull {
		*figurePatchInputValue = nil
	} else if value != nil {
		*figurePatchInputValue = value
	}
}

func setStringField(value *string, figurePatchInputValue **string, isNull bool) {
	if isNull {
		*figurePatchInputValue = nil
	} else if value != nil {
		*figurePatchInputValue = value
	}
}

func setBoolField(value *bool, figurePatchInputValue **bool, isNull bool) {
	if isNull {
		*figurePatchInputValue = nil
	} else if value != nil {
		*figurePatchInputValue = value
	}
}

func NewPatchFigureInput(rawInput []byte, existingFigure FigureItem) FigurePatchInput {
	figurePatchInput := FigurePatchInput{
		FigureCreateInput: FigureCreateInput{
			Name:           existingFigure.Name,
			MaximumHP:      existingFigure.MaximumHP,
			Damage:         existingFigure.Damage,
			Class:          existingFigure.Class,
			Shield:         existingFigure.Shield,
			Retaliate:      existingFigure.Retaliate,
			Rank:           existingFigure.Rank,
			Number:         existingFigure.Number,
			Move:           existingFigure.Move,
			Attack:         existingFigure.Attack,
			Pierce:         existingFigure.Pierce,
			Range:          existingFigure.Range,
			Target:         existingFigure.Target,
			XP:             existingFigure.XP,
			InnateDefenses: existingFigure.InnateDefenses,
			InnateOffenses: existingFigure.InnateOffenses,
			Statuses:       existingFigure.Statuses,
			Special:        existingFigure.Special,
			Alignment:      existingFigure.Alignment,
			AttackPlusC:    existingFigure.AttackPlusC,
			Flying:         existingFigure.Flying,
		}}
	var rawData FigurePatchUnmarshalInput
	json.Unmarshal(rawInput, &rawData)
	setStringField(rawData.Name.Value, &figurePatchInput.Name, rawData.Name.IsNull)
	setNumberField(rawData.MaximumHP.Value, &figurePatchInput.MaximumHP, rawData.MaximumHP.IsNull)
	setNumberField(rawData.Damage.Value, &figurePatchInput.Damage, rawData.Damage.IsNull)
	setStringField(rawData.Class.Value, &figurePatchInput.Class, rawData.Class.IsNull)
	setNumberField(rawData.Shield.Value, &figurePatchInput.Shield, rawData.Shield.IsNull)
	setNumberField(rawData.Retaliate.Value, &figurePatchInput.Retaliate, rawData.Retaliate.IsNull)
	setStringField(rawData.Rank.Value, &figurePatchInput.Rank, rawData.Rank.IsNull)
	setNumberField(rawData.Number.Value, &figurePatchInput.Number, rawData.Number.IsNull)
	setNumberField(rawData.Move.Value, &figurePatchInput.Move, rawData.Move.IsNull)
	setNumberField(rawData.Attack.Value, &figurePatchInput.Attack, rawData.Attack.IsNull)
	setNumberField(rawData.Pierce.Value, &figurePatchInput.Pierce, rawData.Pierce.IsNull)
	setNumberField(rawData.Range.Value, &figurePatchInput.Range, rawData.Range.IsNull)
	setNumberField(rawData.Target.Value, &figurePatchInput.Target, rawData.Target.IsNull)
	setNumberField(rawData.XP.Value, &figurePatchInput.XP, rawData.XP.IsNull)
	setStringField(rawData.InnateDefenses.Value, &figurePatchInput.InnateDefenses, rawData.InnateDefenses.IsNull)
	setStringField(rawData.InnateOffenses.Value, &figurePatchInput.InnateOffenses, rawData.InnateOffenses.IsNull)
	setStringField(rawData.Statuses.Value, &figurePatchInput.Statuses, rawData.Statuses.IsNull)
	setStringField(rawData.Special.Value, &figurePatchInput.Special, rawData.Special.IsNull)
	setStringField(rawData.Alignment.Value, &figurePatchInput.Alignment, rawData.Alignment.IsNull)
	setBoolField(rawData.AttackPlusC.Value, &figurePatchInput.AttackPlusC, rawData.AttackPlusC.IsNull)
	setBoolField(rawData.Flying.Value, &figurePatchInput.Flying, rawData.Flying.IsNull)
	setStringField(rawData.UpdatedAt.Value, &figurePatchInput.UpdatedAt, rawData.UpdatedAt.IsNull)
	return figurePatchInput
}
