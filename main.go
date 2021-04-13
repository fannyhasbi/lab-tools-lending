package main

import (
	"log"

	"github.com/fannyhasbi/lab-tools-lending/config"
	"github.com/fannyhasbi/lab-tools-lending/handler"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/labstack/echo/v4"
)

func main() {
	godotenv.Load()

	e := echo.New()
	e.POST("/", handler.WebhookHandler)

	log.Printf("Server running on port %s\n", config.GetPort())
	e.Logger.Fatal(e.Start(":" + config.GetPort()))

	// http.ListenAndServe(":"+config.GetPort(), http.HandlerFunc(handler.WebhookHandler))
}
