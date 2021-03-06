package main

import (
	"log"
	"os"

	"github.com/fannyhasbi/lab-tools-lending/config"
	"github.com/fannyhasbi/lab-tools-lending/handler"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	godotenv.Load()
	environment := os.Getenv("ENVIRONMENT")
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	if environment == "development" || environment == "" || environment == "production" {
		e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
			log.Println("REQUEST", string(reqBody))
			log.Println("RESPONSE", string(resBody))
		}))
	}

	e.POST("/", handler.WebhookHandler)

	log.Printf("Server running on port %s\n", config.GetPort())
	e.Logger.Fatal(e.Start(":" + config.GetPort()))
}
