package middlerwares

import (
	"strings"
	"time"

	"github.com/GoBootCamp-Group1/Task-Management/api/http/handlers"
	"github.com/GoBootCamp-Group1/Task-Management/pkg/jwt"
	"github.com/GoBootCamp-Group1/Task-Management/pkg/log"

	"github.com/gofiber/fiber/v2"
)

var (
	ErrNoAuthToken        = fiber.NewError(fiber.StatusUnauthorized, "authorization token not specified")
	ErrMalformedAuthToken = fiber.NewError(fiber.StatusUnauthorized, "authorization token is malformed")
	ErrTokenExpired       = fiber.NewError(fiber.StatusUnauthorized, "token expired")
)

func Auth(secret []byte) fiber.Handler {
	return func(c *fiber.Ctx) error {
		h := c.GetReqHeaders()["Authorization"]
		if len(h) == 0 {
			log.ErrorLog.Printf("Error authenticating: %v\n", ErrNoAuthToken)
			return handlers.SendError(c, ErrNoAuthToken)
		}

		// Check if the Authorization header starts with "Bearer "
		if !strings.HasPrefix(h[0], "Bearer ") {
			log.ErrorLog.Printf("Error malformed authentication token: %v\n", ErrMalformedAuthToken)
			return handlers.SendError(c, ErrMalformedAuthToken)
		}

		// Extract the token part
		tokenString := strings.TrimPrefix(h[0], "Bearer ")

		claims, err := jwt.ParseToken(tokenString, secret)
		if err != nil {
			log.ErrorLog.Printf("Error unathorized: %v\n", err)
			return handlers.SendError(c, ErrNoAuthToken)
		}

		c.Locals(jwt.UserClaimKey, claims)

		if claims.ExpiresAt.Before(time.Now()) {
			log.ErrorLog.Printf("Error expired token: %v\n", ErrTokenExpired)
			return handlers.SendError(c, ErrTokenExpired)
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
			return handlers.SendError(c, fiber.NewError(fiber.StatusForbidden, "You don't have access to this section"))
		}

		return c.Next()
	}
}
