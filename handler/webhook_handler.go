package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/fannyhasbi/lab-tools-lending/service"
	"github.com/fannyhasbi/lab-tools-lending/types"
)

func WebhookHandler(res http.ResponseWriter, req *http.Request) {
	// First, decode the JSON response body
	if req.Method != http.MethodPost {
		fmt.Println("Method not allowed")
		return
	}

	body := &types.WebhookRequest{}
	if err := json.NewDecoder(req.Body).Decode(body); err != nil {
		fmt.Println("could not decode request body", err)
		return
	}

	match, err := regexp.MatchString("^/", body.Message.Text)
	if err != nil {
		fmt.Println("regex error")
		return
	}

	if match {
		commandHandler(body)
		return
	}

	// if err := service.SendLocation(body.Message.Chat.ID); err != nil {
	// 	fmt.Println("error in sending reply:", err)
	// 	return
	// }

	fmt.Println("reply sent")
}

func getCommand(message string) string {
	return message[1:]
}

func commandHandler(body *types.WebhookRequest) {
	commandStr := getCommand(body.Message.Text)
	fmt.Printf("The command is : %s\n", commandStr)

	switch commandStr {
	case types.Command().Help:
		helpHandler(body)
	}
}

func helpHandler(body *types.WebhookRequest) {
	if err := service.SayPolo(body.Message.Chat.ID); err != nil {
		fmt.Println("error in sending reply:", err)
		return
	}
}
