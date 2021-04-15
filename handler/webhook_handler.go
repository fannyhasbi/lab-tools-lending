package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"regexp"

	"github.com/fannyhasbi/lab-tools-lending/helper"
	"github.com/fannyhasbi/lab-tools-lending/service"
	"github.com/fannyhasbi/lab-tools-lending/types"
	"github.com/labstack/echo/v4"
)

func WebhookHandler(c echo.Context) error {
	var senderID int64
	var messageText string
	var bodyBytes []byte

	var body *types.WebhookRequest
	var callbackBody *types.InlineCallbackQuery

	var messageService *service.MessageService
	var userService *service.UserService
	var chatSessionService *service.ChatSessionService

	bodyBytes, _ = ioutil.ReadAll(c.Request().Body)
	c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	callbackBody = new(types.InlineCallbackQuery)
	if err := json.Unmarshal(bodyBytes, callbackBody); err != nil {
		log.Println("could not decode request body", err)
		return err
	}

	// check whether it is an inline callback query request or common
	if callbackBody.CallbackQuery.From.ID == 0 {
		body = new(types.WebhookRequest)
		if err := json.Unmarshal(bodyBytes, body); err != nil {
			log.Println("could not decode request body", err)
			return err
		}

		senderID = body.Message.Chat.ID
		messageText = body.Message.Text
		messageService = service.NewMessageService(senderID, messageText, types.RequestTypeCommon)
	} else {
		senderID = callbackBody.CallbackQuery.From.ID
		messageText = callbackBody.CallbackQuery.Data
		messageService = service.NewMessageService(senderID, messageText, types.RequestTypeInlineCallback)
	}

	user := types.User{ID: senderID}

	userService = service.NewUserService()
	user, err := userService.FindByID(senderID)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return messageService.Error()
	}

	chatSessionService = service.NewChatSessionService()
	chatSessions, err := chatSessionService.GetChatSessions(user)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return messageService.Error()
	}

	if err != sql.ErrNoRows && len(chatSessions) > 0 && chatSessions[0].Status != types.ChatSessionStatus["complete"] {
		return sessionProcess(messageText, chatSessions[0], messageService, chatSessionService)
	}

	match, err := regexp.MatchString("^/", messageText)
	if err != nil {
		log.Println("regex error", err)
		return err
	}
	if match {
		return commandHandler(messageText, messageService)
	}

	return messageService.Unknown()
}

func commandHandler(message string, ms *service.MessageService) error {
	commandStr := helper.GetCommand(message)
	log.Printf("the command is : %s\n", commandStr)

	switch commandStr {
	case types.Command().Help:
		return ms.Help()
	case types.Command().Register:
		return ms.Register()
	case types.Command().Check:
		return ms.Check()
	case types.Command().Borrow:
		return ms.Borrow()
	default:
		return ms.Unknown()
	}
}

func sessionProcess(message string, chatSession types.ChatSession, messageService *service.MessageService, chatSessionService *service.ChatSessionService) error {
	var chatSessionDetails []types.ChatSessionDetail
	chatSessionDetails, err := chatSessionService.GetChatSessionDetails(chatSession)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return messageService.Error()
	}

	if err != sql.ErrNoRows {
		messageService.ChangeChatSessionDetails(chatSessionDetails)
		return sessionHandler(chatSessionDetails[0].Topic, message, messageService)
	}

	return nil
}

func sessionHandler(topic types.TopicType, message string, ms *service.MessageService) error {
	switch topic {
	case types.Topic["register_init"], types.Topic["register_confirm"], types.Topic["register_complete"]:
		return ms.Register()
	case types.Topic["borrow_init"], types.Topic["borrow_confirm"]:
		return ms.Borrow()
	default:
		return ms.Unknown()
	}
}
