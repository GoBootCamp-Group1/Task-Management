package utils

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/GoBootCamp-Group1/Task-Management/pkg/jwt"
	"github.com/gofiber/fiber/v2"
)

func HashPassword(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil))
}

func GetUserID(c *fiber.Ctx) (uint, error) {
	userClaims, ok := c.Locals(jwt.UserClaimKey).(*jwt.UserClaims)
	if !ok {
		return 0, fiber.NewError(fiber.StatusUnauthorized, "Invalid user claims")
	}
	return userClaims.UserID, nil
}
