package middlewares

import (
	errors "gloomhaven-companion-service/internal/errors"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gofiber/fiber/v2"
)

func HasScope(requiredScope string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
		claims := token.CustomClaims.(*CustomClaims)
		if !claims.HasScope(requiredScope) {
			return errors.NewForbiddenError()
		}
		return c.Next()
	}
}

func HasOneOfScopes(requiredScopes []string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
		claims := token.CustomClaims.(*CustomClaims)
		for _, requiredScope := range requiredScopes {
			if claims.HasScope(requiredScope) {
				return c.Next()
			}
		}
		return errors.NewForbiddenError()
	}
}
