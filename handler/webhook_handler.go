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
	var chatID int64
	var senderID int64
	var messageText string
	var teleMessage types.TeleMessage
	var requestType types.RequestType
	var bodyBytes []byte

	var body *types.WebhookRequest
	var callbackBody *types.InlineCallbackQuery

	var messageService *service.MessageService
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

		chatID = body.Message.Chat.ID
		messageText = body.Message.Text
		teleMessage = body.Message
		if body.Message.Chat.Type == "group" {
			senderID = body.Message.From.ID
			requestType = types.RequestTypeGroup
		} else {
			senderID = chatID
			requestType = types.RequestTypePrivate
		}
	} else {
		chatID = callbackBody.CallbackQuery.Message.Chat.ID
		messageText = callbackBody.CallbackQuery.Data
		teleMessage = callbackBody.CallbackQuery.Message
		if callbackBody.CallbackQuery.Message.Chat.Type == "group" {
			senderID = callbackBody.CallbackQuery.From.ID
			requestType = types.RequestTypeGroup
		} else {
			senderID = chatID
			requestType = types.RequestTypePrivate
		}
	}

	messageService = service.NewMessageService(chatID, senderID, messageText, requestType, teleMessage)

	user := types.User{ID: senderID}

	chatSessionService = service.NewChatSessionService()
	chatSessions, err := chatSessionService.GetChatSessions(user)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return messageService.Error()
	}

	if err != sql.ErrNoRows && len(chatSessions) > 0 && chatSessions[0].Status != types.ChatSessionStatus["complete"] {
		return sessionProcess(chatSessions[0], messageService, chatSessionService)
	}

	match, err := regexp.MatchString("^/", messageText)
	if err != nil {
		log.Println("regex error", err)
		return err
	}
	if match {
		return commandHandler(messageText, messageService)
	}

	// not interacting with chatbot in group (see: privacy mode in Telegrabm bot)
	if requestType == types.RequestTypeGroup {
		return nil
	}

	if callbackBody.CallbackQuery.From.ID == 0 {
		// on message request (join group, added to group, etc still don't know)
		if len(body.Message.Text) == 0 {
			return nil
		}
	}

	return messageService.Unknown()
}

func commandHandler(message string, ms *service.MessageService) error {
	commandStr := helper.GetCommand(message)
	log.Printf("the command is : %s\n", commandStr)

	switch commandStr {
	case types.CommandStart:
		return ms.FirstStart()
	case types.CommandHelp:
		return ms.Help()
	case types.CommandRegister:
		return ms.Register()
	case types.CommandCheck:
		return ms.Check()
	case types.CommandBorrow:
		return ms.Borrow()
	case types.CommandReturn:
		return ms.ReturnTool()
	case types.CommandAdmin:
		return ms.BeAdmin()
	case types.CommandRespond:
		return ms.Respond()
	case types.CommandManage:
		return ms.Manage()
	case types.CommandReport:
		return ms.Report()
	default:
		return ms.Unknown()
	}
}

func sessionProcess(chatSession types.ChatSession, messageService *service.MessageService, chatSessionService *service.ChatSessionService) error {
	var chatSessionDetails []types.ChatSessionDetail
	chatSessionDetails, err := chatSessionService.GetChatSessionDetails(chatSession)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return messageService.Error()
	}

	if len(chatSessionDetails) > 0 {
		messageService.ChangeChatSessionDetails(chatSessionDetails)
		return sessionHandler(chatSessionDetails[0].Topic, messageService)
	}

	return nil
}

func sessionHandler(topic types.TopicType, ms *service.MessageService) error {
	switch topic {
	case types.Topic["register_init"], types.Topic["register_confirm"], types.Topic["register_complete"]:
		return ms.Register()
	case types.Topic["borrow_init"], types.Topic["borrow_amount"], types.Topic["borrow_date"], types.Topic["borrow_reason"], types.Topic["borrow_confirm"]:
		return ms.Borrow()
	case types.Topic["tool_returning_init"], types.Topic["tool_returning_confirm"]:
		return ms.ReturnTool()
	case types.Topic["respond_borrow_init"]:
		return ms.RespondBorrow()
	case types.Topic["respond_tool_returning_init"]:
		return ms.RespondToolReturning()
	case types.Topic["manage_add_init"], types.Topic["manage_add_name"], types.Topic["manage_add_brand"], types.Topic["manage_add_type"], types.Topic["manage_add_weight"], types.Topic["manage_add_stock"], types.Topic["manage_add_info"], types.Topic["manage_add_photo"], types.Topic["manage_add_confirm"]:
		return ms.ManageAdd()
	case types.Topic["manage_edit_init"], types.Topic["manage_edit_field"], types.Topic["manage_edit_complete"]:
		return ms.ManageEdit()
	case types.Topic["manage_delete_init"], types.Topic["manage_delete_complete"]:
		return ms.ManageDelete()
	case types.Topic["manage_photo_init"], types.Topic["manage_photo_upload"], types.Topic["manage_photo_confirm"]:
		return ms.ManagePhoto()
	default:
		return ms.Unknown()
	}
}
