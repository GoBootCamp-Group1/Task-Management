package handlers

import (
	"fmt"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/services"
	"github.com/GoBootCamp-Group1/Task-Management/pkg/utils"
	"github.com/GoBootCamp-Group1/Task-Management/pkg/validation"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type CreateBoardRequest struct {
	Name      string `json:"name" validate:"required,min=3,max=50,excludesall=;" example:"new board"`
	IsPrivate bool   `json:"is_private" example:"false"`
}

// CreateBoard creates a new board
// @Summary Create Board
// @Description creates a board
// @Tags Board
// @Accept  json
// @Produce json
// @Param   body  body      CreateBoardRequest  true  "Create Board"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /boards [post]
// @Security ApiKeyAuth
func CreateBoard(boardService *services.BoardService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		validate := validation.NewValidator()
		var input CreateBoardRequest

		if err := c.BodyParser(&input); err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}

		err := validate.Struct(input)
		if err != nil {
			fmt.Printf("%+v\n", err)
			return SendError(c, err, fiber.StatusBadRequest)
		}

		userId, err := utils.GetUserID(c)
		if err != nil {
			return SendError(c, err, fiber.StatusInternalServerError)
		}

		boardModel := domains.Board{
			CreatedBy: userId,
			Name:      input.Name,
			IsPrivate: input.IsPrivate,
		}

		err = boardService.CreateBoard(c.Context(), &boardModel)
		if err != nil {
			return SendError(c, err, fiber.StatusInternalServerError)
		}

		return SendSuccessResponse(c, "board")
	}
}

// GetBoardByID get a board
// @Summary Get Board
// @Description gets a board
// @Tags Board
// @Produce json
// @Param   id      path     string  true  "Board ID"
// @Success 200 {object} domains.Board
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /boards/{id} [get]
// @Security ApiKeyAuth
func GetBoardByID(boardService *services.BoardService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseUint(c.Params("id"), 10, 32)
		if err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}

		board, err := boardService.GetBoardByID(c.Context(), uint(id))
		if err != nil {
			return SendError(c, err, fiber.StatusInternalServerError)
		}

		if board == nil {
			return SendError(c, fiber.NewError(fiber.StatusNotFound, "Board not found"), fiber.StatusNotFound)
		}

		return c.JSON(board)
	}
}

type UpdateBoardRequest struct {
	Name      string `json:"name" validate:"required,min=3,max=50,excludesall=;" example:"new board"`
	IsPrivate bool   `json:"is_private" example:"false"`
}

// UpdateBoard updates a new board
// @Summary Update Board
// @Description updates a board
// @Tags Board
// @Accept  json
// @Produce json
// @Param   id      path     string  true  "Board ID"
// @Param   body  body      UpdateBoardRequest  true  "Update Board"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /boards/{id} [put]
// @Security ApiKeyAuth
func UpdateBoard(boardService *services.BoardService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		validate := validation.NewValidator()
		id, err := strconv.ParseUint(c.Params("id"), 10, 32)
		if err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}
		var input UpdateBoardRequest

		if err = c.BodyParser(&input); err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}

		err = validate.Struct(input)
		if err != nil {
			fmt.Printf("%+v\n", err)
			return SendError(c, err, fiber.StatusBadRequest)
		}

		boardModel := domains.Board{
			ID:        uint(id),
			Name:      input.Name,
			IsPrivate: input.IsPrivate,
		}

		err = boardService.UpdateBoard(c.Context(), &boardModel)
		if err != nil {
			return SendError(c, err, fiber.StatusInternalServerError)
		}

		return SendSuccessResponse(c, "board")
	}
}

// DeleteBoard delete a board
// @Summary Delete Board
// @Description deleted a board
// @Tags Board
// @Produce json
// @Param   id      path     string  true  "Board ID"
// @Success 204
// @Failure 400
// @Failure 500
// @Router /boards/{id} [delete]
// @Security ApiKeyAuth
func DeleteBoard(boardService *services.BoardService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseUint(c.Params("id"), 10, 32)
		if err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}

		err = boardService.DeleteBoard(c.Context(), uint(id))
		if err != nil {
			return SendError(c, err, fiber.StatusInternalServerError)
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}
