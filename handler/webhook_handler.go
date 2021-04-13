package handler

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/fannyhasbi/lab-tools-lending/helper"
	"github.com/fannyhasbi/lab-tools-lending/service"
	"github.com/fannyhasbi/lab-tools-lending/types"
	"github.com/labstack/echo/v4"
)

func WebhookHandler(c echo.Context) error {
	var messageService *service.MessageService
	var chatSessionService *service.ChatSessionService

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

	messageService = service.NewMessageService(*body)

	user := types.User{
		ChatID: body.Message.Chat.ID,
	}

	// cek apakah chatID sudah ada di db, jika belum maka harus daftar

	chatSessionService = service.NewChatSessionService()
	chatSessions, err := chatSessionService.GetChatSessions(user)
	if err != nil {
		log.Println(err)
		messageService.Error()
	}

	var chatSessionDetails []types.ChatSessionDetail
	if len(chatSessions) > 0 {
		if chatSessions[0].Status == types.ChatSessionStatus["progress"] {
			chatSessionDetails, err = chatSessionService.GetChatSessionDetails(chatSessions[0])
			if err != nil {
				log.Println(err)
				messageService.Error()
			}
		}

		if len(chatSessionDetails) > 0 {
			sessionHandler(chatSessionDetails[0].Topic, body.Message.Text, messageService)
		}
	}

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

func sessionHandler(topic, message string, ms *service.MessageService) {
	switch topic {
	case types.ChatSessionTopic["register"]:
		fmt.Println("Yoyoy")
	default:
		ms.Unknown()
	}
}
