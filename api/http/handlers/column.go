package handlers

import (
	"fmt"

	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domain"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/service"
	"github.com/GoBootCamp-Group1/Task-Management/pkg/utils"
	"github.com/GoBootCamp-Group1/Task-Management/pkg/validation"
	"github.com/gofiber/fiber/v2"
)

type UpdateColumnRequest struct {
	Name string `json:"name" validate:"required,min=3,max=20,excludesall=;" example:"new column"`
}

type CreateColumnRequest struct {
	Name    string `json:"name" validate:"required,min=3,max=20,excludesall=;" example:"new column"`
	IsFinal bool   `json:"is_final" example:"false"`
}

// @Summary Create Column
// @Description create a column
// @Accept json
// @Success 200
// @Failure 400
// @Failure 500
// @Router /columns/{boardId} [post]
// @Security ApiKeyAuth
func CreateColumn(columnService *service.ColumnService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		boardId, errParam := c.ParamsInt("boardId")
		if errParam != nil {
			// todo: sending error should handle better
			return SendError(c, errParam, fiber.StatusBadRequest)
		}
		validate := validation.NewValidator()
		var input CreateColumnRequest

		if err := c.BodyParser(&input); err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}

		if err := validate.Struct(input); err != nil {
			fmt.Printf("%+v\n", err)
			return SendError(c, err, fiber.StatusBadRequest)
		}

		userId, err := utils.GetUserID(c)
		if err != nil {
			return SendError(c, err, 0)
		}

		columnModel := domain.Column{
			CreatedBy: userId,
			Name:      input.Name,
			IsFinal:   input.IsFinal,
			BoardID:   uint(boardId),
		}

		if err := columnService.CreateColumn(c.Context(), &columnModel); err != nil {
			return SendError(c, err, 0)
		}

		// todo: sending response should handle better
		return SendSuccessResponse(c, "column")
	}
}

// GetColumnByID get a column
// @Summary Get Column
// @Description gets a column
// @Tags Get Column
// @Produce json
// @Param   id      path     string  true  "Column ID"
// @Success 200 {object} domain.Column
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /columns/{id} [get]
// @Security ApiKeyAuth
func GetColumnByID(columnService *service.ColumnService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, errParam := c.ParamsInt("id")
		if errParam != nil {
			return SendError(c, errParam, fiber.StatusBadRequest)
		}

		column, err := columnService.GetColumnById(c.Context(), uint(id))
		if err != nil {
			return SendError(c, err, fiber.StatusInternalServerError)
		}

		if column == nil {
			return SendError(c, fiber.NewError(fiber.StatusNotFound, "Column not found"), fiber.StatusNotFound)
		}

		return c.JSON(column)
	}
}

func GetAllColumns(columnService *service.ColumnService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		boardId, errParam := c.ParamsInt("boardId")
		if errParam != nil {
			return SendError(c, errParam, fiber.StatusBadRequest)
		}

		columns, err := columnService.GetAllColumns(c.Context(), uint(boardId))
		if err != nil {
			return SendError(c, err, 0)
		}

		return c.JSON(columns)
	}
}

func UpdateColumn(columnService *service.ColumnService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		validate := validation.NewValidator()
		id, errParam := c.ParamsInt("id")
		if errParam != nil {
			return SendError(c, errParam, fiber.StatusBadRequest)
		}

		var input domain.ColumnUpdate

		if err := c.BodyParser(&input); err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}

		if err := validate.Struct(input); err != nil {
			fmt.Printf("%+v\n", err)
			return SendError(c, err, fiber.StatusBadRequest)
		}

		input.ID = uint(id)

		if err := columnService.Update(c.Context(), &input); err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}

		return SendSuccessResponse(c, "column")
	}
}

func MoveColumn(columnService *service.ColumnService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		validate := validation.NewValidator()
		id, errParam := c.ParamsInt("id")
		if errParam != nil {
			return SendError(c, errParam, fiber.StatusBadRequest)
		}

		var input domain.ColumnMove

		if err := c.BodyParser(&input); err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}

		if err := validate.Struct(input); err != nil {
			fmt.Printf("%+v\n", err)
			return SendError(c, err, fiber.StatusBadRequest)
		}

		input.ID = uint(id)

		if err := columnService.Move(c.Context(), &input); err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}

		return SendSuccessResponse(c, "column")
	}
}

func ChangeFinalColumn(columnService *service.ColumnService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, errParam := c.ParamsInt("id")
		if errParam != nil {
			return SendError(c, errParam, fiber.StatusBadRequest)
		}

		if err := columnService.Final(c.Context(), uint(id)); err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}

		return SendSuccessResponse(c, "column")
	}
}

func DeleteColumn(columnService *service.ColumnService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, errParam := c.ParamsInt("id")
		if errParam != nil {
			return SendError(c, errParam, fiber.StatusBadRequest)
		}

		if err := columnService.Delete(c.Context(), uint(id)); err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}

		return SendSuccessResponse(c, "column")
	}
}
