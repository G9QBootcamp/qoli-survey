package handler

import (
	"net/http"
	"strconv"

	"github.com/G9QBootcamp/qoli-survey/internal/config"
	"github.com/G9QBootcamp/qoli-survey/internal/db"
	"github.com/G9QBootcamp/qoli-survey/internal/user/dto"
	"github.com/G9QBootcamp/qoli-survey/internal/user/repository"
	"github.com/G9QBootcamp/qoli-survey/internal/user/service"
	"github.com/G9QBootcamp/qoli-survey/pkg/logging"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	conf    *config.Config
	db      db.DbService
	service service.IUserService
	logger  logging.Logger
}

func NewHandler(conf *config.Config, db db.DbService, logger logging.Logger) *UserHandler {
	return &UserHandler{conf: conf, db: db, service: service.New(conf, repository.NewUserRepository(db, logger), logger), logger: logger}
}

func (h *UserHandler) GetUsers(c echo.Context) error {
	users, err := h.service.GetUsers(c.Request().Context(), dto.UserGetRequest{Name: "aa", Page: 1})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})

	}
	return c.JSON(http.StatusOK, users)
}

func (h *UserHandler) Login(c echo.Context) error {
	var req dto.LoginRequest

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	token, expiresAt, err := h.service.Login(c.Request().Context(), req)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	return c.JSON(http.StatusOK, dto.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
	})
}
func (h *UserHandler) UpdateUserProfile(c echo.Context) error {
	userID, ok := c.Get("userID").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "userID not found"})
	}

	var req dto.UpdateUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid input"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "validation failed"})
	}

	updatedUser, err := h.service.UpdateUserProfile(c.Request().Context(), userID, req)
	if err != nil {
		if err.Error() == "date of birth cannot be updated after 24 hours of registration" {
			return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Unable to update profile"})
	}

	return c.JSON(http.StatusOK, updatedUser)
}
func (h *UserHandler) GetProfile(c echo.Context) error {
	userID, ok := c.Get("userID").(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "userID not found"})
	}

	response, err := h.service.GetUser(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "error in get user profile"})
	}
	return c.JSON(http.StatusOK, response)
}

func (h *UserHandler) RestrictUserSurveys(c echo.Context) error {
	id := c.Param("user_id")
	var req struct {
		MaxSurveys int `json:"max_surveys" validate:"required,min=0"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, map[string]string{"error": "validation failed"})
	}

	if err := h.service.SetMaxSurveys(c.Request().Context(), id, req.MaxSurveys); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Max surveys updated successfully"})
}

func (h *UserHandler) Deposit(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	amount, _ := strconv.ParseFloat(c.FormValue("amount"), 64)

	err := h.service.Deposit(c.Request().Context(), userID, amount)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "deposit successful"})
}

func (h *UserHandler) Withdraw(c echo.Context) error {
	userID := c.Get("user_id").(uint)
	amount, _ := strconv.ParseFloat(c.FormValue("amount"), 64)

	err := h.service.Withdraw(c.Request().Context(), userID, amount)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "withdraw successful"})
}

func (h *UserHandler) Transfer(c echo.Context) error {
	senderID := c.Get("user_id").(uint) // Authenticated user
	receiverID, _ := strconv.Atoi(c.Param("user_id"))
	amount, _ := strconv.ParseFloat(c.FormValue("amount"), 64)

	err := h.service.Transfer(c.Request().Context(), senderID, uint(receiverID), amount)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "transfer successful"})
}

func (h *UserHandler) BuyVote(c echo.Context) error {
	buyerID := c.Get("user_id").(uint)
	sellerID, _ := strconv.Atoi(c.Param("seller_id"))
	amount, _ := strconv.ParseFloat(c.FormValue("amount"), 64)

	err := h.service.BuyVote(c.Request().Context(), buyerID, uint(sellerID), amount)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "vote purchase successful"})
}

func (h *UserHandler) SellVote(c echo.Context) error {
	sellerID := c.Get("user_id").(uint)
	buyerID, _ := strconv.Atoi(c.Param("buyer_id"))
	amount, _ := strconv.ParseFloat(c.FormValue("amount"), 64)

	err := h.service.SellVote(c.Request().Context(), sellerID, uint(buyerID), amount)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "vote sold successfully"})
}

func (h *UserHandler) GetBalance(c echo.Context) error {
	userID, ok := c.Get("userID").(uint)
	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "user not found"})
	}

	balance, err := h.service.GetBalance(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]float64{"balance": balance})
}
