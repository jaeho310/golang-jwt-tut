package handler

import (
	"github.com/labstack/echo"
	"server/auth"
)

type AuthHandler struct{}

func (AuthHandler) Init(e *echo.Echo) {
	e.POST("/auth/login", AuthHandler{}.login)
}

func (AuthHandler) login(c echo.Context) error {
	return auth.Login(c)
}
