package handler

import (
	"database/sql"
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
	var userService *service.UserService
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

	userService = service.NewUserService()
	user, err = userService.FindByChatID(body.Message.Chat.ID)
	if err == sql.ErrNoRows {

	} else if err != nil {
		log.Println(err)
		messageService.Error()
	}

	log.Println(user)

	chatSessionService = service.NewChatSessionService()
	chatSessions, err := chatSessionService.GetChatSessions(user)
	if err != nil {
		log.Println(err)
		return messageService.Error()
	}

	var chatSessionDetails []types.ChatSessionDetail
	if len(chatSessions) > 0 {
		if chatSessions[0].Status == types.ChatSessionStatus["progress"] {
			chatSessionDetails, err = chatSessionService.GetChatSessionDetails(chatSessions[0])
			if err != nil {
				log.Println(err)
				return messageService.Error()
			}
		}

		if len(chatSessionDetails) > 0 {
			return sessionHandler(chatSessionDetails[0].Topic, body.Message.Text, messageService)
		}
	}

	if match {
		return commandHandler(body.Message.Text, messageService)
	}

	return c.String(http.StatusOK, "OK")
}

func commandHandler(message string, ms *service.MessageService) error {
	commandStr := helper.GetCommand(message)
	log.Printf("The command is : %s\n", commandStr)

	switch commandStr {
	case types.Command().Help:
		return ms.Help()
	case types.Command().Check:
		return ms.Check()
	default:
		return ms.Unknown()
	}
}

func sessionHandler(topic, message string, ms *service.MessageService) error {
	switch topic {
	case types.ChatSessionTopic["register"]:
		fmt.Println("Yoyoy")
		return nil
	default:
		return ms.Unknown()
	}
}
