package middlerwares

import (
	"errors"
	"github.com/GoBootCamp-Group1/Task-Management/api/http/handlers"
	"github.com/GoBootCamp-Group1/Task-Management/pkg/jwt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func Auth(secret []byte) fiber.Handler {
	return func(c *fiber.Ctx) error {
		h := c.GetReqHeaders()["Authorization"]
		if len(h) == 0 {
			return handlers.SendError(c, errors.New("authorization token not specified"), fiber.StatusUnauthorized)
		}

		// Check if the Authorization header starts with "Bearer "
		if !strings.HasPrefix(h[0], "Bearer ") {
			return handlers.SendError(c, errors.New("authorization token is malformed"), fiber.StatusUnauthorized)
		}

		// Extract the token part
		tokenString := strings.TrimPrefix(h[0], "Bearer ")

		claims, err := jwt.ParseToken(tokenString, secret)
		if err != nil {
			return handlers.SendError(c, err, fiber.StatusUnauthorized)
		}

		c.Locals(jwt.UserClaimKey, claims)

		if claims.ExpiresAt.Before(time.Now()) {
			return handlers.SendError(c, errors.New("token expired"), fiber.StatusUnauthorized)
		}

		return c.Next()
	}
}

func RoleChecker(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims := c.Locals(jwt.UserClaimKey).(*jwt.UserClaims)
		hasAccess := false
		for _, role := range roles {
			if claims.Role == role {
				hasAccess = true
				break
			}
		}

		if !hasAccess {
			return handlers.SendError(c, errors.New("you don't have access to this section"), fiber.StatusForbidden)
		}

		return c.Next()
	}
}
