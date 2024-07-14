package handlers

import (
	"strconv"

	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/services"
	"github.com/GoBootCamp-Group1/Task-Management/pkg/log"
	"github.com/GoBootCamp-Group1/Task-Management/pkg/validation"
	"github.com/gofiber/fiber/v2"
)

var (
	ErrRoleNotFound = fiber.NewError(fiber.StatusNotFound, "Role not found")
)

type CreateRoleRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=50" example:"new role"`
	Description string `json:"description" validate:"required,max=255" example:"role description"`
	Weight      int    `json:"weight" validate:"required" example:"1"`
}

// CreateRole creates a new role
// @Summary Create Role
// @Description creates a role
// @Tags Role
// @Accept  json
// @Produce json
// @Param   body  body      CreateRoleRequest  true  "Create Role"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /roles [post]
// @Security ApiKeyAuth
func CreateRole(roleService *services.RoleService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		validate := validation.NewValidator()
		var input CreateRoleRequest

		if err := c.BodyParser(&input); err != nil {
			log.ErrorLog.Printf("Error parsing role creation request body: %v\n", err)
			return SendError(c, &fiber.Error{Code: fiber.StatusBadRequest, Message: "Error parsing role creation request body"})
		}

		err := validate.Struct(input)
		if err != nil {
			log.ErrorLog.Printf("Error validating role creation request body: %v\n", err)
			return SendError(c, &fiber.Error{Code: fiber.StatusBadRequest, Message: "Error validating role creation request body"})
		}

		roleModel := domains.Role{
			Name:        input.Name,
			Description: input.Description,
			Weight:      input.Weight,
		}

		err = roleService.CreateRole(c.Context(), &roleModel)
		if err != nil {
			log.ErrorLog.Printf("Error creating role: %v\n", err)
			return SendError(c, err)
		}
		msg := "Role created successfully"
		log.InfoLog.Println(msg)

		return SendSuccessResponse(c, msg, roleModel)
	}
}

// GetRoleByID get a role
// @Summary Get Role
// @Description gets a role
// @Tags Role
// @Produce json
// @Param   id      path     string  true  "Role ID"
// @Success 200 {object} domains.Role
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /roles/{id} [get]
// @Security ApiKeyAuth
func GetRoleByID(roleService *services.RoleService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseUint(c.Params("id"), 10, 32)
		if err != nil {
			log.ErrorLog.Printf("Error parsing role id: %v\n", err)
			return SendError(c, &fiber.Error{Code: fiber.StatusBadRequest, Message: "Error parsing role id"})
		}

		role, err := roleService.GetRoleById(c.Context(), uint(id))
		if err != nil {
			log.ErrorLog.Printf("Error getting role: %v\n", err)
			return SendError(c, err)
		}

		if role == nil {
			log.ErrorLog.Printf("Error getting role: %v\n", ErrRoleNotFound)
			return SendError(c, ErrRoleNotFound)
		}

		log.InfoLog.Println("Role loaded successfully")
		return c.JSON(role)
	}
}

type UpdateRoleRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=50" example:"updated role"`
	Description string `json:"description" validate:"required,max=255" example:"updated description"`
	Weight      int    `json:"weight" validate:"required" example:"1"`
}

// UpdateRole updates an existing role
// @Summary Update Role
// @Description updates a role
// @Tags Role
// @Accept  json
// @Produce json
// @Param   id      path     string  true  "Role ID"
// @Param   body  body      UpdateRoleRequest  true  "Update Role"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /roles/{id} [put]
// @Security ApiKeyAuth
func UpdateRole(roleService *services.RoleService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		validate := validation.NewValidator()
		id, err := strconv.ParseUint(c.Params("id"), 10, 32)
		if err != nil {
			log.ErrorLog.Printf("Error parsing role id: %v\n", err)
			return SendError(c, &fiber.Error{Code: fiber.StatusBadRequest, Message: "Error parsing role id"})
		}
		var input UpdateRoleRequest

		if err = c.BodyParser(&input); err != nil {
			log.ErrorLog.Printf("Error parsing role update request body: %v\n", err)
			return SendError(c, &fiber.Error{Code: fiber.StatusBadRequest, Message: "Error parsing role update request body"})
		}

		err = validate.Struct(input)
		if err != nil {
			log.ErrorLog.Printf("Error validating role update request body: %v\n", err)
			return SendError(c, &fiber.Error{Code: fiber.StatusBadRequest, Message: "Error validating role update request body"})
		}

		roleModel := domains.Role{
			ID:          uint(id),
			Name:        input.Name,
			Description: input.Description,
			Weight:      input.Weight,
		}

		err = roleService.UpdateRole(c.Context(), &roleModel)
		if err != nil {
			log.ErrorLog.Printf("Error updating role: %v\n", err)
			return SendError(c, err)
		}
		msg := "Role updated successfully"
		log.InfoLog.Println(msg)

		return SendSuccessResponse(c, msg, roleModel)
	}
}

// DeleteRole delete a role
// @Summary Delete Role
// @Description deletes a role
// @Tags Role
// @Produce json
// @Param   id      path     string  true  "Role ID"
// @Success 204
// @Failure 400
// @Failure 500
// @Router /roles/{id} [delete]
// @Security ApiKeyAuth
func DeleteRole(roleService *services.RoleService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := strconv.ParseUint(c.Params("id"), 10, 32)
		if err != nil {
			log.ErrorLog.Printf("Error parsing role id: %v\n", err)
			return SendError(c, &fiber.Error{Code: fiber.StatusBadRequest, Message: "Error parsing role id"})
		}

		err = roleService.DeleteRole(c.Context(), uint(id))
		if err != nil {
			log.ErrorLog.Printf("Error deleting role: %v\n", err)
			return SendError(c, err)
		}
		log.InfoLog.Println("Role deleted successfully")

		return c.SendStatus(fiber.StatusNoContent)
	}
}
