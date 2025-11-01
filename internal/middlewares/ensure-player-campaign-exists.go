package middlewares

import (
	"gloomhaven-companion-service/internal/constants"
	errors "gloomhaven-companion-service/internal/errors"
	"gloomhaven-companion-service/internal/types"
	"gloomhaven-companion-service/internal/utils"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gofiber/fiber/v2"
)

func EnsurePlayerCampaignExists(s *utils.DynamoDB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
		playerId := token.RegisteredClaims.Subject
		campaignId := c.Params("campaignId")
		playerCampaign := types.PlayerCampaignItem{}
		if err := s.GetItem(
			constants.PARENT,
			constants.PLAYER+constants.SEPERATOR+playerId,
			constants.ENTITY,
			constants.CAMPAIGN+constants.SEPERATOR+campaignId,
			&playerCampaign,
		); err != nil {
			return errors.NewForbiddenError()
		}
		return c.Next()
	}
}
