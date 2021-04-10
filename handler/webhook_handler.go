package handler

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/fannyhasbi/lab-tools-lending/helper"
	"github.com/fannyhasbi/lab-tools-lending/service"
	"github.com/fannyhasbi/lab-tools-lending/types"
	"github.com/labstack/echo/v4"
)

func WebhookHandler(c echo.Context) error {
	body := new(types.WebhookRequest)
	if err := c.Bind(body); err != nil {
		fmt.Println("could not decode request body", err)
		return err
	}

	match, err := regexp.MatchString("^/", body.Message.Text)
	if err != nil {
		fmt.Println("regex error")
		return err
	}

	if match {
		commandHandler(body)
	}

	return c.String(http.StatusOK, "OK")
}

func commandHandler(body *types.WebhookRequest) {
	commandStr := helper.GetCommand(body.Message.Text)
	fmt.Printf("The command is : %s\n", commandStr)

	switch commandStr {
	case types.Command().Help:
		helpHandler(body)
	default:
		unknownCommandHandler(body.Message.Chat.ID)
	}
}

func helpHandler(body *types.WebhookRequest) {
	if err := service.SayPolo(body.Message.Chat.ID); err != nil {
		fmt.Println("error in sending reply:", err)
		return
	}
}

func unknownCommandHandler(chatID int64) {
	if err := service.UnknownCommand(chatID); err != nil {
		fmt.Println("erro in sending reply: ")
	}
}
