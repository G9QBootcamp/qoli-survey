package handler

import (
	"net/http"

	"github.com/G9QBootcamp/qoli-survey/internal/auth/service"
	"github.com/G9QBootcamp/qoli-survey/internal/user/dto"
	"github.com/G9QBootcamp/qoli-survey/internal/user/models"
	"github.com/G9QBootcamp/qoli-survey/internal/user/repository"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func Signup(repo repository.UserRepository, jwtService service.JWTService) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Bind the request to UserCreateRequest DTO
		req := new(dto.UserCreateRequest)
		if err := c.Bind(req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		}

		// Validate the input
		if err := c.Validate(req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
		}

		// Create the user object
		user := models.User{
			NationalID:    req.NationalID,
			Email:         req.Email,
			PasswordHash:  string(hashedPassword),
			WalletBalance: 0, // Default value
		}

		// Save the user to the database
		if err := repo.CreateUser(&user); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
		}

		// Generate JWT token
		token, err := jwtService.GenerateToken(user.ID, user.Email, user.Email)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
		}

		// Prepare the response
		response := dto.UserResponse{
			ID:          user.ID,
			NationalID:  user.NationalID,
			Email:       user.Email,
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			City:        user.City,
			DateOfBirth: user.DateOfBirth,
		}

		// Respond with the user data and token
		return c.JSON(http.StatusOK, map[string]interface{}{
			"user":  response,
			"token": token,
		})
	}
}
