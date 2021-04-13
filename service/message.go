package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/fannyhasbi/lab-tools-lending/config"
	"github.com/fannyhasbi/lab-tools-lending/types"
)

type MessageService struct {
	WebhookRequest types.WebhookRequest
}

func NewMessageService(wr types.WebhookRequest) *MessageService {
	return &MessageService{
		WebhookRequest: wr,
	}
}

func (ms *MessageService) sendMessage(message string) error {
	reqBody := &types.MessageRequest{
		ChatID: ms.WebhookRequest.Message.Chat.ID,
		Text:   message,
	}

	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	res, err := http.Post(fmt.Sprintf("%s/sendMessage", config.WebhookUrl), "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return errors.New("unexpected status" + res.Status)
	}

	return nil
}

func (ms *MessageService) Error() error {
	if err := ms.sendMessage("Maaf, sedang terjadi kesalahan. Silahkan coba beberapa saat lagi."); err != nil {
		log.Println("error in sending reply:", err)
		return err
	}
	return nil
}

func (ms *MessageService) Help() error {
	if err := ms.sendMessage("Halo ini adalah pesan bantuan!"); err != nil {
		log.Println("error in sending reply:", err)
		return err
	}

	return nil
}

func (ms *MessageService) Unknown() error {
	if err := ms.sendMessage("Maaf, perintah tidak dikenali."); err != nil {
		log.Println("error in sending reply:", err)
		return err
	}

	return nil
}

func (ms *MessageService) Check() error {
	if err := ms.sendMessage("Mantap"); err != nil {
		log.Println("error in sending reply", err)
		return err
	}

	return nil
}
