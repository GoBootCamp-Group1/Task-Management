package handlers

import (
	"errors"
	"fmt"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/services"
	"github.com/GoBootCamp-Group1/Task-Management/pkg/validation"
	"github.com/gofiber/fiber/v2"
	"time"
)

type SignUpInput struct {
	Email    string `json:"email" validate:"required,email,excludesall=;" example:"test@example.com"`
	Name     string `json:"name" validate:"required,min=3,max=20,excludesall=;" example:"test"`
	Password string `json:"password" validate:"required,min=8,excludesall=;,password" example:"1234Test@"`
}

// SignUpUser handles the registration of a new user
// @Summary User registration
// @Description Register a user with email, name and password
// @Tags SignUp
// @Accept  json
// @Produce json
// @Param   body  body      SignUpInput  true  "User Registration"
// @Success 201
// @Failure 400
// @Failure 500
// @Router /signup [post]
func SignUpUser(userService *services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		validate := validation.NewValidator()

		var input SignUpInput

		if err := c.BodyParser(&input); err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}

		err := validate.Struct(input)
		if err != nil {
			fmt.Printf("%+v\n", err)
			return SendError(c, err, fiber.StatusBadRequest)
		}

		userModel := domains.User{
			Email:    input.Email,
			Name:     input.Name,
			Password: input.Password,
		}

		err = userService.CreateUser(c.Context(), &userModel)
		if err != nil {
			return SendError(c, err, fiber.StatusInternalServerError)
		}

		return SendSuccessResponse(c, "user")
	}
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email,excludesall=;" example:"test1@test.com"`
	Password string `json:"password" validate:"required,min=8,excludesall=;,password" example:"1234Test@"`
}

// LoginUser handles the login of a user
// @Summary User login
// @Description login user with email and password
// @Tags Login
// @Accept  json
// @Produce json
// @Param   body  body      LoginInput  true  "User Login"
// @Success 200
// @Failure 400
// @Failure 500
// @Router /login [post]
func LoginUser(authService *services.AuthService) fiber.Handler {
	validate := validation.NewValidator()

	return func(c *fiber.Ctx) error {
		var input LoginInput

		c.Cookie(&fiber.Cookie{
			Name:        "X-Session-ID",
			Value:       fmt.Sprint(time.Now().UnixNano()),
			HTTPOnly:    true,
			SessionOnly: true,
		})

		if err := c.BodyParser(&input); err != nil {
			return SendError(c, err, fiber.StatusBadRequest)
		}

		err := validate.Struct(input)
		if err != nil {
			fmt.Printf("%+v\n", err)
			return SendError(c, err, fiber.StatusBadRequest)
		}

		authToken, err := authService.Login(c.Context(), input.Email, input.Password)
		if err != nil {
			return SendError(c, err, fiber.StatusInternalServerError)
		}

		return SendUserToken(c, authToken)
	}
}

func RefreshCreds(authService *services.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		refToken := c.GetReqHeaders()["Authorization"]
		if len(refToken[0]) == 0 {
			return SendError(c, errors.New("token should be provided"), fiber.StatusBadRequest)
		}

		authToken, err := authService.RefreshAuth(c.UserContext(), refToken[0])
		if err != nil {
			return SendError(c, err, fiber.StatusUnauthorized)
		}

		return SendUserToken(c, authToken)
	}
}
