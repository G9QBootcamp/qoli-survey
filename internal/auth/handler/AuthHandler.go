package handler

import (
	"net/http"

	"github.com/G9QBootcamp/qoli-survey/internal/auth/dto"
	"github.com/G9QBootcamp/qoli-survey/internal/auth/repository"
	"github.com/G9QBootcamp/qoli-survey/internal/auth/service"
	"github.com/G9QBootcamp/qoli-survey/internal/config"
	"github.com/G9QBootcamp/qoli-survey/internal/db"
	userRepository "github.com/G9QBootcamp/qoli-survey/internal/user/repository"
	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
	"github.com/labstack/echo/v4"
)

type AuthHandler struct {
	conf    *config.Config
	db      db.DbService
	service service.IAuthService
	logger  logging.Logger
}

func NewHandler(conf *config.Config, db db.DbService, logger logging.Logger) *AuthHandler {
	return &AuthHandler{conf: conf, db: db,
		service: service.New(conf,
			repository.NewAuthRepository(db, logger),
			userRepository.NewUserRepository(db, logger),
			logger),
		logger: logger}
}

func (h *AuthHandler) GenerateOTP(c echo.Context) error {
	userID, ok := c.Get("userID").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "userID not found"})
	}

	otp, err := h.service.SaveOTP(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to save OTP"})
	}

	err = h.service.SendOTPEmail(c.Request().Context(), userID, otp)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to send OTP email"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "OTP sent successfully"})
}

func (h *AuthHandler) VerifyOTP(c echo.Context) error {
	userID, ok := c.Get("userID").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "userID not found"})
	}

	var req dto.VerifyOTPRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	valid, err := h.service.VerifyOTP(c.Request().Context(), userID, req.OTP)
	if err != nil || !valid {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid or expired OTP"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "OTP verified successfully"})
}
