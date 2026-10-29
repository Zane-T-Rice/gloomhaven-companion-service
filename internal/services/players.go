package services

import (
	"gloomhaven-companion-service/internal/constants"
	"gloomhaven-companion-service/internal/types"
	"gloomhaven-companion-service/internal/utils"
)

// The parent of the items that belong to a player.
func playerOwner(playerId string) string {
	return constants.PLAYER + constants.SEPERATOR + playerId
}

// The owners whose templates and editions a player can see, as the parent of their items: the system,
// the player, and every other player in any of the player's campaigns.
func visibleOwners(db utils.DynamoDB, playerId string) ([]string, error) {
	me := playerOwner(playerId)
	owners := []string{constants.TEMPLATE, me}
	seen := map[string]bool{me: true}

	playerCampaigns := []types.PlayerCampaignItem{}
	if err := db.Query(
		constants.PARENT,
		me,
		constants.ENTITY,
		constants.CAMPAIGN,
		nil,
		&playerCampaigns,
	); err != nil {
		return nil, err
	}
	for _, playerCampaign := range playerCampaigns {
		fellows := []types.PlayerCampaignItem{}
		if err := db.Query(
			constants.ENTITY,
			playerCampaign.Entity,
			constants.PARENT,
			constants.PLAYER,
			&constants.ENTITY_INDEX,
			&fellows,
		); err != nil {
			return nil, err
		}
		for _, fellow := range fellows {
			if !seen[fellow.Parent] {
				seen[fellow.Parent] = true
				owners = append(owners, fellow.Parent)
			}
		}
	}
	return owners, nil
}
