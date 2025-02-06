package main

import (
	"payment-service/config"
	"payment-service/internal/handler"
	internal "payment-service/internal/middleware"

	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/joho/godotenv"

	"fmt"
	"os"

	echoswagger "github.com/swaggo/echo-swagger"
)

// @title mandaya project API order-service
// @version 1.0
// @description system booking hotels
// @host localhost:8081
// @BasePath /
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	config.InitDB()

	e := echo.New()
	e.GET("/swagger/*", echoswagger.WrapHandler)

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	private := e.Group("")
	private.Use(internal.CustomJwtMiddleware)
	private.POST("/rooms/payment/:id", handler.PayBooking)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	fmt.Println("Server running on port:", port)
	e.Logger.Fatal(e.Start(":" + port))
}
