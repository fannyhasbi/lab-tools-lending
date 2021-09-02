package service

import (
	"bytes"
	"context"
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
	"golang.org/x/sync/errgroup"
)

type MessageService struct {
	chatID             int64
	messageText        string
	requestType        types.RequestType
	user               types.User
	chatSessionDetails []types.ChatSessionDetail

	chatSessionService   *ChatSessionService
	userService          *UserService
	toolService          *ToolService
	borrowService        *BorrowService
	toolReturningService *ToolReturningService
}

func NewMessageService(chatID, senderID int64, text string, requestType types.RequestType) *MessageService {
	ms := &MessageService{
		chatID:      chatID,
		messageText: text,
		requestType: requestType,
		user:        types.User{ID: senderID},
	}

	ms.initChatSessionService()
	ms.initUserService()
	ms.initToolService()
	ms.initBorrowService()
	ms.initToolReturningService()

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

func (ms *MessageService) initToolReturningService() {
	ms.toolReturningService = NewToolReturningService()
}

func (ms *MessageService) sendMessage(reqBody types.MessageRequest) error {
	if reqBody.ChatID == 0 {
		reqBody.ChatID = ms.chatID
	}

	helper.BuildMessageRequest(&reqBody)

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
	return ms.sendMessage(reqBody)
}

func (ms *MessageService) RecommendRegister() error {
	reqBody := types.MessageRequest{
		Text: fmt.Sprintf("Silahkan registrasi dengan mengetik `/%s` untuk dapat menggunakan sistem ini secara penuh.", types.CommandRegister),
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
	tools, err := ms.toolService.GetAvailableTools()
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return err
	}

	toolID, ok := isIDWithinCommand(ms.messageText)
	if ok && toolID > 0 {
		return ms.checkDetail(toolID)
	}

	message := "Berikut ini daftar alat yang masih tersedia.\n"
	message += fmt.Sprintf("untuk melihat detail alat, ketik perintah \"/%s [id]\"\n\n", types.CommandCheck)
	message += helper.BuildToolListMessage(tools)

	reqBody := types.MessageRequest{
		Text: message,
	}

	return ms.sendMessage(reqBody)
}

func (ms *MessageService) checkDetail(toolID int64) error {
	tool, err := ms.toolService.FindByID(toolID)
	if err != nil || tool.Stock < 1 {
		log.Println("[ERR][Borrow][FindByID]", err)
		return ms.sendMessage(types.MessageRequest{
			Text: "Maaf, nomor alat yang Anda pilih tidak tersedia.",
		})
	}

	message := fmt.Sprintf(`Nama: %s
	Brand: %s
	Tipe: %s
	Berat: %.2f gram
	Stok: %d

	Keterangan:
	%s
	`, tool.Name, tool.Brand, tool.ProductType, tool.Weight, tool.Stock, tool.AdditionalInformation)
	message = helper.RemoveTab(message)

	return ms.sendMessage(types.MessageRequest{
		Text: message,
		ReplyMarkup: types.InlineKeyboardMarkup{
			InlineKeyboard: [][]types.InlineKeyboardButton{
				{
					{
						Text:         "Pinjam",
						CallbackData: fmt.Sprintf("/%s %d", types.CommandBorrow, tool.ID),
					},
				},
			},
		},
	})
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
			Text: "Tidak bisa melakukan registrasi, Anda sudah terdaftar ke dalam sistem pada " + helper.TranslateDateStringToBahasa(user.CreatedAt),
		}
		return ms.sendMessage(reqBody)
	}

	if len(ms.chatSessionDetails) == 0 {
		return ms.registerInit()
	}

	switch ms.chatSessionDetails[0].Topic {
	case types.Topic["register_init"]:
		return ms.registerConfirm()
	case types.Topic["register_confirm"]:
		return ms.registerComplete()
	}

	return nil
}

func (ms *MessageService) registerInit() error {
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
	sessionDataGenerator := helper.NewSessionDataGenerator()
	generatedSessionData := sessionDataGenerator.RegisterComplete(true)

	if err := ms.saveChatSessionDetail(types.Topic["register_complete"], generatedSessionData); err != nil {
		return err
	}

	chatSessionID := ms.chatSessionDetails[0].ChatSessionID

	if err := ms.chatSessionService.UpdateChatSessionStatus(chatSessionID, types.ChatSessionStatus["complete"]); err != nil {
		return err
	}

	reqBody := types.MessageRequest{
		Text: fmt.Sprintf("Selamat! Anda telah terdaftar dan dapat menggunakan sistem ini.\n\nSilahkan ketik `/%s` untuk bantuan.", types.CommandHelp),
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
		Text: "Registrasi berhasil dibatalkan.",
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
	currentYear := time.Now().Year()
	if batch < 2008 || batch > currentYear {
		return fmt.Errorf("batch is beyond the limit")
	}
	return nil
}

func (ms *MessageService) notRegistered() error {
	reqBody := types.MessageRequest{
		Text: fmt.Sprintf("Maaf, Anda belum terdaftar kedalam sistem. Silahkan registrasi dengan cara ketik `/%s`.", types.CommandRegister),
	}
	return ms.sendMessage(reqBody)
}

func (ms *MessageService) Borrow() error {
	user, err := ms.userService.FindByID(ms.user.ID)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ERR][Borrow]", err)
		return err
	}

	ms.user = user
	if !ms.user.IsRegistered() {
		return ms.notRegistered()
	}

	toolID, ok := isIDWithinCommand(ms.messageText)
	if ok && toolID > 0 {
		return ms.borrowInit(toolID)
	}

	if len(ms.chatSessionDetails) > 0 {
		switch ms.chatSessionDetails[0].Topic {
		case types.Topic["borrow_init"]:
			return ms.borrowAskDateRange()
		case types.Topic["borrow_date"]:
			return ms.borrowReason()
		case types.Topic["borrow_reason"]:
			return ms.borrowConfirm()
		}
	}

	return ms.borrowMechanism()
}

func isIDWithinCommand(s string) (int64, bool) {
	ss := strings.Split(s, " ")
	if len(ss) != 2 {
		return 0, false
	}

	i, err := strconv.ParseInt(ss[1], 10, 64)
	if err != nil {
		return 0, false
	}

	return i, true
}

func (ms *MessageService) borrowMechanism() error {
	var message string
	var reqBody types.MessageRequest

	message = fmt.Sprintf(`*Mekanisme Peminjaman*

	1\. Cek ketersediaan alat dengan mengetik /%s
	2\. Ketik perintah "*/%s \[id\]*", dimana *id* adalah nomor unik alat yang akan dipinjam

	Contoh : "*/%s 321*"`, types.CommandCheck, types.CommandBorrow, types.CommandBorrow)
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
		return ms.sendMessage(types.MessageRequest{
			Text: "Maaf, nomor alat yang Anda pilih tidak tersedia.",
		})
	}

	if tool.Stock < 1 {
		return ms.sendMessage(types.MessageRequest{
			Text: "Stok barang sudah habis. Tidak dapat melakukan pengajuan peminjaman.",
		})
	}

	borrows, err := ms.borrowService.GetCurrentlyBeingBorrowedAndRequestedByUserID(ms.user.ID)
	if err != nil {
		log.Println(err)
		return ms.Error()
	}

	if len(borrows) > 0 {
		var message string
		requestingBorrows := helper.GetBorrowsByStatus(borrows, types.GetBorrowStatus("request"))
		if len(requestingBorrows) > 0 {
			message = "Maaf, saat ini status Anda sedang mengajukan peminjaman, silahkan tunggu hingga pengurus menanggapi pengajuan tersebut."
		} else {
			message = "Maaf, saat ini status Anda sedang meminjam barang sehingga tidak dapat mengajukan peminjaman.\n"
			message += fmt.Sprintf(`Untuk melakukan pengembalian silahkan ketik "/%s"`, types.CommandReturn)
		}

		reqBody := types.MessageRequest{
			Text: message,
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
		return ms.Error()
	}

	sessionDataGenerator := helper.NewSessionDataGenerator()
	generatedSessionData := sessionDataGenerator.BorrowInit(tool.ID)

	if err = ms.saveChatSessionDetail(types.Topic["borrow_init"], generatedSessionData); err != nil {
		log.Println("[ERR][Borrow][saveChatSessionDetail]", err)
		return err
	}

	reqBody := types.MessageRequest{
		Text: fmt.Sprintf("Berapa lama waktu peminjaman?\n\nJika tidak ada dalam pilihan, maka sebutkan jumlah hari. Minimal durasi peminjaman adalah %d hari.", types.BorrowMinimalDuration),
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
	borrowDateRange, err := helper.GetBorrowTimeRangeValue(ms.messageText)
	if err != nil {
		log.Println("[ERR][Borrow][getBorrowDateRange]", err)
		return ms.sendMessage(types.MessageRequest{
			Text: "Mohon sebutkan jumlah hari.",
		})
	}

	if borrowDateRange < types.BorrowMinimalDuration {
		return ms.sendMessage(types.MessageRequest{
			Text: fmt.Sprintf("Minimal durasi peminjaman adalah %d hari", types.BorrowMinimalDuration),
		})
	}

	returnDate := time.Now().AddDate(0, 0, borrowDateRange)

	borrow, err := ms.borrowService.FindInitialByUserID(ms.user.ID)
	if err != nil {
		log.Println("[ERR][Borrow][FindInitialByUserID]", err)
		return ms.Error()
	}

	borrow.ReturnDate = sql.NullString{
		Valid:  true,
		String: returnDate.Format("2006-01-02"),
	}

	borrow, err = ms.borrowService.UpdateBorrow(borrow)
	if err != nil {
		log.Println("[ERR][Borrow][UpdateBorrow]", err)
		return ms.Error()
	}

	sessionDataGenerator := helper.NewSessionDataGenerator()
	generatedSessionData := sessionDataGenerator.BorrowDateRange(borrowDateRange)

	if err = ms.saveChatSessionDetail(types.Topic["borrow_date"], generatedSessionData); err != nil {
		log.Println("[ERR][Borrow][saveChatSessionDetail]", err)
		return ms.Error()
	}

	return ms.sendMessage(types.MessageRequest{
		Text: "Apa alasan Anda meminjam barang ini?",
	})
}

func (ms *MessageService) borrowReason() error {
	borrow, err := ms.borrowService.FindInitialByUserID(ms.user.ID)
	if err != nil {
		log.Println("[ERR][borrowReason][FindInitialByUserID]", err)
		return ms.Error()
	}

	tool, err := ms.toolService.FindByID(borrow.ToolID)
	if err != nil {
		log.Println("[ERR][Borrow][FindByID]", err)
		return ms.Error()
	}

	if err = ms.borrowService.UpdateBorrowReason(borrow.ID, ms.messageText); err != nil {
		log.Println("[ERR][borrowReason][UpdateBorrowReason]", err)
		return ms.Error()
	}

	if err = ms.saveChatSessionDetail(types.Topic["borrow_reason"], ""); err != nil {
		log.Println("[ERR][borrowReason][saveChatSessionDetail]", err)
		return ms.Error()
	}

	returnDateStr := helper.ChangeDateStringFormat(borrow.ReturnDate.String)
	returnDateTime, err := time.Parse(helper.BasicDateLayout, returnDateStr)
	if err != nil {
		log.Println("[ERR][borrowReason][Parse]")
		return ms.Error()
	}

	currentDateTime := time.Now()
	borrowDateRange := int(returnDateTime.Sub(currentDateTime).Hours()/24) + 1 // because current day is counted

	message := fmt.Sprintf(`Nama alat : %s
		Jumlah : %d
		Tanggal Pengembalian : %s (%d hari)
		Alamat peminjam : %s

		Pastikan data sudah benar. Tekan "Lanjutkan" untuk mengajukan ke pengurus.
	`, tool.Name, 1, helper.TranslateDateStringToBahasa(borrow.ReturnDate.String), borrowDateRange, ms.user.Address)
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
		return ms.Error()
	}

	var message string
	if userResponse {
		borrow.Status = types.GetBorrowStatus("request")
		message = "Pengajuan peminjaman berhasil, silahkan tunggu hingga pengurus menanggapi pengajuan."
	} else {
		borrow.Status = types.GetBorrowStatus("cancel")
		message = "Pengajuan berhasil dibatalkan"
	}

	_, err = ms.borrowService.UpdateBorrow(borrow)
	if err != nil {
		log.Println("[ERR][borrowConfirm][UpdateBorrow]", err)
		return ms.Error()
	}

	if userResponse {
		go ms.sendBorrowToAdmin(borrow)
	}

	return ms.sendMessage(types.MessageRequest{
		Text: message,
	})
}

func (ms *MessageService) sendBorrowToAdmin(borrow types.Borrow) error {
	message := fmt.Sprintf(`Seseorang baru saja mengajukan peminjaman barang

	ID Pengajuan: %d
	Nama Barang: %s
	
	Anda dapat menanggapi pengajuan ini dengan mengetik perintah "/%s %s %d"`,
		borrow.ID, borrow.Tool.Name, types.CommandRespond, types.RespondTypeBorrow, borrow.ID)
	message = helper.RemoveTab(message)

	return ms.sendMessage(types.MessageRequest{
		ChatID: types.AdminGroupIDs[0],
		Text:   message,
	})
}

func (ms *MessageService) ReturnTool() error {
	user, err := ms.userService.FindByID(ms.user.ID)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ERR][Borrow]", err)
		return err
	}

	ms.user = user
	if !ms.user.IsRegistered() {
		return ms.notRegistered()
	}

	if ok := isFlagWithinReturningCommand(ms.messageText); ok {
		return ms.toolReturningInit()
	}

	if len(ms.chatSessionDetails) > 0 {
		switch ms.chatSessionDetails[0].Topic {
		case types.Topic["tool_returning_init"]:
			return ms.toolReturningConfirm()
		case types.Topic["tool_returning_confirm"]:
			return ms.toolReturningComplete()
		}
	}

	return ms.currentlyBorrowedTools()
}

func isFlagWithinReturningCommand(s string) bool {
	ss := strings.Split(s, " ")
	if len(ss) != 2 {
		return false
	}

	if strings.Contains(ss[1], types.ToolReturningFlag) {
		return true
	}

	return false
}

/* func (ms *MessageService) borrowedTools() error {
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
} */

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
		ReplyMarkup: types.InlineKeyboardMarkup{
			InlineKeyboard: [][]types.InlineKeyboardButton{
				{
					{
						Text:         "Ajukan Pengembalian",
						CallbackData: fmt.Sprintf("/%s %s", types.CommandReturn, types.ToolReturningFlag),
					},
				},
			},
		},
	}

	return ms.sendMessage(reqBody)
}

func (ms *MessageService) toolReturningInit() error {
	_, err := ms.borrowService.FindCurrentlyBeingBorrowedByUserID(ms.user.ID)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return ms.Error()
	}

	if err == sql.ErrNoRows {
		return ms.sendMessage(types.MessageRequest{
			Text: "Saat ini tidak ada alat yang sedang Anda pinjam.",
		})
	}

	_, err = ms.toolReturningService.FindOnProgressByUserID(ms.user.ID)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return ms.Error()
	}

	if err != sql.ErrNoRows {
		return ms.sendMessage(types.MessageRequest{
			Text: "Maaf, Anda sudah mengajukan pengembalian sebelumnya. Silahkan tunggu hingga pengurus menanggapi pengajuan tersebut.",
		})
	}

	if err := ms.saveChatSessionDetail(types.Topic["tool_returning_init"], ""); err != nil {
		log.Println("[ERR][toolReturningInit][saveChatSessionDetail]", err)
		return err
	}

	return ms.sendMessage(types.MessageRequest{
		Text: "Tulis keterangan pengembalian. Dapat berupa kondisi barang, alasan pengembalian, dsb.",
	})
}

func (ms *MessageService) toolReturningConfirm() error {
	borrow, err := ms.borrowService.FindCurrentlyBeingBorrowedByUserID(ms.user.ID)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err)
		return err
	}

	if err == sql.ErrNoRows {
		reqBody := types.MessageRequest{
			Text: "Saat ini tidak ada alat yang sedang Anda pinjam.",
		}

		return ms.sendMessage(reqBody)
	}

	message := fmt.Sprintf(`
		Apakah Anda yakin data ini sudah benar?

		Nama peminjam: %s
		Nama barang: %s
		Dipinjam sejak: %s
		Tanggal pengembalian: %s
		Keterangan:
		%s
	`, ms.user.Name, borrow.Tool.Name, helper.TranslateDateStringToBahasa(borrow.CreatedAt), helper.TranslateDateToBahasa(time.Now()), ms.messageText)
	message = helper.RemoveTab(message)

	errs, _ := errgroup.WithContext(context.Background())
	errs.Go(func() error {
		sessionDataGenerator := helper.NewSessionDataGenerator()
		generatedSessionData := sessionDataGenerator.ToolReturningConfirm(ms.messageText)
		return ms.saveChatSessionDetail(types.Topic["tool_returning_confirm"], generatedSessionData)
	})

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

	err = errs.Wait()
	if err != nil {
		log.Println("[ERR][toolReturningConfirm][saveChatSessionDetail]", err)
		return err
	}

	return ms.sendMessage(reqBody)
}

func (ms *MessageService) toolReturningComplete() error {
	var userResponse bool

	chatSessionID := ms.chatSessionDetails[0].ChatSessionID

	if ms.messageText == "yes" {
		userResponse = true
	} else {
		userResponse = false
	}

	sessionDataGenerator := helper.NewSessionDataGenerator()
	generatedSessionData := sessionDataGenerator.ToolReturningComplete(userResponse)
	err := ms.saveChatSessionDetail(types.Topic["tool_returning_complete"], generatedSessionData)
	if err != nil {
		log.Println("[ERR][ReturnTool[saveChatSessionDetail]", err)
		return err
	}

	if err := ms.chatSessionService.UpdateChatSessionStatus(chatSessionID, types.ChatSessionStatus["complete"]); err != nil {
		log.Println("[ERR][ReturnToole][UpdateChatSessionStatus]", err)
		return err
	}

	errs, _ := errgroup.WithContext(context.Background())
	if userResponse {
		errs.Go(func() error {
			return ms.toolReturningCompletePositive()
		})
	} else {
		errs.Go(func() error {
			return ms.toolReturningCompleteNegative()
		})
	}

	err = errs.Wait()
	if err != nil {
		log.Println("[ERR][ReturnTool][toolReturningComplete]", err)
		return err
	}

	return nil
}

func (ms *MessageService) toolReturningCompletePositive() error {
	borrow, err := ms.borrowService.FindCurrentlyBeingBorrowedByUserID(ms.user.ID)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	var additionalInfo string
	confirmationChatSessionDetail, found := helper.GetChatSessionDetailByTopic(ms.chatSessionDetails, types.Topic["tool_returning_confirm"])
	if found {
		dataParsed, err := gabs.ParseJSON([]byte(confirmationChatSessionDetail.Data))
		if err != nil {
			return err
		}
		value, ok := dataParsed.Path("additional_info").Data().(string)
		if ok {
			additionalInfo = value
		}
	}

	toolReturning := types.ToolReturning{
		UserID:         ms.user.ID,
		ToolID:         borrow.ToolID,
		Status:         types.GetToolReturningStatus("request"),
		AdditionalInfo: additionalInfo,
	}

	toolReturning, err = ms.toolReturningService.SaveToolReturning(toolReturning)
	if err != nil {
		return err
	}

	// get all value along with tools and users value in ToolReturning
	toolReturning, err = ms.toolReturningService.FindToolReturningByID(toolReturning.ID)
	if err != nil {
		return err
	}

	go ms.sendToolReturningToAdmin(toolReturning)

	reqBody := types.MessageRequest{
		Text: "Pengajuan pengembalian berhasil, silahkan tunggu hingga pengurus menanggapi pengajuan tersebut.",
	}

	return ms.sendMessage(reqBody)
}

func (ms *MessageService) toolReturningCompleteNegative() error {
	reqBody := types.MessageRequest{
		Text: "Pengajuan pengembalian berhasil dibatalkan.",
	}

	return ms.sendMessage(reqBody)
}

func (ms *MessageService) sendToolReturningToAdmin(toolReturning types.ToolReturning) error {
	message := fmt.Sprintf(`Seseorang baru saja mengajukan pengembalian barang
	
	ID Pengajuan: %d
	Nama Barang: %s

	Anda dapat menanggapi pengajuan ini dengan mengetik perintah "/%s %s %d"`, toolReturning.ID, toolReturning.Tool.Name, types.CommandRespond, types.RespondTypeToolReturning, toolReturning.ID)

	return ms.sendMessage(types.MessageRequest{
		ChatID: types.AdminGroupIDs[0],
		Text:   message,
	})
}

/**
*
* Admin handlers
*
 */

func (ms *MessageService) isEligibleAdmin() bool {
	if ms.requestType != types.RequestTypeGroup {
		return false
	}

	adminIDs := []int64{284324420}
	for _, id := range adminIDs {
		if ms.user.ID == id {
			return true
		}
	}
	return false
}

func (ms *MessageService) Respond() error {
	if ok := ms.isEligibleAdmin(); !ok {
		log.Println("[INFO] Not eligible user accessing admin command", ms.messageText)
		return ms.Unknown()
	}

	respCommands, ok := helper.GetRespondCommands(ms.messageText)
	if !ok {
		return ms.ListToRespond()
	}

	if respCommands.Text != "yes" && respCommands.Text != "no" && respCommands.Text != "" {
		return ms.sendMessage(types.MessageRequest{
			Text: "Maaf, perintah tidak dikenali. Pilihan yang tersedia adalah \"yes\" dan \"no\"",
		})
	}

	if respCommands.Type == types.RespondTypeBorrow {
		return ms.respondBorrow(respCommands)
	} else if respCommands.Type == types.RespondTypeToolReturning {
		return ms.respondToolReturning(respCommands)
	}

	return ms.Unknown()
}

func (ms *MessageService) ListToRespond() error {
	var message string

	borrows, err := ms.borrowService.GetBorrowRequests()
	if err != nil {
		log.Println("[ERR][ListToRespond][GetBorrowRequests]", err)
		return ms.Error()
	}

	toolRets, err := ms.toolReturningService.GetToolReturningRequests()
	if err != nil {
		log.Println("[ERR][ListToRespond][GetToolReturningRequests]", err)
		return ms.Error()
	}

	message += "Daftar Pengajuan Peminjaman\n"
	if len(borrows) > 0 {
		message += helper.BuildBorrowRequestListMessage(borrows)
	} else {
		message += "- tidak ada\n"
	}

	message += "\nDaftar Pengajuan Pengembalian\n"
	if len(toolRets) > 0 {
		message += helper.BuildToolReturningRequestListMessage(toolRets)
	} else {
		message += "- tidak ada\n"
	}

	message += fmt.Sprintf("\n\nUntuk menanggapi pengajuan ketik perintah \"/%s [pinjam/kembali] [id]\"", types.CommandRespond)
	message += fmt.Sprintf("\ncontoh: \"/%s pinjam 173\"", types.CommandRespond)

	return ms.sendMessage(types.MessageRequest{
		Text: message,
	})
}

func (ms *MessageService) respondBorrow(commands types.RespondCommands) error {
	borrow, err := ms.borrowService.FindBorrowByID(commands.ID)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ERR][respondBorrow][FindBorrowByID]", err)
		return ms.Error()
	}

	if err == sql.ErrNoRows || borrow.Status != types.GetBorrowStatus("request") {
		return ms.sendMessage(types.MessageRequest{
			Text: "Gagal menanggapi, ID tidak ditemukan.",
		})
	}

	if commands.Text == "yes" {
		return ms.respondBorrowPositive(borrow)
	} else if commands.Text == "no" {
		return ms.respondBorrowNegative(borrow)
	}

	return ms.respondBorrowDetail(borrow)
}

func (ms *MessageService) respondBorrowDetail(borrow types.Borrow) error {
	message := fmt.Sprintf(`
		ID: %d
		Nama pemohon: %s (%s)
		Barang: %s
		Diajukan pada: %s
		Estimasi pengembalian: %s
		Alasan peminjaman:
		%s
	`, borrow.ID, borrow.User.Name, borrow.User.NIM, borrow.Tool.Name, helper.TranslateDateStringToBahasa(borrow.CreatedAt), helper.TranslateDateStringToBahasa(borrow.ReturnDate.String), borrow.Reason.String)
	message = helper.RemoveTab(message)

	return ms.sendMessage(types.MessageRequest{
		Text: message,
		ReplyMarkup: types.InlineKeyboardMarkup{
			InlineKeyboard: [][]types.InlineKeyboardButton{
				{
					{
						Text:         "Setujui",
						CallbackData: fmt.Sprintf("/%s %s %d yes", types.CommandRespond, types.RespondTypeBorrow, borrow.ID),
					},
					{
						Text:         "Tolak",
						CallbackData: fmt.Sprintf("/%s %s %d no", types.CommandRespond, types.RespondTypeBorrow, borrow.ID),
					},
				},
			},
		},
	})
}

func (ms *MessageService) respondBorrowPositive(borrow types.Borrow) error {
	if borrow.Tool.Stock < 1 {
		return ms.sendMessage(types.MessageRequest{
			Text: "Stok barang sudah habis. Tidak dapat menyetujui peminjaman.",
		})
	}

	if err := ms.toolService.DecreaseStock(borrow.ToolID); err != nil {
		log.Println("[ERR][respondBorrowPositive][DecreaseStock]", err)
		return ms.Error()
	}

	borrow.Status = types.GetBorrowStatus("progress")
	if _, err := ms.borrowService.UpdateBorrow(borrow); err != nil {
		log.Println("[ERR][respondBorrowPositive][UpdateBorrow]", err)
		return ms.Error()
	}

	reqBody := types.MessageRequest{
		ChatID: borrow.UserID,
		Text:   fmt.Sprintf("Pengajuan peminjaman \"%s\" telah disetujui oleh pengurus.", borrow.Tool.Name),
	}
	if err := ms.sendMessage(reqBody); err != nil {
		log.Println("error in sending reply:", err)
		return err
	}

	return ms.sendMessage(types.MessageRequest{
		Text: "Pengajuan peminjaman berhasil disetujui.",
	})
}

func (ms *MessageService) respondBorrowNegative(borrow types.Borrow) error {
	borrow.Status = types.GetBorrowStatus("reject")
	if _, err := ms.borrowService.UpdateBorrow(borrow); err != nil {
		log.Println("[ERR][respondBorrowNegative][UpdateBorrow]", err)
		return ms.Error()
	}

	reqBody := types.MessageRequest{
		ChatID: borrow.UserID,
		Text:   fmt.Sprintf("Pengajuan peminjaman \"%s\" telah ditolak oleh pengurus.", borrow.Tool.Name),
	}
	if err := ms.sendMessage(reqBody); err != nil {
		log.Println("error in sending reply:", err)
		return err
	}

	return ms.sendMessage(types.MessageRequest{
		Text: "Pengajuan peminjaman berhasil ditolak.",
	})
}

func (ms *MessageService) respondToolReturning(commands types.RespondCommands) error {
	toolReturning, err := ms.toolReturningService.FindToolReturningByID(commands.ID)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ERR][respondToolReturning][FindToolReturningByID]", err)
		return ms.Error()
	}

	if err == sql.ErrNoRows || toolReturning.Status != types.GetToolReturningStatus("request") {
		return ms.sendMessage(types.MessageRequest{
			Text: "Gagal menanggapi, ID tidak ditemukan.",
		})
	}

	if commands.Text == "yes" {
		return ms.respondToolReturningPositive(toolReturning)
	} else if commands.Text == "no" {
		return ms.respondToolReturningNegative(toolReturning)
	}

	return ms.respondToolReturningDetail(toolReturning)
}

func (ms *MessageService) respondToolReturningDetail(toolReturning types.ToolReturning) error {
	message := fmt.Sprintf(`
		ID: %d
		Nama pemohon: %s (%s)
		Barang: %s
		Diajukan pada: %s
		Keterangan:
		%s
	`, toolReturning.ID, toolReturning.User.Name, toolReturning.User.NIM, toolReturning.Tool.Name, helper.TranslateDateStringToBahasa(toolReturning.CreatedAt), toolReturning.AdditionalInfo)
	message = helper.RemoveTab(message)

	return ms.sendMessage(types.MessageRequest{
		Text: message,
		ReplyMarkup: types.InlineKeyboardMarkup{
			InlineKeyboard: [][]types.InlineKeyboardButton{
				{
					{
						Text:         "Setujui",
						CallbackData: fmt.Sprintf("/%s %s %d yes", types.CommandRespond, types.RespondTypeToolReturning, toolReturning.ID),
					},
					{
						Text:         "Tolak",
						CallbackData: fmt.Sprintf("/%s %s %d no", types.CommandRespond, types.RespondTypeToolReturning, toolReturning.ID),
					},
				},
			},
		},
	})
}

func (ms *MessageService) respondToolReturningPositive(toolReturning types.ToolReturning) error {
	if err := ms.toolService.IncreaseStock(toolReturning.ToolID); err != nil {
		log.Println("[ERR][respondToolReturningPositive][IncreaseStock]", err)
		return ms.Error()
	}

	borrow, err := ms.borrowService.FindCurrentlyBeingBorrowedByUserID(toolReturning.UserID)
	if err != nil {
		log.Println("[ERR][respondToolReturningPositive][FindCurrentlyBeingBorrowedByUserID]", err)
		return ms.Error()
	}

	borrow.Status = types.GetBorrowStatus("returned")
	if _, err := ms.borrowService.UpdateBorrow(borrow); err != nil {
		log.Println("[ERR][sendToolReturningToAdmin][UpdateBorrow]", err)
		return ms.Error()
	}

	if err := ms.toolReturningService.UpdateToolReturningStatus(toolReturning.ID, types.GetToolReturningStatus("complete")); err != nil {
		log.Println("[ERR][sendToolReturningToAdmin][UpdateToolReturningStatus]", err)
		return ms.Error()
	}

	if err := ms.toolReturningService.UpdateToolReturningConfirmedAt(toolReturning.ID, time.Now()); err != nil {
		log.Println("[ERR][respondToolReturningNegative][UpdateToolReturningConfirmedAt]", err)
		return ms.Error()
	}

	reqBody := types.MessageRequest{
		ChatID: toolReturning.UserID,
		Text:   fmt.Sprintf("Pengajuan pengembalian \"%s\" telah disetujui oleh pengurus.", toolReturning.Tool.Name),
	}
	if err := ms.sendMessage(reqBody); err != nil {
		log.Println("error in sending reply:", err)
		return err
	}

	return ms.sendMessage(types.MessageRequest{
		Text: "Pengajuan pengembalian berhasil disetujui.",
	})
}

func (ms *MessageService) respondToolReturningNegative(toolReturning types.ToolReturning) error {
	if err := ms.toolReturningService.UpdateToolReturningStatus(toolReturning.ID, types.GetToolReturningStatus("reject")); err != nil {
		log.Println("[ERR][respondToolReturningNegative][UpdateToolReturningStatus]", err)
		return ms.Error()
	}

	if err := ms.toolReturningService.UpdateToolReturningConfirmedAt(toolReturning.ID, time.Now()); err != nil {
		log.Println("[ERR][respondToolReturningNegative][UpdateToolReturningConfirmedAt]", err)
		return ms.Error()
	}

	reqBody := types.MessageRequest{
		ChatID: toolReturning.UserID,
		Text:   fmt.Sprintf("Pengajuan pengembalian \"%s\" telah ditolak oleh pengurus.", toolReturning.Tool.Name),
	}
	if err := ms.sendMessage(reqBody); err != nil {
		log.Println("error in sending reply:", err)
		return err
	}

	return ms.sendMessage(types.MessageRequest{
		Text: "Pengajuan pengembalian berhasil ditolak",
	})
}
