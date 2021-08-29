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
	borrowService      *BorrowService
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
	ms.initBorrowService()

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

func (ms *MessageService) initBorrowService() {
	ms.borrowService = NewBorrowService()
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
	var message string

	tools, err := ms.toolService.GetAvailableTools()
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

	msg := fmt.Sprintf(`Apakah Anda yakin data ini sudah benar?

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
						Text:         "Lanjutkan",
						CallbackData: "yes",
					},
					{
						Text:         "Batalkan",
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
		Text: fmt.Sprintf("Selamat! Anda telah terdaftar dan dapat menggunakan sistem ini.\n\nSilahkan ketik `/%s` untuk bantuan.", types.Command().Help),
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
	case types.Topic["borrow_date"]:
		return ms.borrowConfirm()
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

	message = fmt.Sprintf(`*Mekanisme Peminjaman*

	1\. Cek ketersediaan alat dengan mengetik /%s
	2\. Ketik perintah "*/%s \[id\]*", dimana *id* adalah nomor unik alat yang akan dipinjam

	Contoh : "*/%s 321*"`, types.Command().Check, types.Command().Borrow, types.Command().Borrow)
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

	borrow, err = ms.borrowService.SaveBorrow(borrow)
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
						CallbackData: strconv.Itoa(types.BorrowTimeRangeMap["oneweek"]),
					},
					{
						Text:         "2 Minggu",
						CallbackData: strconv.Itoa(types.BorrowTimeRangeMap["twoweek"]),
					},
				},
				{
					{
						Text:         "1 Bulan",
						CallbackData: strconv.Itoa(types.BorrowTimeRangeMap["onemonth"]),
					},
					{
						Text:         "2 Bulan",
						CallbackData: strconv.Itoa(types.BorrowTimeRangeMap["twomonth"]),
					},
				},
			},
		},
	}

	return ms.sendMessage(reqBody)
}

func (ms *MessageService) borrowAskDateRange() error {
	var borrowAmount int = 1

	borrowDateRange, err := helper.GetBorrowTimeRangeValue(ms.messageText)
	if err != nil {
		log.Println("[ERR][Borrow][getBorrowDateRange]", err)
		reqBody := types.MessageRequest{
			Text: "Mohon sebutkan jumlah hari.",
		}

		return ms.sendMessage(reqBody)
	}

	returnDate := time.Now().AddDate(0, 0, borrowDateRange)
	returnDateStr := helper.TranslateDateToBahasa(returnDate)

	borrow, err := ms.borrowService.FindInitialByUserID(ms.user.ID)
	if err != nil {
		log.Println("[ERR][Borrow][FindInitialByUserID]", err)
		reqBody := types.MessageRequest{
			Text: "Maaf, sedang terjadi kesalahan. Silahkan coba beberapa saat lagi.",
		}

		return ms.sendMessage(reqBody)
	}

	tool, err := ms.toolService.FindByID(borrow.ToolID)
	if err != nil {
		log.Println("[ERR][Borrow][FindByID]", err)
		reqBody := types.MessageRequest{
			Text: "Maaf, sedang terjadi kesalahan. Silahkan coba beberapa saat lagi.",
		}

		return ms.sendMessage(reqBody)
	}

	borrow.ReturnDate = sql.NullString{
		Valid:  true,
		String: returnDate.Format("2006-01-02"),
	}

	borrow, err = ms.borrowService.UpdateBorrow(borrow)
	if err != nil {
		log.Println("[ERR][Borrow][UpdateBorrow]", err)
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
		Tanggal Pengembalian : %s (%d hari)
		Alamat peminjam : %s
	`, tool.Name, borrowAmount, returnDateStr, borrowDateRange, ms.user.Address)
	message = helper.RemoveTab(message)

	reqBody := types.MessageRequest{
		Text: message,
		ReplyMarkup: types.InlineKeyboardMarkup{
			InlineKeyboard: [][]types.InlineKeyboardButton{
				{
					{
						Text:         "Lanjutkan",
						CallbackData: "yes",
					},
					{
						Text:         "Batalkan",
						CallbackData: "no",
					},
				},
			},
		},
	}

	return ms.sendMessage(reqBody)
}

func (ms *MessageService) borrowConfirm() error {
	var userResponse bool
	if ms.messageText == "yes" {
		userResponse = true
	} else {
		userResponse = false
	}

	sessionDataGenerator := helper.NewSessionDataGenerator()
	generatedSessionData := sessionDataGenerator.BorrowConfirmation(userResponse)

	if err := ms.saveChatSessionDetail(types.Topic["borrow_confirm"], generatedSessionData); err != nil {
		return err
	}

	chatSessionID := ms.chatSessionDetails[0].ChatSessionID
	if err := ms.chatSessionService.UpdateChatSessionStatus(chatSessionID, types.ChatSessionStatus["complete"]); err != nil {
		return err
	}

	borrow, err := ms.borrowService.FindInitialByUserID(ms.user.ID)
	if err != nil {
		log.Println("[ERR][borrowConfirm][FindInitialByUserID]", err)
		reqBody := types.MessageRequest{
			Text: "Maaf, sedang terjadi kesalahan. Silahkan coba beberapa saat lagi.",
		}

		return ms.sendMessage(reqBody)
	}

	var message string
	if userResponse {
		message = "Pengajuan peminjaman berhasil, silahkan tunggu hingga pengurus mengkonfirmasi pengajuan."
	} else {
		borrow.Status = types.GetBorrowStatus("cancel")
		message = "Pengajuan berhasil dibatalkan"
	}

	_, err = ms.borrowService.UpdateBorrow(borrow)
	if err != nil {
		log.Println("[ERR][borrowConfirm][UpdateBorrow]", err)
		reqBody := types.MessageRequest{
			Text: "Maaf, sedang terjadi kesalahan. Silahkan coba beberapa saat lagi.",
		}

		return ms.sendMessage(reqBody)
	}

	if userResponse {
		go func() {
			time.Sleep(2 * time.Second)
			ms.sendToAdmin(borrow)
		}()
	}

	return ms.sendMessage(types.MessageRequest{
		Text: message,
	})
}

// todo: Notif to admin
func (ms *MessageService) sendToAdmin(borrow types.Borrow) error {
	/**
	* todo: update trigger when admin confirm
	* but for this time it update here
	**/
	borrow.Status = types.GetBorrowStatus("progress")
	if _, err := ms.borrowService.UpdateBorrow(borrow); err != nil {
		log.Println("[ERR][sendToAdmin][UpdateBorrow]", err)
		reqBody := types.MessageRequest{
			Text: "Maaf, sedang terjadi kesalahan. Silahkan coba beberapa saat lagi.",
		}

		return ms.sendMessage(reqBody)
	}

	return ms.sendMessage(types.MessageRequest{
		Text: "Permintaan Anda telah disetujui.",
	})
}

func (ms *MessageService) ReturnTool() error {
	// return ms.borrowedTools()
	return ms.currentlyBorrowedTools()
}

func (ms *MessageService) borrowedTools() error {
	var message string

	borrows, err := ms.borrowService.FindByUserID(ms.user.ID)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return err
	}

	borrowGroup := helper.GroupBorrowStatus(borrows)
	progressLen := len(borrowGroup[types.GetBorrowStatus("progress")])
	returnedLen := len(borrowGroup[types.GetBorrowStatus("returned")])

	if progressLen > 0 || returnedLen > 0 {
		message += "Alat yang sedang Anda pinjam:\n"
		if progressLen > 0 {
			message += helper.BuildBorrowedMessage(borrowGroup[types.GetBorrowStatus("progress")])
		} else {
			message += "Saat ini tidak ada alat yang sedang Anda pinjam."
		}

		message += "\n\nAlat yang sudah pernah Anda pinjam:\n"
		if returnedLen > 0 {
			message += helper.BuildBorrowedMessage(borrowGroup[types.GetBorrowStatus("returned")])
		} else {
			message += "Belum ada alat yang pernah Anda kembalikan."
		}
	} else {
		message += "Anda belum pernah meminjam alat."
	}

	reqBody := types.MessageRequest{
		Text: message,
	}

	return ms.sendMessage(reqBody)
}

func (ms *MessageService) currentlyBorrowedTools() error {
	var message string
	var returnDateLayout string = "2006-01-02T15:04:05Z"

	borrow, err := ms.borrowService.FindCurrentlyBeingBorrowedByUserID(ms.user.ID)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return err
	}

	if err == sql.ErrNoRows {
		message += "Saat ini tidak ada alat yang sedang Anda pinjam."
		reqBody := types.MessageRequest{
			Text: message,
		}

		return ms.sendMessage(reqBody)
	}

	message += fmt.Sprintf("Anda sedang meminjam %s sejak %s.\nJadwal pengembalian pada %s\n\n", borrow.Tool.Name, helper.TranslateDateStringToBahasa(borrow.CreatedAt), helper.TranslateDateStringToBahasa(borrow.ReturnDate.String))

	returnDateTime, err := time.Parse(returnDateLayout, borrow.ReturnDate.String)
	if err != nil {
		log.Println(err)
		return err
	}

	currentDateTime := time.Now()
	dayDifference := int(returnDateTime.Sub(currentDateTime).Hours() / 24)

	if dayDifference < 0 {
		message += fmt.Sprintf("Anda sudah lewat %d hari dari jadwal pengembalian barang.", -dayDifference)
	} else if dayDifference == 0 {
		message += "Hari ini adalah jadwal pengembalian barang."
	} else {
		message += fmt.Sprintf("Durasi peminjaman tersisa %d hari lagi.", dayDifference)
	}

	reqBody := types.MessageRequest{
		Text: message,
	}

	return ms.sendMessage(reqBody)
}
