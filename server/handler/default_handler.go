package handler

import (
	"github.com/labstack/echo"
	"net/http"
)

type DefaultHandler struct{}

func (DefaultHandler) Init(e *echo.Echo) {
	e.GET("/health", DefaultHandler{}.health)
}

func (DefaultHandler) health(c echo.Context) error {
	return c.JSON(http.StatusOK, "ok")
}
