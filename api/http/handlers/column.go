package handlers

import (
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/services"
	"github.com/GoBootCamp-Group1/Task-Management/pkg/log"
	"github.com/GoBootCamp-Group1/Task-Management/pkg/utils"
	"github.com/GoBootCamp-Group1/Task-Management/pkg/validation"
	"github.com/gofiber/fiber/v2"
)

var (
	ErrColumnNotFound = fiber.NewError(fiber.StatusNotFound, "Column not found")
)

type UpdateColumnRequest struct {
	Name string `json:"name" validate:"required,min=3,max=20,excludesall=;" example:"new column"`
}

type MoveColumnRequest struct {
	OrderPosition int `json:"position" validate:"required,min=0" example:"2"`
}
type CreateColumnRequest struct {
	Name    string `json:"name" validate:"required,min=3,max=20,excludesall=;" example:"new column"`
	IsFinal bool   `json:"is_final" example:"false"`
}

// CreateColumn creates a column
// @Summary Create Column
// @Description create a column
// @Tags Column
// @Accept json
// @Produce json
// @Param   body      body     CreateColumnRequest  true  "Create Column"
// @Param   boardId      path     string  true  "Board ID"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /boards/{boardId}/columns [post]
// @Security ApiKeyAuth
func CreateColumn(columnService *services.ColumnService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		boardId, errParam := c.ParamsInt("boardId")
		if errParam != nil {
			log.ErrorLog.Printf("Error parsing board id: %v\n", errParam)
			// todo: sending error should handle better
			return SendError(c, errParam)
		}
		validate := validation.NewValidator()
		var input CreateColumnRequest

		if err := c.BodyParser(&input); err != nil {
			log.ErrorLog.Printf("Error parsing column creation request body: %v\n", err)
			return SendError(c, err)
		}

		if err := validate.Struct(input); err != nil {
			log.ErrorLog.Printf("Error validating column creation request body: %v\n", err)
			return SendError(c, err)
		}

		userId, err := utils.GetUserID(c)
		if err != nil {
			log.ErrorLog.Printf("Error loading user: %v\n", err)
			return SendError(c, err)
		}

		columnModel := domains.Column{
			CreatedBy: userId,
			Name:      input.Name,
			IsFinal:   input.IsFinal,
			BoardID:   uint(boardId),
		}

		if err = columnService.CreateColumn(c.Context(), &columnModel); err != nil {
			log.ErrorLog.Printf("Error creating column: %v\n", err)
			return SendError(c, err)
		}
		MsgColumnCreation := "Column created successfully"
		log.InfoLog.Println(MsgColumnCreation)

		return SendSuccessResponse(
			c,
			MsgColumnCreation,
			columnModel)
	}
}

// GetColumnByID get a column
// @Summary Get Column
// @Description gets a column
// @Tags Column
// @Produce json
// @Param   id      path     string  true  "Column ID"
// @Param   boardId      path     string  true  "Board ID"
// @Success 200 {object} Response
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /boards/{boardId}/columns/{id} [get]
// @Security ApiKeyAuth
func GetColumnByID(columnService *services.ColumnService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, errParam := c.ParamsInt("id")
		if errParam != nil {
			log.ErrorLog.Printf("Error parsing column id: %v\n", errParam)
			return SendError(c, errParam)
		}

		userId, err := utils.GetUserID(c)
		if err != nil {
			log.ErrorLog.Printf("Error loading user: %v\n", err)
			return SendError(c, err)
		}

		column, err := columnService.GetColumnById(c.Context(), userId, uint(id))
		if err != nil {
			log.ErrorLog.Printf("Error getting column: %v\n", err)
			return SendError(c, err)
		}

		if column == nil {
			log.ErrorLog.Printf("Error getting column: %v\n", err)
			return SendError(c, ErrColumnNotFound)
		}
		msg := "Column loaded successfully"
		log.InfoLog.Println(msg)

		return SendSuccessResponse(c, msg, column)
	}
}

// GetAllColumns get all columns of a board
// @Summary Get Columns
// @Description gets all columns of a bard
// @Tags Column
// @Produce json
// @Param   boardId      path     string  true  "Board ID"
// @Success 200 {object} Response
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /boards/{boardId}/columns [get]
// @Security ApiKeyAuth
func GetAllColumns(columnService *services.ColumnService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		boardId, errParam := c.ParamsInt("boardId")
		if errParam != nil {
			log.ErrorLog.Printf("Error parsing board id: %v\n", errParam)
			return SendError(c, &fiber.Error{Code: 400, Message: "Board Id is not valid or didn't send"})
		}

		page, pageSize := PageAndPageSize(c)

		userId, err := utils.GetUserID(c)
		if err != nil {
			log.ErrorLog.Printf("Error loading user: %v\n", err)
			return SendError(c, err)
		}

		columns, err := columnService.GetAllColumns(c.Context(), userId, uint(boardId), page, pageSize)
		if err != nil {
			log.ErrorLog.Printf("Error getting all columns: %v\n", err)
			return SendError(c, err)
		}
		msg := "Columns loaded successfully"
		log.InfoLog.Println(msg)

		return SendSuccessPaginateResponse(c, "Columns loaded successfully", columns.Data, uint(columns.Page), uint(columns.PageSize), uint(columns.Total))
	}
}

// UpdateColumn update a column
// @Summary Update Column
// @Description update a column
// @Tags Column
// @Accept json
// @Produce json
// @Param   id      path     string  true  "Column ID"
// @Param   boardId      path     string  true  "Board ID"
// @Param   name      body     UpdateColumnRequest  true  "New Column name"
// @Success 200 {object} Response
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /boards/{boardId}/columns/{id} [put]
// @Security ApiKeyAuth
func UpdateColumn(columnService *services.ColumnService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		validate := validation.NewValidator()
		id, errParam := c.ParamsInt("id")
		if errParam != nil {
			log.ErrorLog.Printf("Error parsing column id: %v\n", errParam)
			return SendError(c, errParam)
		}

		boardId, errParam := c.ParamsInt("boardId")
		if errParam != nil {
			log.ErrorLog.Printf("Error parsing board id: %v\n", errParam)
			return SendError(c, errParam)
		}

		var input UpdateColumnRequest

		if err := c.BodyParser(&input); err != nil {
			log.ErrorLog.Printf("Error parsing column update request body: %v\n", err)
			return SendError(c, err)
		}

		if err := validate.Struct(input); err != nil {
			log.ErrorLog.Printf("Error validating column creation request body: %v\n", err)
			return SendError(c, err)
		}

		column := domains.ColumnUpdate{
			ID:   uint(id),
			Name: input.Name,
		}

		userId, err := utils.GetUserID(c)
		if err != nil {
			log.ErrorLog.Printf("Error loading user: %v\n", err)
			return SendError(c, err)
		}

		if err := columnService.Update(c.Context(), uint(boardId), userId, &column); err != nil {
			log.ErrorLog.Printf("Error updating column: %v\n", err)
			return SendError(c, err)
		}
		MsgColumnUpdate := "Column updated successfully"
		log.InfoLog.Println(MsgColumnUpdate)

		return SendSuccessResponse(
			c,
			MsgColumnUpdate,
			input)
	}
}

// MoveColumn move a column
// @Summary Move Column
// @Description move a column
// @Tags Column
// @Accept json
// @Produce json
// @Param   id      path     string  true  "Column ID"
// @Param   boardId      path     string  true  "Board ID"
// @Param   body     body     MoveColumnRequest  true  "new position of column [id]"
// @Success 200 {object} Response
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /boards/{boardId}/columns/{id}/move [put]
// @Security ApiKeyAuth
func MoveColumn(columnService *services.ColumnService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		validate := validation.NewValidator()
		id, errParam := c.ParamsInt("id")
		if errParam != nil {
			log.ErrorLog.Printf("Error parsing column id: %v\n", errParam)
			return SendError(c, errParam)
		}

		boardId, errParam := c.ParamsInt("boardId")
		if errParam != nil {
			log.ErrorLog.Printf("Error parsing board id: %v\n", errParam)
			return SendError(c, errParam)
		}

		userId, err := utils.GetUserID(c)
		if err != nil {
			log.ErrorLog.Printf("Error loading user: %v\n", err)
			return SendError(c, err)
		}

		var input MoveColumnRequest

		if err := c.BodyParser(&input); err != nil {
			log.ErrorLog.Printf("Error parsing column move request body: %v\n", err)
			return SendError(c, err)
		}

		if err := validate.Struct(input); err != nil {
			log.ErrorLog.Printf("Error validating column move request body: %v\n", err)
			return SendError(c, err)
		}

		column := domains.ColumnMove{ID: uint(id), OrderPosition: input.OrderPosition}

		if err := columnService.Move(c.Context(), uint(boardId), userId, &column); err != nil {
			log.ErrorLog.Printf("Error moving column: %v\n", err)
			return SendError(c, err)
		}
		MsgColumnMove := "Column moved successfully"
		log.InfoLog.Println(MsgColumnMove)

		return SendSuccessResponse(
			c,
			MsgColumnMove,
			input)
	}
}

// ChangeFinalColumn make a column as a final column
// @Summary Final Column
// @Description make a column as a final column
// @Tags Column
// @Accept json
// @Produce json
// @Param   id      path     string  true  "Column ID"
// @Param   boardId      path     string  true  "Board ID"
// @Success 200 {object} Response
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /boards/{boardId}/columns/{id}/final [put]
// @Security ApiKeyAuth
func ChangeFinalColumn(columnService *services.ColumnService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, errParam := c.ParamsInt("id")
		if errParam != nil {
			log.ErrorLog.Printf("Error parsing column id: %v\n", errParam)
			return SendError(c, errParam)
		}

		boardId, errParam := c.ParamsInt("boardId")
		if errParam != nil {
			log.ErrorLog.Printf("Error parsing board id: %v\n", errParam)
			return SendError(c, errParam)
		}

		userId, err := utils.GetUserID(c)
		if err != nil {
			log.ErrorLog.Printf("Error loading user: %v\n", err)
			return SendError(c, err)
		}

		if err := columnService.Final(c.Context(), uint(boardId), userId, uint(id)); err != nil {
			log.ErrorLog.Printf("Error changing final column: %v\n", err)
			return SendError(c, err)
		}
		MsgFinalColumnChange := "Final column changed successfully"
		log.InfoLog.Println(MsgFinalColumnChange)

		return SendSuccessResponse(
			c,
			MsgFinalColumnChange,
			id)
	}
}

// DeleteColumn delete a column
// @Summary Delete Column
// @Description delete a column
// @Tags Column
// @Accept json
// @Produce json
// @Param   id      path     string  true  "Column ID"
// @Param   boardId      path     string  true  "Board ID"
// @Success 200
// @Failure 404
// @Failure 500
// @Router /boards/{boardId}/columns/{id} [delete]
// @Security ApiKeyAuth
func DeleteColumn(columnService *services.ColumnService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, errParam := c.ParamsInt("id")
		if errParam != nil {
			log.ErrorLog.Printf("Error parsing column id: %v\n", errParam)
			return SendError(c, errParam)
		}

		boardId, errParam := c.ParamsInt("boardId")
		if errParam != nil {
			log.ErrorLog.Printf("Error parsing board id: %v\n", errParam)
			return SendError(c, errParam)
		}

		userId, err := utils.GetUserID(c)
		if err != nil {
			log.ErrorLog.Printf("Error loading user: %v\n", err)
			return SendError(c, err)
		}

		if err := columnService.Delete(c.Context(), uint(boardId), userId, uint(id)); err != nil {
			log.ErrorLog.Printf("Error deleting column: %v\n", err)
			return SendError(c, err)
		}
		MsgColumnDelete := "Column deleted successfully"
		log.InfoLog.Println(MsgColumnDelete)

		return SendSuccessResponse(
			c,
			MsgColumnDelete,
			id)
	}
}
