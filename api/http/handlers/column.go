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
		log.InfoLog.Println("Column created successfully")

		// todo: sending response should handle better
		return SendSuccessResponse(c, "column")
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

func GetAllColumns(columnService *services.ColumnService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		boardId, errParam := c.ParamsInt("boardId")
		if errParam != nil {
			log.ErrorLog.Printf("Error parsing board id: %v\n", errParam)
			return SendError(c, errParam, fiber.StatusBadRequest)
		}

		columns, err := columnService.GetAllColumns(c.Context(), uint(boardId))
		if err != nil {
			log.ErrorLog.Printf("Error getting all columns: %v\n", err)
			return SendError(c, err, 0)
		}
		log.InfoLog.Println("Columns loaded successfully")

		return c.JSON(columns)
	}
}

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
		log.InfoLog.Println("Column updated successfully")

		return SendSuccessResponse(c, "column")
	}
}

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
		log.InfoLog.Println("Column moved successfully")

		return SendSuccessResponse(c, "column")
	}
}

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
		log.InfoLog.Println("Final column changed successfully")

		return SendSuccessResponse(c, "column")
	}
}

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
		log.InfoLog.Println("Column deleted successfully")

		return SendSuccessResponse(c, "column")
	}
}
