package handlers

import (
	"fmt"
	"time"

	"github.com/GoBootCamp-Group1/Task-Management/internal/core/domains"
	"github.com/GoBootCamp-Group1/Task-Management/internal/core/services"
	"github.com/GoBootCamp-Group1/Task-Management/pkg/log"
	"github.com/GoBootCamp-Group1/Task-Management/pkg/validation"
	"github.com/gofiber/fiber/v2"
)

var (
	ErrRefreshTokenNotProvided = &fiber.Error{Code: fiber.StatusBadRequest, Message: "Refresh token not provided"}
)

type SignUpInput struct {
	Email    string `json:"email" validate:"required,email,excludesall=;" example:"test@example.com"`
	Name     string `json:"name" validate:"required,min=3,max=20,excludesall=;" example:"test"`
	Password string `json:"password" validate:"required,min=8,excludesall=;,password" example:"1234Test@"`
}

// SignUpUser handles the registration of a new user
// @Summary User registration
// @Description Register a user with email, name and password
// @Tags Authentication
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
			log.ErrorLog.Printf("Error parsing user sign-up request body: %v\n", err)
			return SendError(c, err)
		}

		err := validate.Struct(input)
		if err != nil {
			log.ErrorLog.Printf("Error validating user sign-up request body: %v\n", err)
			return SendError(c, err)
		}

		userModel := domains.User{
			Email:    input.Email,
			Name:     input.Name,
			Password: input.Password,
		}

		err = userService.CreateUser(c.Context(), &userModel)
		if err != nil {
			log.ErrorLog.Printf("Error creating user: %v\n", err)
			return SendError(c, err)
		}
		msg := "User signed up (created) successfully"
		log.InfoLog.Println(msg)

		return SendSuccessResponse(c, msg, userModel)
	}
}

type LoginInput struct {
	Email    string `json:"email" validate:"required,email,excludesall=;" example:"test1@test.com"`
	Password string `json:"password" validate:"required,min=8,excludesall=;,password" example:"1234Test@"`
}

// LoginUser handles the login of a user
// @Summary User login
// @Description login user with email and password
// @Tags Authentication
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
			log.ErrorLog.Printf("Error parsing : %v\n", err)
			return SendError(c, &fiber.Error{Code: fiber.StatusBadRequest, Message: "Error parsing user login request body"})
		}

		err := validate.Struct(input)
		if err != nil {
			log.ErrorLog.Printf("Error vaidating user login request body: %v\n", err)
			return SendError(c, &fiber.Error{Code: fiber.StatusBadRequest, Message: "Error validating user login request body"})
		}

		authToken, err := authService.Login(c.Context(), input.Email, input.Password)
		if err != nil {
			log.ErrorLog.Printf("Error logging in user: %v\n", err)
			return SendError(c, err)
		}

		log.InfoLog.Println("User logged in successfully")
		return SendUserToken(c, authToken)
	}
}

func RefreshCreds(authService *services.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		refToken := c.GetReqHeaders()["Authorization"]
		if len(refToken[0]) == 0 {
			log.ErrorLog.Printf("Error refreshing token: %v\n", ErrRefreshTokenNotProvided)
			return SendError(c, ErrRefreshTokenNotProvided)
		}

		authToken, err := authService.RefreshAuth(c.UserContext(), refToken[0])
		if err != nil {
			log.ErrorLog.Printf("Error refreshing token, not authorized: %v\n", err)
			return SendError(c, err)
		}
		log.InfoLog.Println("Token refreshed successfully")

		return SendUserToken(c, authToken)
	}
}
