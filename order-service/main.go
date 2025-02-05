package main

import (
	"order-service/config"
	"order-service/internal/handler"
	internal "order-service/internal/middleware"

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

	public := e.Group("")
	public.GET("/rooms", handler.GetRooms)

	private := e.Group("")
	private.Use(internal.CustomJwtMiddleware)
	private.GET("/rooms/users", handler.GetUserRooms)
	private.POST("/rooms/booking", handler.BookRoom)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	fmt.Println("Server running on port:", port)
	e.Logger.Fatal(e.Start(":" + port))
}
