package main

import (
	"user-service/config"
	"user-service/internal/handler"

	//"user-service/internal/middleware"

	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/joho/godotenv"

	"fmt"
	"os"
	//echoswagger "github.com/swaggo/echo-swagger"
)

// @title mandaya project API
// @version 1.0
// @description system booking hotels
// @host localhost:8080
// @BasePath /
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	config.InitDB()

	e := echo.New()
	//e.GET("/swagger/*", echoswagger.WrapHandler)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	public := e.Group("")
	public.POST("/users/register", handler.Register)
	public.POST("/users/login", handler.LoginUser)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Server running on port:", port)
	e.Logger.Fatal(e.Start(":" + port))
}
