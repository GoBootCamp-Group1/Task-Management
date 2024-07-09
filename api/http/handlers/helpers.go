package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/services"
	"github.com/GoBootCamp-Group1/Task-Management/pkg/jwt"
	"github.com/GoBootCamp-Group1/Task-Management/pkg/validation"
	"github.com/GoBootCamp-Group1/Task-Management/pkg/valuecontext"
	"github.com/gofiber/fiber/v2"
)

const UserClaimKey = jwt.UserClaimKey

var (
	errWrongClaimType = errors.New("wrong claim type")
)

type ServiceFactory[T any] func(context.Context) T

func SendSuccessResponse(c *fiber.Ctx, entity string) error {
	return c.Status(fiber.StatusOK).JSON(map[string]any{
		"message": fmt.Sprintf("%s created successfully", entity),
	})
}

func SendError(c *fiber.Ctx, err error, status int) error {
	if status == 0 {
		status = fiber.StatusInternalServerError
	}

	c.Locals(valuecontext.IsTxError, err)

	return c.Status(status).JSON(map[string]any{
		"error_msg": err.Error(),
	})
}

func SendUserToken(c *fiber.Ctx, authToken *services.UserToken) error {
	return c.Status(fiber.StatusOK).JSON(map[string]any{
		"auth":    authToken.AuthorizationToken,
		"refresh": authToken.RefreshToken,
		"exp":     authToken.ExpiresAt,
	})
}

func PageAndPageSize(c *fiber.Ctx) (int, int) {
	page, pageSize := c.QueryInt("page"), c.QueryInt("page_size")
	if page <= 0 {
		page = 1
	}

	if pageSize <= 0 {
		pageSize = 20
	}

	return page, pageSize
}

func ValidateAndFill(c *fiber.Ctx, input any) error {
	validate := validation.NewValidator()

	if err := c.BodyParser(&input); err != nil {
		return SendError(c, err, fiber.StatusBadRequest)
	}

	err := validate.Struct(input)
	if err != nil {
		fmt.Printf("%+v\n", err)
		return SendError(c, err, fiber.StatusBadRequest)
	}

	return nil
}
