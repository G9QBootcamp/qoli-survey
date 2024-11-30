package server

import (
	"log"
	"net/http"

	"github.com/G9QBootcamp/qoli-survey/pkg/validation"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type Server struct {
	Echo *echo.Echo
}

func NewHttpServer() *Server {
	e := echo.New()

	validate := validator.New()
	validation.RegisterCustomValidation(validate)
	e.Validator = &validation.CustomValidator{Validator: validate}

	return &Server{Echo: e}
}

func (s *Server) Start(address string) {
	log.Printf("Starting server on %s\n", address)
	if err := s.Echo.Start(address); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Shutting down the server: %v", err)
	}
}
