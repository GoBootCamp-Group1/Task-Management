package handlers

import (
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domain"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/service"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func CreateBoard(boardService *service.BoardService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input struct {
			CreatedBy uint   `json:"created_by"`
			Name      string `json:"name"`
			IsPrivate bool   `json:"is_private"`
		}

		if err := c.BodyParser(&input); err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}

		boardModel := domain.Board{
			CreatedBy: input.CreatedBy,
			Name:      input.Name,
			IsPrivate: input.IsPrivate,
		}

		err := boardService.CreateBoard(c.Context(), &boardModel)
		if err != nil {
			return SendError(c, err, fiber.StatusInternalServerError)
		}

		return SendSuccessResponse(c, "board")
	}
}

func GetBoardByID(boardService *service.BoardService) fiber.Handler {
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

func UpdateBoard(boardService *service.BoardService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var input struct {
			ID        uint   `json:"id"`
			Name      string `json:"name"`
			IsPrivate bool   `json:"is_private"`
		}

		if err := c.BodyParser(&input); err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}

		boardModel := domain.Board{
			ID:        input.ID,
			Name:      input.Name,
			IsPrivate: input.IsPrivate,
		}

		err := boardService.UpdateBoard(c.Context(), &boardModel)
		if err != nil {
			return SendError(c, err, fiber.StatusInternalServerError)
		}

		return SendSuccessResponse(c, "board")
	}
}

func DeleteBoard(boardService *service.BoardService) fiber.Handler {
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
