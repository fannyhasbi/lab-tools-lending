package handler

import (
	"log"
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
		log.Println("could not decode request body", err)
		return err
	}

	match, err := regexp.MatchString("^/", body.Message.Text)
	if err != nil {
		log.Println("regex error", err)
		return err
	}

	messageService := service.NewMessageService(*body)

	if match {
		commandHandler(body.Message.Text, messageService)
	}

	return c.String(http.StatusOK, "OK")
}

func commandHandler(message string, ms *service.MessageService) {
	commandStr := helper.GetCommand(message)
	log.Printf("The command is : %s\n", commandStr)

	switch commandStr {
	case types.Command().Help:
		ms.Help()
	case types.Command().Check:
		ms.Check()
	default:
		ms.Unknown()
	}
}
