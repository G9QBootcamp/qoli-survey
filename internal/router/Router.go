package router

import (
	"github.com/G9QBootcamp/qoli-survey/internal/config"
	"github.com/G9QBootcamp/qoli-survey/internal/db"
	"github.com/G9QBootcamp/qoli-survey/internal/user/handler"
	"github.com/G9QBootcamp/qoli-survey/internal/user/repository"
	"github.com/G9QBootcamp/qoli-survey/internal/user/service"
	"github.com/labstack/echo/v4"
)

func NewUserRouter(conf *config.Config, dbService db.Service, e *echo.Echo) {
	userRepo := repository.NewUserRepository(dbService.GetDb())
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	// Register routes
	e.POST("/signup", userHandler.Signup)
}
