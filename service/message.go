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
	"strings"
	"time"

	"github.com/Jeffail/gabs"
	"github.com/fannyhasbi/lab-tools-lending/config"
	"github.com/fannyhasbi/lab-tools-lending/helper"
	"github.com/fannyhasbi/lab-tools-lending/types"
)

type MessageService struct {
	messageText        string
	user               types.User
	requestType        string
	chatSessionDetails []types.ChatSessionDetail

	chatSessionService *ChatSessionService
	userService        *UserService
	toolService        *ToolService
}

func NewMessageService(senderID int64, text string, requestType string) *MessageService {
	ms := &MessageService{
		messageText: text,
		requestType: requestType,
		user:        types.User{ID: senderID},
	}

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
		Text: fmt.Sprintf("Silahkan registrasi dengan mengetik `/%s` untuk dapat menggunakan sistem ini secara penuh.", types.Command().Register),
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

func (ms *MessageService) saveChatSessionDetail(topic types.TopicType, sessionData string) error {
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

	if len(sessionData) == 0 {
		sessionData = "{}"
	}

	chatSessionDetail := types.ChatSessionDetail{
		Topic:         topic,
		ChatSessionID: chatSession.ID,
		Data:          sessionData,
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
			Text: "Tidak bisa melakukan registrasi kembali, Anda sudah pernah terdaftar ke dalam sistem pada " + helper.GetDateFromTimestamp(user.CreatedAt),
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

	if err = ms.saveChatSessionDetail(types.Topic["register_init"], ""); err != nil {
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
			Text: "Format registrasi salah, mohon cek format kembali kemudian kirim ulang.",
		}
		return ms.sendMessage(reqBody)
	}

	err = validateRegisterConfirmation(registrationMessage)
	if err != nil {
		log.Println("[ERR][Registration][Validation]", err)
		reqBody := types.MessageRequest{
			Text: "Format registrasi salah. Mohon cek format kembali, kemudian kirim ulang.",
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

	if err = ms.saveChatSessionDetail(types.Topic["register_confirm"], ""); err != nil {
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
	if err := ms.saveChatSessionDetail(types.Topic["register_complete"], ""); err != nil {
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
		Text: "Registrasi berhasil dibatalkan",
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
			Text: fmt.Sprintf("Maaf, Anda belum terdaftar kedalam sistem. Silahkan registrasi dengan cara ketik `/%s`.", types.Command().Register),
		}
		return ms.sendMessage(reqBody)
	}

	ok, toolID := isToolIDWithinBorrowMessage(ms.messageText)
	if ok && toolID > 0 {
		return ms.borrowInit(toolID)
	}

	if len(ms.chatSessionDetails) == 0 {
		return ms.borrowMechanism()
	}

	switch ms.chatSessionDetails[0].Topic {
	case types.Topic["borrow_init"]:
		return ms.borrowAskDateRange()
	case types.Topic["borrow_confirm"]:
		return ms.borrowComplete()
	}

	return nil
}

func isToolIDWithinBorrowMessage(s string) (bool, int64) {
	ss := strings.Split(s, " ")
	if len(ss) != 2 {
		return false, 0
	}

	i, err := strconv.ParseInt(ss[1], 10, 64)
	if err != nil {
		return false, 0
	}

	return true, i
}

func (ms *MessageService) borrowMechanism() error {
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

func (ms *MessageService) borrowInit(toolID int64) error {
	tool, err := ms.toolService.FindByID(toolID)
	if err != nil {
		log.Println("[ERR][Borrow][FindByID]", err)
		reqBody := types.MessageRequest{
			Text: "Maaf, nomor alat yang Anda pilih tidak tersedia.",
		}

		return ms.sendMessage(reqBody)
	}

	borrow := types.Borrow{
		Amount: 1,
		Status: types.GetBorrowStatus("init"),
		UserID: ms.user.ID,
		ToolID: tool.ID,
	}

	borrowService := NewBorrowService()
	borrow, err = borrowService.SaveBorrow(borrow)
	if err != nil {
		log.Println("[ERR][Borrow][SaveBorrow]", err)
		reqBody := types.MessageRequest{
			Text: "Maaf, sedang terjadi kesalahan. Silahkan coba beberapa saat lagi.",
		}

		return ms.sendMessage(reqBody)
	}

	sessionDataGenerator := helper.NewSessionDataGenerator()
	generatedSessionData := sessionDataGenerator.BorrowInit(tool.ID)

	if err = ms.saveChatSessionDetail(types.Topic["borrow_init"], generatedSessionData); err != nil {
		log.Println("[ERR][Borrow][saveChatSessionDetail]", err)
		return err
	}

	reqBody := types.MessageRequest{
		Text: "Berapa lama waktu peminjaman?\n\nJika tidak ada dalam pilihan, maka sebutkan jumlah hari.",
		ReplyMarkup: types.InlineKeyboardMarkup{
			InlineKeyboard: [][]types.InlineKeyboardButton{
				{
					{
						Text:         "1 Minggu",
						CallbackData: types.GetBorrowTimeRange().OneWeek,
					},
					{
						Text:         "2 Minggu",
						CallbackData: types.GetBorrowTimeRange().TwoWeek,
					},
				},
				{
					{
						Text:         "1 Bulan",
						CallbackData: types.GetBorrowTimeRange().OneMonth,
					},
					{
						Text:         "2 Bulan",
						CallbackData: types.GetBorrowTimeRange().TwoMonth,
					},
				},
			},
		},
	}

	return ms.sendMessage(reqBody)
}

func getBorrowTimeRange(message string) (r int, err error) {
	btr := types.GetBorrowTimeRange()
	switch message {
	case btr.OneWeek:
		r = 7
	case btr.TwoWeek:
		r = 14
	case btr.OneMonth:
		r = 30
	case btr.TwoMonth:
		r = 60
	default:
		r, err = strconv.Atoi(message)
	}

	return r, err
}

func (ms *MessageService) borrowAskDateRange() error {
	var borrowAmount int = 1

	borrowDateRange, err := getBorrowTimeRange(ms.messageText)
	if err != nil {
		log.Println("[ERR][Borrow][getBorrowDateRange]", err)
		reqBody := types.MessageRequest{
			Text: "Mohon sebutkan jumlah hari.",
		}

		return ms.sendMessage(reqBody)
	}

	returnDate := time.Now().AddDate(0, 0, borrowDateRange)
	returnDateStr := helper.TranslateDateToBahasa(returnDate)

	sessionData := ms.chatSessionDetails[0].Data
	sessionDataParsed, err := gabs.ParseJSON([]byte(sessionData))
	if err != nil {
		log.Println("[ERR][Borrow][ParseJSON]", err)
		reqBody := types.MessageRequest{
			Text: "Maaf, sedang terjadi kesalahan. Silahkan coba beberapa saat lagi.",
		}

		return ms.sendMessage(reqBody)
	}

	tempToolID, ok := sessionDataParsed.Path("tool_id").Data().(float64)
	if !ok {
		log.Println("[ERR][Borrow][Path] tool_id does not exist in json")
		reqBody := types.MessageRequest{
			Text: "Maaf, sedang terjadi kesalahan. Silahkan coba beberapa saat lagi.",
		}

		return ms.sendMessage(reqBody)
	}

	toolID := int64(tempToolID)

	tool, err := ms.toolService.FindByID(toolID)
	if err != nil {
		log.Println("[ERR][Borrow][FindByID]", err)
		reqBody := types.MessageRequest{
			Text: "Maaf, nomor alat yang Anda pilih tidak tersedia.",
		}

		return ms.sendMessage(reqBody)
	}

	borrow := types.Borrow{
		Amount: borrowAmount,
		Status: types.GetBorrowStatus("progress"),
		UserID: ms.user.ID,
		ToolID: tool.ID,
	}

	borrowService := NewBorrowService()
	borrow, err = borrowService.SaveBorrow(borrow)
	if err != nil {
		log.Println("[ERR][Borrow][SaveBorrow]", err)
		reqBody := types.MessageRequest{
			Text: "Maaf, sedang terjadi kesalahan. Silahkan coba beberapa saat lagi.",
		}

		return ms.sendMessage(reqBody)
	}

	sessionDataGenerator := helper.NewSessionDataGenerator()
	generatedSessionData := sessionDataGenerator.BorrowDateRange(borrowDateRange)

	if err = ms.saveChatSessionDetail(types.Topic["borrow_date"], generatedSessionData); err != nil {
		log.Println("[ERR][Borrow][saveChatSessionDetail]", err)
		return err
	}

	message := fmt.Sprintf(`Nama alat : %s
		Jumlah : %d
		Tanggal Pengembalian : %s
		Alamat peminjam : %s
	`, tool.Name, borrowAmount, returnDateStr, ms.user.Address)
	message = helper.RemoveTab(message)

	reqBody := types.MessageRequest{
		Text: message,
		ReplyMarkup: types.InlineKeyboardMarkup{
			InlineKeyboard: [][]types.InlineKeyboardButton{
				{
					{
						Text:         "Benar",
						CallbackData: "yes",
					},
					{
						Text:         "Salah",
						CallbackData: "no",
					},
				},
			},
		},
	}

	return ms.sendMessage(reqBody)
}

func (ms *MessageService) borrowComplete() error {
	reqBody := types.MessageRequest{
		Text: "Berhasil meminjam!",
	}
	return ms.sendMessage(reqBody)
}
