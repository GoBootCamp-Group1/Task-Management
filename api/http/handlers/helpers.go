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

type Response struct {
	Success bool   `json:"success"`
	Status  int    `json:"status"`
	Data    any    `json:"data,omitempty"`
	Message string `json:"message,omitempty"`
}

func SendSuccessResponse(c *fiber.Ctx, message string, data any) error {
	response := Response{
		Success: true,
		Status:  fiber.StatusOK,
		Data:    data,
		Message: message,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

func SendError(c *fiber.Ctx, err error, status int) error {
	if status == 0 {
		status = fiber.StatusInternalServerError
	}

	c.Locals(valuecontext.IsTxError, err)

	response := Response{
		Success: false,
		Status:  status,
		Message: err.Error(),
	}
	return c.Status(status).JSON(response)
}

func SendUserToken(c *fiber.Ctx, authToken *services.UserToken) error {
	response := Response{
		Success: true,
		Status:  fiber.StatusOK,
		Data: map[string]interface{}{
			"auth":    authToken.AuthorizationToken,
			"refresh": authToken.RefreshToken,
			"exp":     authToken.ExpiresAt,
		},
	}
	return c.Status(fiber.StatusOK).JSON(response)
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
