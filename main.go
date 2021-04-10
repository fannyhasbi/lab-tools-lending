package main

import (
	"fmt"

	"github.com/fannyhasbi/lab-tools-lending/config"
	"github.com/fannyhasbi/lab-tools-lending/handler"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.POST("/", handler.WebhookHandler)

	fmt.Printf("Server running on port %s\n", config.GetPort())
	e.Logger.Fatal(e.Start(":" + config.GetPort()))

	// http.ListenAndServe(":"+config.GetPort(), http.HandlerFunc(handler.WebhookHandler))
}
