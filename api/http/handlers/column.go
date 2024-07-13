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

type CreateColumnRequest struct {
	Name    string `json:"name" validate:"required,min=3,max=20,excludesall=;" example:"new column"`
	IsFinal bool   `json:"is_final" example:"false"`
}

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
// @Router /columns/{boardId} [post]
// @Security ApiKeyAuth
func CreateColumn(columnService *services.ColumnService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		boardId, errParam := c.ParamsInt("boardId")
		if errParam != nil {
			log.ErrorLog.Printf("Error parsing board id: %v\n", errParam)
			// todo: sending error should handle better
			return SendError(c, errParam, fiber.StatusBadRequest)
		}
		validate := validation.NewValidator()
		var input CreateColumnRequest

		if err := c.BodyParser(&input); err != nil {
			log.ErrorLog.Printf("Error parsing column creation request body: %v\n", err)
			return SendError(c, err, fiber.StatusBadRequest)
		}

		if err := validate.Struct(input); err != nil {
			log.ErrorLog.Printf("Error validating column creation request body: %v\n", err)
			return SendError(c, err, fiber.StatusBadRequest)
		}

		userId, err := utils.GetUserID(c)
		if err != nil {
			log.ErrorLog.Printf("Error loading user: %v\n", err)
			return SendError(c, err, 0)
		}

		columnModel := domains.Column{
			CreatedBy: userId,
			Name:      input.Name,
			IsFinal:   input.IsFinal,
			BoardID:   uint(boardId),
		}

		if err = columnService.CreateColumn(c.Context(), &columnModel); err != nil {
			log.ErrorLog.Printf("Error creating column: %v\n", err)
			return SendError(c, err, 0)
		}
		MsgColumnCreation := "Column created successfully"
		log.InfoLog.Println(MsgColumnCreation)

		// todo: sending response should handle better
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
// @Success 200 {object} domains.Column
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /columns/{id} [get]
// @Security ApiKeyAuth
func GetColumnByID(columnService *services.ColumnService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, errParam := c.ParamsInt("id")
		if errParam != nil {
			log.ErrorLog.Printf("Error parsing column id: %v\n", errParam)
			return SendError(c, errParam, fiber.StatusBadRequest)
		}

		column, err := columnService.GetColumnById(c.Context(), uint(id))
		if err != nil {
			log.ErrorLog.Printf("Error getting column: %v\n", err)
			return SendError(c, err, fiber.StatusInternalServerError)
		}

		if column == nil {
			log.ErrorLog.Printf("Error getting column: %v\n", err)
			return SendError(c, ErrColumnNotFound, fiber.StatusNotFound)
		}
		log.InfoLog.Println("Column loaded successfully")

		return c.JSON(column)
	}
}

// GetAllColumns get all columns of a board
// @Summary Get Columns
// @Description gets all columns of a bard
// @Tags Column
// @Produce json
// @Success 200 {object} domains.Column
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /columns [get]
// @Security ApiKeyAuth
func GetAllColumns(columnService *services.ColumnService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		boardId, errParam := c.ParamsInt("boardId")
		if errParam != nil {
			log.ErrorLog.Printf("Error parsing board id: %v\n", errParam)
			return SendError(c, errParam, fiber.StatusBadRequest)
		}

		page, pageSize := PageAndPageSize(c)

		columns, err := columnService.GetAllColumns(c.Context(), uint(boardId), page, pageSize)
		if err != nil {
			log.ErrorLog.Printf("Error getting all columns: %v\n", err)
			return SendError(c, err, 0)
		}
		log.InfoLog.Println("Columns loaded successfully")

		return c.JSON(columns)
	}
}

// UpdateColumn update a column
// @Summary Update Column
// @Description update a column
// @Tags Column
// @Accept json
// @Produce json
// @Param   id      path     string  true  "Column ID"
// @Param   name      body     string  true  "New Column name"
// @Success 200 {object} domains.Column
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /columns/{id} [put]
// @Security ApiKeyAuth
func UpdateColumn(columnService *services.ColumnService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		validate := validation.NewValidator()
		id, errParam := c.ParamsInt("id")
		if errParam != nil {
			log.ErrorLog.Printf("Error parsing column id: %v\n", errParam)
			return SendError(c, errParam, fiber.StatusBadRequest)
		}

		var input domains.ColumnUpdate

		if err := c.BodyParser(&input); err != nil {
			log.ErrorLog.Printf("Error parsing column update request body: %v\n", err)
			return SendError(c, err, fiber.StatusBadRequest)
		}

		if err := validate.Struct(input); err != nil {
			log.ErrorLog.Printf("Error validating column creation request body: %v\n", err)
			return SendError(c, err, fiber.StatusBadRequest)
		}

		input.ID = uint(id)

		if err := columnService.Update(c.Context(), &input); err != nil {
			log.ErrorLog.Printf("Error updating column: %v\n", err)
			return SendError(c, err, fiber.StatusBadRequest)
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
// @Param   order_position      body     number  true  "new position of column [id]"
// @Success 200 {object} domains.Column
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /columns/{id}/move [put]
// @Security ApiKeyAuth
func MoveColumn(columnService *services.ColumnService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		validate := validation.NewValidator()
		id, errParam := c.ParamsInt("id")
		if errParam != nil {
			log.ErrorLog.Printf("Error parsing column id: %v\n", errParam)
			return SendError(c, errParam, fiber.StatusBadRequest)
		}

		var input domains.ColumnMove

		if err := c.BodyParser(&input); err != nil {
			log.ErrorLog.Printf("Error parsing column move request body: %v\n", err)
			return SendError(c, err, fiber.StatusBadRequest)
		}

		if err := validate.Struct(input); err != nil {
			log.ErrorLog.Printf("Error validating column move request body: %v\n", err)
			return SendError(c, err, fiber.StatusBadRequest)
		}

		input.ID = uint(id)

		if err := columnService.Move(c.Context(), &input); err != nil {
			log.ErrorLog.Printf("Error moving column: %v\n", err)
			return SendError(c, err, fiber.StatusBadRequest)
		}
		MsgColumnMove := "Column moved successfully"
		log.InfoLog.Println(MsgColumnMove)

		return SendSuccessResponse(
			c,
			MsgColumnMove,
			input)
	}
}

// FinalColumn make a column as a final column
// @Summary Final Column
// @Description make a column as a final column
// @Tags Column
// @Accept json
// @Produce json
// @Param   id      path     string  true  "Column ID"
// @Success 200 {object} domains.Column
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /columns/{id}/final [put]
// @Security ApiKeyAuth
func ChangeFinalColumn(columnService *services.ColumnService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, errParam := c.ParamsInt("id")
		if errParam != nil {
			log.ErrorLog.Printf("Error parsing column id: %v\n", errParam)
			return SendError(c, errParam, fiber.StatusBadRequest)
		}

		if err := columnService.Final(c.Context(), uint(id)); err != nil {
			log.ErrorLog.Printf("Error changing final column: %v\n", err)
			return SendError(c, err, fiber.StatusBadRequest)
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
// @Success 200
// @Failure 404
// @Failure 500
// @Router /columns/{id} [delete]
// @Security ApiKeyAuth
func DeleteColumn(columnService *services.ColumnService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, errParam := c.ParamsInt("id")
		if errParam != nil {
			log.ErrorLog.Printf("Error parsing column id: %v\n", errParam)
			return SendError(c, errParam, fiber.StatusBadRequest)
		}

		if err := columnService.Delete(c.Context(), uint(id)); err != nil {
			log.ErrorLog.Printf("Error deleting column: %v\n", err)
			return SendError(c, err, fiber.StatusBadRequest)
		}
		MsgColumnDelete := "Column deleted successfully"
		log.InfoLog.Println(MsgColumnDelete)

		return SendSuccessResponse(
			c,
			MsgColumnDelete,
			id)
	}
}
