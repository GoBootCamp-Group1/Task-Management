package handlers

import (
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/services"
	"github.com/GoBootCamp-Group1/Task-Management/pkg/log"
	"github.com/GoBootCamp-Group1/Task-Management/pkg/utils"
	"github.com/GoBootCamp-Group1/Task-Management/pkg/validation"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

var (
	ErrBoardNotFound = fiber.NewError(fiber.StatusNotFound, "Board not found")
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
			log.ErrorLog.Printf("Error parsing board creation request body: %v\n", err)
			return SendError(c, err, fiber.StatusBadRequest)
		}

		err := validate.Struct(input)
		if err != nil {
			log.ErrorLog.Printf("Error validating board creation request body: %v\n", err)
			return SendError(c, err, fiber.StatusBadRequest)
		}

		userId, err := utils.GetUserID(c)
		if err != nil {
			log.ErrorLog.Printf("Error loading user: %v\n", err)
			return SendError(c, err, fiber.StatusInternalServerError)
		}

		boardModel := domains.Board{
			CreatedBy: userId,
			Name:      input.Name,
			IsPrivate: input.IsPrivate,
		}

		err = boardService.CreateBoard(c.Context(), &boardModel)
		if err != nil {
			log.ErrorLog.Printf("Error creating board: %v\n", err)
			return SendError(c, err, fiber.StatusInternalServerError)
		}
		msg := "Board created successfully"
		log.InfoLog.Println(msg)

		return SendSuccessResponse(c, msg, boardModel)
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
			log.ErrorLog.Printf("Error parsing board id: %v\n", err)
			return SendError(c, err, fiber.StatusBadRequest)
		}

		board, err := boardService.GetBoardByID(c.Context(), uint(id))
		if err != nil {
			log.ErrorLog.Printf("Error getting board: %v\n", err)
			return SendError(c, err, fiber.StatusInternalServerError)
		}

		if board == nil {
			log.ErrorLog.Printf("Error getting board: %v\n", ErrBoardNotFound)
			return SendError(c, ErrBoardNotFound, fiber.StatusNotFound)
		}

		log.InfoLog.Println("Board loaded successfully")
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
			log.ErrorLog.Printf("Error parsing board id: %v\n", err)
			return SendError(c, err, fiber.StatusBadRequest)
		}
		var input UpdateBoardRequest

		if err = c.BodyParser(&input); err != nil {
			log.ErrorLog.Printf("Error parsing board update request body: %v\n", err)
			return SendError(c, err, fiber.StatusBadRequest)
		}

		err = validate.Struct(input)
		if err != nil {
			log.ErrorLog.Printf("Error validating board update request body: %v\n", err)
			return SendError(c, err, fiber.StatusBadRequest)
		}

		boardModel := domains.Board{
			ID:        uint(id),
			Name:      input.Name,
			IsPrivate: input.IsPrivate,
		}

		err = boardService.UpdateBoard(c.Context(), &boardModel)
		if err != nil {
			log.ErrorLog.Printf("Error updating board: %v\n", err)
			return SendError(c, err, fiber.StatusInternalServerError)
		}
		msg := "Board updated successfully"
		log.InfoLog.Println(msg)

		return SendSuccessResponse(c, msg, boardModel)
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
			log.ErrorLog.Printf("Error parsing board id: %v\n", err)
			return SendError(c, err, fiber.StatusBadRequest)
		}

		err = boardService.DeleteBoard(c.Context(), uint(id))
		if err != nil {
			log.ErrorLog.Printf("Error deleting board: %v\n", err)
			return SendError(c, err, fiber.StatusInternalServerError)
		}
		log.InfoLog.Println("Board deleted successfully")

		return c.SendStatus(fiber.StatusNoContent)
	}
}

type InviteUserRequest struct {
	UserId   uint   `json:"user_id"`
	RoleName string `json:"role_name"`
}

// InviteUserToBoard invite a user to the board
// @Summary Invite User to Board
// @Description invites a user to board
// @Tags Board
// @Accept  json
// @Produce json
// @Param   id      path     string  true  "Board ID"
// @Param   body  body      InviteUserRequest  true  "Create Board"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /boards/{id}/add-user [post]
// @Security ApiKeyAuth
func InviteUserToBoard(boardService *services.BoardService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		boardId, err := strconv.ParseUint(c.Params("id"), 10, 32)
		if err != nil {
			log.ErrorLog.Printf("Error parsing board id: %v\n", err)
			return SendError(c, err, fiber.StatusBadRequest)
		}

		validate := validation.NewValidator()
		var input InviteUserRequest

		if err = c.BodyParser(&input); err != nil {
			log.ErrorLog.Printf("Error parsing user invitation request body: %v\n", err)
			return SendError(c, err, fiber.StatusBadRequest)
		}

		if err = validate.Struct(input); err != nil {
			log.ErrorLog.Printf("Error validating user invitation request body: %v\n", err)
			return SendError(c, err, fiber.StatusBadRequest)
		}

		if err = boardService.InviteUserToBoard(c.Context(), input.UserId, uint(boardId), input.RoleName); err != nil {
			log.ErrorLog.Printf("Error inviting user: %v\n", err)
			return SendError(c, err, fiber.StatusInternalServerError)
		}
		msg := "User invited successfully"
		log.InfoLog.Println(msg)

		return SendSuccessResponse(c, msg, nil)
	}
}

// RemoveUserFromBoard remove user from the board
// @Summary Remove User from Board
// @Description removes a user from board
// @Tags Board
// @Accept  json
// @Produce json
// @Param   board_id      path     string  true  "Board ID"
// @Param   user_id      path     string  true  "User ID"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /boards/{board_id}/users/{user_id} [delete]
// @Security ApiKeyAuth
func RemoveUserFromBoard(boardService *services.BoardService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		boardId, err := strconv.ParseUint(c.Params("board_id"), 10, 32)
		if err != nil {
			log.ErrorLog.Printf("Error parsing board id: %v\n", err)
			return SendError(c, err, fiber.StatusBadRequest)
		}

		userId, err := strconv.ParseUint(c.Params("user_id"), 10, 32)
		if err != nil {
			log.ErrorLog.Printf("Error parsing user id: %v\n", err)
			return SendError(c, err, fiber.StatusBadRequest)
		}

		if err = boardService.RemoveUserFromBoard(c.Context(), uint(userId), uint(boardId)); err != nil {
			log.ErrorLog.Printf("Error removing user from board: %v\n", err)
			return SendError(c, err, fiber.StatusInternalServerError)
		}
		msg := "User removed from board successfully"
		log.InfoLog.Println(msg)

		return SendSuccessResponse(c, msg, nil)
	}
}

type ChangeUserRoleRequest struct {
	RoleName string `json:"role_name"`
}

// ChangeUserRoleInBoard change user role in board
// @Summary Change User Role in Board
// @Description changes user role in board
// @Tags Board
// @Accept  json
// @Produce json
// @Param   board_id      path     string  true  "Board ID"
// @Param   user_id      path     string  true  "User ID"
// @Param   body  body      ChangeUserRoleRequest  true  "Create Board"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /boards/{board_id}/users/{user_id}  [put]
// @Security ApiKeyAuth
func ChangeUserRoleInBoard(boardService *services.BoardService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		boardId, err := strconv.ParseUint(c.Params("board_id"), 10, 32)
		if err != nil {
			log.ErrorLog.Printf("Error parsing board id: %v\n", err)
			return SendError(c, err, fiber.StatusBadRequest)
		}

		userId, err := strconv.ParseUint(c.Params("user_id"), 10, 32)
		if err != nil {
			log.ErrorLog.Printf("Error parsing user id: %v\n", err)
			return SendError(c, err, fiber.StatusBadRequest)
		}

		validate := validation.NewValidator()
		var input ChangeUserRoleRequest

		if err = c.BodyParser(&input); err != nil {
			log.ErrorLog.Printf("Error parsing user role change request body: %v\n", err)
			return SendError(c, err, fiber.StatusBadRequest)
		}

		if err = validate.Struct(input); err != nil {
			log.ErrorLog.Printf("Error validating user role change request body: %v\n", err)
			return SendError(c, err, fiber.StatusBadRequest)
		}

		if err = boardService.ChangeUserRole(c.Context(), uint(userId), uint(boardId), input.RoleName); err != nil {
			log.ErrorLog.Printf("Error changeing user role: %v\n", err)
			return SendError(c, err, fiber.StatusInternalServerError)
		}
		msg := "User role changed successfully"
		log.InfoLog.Println(msg)

		return SendSuccessResponse(c, msg, nil)
	}
}
