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
	userID, ok := c.Get("userID").(uint)
	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "user not found"})
	}
	var req struct {
		Amount float64 `json:"amount" validate:"required,min=0"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	err := h.service.Deposit(c.Request().Context(), userID, req.Amount)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "deposit successful"})
}

func (h *UserHandler) Withdraw(c echo.Context) error {
	userID, ok := c.Get("userID").(uint)
	if !ok || userID == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "user not found"})
	}
	var req struct {
		Amount float64 `json:"amount" validate:"required,min=0"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	err := h.service.Withdraw(c.Request().Context(), userID, req.Amount)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "withdraw successful"})
}

func (h *UserHandler) Transfer(c echo.Context) error {
	senderID, ok := c.Get("userID").(uint)
	if !ok || senderID == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "user not found"})
	}
	receiverID, _ := strconv.Atoi(c.Param("user_id"))
	var req struct {
		Amount float64 `json:"amount" validate:"required,min=0"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	err := h.service.Transfer(c.Request().Context(), senderID, uint(receiverID), req.Amount)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "transfer successful"})
}

func (h *UserHandler) BuyVote(c echo.Context) error {
	buyerID, ok := c.Get("userID").(uint)
	if !ok || buyerID == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "user not found"})
	}

	var req struct {
		Amount float64 `json:"amount" validate:"required,min=0"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	sellerID, err1 := strconv.Atoi(c.Param("seller_id"))
	voteID, err2 := strconv.Atoi(c.Param("vote_id"))

	if err1 != nil || err2 != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request parameters"})
	}

	voterId, err := h.service.GetVoterID(c.Request().Context(), uint(voteID))
	if err != nil || voterId != uint(sellerID) {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "user is not voter of the vote"})
	}

	err = h.service.BuyVote(c.Request().Context(), buyerID, uint(sellerID), uint(voteID), req.Amount)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "vote purchase successful"})
}

func (h *UserHandler) SellVote(c echo.Context) error {
	sellerID, ok := c.Get("userID").(uint)
	if !ok || sellerID == 0 {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "user not found"})
	}

	var req struct {
		Amount float64 `json:"amount" validate:"required,min=0"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	buyerID, err1 := strconv.Atoi(c.Param("buyer_id"))
	voteID, err2 := strconv.Atoi(c.Param("vote_id"))

	if err1 != nil || err2 != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request parameters"})
	}

	voterId, err := h.service.GetVoterID(c.Request().Context(), uint(voteID))
	if err != nil || voterId != sellerID {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "user is not voter of the vote"})
	}

	err = h.service.SellVote(c.Request().Context(), sellerID, uint(buyerID), uint(voteID), req.Amount)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to sell votes"})
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
