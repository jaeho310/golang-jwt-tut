package main

import (
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"os"
	"server/handler"
)

func main() {
	loadEnv()
	e := echo.New()
	initMiddleware(e)
	initHandler(e)
	e.Logger.Fatal(e.Start(":8395"))
}

func initMiddleware(e *echo.Echo) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
}

func initHandler(e *echo.Echo) {
	handler.DefaultHandler{}.Init(e)
	handler.AuthHandler{}.Init(e)
	handler.UserHandler{}.Init(e)
}

func loadEnv() {
	err := godotenv.Load("./env/.env")
	if err != nil {
		panic(err)
	}
	if len(os.Getenv("JWT_SECRET")) < 1 {
		panic(err)
	}
}
