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
// @description system booking mandaya hotels
// @host localhost:8082
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
	private.POST("/create_invoice", handler.CreateTopUpInvoice)
	private.POST("/get_price/:id", handler.GetPrice)

	e.POST("/xendit_webhook", handler.XenditWebhook)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	fmt.Println("Server running on port:", port)
	e.Logger.Fatal(e.Start(":" + port))
}
