package main

import (
	"log"

	"github.com/fannyhasbi/lab-tools-lending/config"
	"github.com/fannyhasbi/lab-tools-lending/handler"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	godotenv.Load()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	e := echo.New()

	e.Use(middleware.Logger())

	e.POST("/", handler.WebhookHandler)

	log.Printf("Server running on port %s\n", config.GetPort())
	e.Logger.Fatal(e.Start(":" + config.GetPort()))
}
