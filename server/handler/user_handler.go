package handler

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
	"os"
	"server/auth"
)

type UserHandler struct{}

func (UserHandler) Init(e *echo.Echo) {
	authApiGroup := e.Group("/api/user")
	// 에코 프레임워크를 사용하여 secret과 헤더의 알고리즘을 이용해 인증합니다.
	authApiGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     &auth.MyClaim{},
		SigningKey: []byte(os.Getenv("JWT_SECRET")),
	}))
	authApiGroup.GET("/me", UserHandler{}.getUserInfo)
}

func (UserHandler) getUserInfo(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*auth.MyClaim)
	name := claims.UserName
	return c.JSON(http.StatusOK, echo.Map{"name": name})
}
