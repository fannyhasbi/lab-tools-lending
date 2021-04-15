package service

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/fannyhasbi/lab-tools-lending/config"
	"github.com/fannyhasbi/lab-tools-lending/helper"
	"github.com/fannyhasbi/lab-tools-lending/types"
)

type MessageService struct {
	messageText        string
	user               types.User
	chatSessionDetails []types.ChatSessionDetail

	chatSessionService *ChatSessionService
	userService        *UserService
	toolService        *ToolService
}

func NewMessageService(senderID int64, text string, requestType string) *MessageService {
	ms := &MessageService{
		messageText: text,
	}

	ms.user = types.User{ID: senderID}

	ms.initChatSessionService()
	ms.initUserService()
	ms.initToolService()

	return ms
}

func (ms *MessageService) initChatSessionService() {
	ms.chatSessionService = NewChatSessionService()
}

func (ms *MessageService) initUserService() {
	ms.userService = NewUserService()
}

func (ms *MessageService) initToolService() {
	ms.toolService = NewToolService()
}

func buildMessageRequest(data *types.MessageRequest) {
	if len(data.ReplyMarkup.InlineKeyboard) == 0 {
		inlineKeyboard := make([][]types.InlineKeyboardButton, 0)
		data.ReplyMarkup.InlineKeyboard = inlineKeyboard
	}
}

func (ms *MessageService) sendMessage(reqBody types.MessageRequest) error {
	if reqBody.ChatID == 0 {
		reqBody.ChatID = ms.user.ID
	}

	buildMessageRequest(&reqBody)

	reqBytes, err := json.Marshal(&reqBody)
	if err != nil {
		return err
	}

	res, err := http.Post(fmt.Sprintf("%s/sendMessage", config.WebhookUrl), "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		defer res.Body.Close()

		var j interface{}
		err = json.NewDecoder(res.Body).Decode(&j)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%#v\n", j)
		return errors.New("unexpected status" + res.Status)
	}

	return nil
}

func (ms *MessageService) ChangeChatSessionDetails(d []types.ChatSessionDetail) {
	ms.chatSessionDetails = d
}

func (ms *MessageService) Error() error {
	reqBody := types.MessageRequest{
		Text: "Maaf, sedang terjadi kesalahan. Silahkan coba beberapa saat lagi.",
	}
	if err := ms.sendMessage(reqBody); err != nil {
		log.Println("error in sending reply:", err)
		return err
	}
	return nil
}

func (ms *MessageService) RecommendRegister() error {
	reqBody := types.MessageRequest{
		Text: fmt.Sprintf("Silahkan daftar dengan mengetik `/%s` untuk dapat menggunakan sistem ini secara penuh.", types.Command().Register),
	}
	if err := ms.sendMessage(reqBody); err != nil {
		log.Println("error in sending reply:", err)
		return err
	}
	return nil
}

func (ms *MessageService) Help() error {
	reqBody := types.MessageRequest{
		Text: "Halo ini adalah pesan bantuan!",
	}
	if err := ms.sendMessage(reqBody); err != nil {
		log.Println("error in sending reply:", err)
		return err
	}

	return nil
}

func (ms *MessageService) Unknown() error {
	reqBody := types.MessageRequest{
		Text: "Maaf, perintah tidak dikenali.",
	}
	if err := ms.sendMessage(reqBody); err != nil {
		log.Println("error in sending reply:", err)
		return err
	}

	return nil
}

func (ms *MessageService) Check() error {
	var toolService *ToolService
	var message string

	toolService = NewToolService()

	tools, err := toolService.GetAvailableTools()
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return err
	}

	message = "Berikut ini daftar alat yang masih tersedia.\n\n"
	message = message + helper.BuildToolListMessage(tools)

	reqBody := types.MessageRequest{
		Text: message,
	}

	return ms.sendMessage(reqBody)
}

func (ms *MessageService) saveChatSessionDetail(topic types.TopicType) error {
	var chatSession types.ChatSession

	chatSessions, err := ms.chatSessionService.GetChatSessions(ms.user)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if len(chatSessions) == 0 {
		chatSessionParam := types.ChatSession{
			Status: types.ChatSessionStatus["progress"],
			UserID: ms.user.ID,
		}

		chatSession, err = ms.chatSessionService.SaveChatSession(chatSessionParam)
		if err != nil {
			return err
		}
	} else {
		chatSession = chatSessions[0]
	}

	chatSessionDetail := types.ChatSessionDetail{
		Topic:         topic,
		ChatSessionID: chatSession.ID,
	}
	_, err = ms.chatSessionService.SaveChatSessionDetail(chatSessionDetail)

	return err
}

func (ms *MessageService) Register() error {
	user, err := ms.userService.FindByID(ms.user.ID)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return err
	}

	if user.IsRegistered() && len(ms.chatSessionDetails) == 0 {
		reqBody := types.MessageRequest{
			Text: "Tidak bisa melakukan pendaftaran kembali, Anda sudah pernah terdaftar ke dalam sistem pada " + helper.GetDateFromTimestamp(user.CreatedAt),
		}
		return ms.sendMessage(reqBody)
	}

	if len(ms.chatSessionDetails) == 0 {
		return ms.registerAsk()
	}

	switch ms.chatSessionDetails[0].Topic {
	case types.Topic["register_init"]:
		return ms.registerConfirm()
	case types.Topic["register_confirm"]:
		return ms.registerComplete()
	}

	return nil
}

func (ms *MessageService) registerAsk() error {
	msg := `Silahkan isi beberapa pertanyaan berikut secara urut (pisahkan dengan baris baru)

		Nama Lengkap
		Nomor Induk Mahasiswa
		Angkatan
		Alamat lengkap tempat tinggal sekarang

		Contoh

		Fanny Hasbi
		211201XXXXXXXX
		2016
		Jalan Jenderal Sudirman No. 189, Pangembon, Brebes, Jawa Tengah`
	msg = helper.RemoveTab(msg)

	reqBody := types.MessageRequest{
		Text: msg,
	}
	if err := ms.sendMessage(reqBody); err != nil {
		log.Println("error in sending reply", err)
		return err
	}

	_, err := ms.userService.SaveUser(ms.user)
	if err != nil {
		return err
	}

	if err = ms.saveChatSessionDetail(types.Topic["register_init"]); err != nil {
		log.Println("[ERR][Registration][saveChatSessionDetail]", err)
		return err
	}

	return nil
}

func (ms *MessageService) registerConfirm() error {
	registrationMessage, err := getRegistrationMessage(ms.messageText)
	if err != nil {
		log.Println("[ERR][Registration][getRegistrationMessage]", err)
		reqBody := types.MessageRequest{
			Text: "Format pendaftaran salah, mohon cek format kembali kemudian kirim ulang.",
		}
		return ms.sendMessage(reqBody)
	}

	err = validateRegisterConfirmation(registrationMessage)
	if err != nil {
		log.Println("[ERR][Registration][Validation]", err)
		reqBody := types.MessageRequest{
			Text: "Format pendaftaran salah. Mohon cek format kembali, kemudian kirim ulang.",
		}
		return ms.sendMessage(reqBody)
	}

	user := types.User{
		ID:      ms.user.ID,
		Name:    registrationMessage.Name,
		NIM:     registrationMessage.NIM,
		Batch:   uint16(registrationMessage.Batch),
		Address: registrationMessage.Address,
	}

	user, err = ms.userService.UpdateUser(user)
	if err != nil {
		log.Println("[ERR][Registration[UpdateUser]", err)
		reqBody := types.MessageRequest{
			Text: "Terjadi kesalahan, mohon coba beberapa saat lagi.",
		}
		return ms.sendMessage(reqBody)
	}

	msg := fmt.Sprintf(`Apakah data ini sudah benar?

		Nama : %s
		NIM : %s
		Angkatan : %d
		Alamat : %s`, user.Name, user.NIM, user.Batch, user.Address)
	msg = helper.RemoveTab(msg)

	reqBody := types.MessageRequest{
		Text: msg,
		ReplyMarkup: types.InlineKeyboardMarkup{
			InlineKeyboard: [][]types.InlineKeyboardButton{
				{
					{
						Text:         "Yakin",
						CallbackData: "yes",
					},
					{
						Text:         "Tidak",
						CallbackData: "no",
					},
				},
			},
		},
	}

	if err = ms.saveChatSessionDetail(types.Topic["register_confirm"]); err != nil {
		log.Println("[ERR][Registration][saveChatSessionDetail]", err)
		return err
	}

	return ms.sendMessage(reqBody)
}

func (ms *MessageService) registerComplete() error {
	var err error

	if ms.messageText == "yes" {
		err = ms.registerCompletePositive()
	} else {
		err = ms.registerCompleteNegative()
	}

	if err != nil {
		log.Println("[ERR][Registration][registerComplete]", err)
		return err
	}

	return nil
}

func (ms *MessageService) registerCompletePositive() error {
	if err := ms.saveChatSessionDetail(types.Topic["register_complete"]); err != nil {
		return err
	}

	chatSessionID := ms.chatSessionDetails[0].ChatSessionID

	if err := ms.chatSessionService.UpdateChatSessionStatus(chatSessionID, types.ChatSessionStatus["complete"]); err != nil {
		return err
	}

	reqBody := types.MessageRequest{
		Text: "Selamat! Anda telah terdaftar dan dapat menggunakan sistem ini.\n\nSilahkan ketik `/help` untuk bantuan.",
	}

	return ms.sendMessage(reqBody)
}

func (ms *MessageService) registerCompleteNegative() error {
	sessionID := ms.chatSessionDetails[0].ChatSessionID
	if err := ms.chatSessionService.DeleteChatSessionDetailByChatSessionID(sessionID); err != nil {
		return err
	}

	if err := ms.chatSessionService.DeleteChatSession(sessionID); err != nil {
		return err
	}

	if err := ms.userService.DeleteUser(ms.user.ID); err != nil {
		return err
	}

	reqBody := types.MessageRequest{
		Text: "Pendaftaran berhasil dibatalkan",
	}

	return ms.sendMessage(reqBody)
}

func getRegistrationMessage(message string) (types.QuestionRegistration, error) {
	registrationMessage := types.QuestionRegistration{}

	splittedMessage := helper.SplitNewLine(message)
	if len(splittedMessage) != 4 {
		return registrationMessage, fmt.Errorf("invalid registration format")
	}

	batch, err := strconv.Atoi(splittedMessage[2])
	if err != nil {
		return registrationMessage, err
	}

	registrationMessage.Name = splittedMessage[0]
	registrationMessage.NIM = splittedMessage[1]
	registrationMessage.Batch = batch
	registrationMessage.Address = splittedMessage[3]

	return registrationMessage, nil
}

func validateRegisterConfirmation(reg types.QuestionRegistration) error {
	if err := validateRegisterMessageBatch(reg.Batch); err != nil {
		return err
	}

	if len(reg.Name) < 4 {
		return fmt.Errorf("invalid name length")
	}

	if len(reg.NIM) < 9 || len(reg.NIM) > 14 {
		return fmt.Errorf("invalid NIM length")
	}

	if len(reg.Address) < 5 {
		return fmt.Errorf("invalid address length")
	}

	return nil
}

func validateRegisterMessageBatch(batch int) error {
	currentYear, _, _ := time.Now().Date()
	if batch < 2008 || batch > currentYear {
		return fmt.Errorf("batch is beyond the limit")
	}
	return nil
}

func (ms *MessageService) Borrow() error {
	user, err := ms.userService.FindByID(ms.user.ID)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ERR][Borrow]", err)
		return err
	}

	ms.user = user
	if !ms.user.IsRegistered() {
		reqBody := types.MessageRequest{
			Text: "Maaf, Anda belum terdaftar dalam sistem. Silahkan daftar dengan cara ketik `/daftar`.",
		}
		return ms.sendMessage(reqBody)
	}

	if len(ms.chatSessionDetails) == 0 {
		return ms.borrowAsk()
	}

	switch ms.chatSessionDetails[0].Topic {
	case types.Topic["borrow_init"]:
		return ms.borrowConfirm()
	case types.Topic["borrow_confirm"]:
		return ms.borrowComplete()
	}

	return nil
}

func (ms *MessageService) borrowAsk() error {
	var message string
	var reqBody types.MessageRequest

	message = `*Mekanisme Peminjaman*

	1\. Cek ketersediaan alat dengan mengetik /cek
	2\. Ketik perintah "*/pinjam \[id\]*", dimana *id* adalah nomor unik alat yang akan dipinjam

	Contoh : "*/pinjam 321*"`
	message = helper.RemoveTab(message)

	reqBody = types.MessageRequest{
		ParseMode: "MarkdownV2",
		Text:      message,
	}

	if err := ms.sendMessage(reqBody); err != nil {
		log.Println("[ERR][Borrow][sendMessage]", err)
		return err
	}

	return ms.Check()
}

func (ms *MessageService) borrowConfirm() error {
	var borrowAmount int
	returnDate := time.Now().AddDate(0, 0, 30).Format("2006-01-02")

	id, err := getToolIDFromBorrowMessage(ms.messageText)
	if err != nil {
		log.Println("[ERR][Borrow][getToolIDFromBorrowMessage]", err)
		reqBody := types.MessageRequest{
			Text: "Maaf, format yang Anda masukkan salah. Silahkan kirimkan nomor alat sesuai dengan daftar.",
		}

		return ms.sendMessage(reqBody)
	}

	tool, err := ms.toolService.FindByID(id)
	if err != nil {
		log.Println("[ERR][Borrow][FindByID]", err)
		reqBody := types.MessageRequest{
			Text: "Maaf, nomor alat yang Anda pilih tidak tersedia.",
		}

		return ms.sendMessage(reqBody)
	}

	borrowService := NewBorrowService()
	borrow := types.Borrow{
		Amount:     borrowAmount,
		ReturnDate: returnDate,
		UserID:     ms.user.ID,
		ToolID:     tool.ID,
	}

	borrow, err = borrowService.SaveBorrow(borrow)
	if err != nil {
		log.Println("[ERR][Borrow][SaveBorrow]", err)
		reqBody := types.MessageRequest{
			Text: "Maaf, sedang terjadi kesalahan. Silahkan coba beberapa saat lagi.",
		}

		return ms.sendMessage(reqBody)
	}

	message := fmt.Sprintf(`Apakah data ini sudah benar?

		Nama alat : %s
		Jumlah : %d
		Tanggal Pengembalian : %s
		Alamat peminjam :
		%s
	`, tool.Name, borrowAmount, returnDate, ms.user.Address)
	message = helper.RemoveTab(message)

	reqBody := types.MessageRequest{
		Text: message,
		ReplyMarkup: types.InlineKeyboardMarkup{
			InlineKeyboard: [][]types.InlineKeyboardButton{
				{
					{
						Text:         "Yakin",
						CallbackData: "yes",
					},
					{
						Text:         "Tidak",
						CallbackData: "no",
					},
				},
			},
		},
	}

	if err = ms.saveChatSessionDetail(types.Topic["borrow_confirm"]); err != nil {
		log.Println("[ERR][Borrow][saveChatSessionDetail]", err)
		return err
	}

	return ms.sendMessage(reqBody)
}

func getToolIDFromBorrowMessage(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func (ms *MessageService) borrowComplete() error {
	reqBody := types.MessageRequest{
		Text: "Berhasil meminjam!",
	}
	return ms.sendMessage(reqBody)
}
