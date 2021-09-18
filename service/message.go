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
	message            types.TeleMessage
	requestType        types.RequestType
	user               types.User
	chatSessionDetails []types.ChatSessionDetail

	chatSessionService   *ChatSessionService
	userService          *UserService
	toolService          *ToolService
	borrowService        *BorrowService
	toolReturningService *ToolReturningService
}

func NewMessageService(chatID, senderID int64, text string, requestType types.RequestType, teleMessage types.TeleMessage) *MessageService {
	ms := &MessageService{
		chatID:      chatID,
		messageText: text,
		requestType: requestType,
		message:     teleMessage,
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

func (ms *MessageService) isRegistered() bool {
	user, err := ms.userService.FindByID(ms.user.ID)
	if err != nil && err != sql.ErrNoRows {
		return false
	}

	ms.user = user
	return ms.user.IsRegistered()
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

	res, err := http.Post(fmt.Sprintf("%s/sendMessage", config.WebhookUrl()), "application/json", bytes.NewBuffer(reqBytes))
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
		return errors.New("unexpected status " + res.Status)
	}

	return nil
}

func (ms *MessageService) sendPhoto(reqBody types.PhotoRequest) error {
	if reqBody.ChatID == 0 {
		reqBody.ChatID = ms.chatID
	}

	reqBytes, err := json.Marshal(&reqBody)
	if err != nil {
		return err
	}

	res, err := http.Post(fmt.Sprintf("%s/sendPhoto", config.WebhookUrl()), "application/json", bytes.NewBuffer(reqBytes))
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
		return errors.New("unexpected status " + res.Status)
	}

	return nil
}

func (ms *MessageService) sendPhotoGroup(reqBody types.PhotoGroupRequest) error {
	if reqBody.ChatID == 0 {
		reqBody.ChatID = ms.chatID
	}

	reqBytes, err := json.Marshal(&reqBody)
	if err != nil {
		return err
	}

	res, err := http.Post(fmt.Sprintf("%s/sendMediaGroup", config.WebhookUrl()), "application/json", bytes.NewBuffer(reqBytes))
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

		return errors.New("unexpected status " + res.Status)
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

func (ms *MessageService) FirstStart() error {
	message := "Selamat Datang di Layanan Chatbot Peminjaman Barang Laboratorium Teknik Komputer!"
	message += fmt.Sprintf("\n\nSilahkan lakukan registrasi dengan mengirim perintah \"/%s\"", types.CommandRegister)
	return ms.sendMessage(types.MessageRequest{
		Text: message,
	})
}

func (ms *MessageService) Help() error {
	message := fmt.Sprintf(`/%s - Mendaftarkan diri agar dapat menggunakan sistem
		/%s - Cek ketersediaan barang
		/%s - Mulai pengajuan peminjaman barang
		/%s - Mulai pengajuan Pengembalian barang
		/%s - Menampilkan panduan penggunaan bot`, types.CommandRegister, types.CommandCheck, types.CommandBorrow, types.CommandReturn, types.CommandHelp)

	if ms.isEligibleAdmin() {
		message = fmt.Sprintf(`/%s - Cek ketersediaan barang
			/%s - Menanggapi pengajuan peminjaman dan pengembalian barang
			/%s - Menambah dan mengubah data barang
			/%s - Melihat laporan bulanan
			/%s - Menampilkan panduan penggunaan bot`, types.CommandCheck, types.CommandRespond, types.CommandManage, types.CommandReport, types.CommandHelp)
	}

	return ms.sendMessage(types.MessageRequest{
		Text: helper.RemoveTab(message),
	})
}

func (ms *MessageService) Unknown() error {
	reqBody := types.MessageRequest{
		Text: "Maaf, perintah tidak dikenali.",
	}
	return ms.sendMessage(reqBody)
}

func (ms *MessageService) Check() error {
	checkCommandOrder, ok := helper.GetCheckCommandOrder(ms.messageText)
	if ok {
		if checkCommandOrder.Text == types.CheckTypePhoto {
			return ms.checkDetailPhoto(checkCommandOrder.ID)
		}

		return ms.checkDetail(checkCommandOrder.ID)
	}

	var tools []types.Tool
	var err error
	if ms.isEligibleAdmin() {
		tools, err = ms.toolService.GetTools()
	} else {
		tools, err = ms.toolService.GetAvailableTools()
	}

	if err != nil {
		log.Println("[ERR][Check][GetTools]", err)
		return ms.Error()
	}

	if len(tools) < 1 {
		return ms.sendMessage(types.MessageRequest{
			Text: "Tidak ada barang yang tersedia.",
		})
	}

	message := "Berikut ini daftar alat yang masih tersedia.\n"
	message += fmt.Sprintf("untuk melihat detail alat, ketik perintah \"/%s [id]\"\n\n", types.CommandCheck)
	message += helper.BuildToolListMessage(tools)

	return ms.sendMessage(types.MessageRequest{
		Text: message,
	})
}

func (ms *MessageService) checkDetail(toolID int64) error {
	tool, err := ms.toolService.FindByID(toolID)
	if err != nil {
		log.Println("[ERR][checkDetail][FindByID]", err)
		return ms.sendMessage(types.MessageRequest{
			Text: "Maaf, nomor alat yang Anda pilih tidak tersedia.",
		})
	}

	if err == sql.ErrNoRows {
		return ms.sendMessage(types.MessageRequest{
			Text: "ID tidak ditemukan.",
		})
	}

	if tool.Stock < 1 && !ms.isEligibleAdmin() {
		return ms.sendMessage(types.MessageRequest{
			Text: "Maaf, nomor alat yang Anda pilih tidak tersedia",
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

	var inlineKeyboard [][]types.InlineKeyboardButton
	if ms.isEligibleAdmin() {
		inlineKeyboard = [][]types.InlineKeyboardButton{
			{{
				Text:         "Lihat Foto",
				CallbackData: fmt.Sprintf("/%s %d %s", types.CommandCheck, tool.ID, types.CheckTypePhoto),
			}},
			{{
				Text:         "Ubah Data",
				CallbackData: fmt.Sprintf("/%s %s %d", types.CommandManage, types.ManageTypeEdit, tool.ID),
			}},
			{{
				Text:         "Hapus",
				CallbackData: fmt.Sprintf("/%s %s %d", types.CommandManage, types.ManageTypeDelete, tool.ID),
			}},
		}
	} else {
		inlineKeyboard = [][]types.InlineKeyboardButton{
			{{
				Text:         "Lihat Foto",
				CallbackData: fmt.Sprintf("/%s %d %s", types.CommandCheck, tool.ID, types.CheckTypePhoto),
			}},
			{{
				Text:         "Pinjam",
				CallbackData: fmt.Sprintf("/%s %d", types.CommandBorrow, tool.ID),
			}},
		}
	}

	reqBody := types.MessageRequest{
		Text: message,
		ReplyMarkup: types.InlineKeyboardMarkup{
			InlineKeyboard: inlineKeyboard,
		},
	}

	return ms.sendMessage(reqBody)
}

func (ms *MessageService) checkDetailPhoto(toolID int64) error {
	tool, err := ms.toolService.FindByID(toolID)
	if err != nil {
		log.Println("[ERR][checkDetailPhoto][FindByID]", err)
		return ms.sendMessage(types.MessageRequest{
			Text: "Maaf, nomor alat yang Anda pilih tidak tersedia.",
		})
	}

	if err == sql.ErrNoRows {
		return ms.sendMessage(types.MessageRequest{
			Text: "ID tidak ditemukan.",
		})
	}

	if tool.Stock < 1 && !ms.isEligibleAdmin() {
		return ms.sendMessage(types.MessageRequest{
			Text: "Maaf, nomor alat yang Anda pilih tidak tersedia.",
		})
	}

	photos, err := ms.toolService.GetPhotos(toolID)
	if err != nil {
		log.Println("[ERR][checkDetailPhoto][GetPhotos]", err)
		return ms.Error()
	}

	if len(photos) > 1 {
		var media []types.InputMediaPhoto
		for _, photo := range photos {
			temp := types.InputMediaPhoto{
				Type:  "photo",
				Media: photo.FileID,
			}
			media = append(media, temp)
		}

		return ms.sendPhotoGroup(types.PhotoGroupRequest{
			Media: media,
		})
	}

	if len(photos) == 1 {
		return ms.sendPhoto(types.PhotoRequest{
			Photo: photos[0].FileID,
		})
	}

	return ms.sendMessage(types.MessageRequest{
		Text: "Foto tidak tersedia untuk barang ini.",
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
		log.Println("[ERR][Register][FindByID]", err)
		return ms.Error()
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
		log.Println("[ERR][registerIniti][saveChatSessionDetail]", err)
		return err
	}

	return nil
}

func (ms *MessageService) registerConfirm() error {
	registrationMessage, err := getRegistrationMessage(ms.messageText)
	if err != nil {
		log.Println("[ERR][registerConfirm][getRegistrationMessage]", err)
		reqBody := types.MessageRequest{
			Text: "Format registrasi salah, mohon cek format kembali kemudian kirim ulang.",
		}
		return ms.sendMessage(reqBody)
	}

	err = validateRegisterConfirmation(registrationMessage)
	if err != nil {
		log.Println("[ERR][registerConfirm][Validation]", err)
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
		log.Println("[ERR][registerConfirm[UpdateUser]", err)
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
		log.Println("[ERR][registerConfirm][saveChatSessionDetail]", err)
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
		log.Println("[ERR][registerComplete][registerComplete]", err)
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
		Text: "Registrasi dibatalkan.",
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
		log.Println("[ERR][Borrow][FindByID]", err)
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
			return ms.borrowAmount()
		case types.Topic["borrow_amount"]:
			return ms.borrowDuration()
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
		log.Println("[ERR][borrowMechanism][sendMessage]", err)
		return err
	}

	return ms.Check()
}

func (ms *MessageService) borrowInit(toolID int64) error {
	tool, err := ms.toolService.FindByID(toolID)
	if err != nil {
		log.Println("[ERR][borrowInit][FindByID]", err)
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
		log.Println("[ERR][borrowInit][GetCurrentlyBeingBorrowedAndRequestedByUserID]", err)
		return ms.Error()
	}

	borrowStatus, same := helper.GetSameBorrow(borrows, toolID)

	if same {
		var message string
		if borrowStatus == types.GetBorrowStatus("request") {
			message = "Maaf, Anda sudah mengajukan peminjaman barang yang sama, silahkan tunggu hingga pengurus menanggapi pengajuan tersebut."
		} else {
			message = "Maaf, Anda sedang meminjam barang yang sama sehingga tidak dapat mengajukan peminjaman.\n"
			message += fmt.Sprintf(`Untuk melakukan pengembalian silahkan ketik "/%s"`, types.CommandReturn)
		}

		return ms.sendMessage(types.MessageRequest{
			Text: message,
		})
	}

	sessionDataGenerator := helper.NewSessionDataGenerator()
	generatedSessionData := sessionDataGenerator.BorrowInit(tool.ID)

	if err = ms.saveChatSessionDetail(types.Topic["borrow_init"], generatedSessionData); err != nil {
		log.Println("[ERR][borrowInit][saveChatSessionDetail]", err)
		return err
	}

	return ms.sendMessage(types.MessageRequest{
		Text: "Berapa jumlah yang ingin dipinjam?\n\nJika tidak ada dalam pilihan, maka sebutkan dalam angka (min. 1).",
		ReplyMarkup: types.InlineKeyboardMarkup{
			InlineKeyboard: [][]types.InlineKeyboardButton{
				{
					{
						Text:         "1",
						CallbackData: "1",
					},
					{
						Text:         "2",
						CallbackData: "2",
					},
					{
						Text:         "3",
						CallbackData: "3",
					},
				},
			},
		},
	})
}

func (ms *MessageService) borrowAmount() error {
	amount, err := strconv.Atoi(ms.messageText)
	if err != nil || amount < 1 {
		log.Println("[ERR][borrowAmount][Atoi]", err)
		return ms.sendMessage(types.MessageRequest{
			Text: "Mohon sebutkan jumlah barang dalam angka.",
		})
	}

	borrowSession := helper.GetBorrowFromChatSessionDetail(ms.chatSessionDetails)
	tool, err := ms.toolService.FindByID(borrowSession.ToolID)
	if err != nil {
		log.Println("[ERR][borrowAmount][FindByID]", err)
		return ms.Error()
	}

	if int64(amount) > tool.Stock {
		return ms.sendMessage(types.MessageRequest{
			Text: fmt.Sprintf("Tidak bisa meminjam barang melebihi stok yang ada. Stok saat ini %d", tool.Stock),
		})
	}

	sessionDataGenerator := helper.NewSessionDataGenerator()
	generatedSessionData := sessionDataGenerator.BorrowAmount(amount)

	if err = ms.saveChatSessionDetail(types.Topic["borrow_amount"], generatedSessionData); err != nil {
		log.Println("[ERR][borrowAmount][saveChatSessionDetail]", err)
		return ms.Error()
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

func (ms *MessageService) borrowDuration() error {
	duration, err := helper.GetDurationValue(ms.messageText)
	if err != nil {
		log.Println("[ERR][borrowDuration][GetDurationValue]", err)
		return ms.sendMessage(types.MessageRequest{
			Text: "Mohon sebutkan jumlah hari.",
		})
	}

	if duration < types.BorrowMinimalDuration {
		return ms.sendMessage(types.MessageRequest{
			Text: fmt.Sprintf("Minimal durasi peminjaman adalah %d hari", types.BorrowMinimalDuration),
		})
	}

	sessionDataGenerator := helper.NewSessionDataGenerator()
	generatedSessionData := sessionDataGenerator.BorrowDuration(duration)

	if err = ms.saveChatSessionDetail(types.Topic["borrow_date"], generatedSessionData); err != nil {
		log.Println("[ERR][borrowDuration][saveChatSessionDetail]", err)
		return ms.Error()
	}

	return ms.sendMessage(types.MessageRequest{
		Text: "Apa alasan Anda meminjam barang ini?",
	})
}

func (ms *MessageService) borrowReason() error {
	borrow := helper.GetBorrowFromChatSessionDetail(ms.chatSessionDetails)

	tool, err := ms.toolService.FindByID(borrow.ToolID)
	if err != nil {
		log.Println("[ERR][borrowReason][FindByID]", err)
		return ms.Error()
	}

	sessionDataGenerator := helper.NewSessionDataGenerator()
	generatedSessionData := sessionDataGenerator.BorrowReason(ms.messageText)

	if err = ms.saveChatSessionDetail(types.Topic["borrow_reason"], generatedSessionData); err != nil {
		log.Println("[ERR][borrowReason][saveChatSessionDetail]", err)
		return ms.Error()
	}

	returnDate := time.Now().AddDate(0, 0, borrow.Duration).Format(types.BasicDateLayout)

	message := fmt.Sprintf(`Nama alat : %s
		Jumlah : %d
		Tanggal Pengembalian : %s (%d hari)
		Alamat peminjam : %s
		Alasan:
		%s

		Pastikan data sudah benar. Tekan "Lanjutkan" untuk mengajukan ke pengurus.
	`, tool.Name, borrow.Amount, helper.TranslateDateStringToBahasa(returnDate), borrow.Duration, ms.user.Address, ms.messageText)
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

	if !userResponse {
		return ms.sendMessage(types.MessageRequest{
			Text: "Pengajuan dibatalkan",
		})
	}

	borrowSession := helper.GetBorrowFromChatSessionDetail(ms.chatSessionDetails)

	borrow := types.Borrow{
		Amount:   borrowSession.Amount,
		Duration: borrowSession.Duration,
		Status:   types.GetBorrowStatus("request"),
		UserID:   ms.user.ID,
		ToolID:   borrowSession.ToolID,
		Reason:   borrowSession.Reason,
	}

	borrowID, err := ms.borrowService.SaveBorrow(borrow)
	if err != nil {
		log.Println("[ERR][borrowConfirm][SaveBorrow]", err)
		return ms.Error()
	}

	go ms.notifyBorrowRequestToAdmin(borrowID)

	return ms.sendMessage(types.MessageRequest{
		Text: "Pengajuan peminjaman berhasil, silahkan tunggu hingga pengurus menanggapi pengajuan.",
	})
}

func (ms *MessageService) notifyBorrowRequestToAdmin(borrowID int64) error {
	borrow, err := ms.borrowService.FindBorrowByID(borrowID)
	if err != nil {
		log.Println("[ERR][notifyBorrowRequestToAdmin][FindBorrowByID]", err)
		return ms.Error()
	}

	message := fmt.Sprintf(`Seseorang baru saja mengajukan peminjaman barang

	Nama Pemohon: %s
	Barang: %s`, borrow.User.Name, borrow.Tool.Name)
	message = helper.RemoveTab(message)

	return ms.sendMessage(types.MessageRequest{
		ChatID: helper.GetAdminGroupID(),
		Text:   message,
		ReplyMarkup: types.InlineKeyboardMarkup{
			InlineKeyboard: [][]types.InlineKeyboardButton{
				{
					{
						Text:         "Tanggapi",
						CallbackData: fmt.Sprintf("/%s %s %d", types.CommandRespond, types.RespondTypeBorrow, borrow.ID),
					},
				},
			},
		},
	})
}

func (ms *MessageService) ReturnTool() error {
	user, err := ms.userService.FindByID(ms.user.ID)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ERR][ReturnTool][FindByID]", err)
		return err
	}

	ms.user = user
	if !ms.user.IsRegistered() {
		return ms.notRegistered()
	}

	borrowID, ok := isIDWithinCommand(ms.messageText)
	if ok && borrowID > 0 {
		return ms.toolReturningInit(borrowID)
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
	borrows, err := ms.borrowService.GetCurrentlyBeingBorrowedByUserID(ms.user.ID)
	if err != nil {
		log.Println("[ERR][currentlyBorrowedTools][GetCurrentlyBeingBorrowedByUserID]", err)
		return err
	}

	if len(borrows) == 0 {
		return ms.sendMessage(types.MessageRequest{
			Text: "Saat ini tidak ada alat yang sedang Anda pinjam.",
		})
	}

	message := "Berikut ini daftar alat yang sedang Anda pinjam.\n\n"
	message += helper.BuildBorrowedMessage(borrows)
	message += fmt.Sprintf("\nUntuk mengajukan pengembalian ketik perintah\n\"/%s [id_peminjaman]\"\n\n", types.CommandReturn)

	return ms.sendMessage(types.MessageRequest{
		Text: message,
	})
}

func (ms *MessageService) toolReturningInit(borrowID int64) error {
	borrow, err := ms.borrowService.FindBorrowByID(borrowID)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ERR][toolReturningInit][FindBorrowByID]", err)
		return ms.Error()
	}

	if err == sql.ErrNoRows || borrow.Status != types.GetBorrowStatus("progress") || borrow.UserID != ms.user.ID {
		return ms.sendMessage(types.MessageRequest{
			Text: "ID peminjaman tidak ditemukan.",
		})
	}

	rets, err := ms.toolReturningService.GetCurrentlyBeingRequested(ms.user.ID, borrowID)
	if err != nil {
		log.Println("[ERR][toolReturningInit][GetCurrentlyBeingRequested]", err)
		return ms.Error()
	}

	if len(rets) > 0 {
		return ms.sendMessage(types.MessageRequest{
			Text: "Maaf, Anda sudah mengajukan pengembalian barang yang sama. Silahkan tunggu hingga pengurus menanggapi pengajuan tersebut.",
		})
	}

	sessionDataGenerator := helper.NewSessionDataGenerator()
	generatedSessionData := sessionDataGenerator.ToolReturningInit(borrowID)

	if err := ms.saveChatSessionDetail(types.Topic["tool_returning_init"], generatedSessionData); err != nil {
		log.Println("[ERR][toolReturningInit][saveChatSessionDetail]", err)
		return err
	}

	return ms.sendMessage(types.MessageRequest{
		Text: "Tulis keterangan pengembalian. Dapat berupa kondisi barang, alasan pengembalian, dsb.",
	})
}

func (ms *MessageService) toolReturningConfirm() error {
	toolReturningSession, found := helper.GetChatSessionDetailByTopic(ms.chatSessionDetails, types.Topic["tool_returning_init"])
	if !found {
		return ms.Unknown()
	}

	dataParsed, err := gabs.ParseJSON([]byte(toolReturningSession.Data))
	if err != nil {
		log.Println("[ERR][toolReturningConfirm][ParseJSON]", err)
		return ms.Error()
	}

	var borrowID int64
	bID, ok := dataParsed.Path("borrow_id").Data().(float64)
	if !ok {
		log.Println("[ERR][toolReturningConfirm][Path] borrow_id not found")
		return ms.Error()
	}

	borrowID = int64(bID)

	borrow, err := ms.borrowService.FindBorrowByID(borrowID)
	if err != nil {
		log.Println("[ERR][toolReturningConfirm][FindBorrowByID]", err)
		return ms.Error()
	}

	message := fmt.Sprintf(`Nama peminjam: %s
		Nama barang: %s
		Jumlah: %d
		Dipinjam sejak: %s
		Tanggal pengembalian: %s
		Keterangan:
		%s
	
	
		Pastikan data sudah benar kemudian tekan "Lanjutkan".`,
		ms.user.Name, borrow.Tool.Name, borrow.Amount, helper.TranslateDateStringToBahasa(borrow.ConfirmedAt.Time.Format(types.BasicDateLayout)), helper.TranslateDateToBahasa(time.Now()), ms.messageText)
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
		log.Println("[ERR][toolReturningComplete[saveChatSessionDetail]", err)
		return err
	}

	if err := ms.chatSessionService.UpdateChatSessionStatus(chatSessionID, types.ChatSessionStatus["complete"]); err != nil {
		log.Println("[ERR][toolReturningComplete][UpdateChatSessionStatus]", err)
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
		log.Println("[ERR][toolReturningComplete][toolReturningComplete]", err)
		return err
	}

	return nil
}

func (ms *MessageService) toolReturningCompletePositive() error {
	toolReturningSession, found := helper.GetChatSessionDetailByTopic(ms.chatSessionDetails, types.Topic["tool_returning_init"])
	if !found {
		return errors.New("session not found")
	}

	dataParsed, err := gabs.ParseJSON([]byte(toolReturningSession.Data))
	if err != nil {
		return err
	}

	var borrowID int64
	bID, ok := dataParsed.Path("borrow_id").Data().(float64)
	if !ok {
		return errors.New("borrow_id not found")
	}

	borrowID = int64(bID)

	_, err = ms.borrowService.FindBorrowByID(borrowID)
	if err != nil {
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
		BorrowID:       borrowID,
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

	go ms.notifyToolReturningRequestToAdmin(toolReturning)

	reqBody := types.MessageRequest{
		Text: "Pengajuan pengembalian berhasil, silahkan tunggu hingga pengurus menanggapi pengajuan tersebut.",
	}

	return ms.sendMessage(reqBody)
}

func (ms *MessageService) toolReturningCompleteNegative() error {
	reqBody := types.MessageRequest{
		Text: "Pengajuan pengembalian dibatalkan.",
	}

	return ms.sendMessage(reqBody)
}

func (ms *MessageService) notifyToolReturningRequestToAdmin(toolReturning types.ToolReturning) error {
	message := fmt.Sprintf(`Seseorang baru saja mengajukan pengembalian barang
	
	Nama Pemohon: %s
	Barang: %s`, toolReturning.Borrow.User.Name, toolReturning.Borrow.Tool.Name)

	return ms.sendMessage(types.MessageRequest{
		ChatID: helper.GetAdminGroupID(),
		Text:   message,
		ReplyMarkup: types.InlineKeyboardMarkup{
			InlineKeyboard: [][]types.InlineKeyboardButton{
				{
					{
						Text:         "Tanggapi",
						CallbackData: fmt.Sprintf("/%s %s %d", types.CommandRespond, types.RespondTypeToolReturning, toolReturning.ID),
					},
				},
			},
		},
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
	if !ms.isEligibleAdmin() {
		log.Println("[INFO] Not eligible user accessing admin command", ms.messageText)
		return ms.Unknown()
	}

	respCommands, ok := helper.GetRespondCommandOrder(ms.messageText)
	if !ok {
		return ms.ListToRespond()
	}

	if respCommands.Text != "yes" && respCommands.Text != "no" && respCommands.Text != "" {
		return ms.sendMessage(types.MessageRequest{
			Text: "Maaf, perintah tidak dikenali. Pilihan yang tersedia adalah \"yes\" dan \"no\"",
		})
	}

	if respCommands.Type == types.RespondTypeBorrow {
		return ms.respondBorrowInit(respCommands)
	} else if respCommands.Type == types.RespondTypeToolReturning {
		return ms.respondToolReturningInit(respCommands)
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

func (ms *MessageService) RespondBorrow() error {
	if !ms.isRegistered() || !ms.isEligibleAdmin() {
		log.Println("[INFO] Not eligible user accessing admin command", ms.messageText)
		return ms.Unknown()
	}

	if len(ms.chatSessionDetails) > 0 {
		switch ms.chatSessionDetails[0].Topic {
		case types.Topic["respond_borrow_init"]:
			return ms.respondBorrowComplete()
		}
	}

	return ms.Unknown()
}

func (ms *MessageService) respondBorrowInit(commands types.RespondCommandOrder) error {
	borrow, err := ms.borrowService.FindBorrowByID(commands.ID)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ERR][respondBorrowInit][FindBorrowByID]", err)
		return ms.Error()
	}

	if err == sql.ErrNoRows || borrow.Status != types.GetBorrowStatus("request") {
		return ms.sendMessage(types.MessageRequest{
			Text: "Gagal menanggapi, ID tidak ditemukan.",
		})
	}

	if commands.Text == "" {
		return ms.respondBorrowDetail(borrow)
	}

	if commands.Text == "yes" && borrow.Tool.Stock < int64(borrow.Amount) {
		return ms.sendMessage(types.MessageRequest{
			Text: fmt.Sprintf("Jumlah yang dipinjam melebihi stok yang ada. Stok saat ini %d", borrow.Tool.Stock),
		})
	}

	sessionDataGenerator := helper.NewSessionDataGenerator()
	generatedSessionData := sessionDataGenerator.RespondBorrowInit(borrow.ID, commands.Text)

	if err = ms.saveChatSessionDetail(types.Topic["respond_borrow_init"], generatedSessionData); err != nil {
		log.Println("[ERR][respondBorrowInit][saveChatSessionDetail]", err)
		return ms.Error()
	}

	return ms.sendMessage(types.MessageRequest{
		Text: "Tuliskan keterangan tambahan.",
	})
}

func (ms *MessageService) respondBorrowComplete() error {
	respondBorrowSession, ok := helper.GetChatSessionDetailByTopic(ms.chatSessionDetails, types.Topic["respond_borrow_init"])
	if !ok {
		return ms.Unknown()
	}

	dataParsed, err := gabs.ParseJSON([]byte(respondBorrowSession.Data))
	if err != nil {
		log.Println("[ERR][respondBorrowComplete][ParseJSON]", err)
		return ms.Error()
	}

	var borrowID int64
	bID, ok := dataParsed.Path("borrow_id").Data().(float64)
	if ok {
		borrowID = int64(bID)
	}

	borrow, err := ms.borrowService.FindBorrowByID(borrowID)
	if err != nil {
		log.Println("[ERR][respondBorrowComplete][FindBorrowByID]", err)
		return ms.Error()
	}

	userResponse, _ := dataParsed.Path("user_response").Data().(string)

	sessionDataGenerator := helper.NewSessionDataGenerator()
	generatedSessionData := sessionDataGenerator.RespondBorrowComplete(ms.messageText)

	if err := ms.saveChatSessionDetail(types.Topic["respond_borrow_complete"], generatedSessionData); err != nil {
		log.Println("[ERR][respondBorrowComplete][saveChatSessionDetail]", err)
		return ms.Error()
	}

	if err := ms.borrowService.UpdateBorrowConfirm(borrow.ID, time.Now(), ms.message.From.FirstName, ms.message.From.LastName); err != nil {
		log.Println("[ERR][respondBorrowComplete][UpdateBorrowConfirmedAt]", err)
		return ms.Error()
	}

	chatSessionID := ms.chatSessionDetails[0].ChatSessionID
	if err := ms.chatSessionService.UpdateChatSessionStatus(chatSessionID, types.ChatSessionStatus["complete"]); err != nil {
		log.Println("[ERR][respondBorrowComplete][UpdateChatSessionStatus]", err)
		return ms.Error()
	}

	if userResponse == "yes" {
		return ms.respondBorrowPositive(borrow)
	}

	return ms.respondBorrowNegative(borrow)
}

func (ms *MessageService) respondBorrowDetail(borrow types.Borrow) error {
	message := fmt.Sprintf(`
		ID: %d
		Nama pemohon: %s (%s)
		Barang: %s
		Jumlah: %d
		Diajukan pada: %s
		Durasi peminjaman: %d hari
		Alamat pemohon:
		%s

		Alasan peminjaman:
		%s
	`, borrow.ID, borrow.User.Name, borrow.User.NIM, borrow.Tool.Name, borrow.Amount, helper.TranslateDateStringToBahasa(borrow.CreatedAt), borrow.Duration, borrow.User.Address, borrow.Reason.String)
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
	if err := ms.borrowService.UpdateBorrowStatus(borrow.ID, types.GetBorrowStatus("progress")); err != nil {
		log.Println("[ERR][respondBorrowPositive][UpdateBorrowStatus]", err)
		return ms.Error()
	}

	if err := ms.toolService.DecreaseStock(borrow.ToolID, borrow.Amount); err != nil {
		log.Println("[ERR][respondBorrowPositive][DecreaseStock]", err)
		return ms.Error()
	}

	returnDate := time.Now().AddDate(0, 0, borrow.Duration)
	message := fmt.Sprintf(`Pengajuan peminjaman "%s" telah disetujui oleh pengurus.
		Batas akhir peminjaman: %s (%d hari)

		Keterangan:
		%s`, borrow.Tool.Name, helper.TranslateDateToBahasa(returnDate), borrow.Duration, ms.messageText)
	message = helper.RemoveTab(message)

	reqBody := types.MessageRequest{
		ChatID: borrow.UserID,
		Text:   message,
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
	if err := ms.borrowService.UpdateBorrowStatus(borrow.ID, types.GetBorrowStatus("reject")); err != nil {
		log.Println("[ERR][respondBorrowNegative][UpdateBorrowStatus]", err)
		return ms.Error()
	}

	reqBody := types.MessageRequest{
		ChatID: borrow.UserID,
		Text:   fmt.Sprintf("Pengajuan peminjaman \"%s\" telah ditolak oleh pengurus.\n\nKeterangan:\n%s", borrow.Tool.Name, ms.messageText),
	}
	if err := ms.sendMessage(reqBody); err != nil {
		log.Println("error in sending reply:", err)
		return err
	}

	return ms.sendMessage(types.MessageRequest{
		Text: "Pengajuan peminjaman berhasil ditolak.",
	})
}

func (ms *MessageService) RespondToolReturning() error {
	if !ms.isRegistered() || !ms.isEligibleAdmin() {
		log.Println("[INFO] Not eligible user accessing admin command", ms.messageText)
		return ms.Unknown()
	}

	if len(ms.chatSessionDetails) > 0 {
		switch ms.chatSessionDetails[0].Topic {
		case types.Topic["respond_tool_returning_init"]:
			return ms.respondToolReturningComplete()
		}
	}

	return ms.Unknown()
}

func (ms *MessageService) respondToolReturningInit(commands types.RespondCommandOrder) error {
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

	if commands.Text == "" {
		return ms.respondToolReturningDetail(toolReturning)
	}

	sessionDataGenerator := helper.NewSessionDataGenerator()
	generatedSessionData := sessionDataGenerator.RespondToolReturningInit(toolReturning.ID, commands.Text)

	if err = ms.saveChatSessionDetail(types.Topic["respond_tool_returning_init"], generatedSessionData); err != nil {
		log.Println("[ERR][respondToolReturningInit][saveChatSessionDetail]", err)
		return ms.Error()
	}

	return ms.sendMessage(types.MessageRequest{
		Text: "Tuliskan keterangan tambahan.",
	})
}

func (ms *MessageService) respondToolReturningComplete() error {
	respondToolReturningSession, ok := helper.GetChatSessionDetailByTopic(ms.chatSessionDetails, types.Topic["respond_tool_returning_init"])
	if !ok {
		return ms.Unknown()
	}

	dataParsed, err := gabs.ParseJSON([]byte(respondToolReturningSession.Data))
	if err != nil {
		log.Println("[ERR][respondToolReturningComplete][ParseJSON]", err)
		return ms.Error()
	}

	var toolReturningID int64
	trID, ok := dataParsed.Path("tool_returning_id").Data().(float64)
	if ok {
		toolReturningID = int64(trID)
	}

	toolReturning, err := ms.toolReturningService.FindToolReturningByID(toolReturningID)
	if err != nil {
		log.Println("[ERR][respondToolReturningComplete][FindToolReturningByID]", err)
		return ms.Error()
	}

	userResponse, _ := dataParsed.Path("user_response").Data().(string)

	sessionDataGenerator := helper.NewSessionDataGenerator()
	generatedSessionData := sessionDataGenerator.RespondToolReturningComplete(ms.messageText)

	if err := ms.saveChatSessionDetail(types.Topic["respond_tool_returning_complete"], generatedSessionData); err != nil {
		log.Println("[ERR][respondToolReturningComplete][saveChatSessionDetail]", err)
		return ms.Error()
	}

	if err := ms.toolReturningService.UpdateToolReturningConfirm(toolReturning.ID, time.Now(), ms.message.From.FirstName, ms.message.From.LastName); err != nil {
		log.Println("[ERR][respondToolReturningComplete][UpdateToolReturningConfirmedAt]", err)
		return ms.Error()
	}

	chatSessionID := ms.chatSessionDetails[0].ChatSessionID
	if err := ms.chatSessionService.UpdateChatSessionStatus(chatSessionID, types.ChatSessionStatus["complete"]); err != nil {
		log.Println("[ERR][respondToolReturningComplete][UpdateChatSessionStatus]", err)
		return ms.Error()
	}

	if userResponse == "yes" {
		return ms.respondToolReturningApprove(toolReturning)
	}

	return ms.respondToolReturningReject(toolReturning)
}

func (ms *MessageService) respondToolReturningDetail(toolReturning types.ToolReturning) error {
	message := fmt.Sprintf(`
		ID: %d
		Diajukan pada: %s
		Nama pemohon: %s (%s)
		Barang: %s
		Jumlah: %d
		Dipinjam sejak: %s
		Durasi peminjaman: %d hari
		Alamat peminjam:
		%s

		Keterangan:
		%s
	`, toolReturning.ID, helper.TranslateDateStringToBahasa(toolReturning.CreatedAt), toolReturning.Borrow.User.Name, toolReturning.Borrow.User.NIM, toolReturning.Borrow.Tool.Name, toolReturning.Borrow.Amount, helper.TranslateDateToBahasa(toolReturning.Borrow.ConfirmedAt.Time), toolReturning.Borrow.Duration, toolReturning.Borrow.User.Address, toolReturning.AdditionalInfo)
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

func (ms *MessageService) respondToolReturningApprove(toolReturning types.ToolReturning) error {
	borrow, err := ms.borrowService.FindBorrowByID(toolReturning.BorrowID)
	if err != nil {
		log.Println("[ERR][respondToolReturningApprove][FindBorrowByID]", err)
		return ms.Error()
	}

	if err := ms.borrowService.UpdateBorrowStatus(borrow.ID, types.GetBorrowStatus("returned")); err != nil {
		log.Println("[ERR][respondToolReturningApprove][UpdateBorrowStatus]", err)
		return ms.Error()
	}

	if err := ms.toolReturningService.UpdateToolReturningStatus(toolReturning.ID, types.GetToolReturningStatus("complete")); err != nil {
		log.Println("[ERR][respondToolReturningToAdmin][UpdateToolReturningStatus]", err)
		return ms.Error()
	}

	if err := ms.toolService.IncreaseStock(toolReturning.Borrow.ToolID, borrow.Amount); err != nil {
		log.Println("[ERR][respondToolReturningApprove][IncreaseStock]", err)
		return ms.Error()
	}

	reqBody := types.MessageRequest{
		ChatID: toolReturning.Borrow.UserID,
		Text:   fmt.Sprintf("Pengajuan pengembalian \"%s\" telah disetujui oleh pengurus.\n\nKeterangan:\n%s", toolReturning.Borrow.Tool.Name, ms.messageText),
	}
	if err := ms.sendMessage(reqBody); err != nil {
		log.Println("error in sending reply:", err)
		return err
	}

	return ms.sendMessage(types.MessageRequest{
		Text: "Pengajuan pengembalian berhasil disetujui.",
	})
}

func (ms *MessageService) respondToolReturningReject(toolReturning types.ToolReturning) error {
	if err := ms.toolReturningService.UpdateToolReturningStatus(toolReturning.ID, types.GetToolReturningStatus("reject")); err != nil {
		log.Println("[ERR][respondToolReturningReject][UpdateToolReturningStatus]", err)
		return ms.Error()
	}

	reqBody := types.MessageRequest{
		ChatID: toolReturning.Borrow.UserID,
		Text:   fmt.Sprintf("Pengajuan pengembalian \"%s\" telah ditolak oleh pengurus.\n\nKeterangan:\n%s", toolReturning.Borrow.Tool.Name, ms.messageText),
	}
	if err := ms.sendMessage(reqBody); err != nil {
		log.Println("error in sending reply:", err)
		return err
	}

	return ms.sendMessage(types.MessageRequest{
		Text: "Pengajuan pengembalian berhasil ditolak",
	})
}

func (ms *MessageService) Manage() error {
	if !ms.isEligibleAdmin() {
		log.Println("[INFO] Not eligible user accessing admin command", ms.messageText)
		return ms.Unknown()
	}

	manageCommands, ok := helper.GetManageCommandOrder(ms.messageText)
	if !ok {
		return ms.manageMenu()
	}

	if manageCommands.ID == 0 {
		if manageCommands.Type == types.ManageTypeEdit {
			return ms.sendMessage(types.MessageRequest{
				Text: fmt.Sprintf(
					"Untuk melakukan pengubahan data silahkan kirim perintah\n\"/%s %s [id_barang]\"\n\nContoh: \"/%s %s 5\"",
					types.CommandManage, types.ManageTypeEdit, types.CommandManage, types.ManageTypeEdit),
			})
		}

		if manageCommands.Type == types.ManageTypeDelete {
			return ms.sendMessage(types.MessageRequest{
				Text: fmt.Sprintf(
					"Untuk melakukan penghapusan barang silahkan kirim perintah\n\"/%s %s [id_barang]\"\n\nContoh: \"/%s %s 5\"",
					types.CommandManage, types.ManageTypeDelete, types.CommandManage, types.ManageTypeDelete),
			})
		}
	}

	switch manageCommands.Type {
	case types.ManageTypeAdd:
		return ms.manageAddInit()
	case types.ManageTypeEdit:
		return ms.manageEditInit(manageCommands.ID)
	case types.ManageTypeDelete:
		return ms.manageDeleteInit(manageCommands.ID)
	case types.ManageTypePhoto:
		return ms.managePhotoInit(manageCommands.ID)
	}

	return ms.Unknown()
}

func (ms *MessageService) manageMenu() error {
	return ms.sendMessage(types.MessageRequest{
		Text: "Silahkan pilih menu pengelolaan.",
		ReplyMarkup: types.InlineKeyboardMarkup{
			InlineKeyboard: [][]types.InlineKeyboardButton{
				{{
					Text:         "Tambah Barang",
					CallbackData: fmt.Sprintf("/%s %s", types.CommandManage, types.ManageTypeAdd),
				}},
				{{
					Text:         "Edit Barang",
					CallbackData: fmt.Sprintf("/%s %s", types.CommandManage, types.ManageTypeEdit),
				}},
			},
		},
	})
}

func (ms *MessageService) manageAddInit() error {
	if err := ms.saveChatSessionDetail(types.Topic["manage_add_init"], ""); err != nil {
		log.Println("[ERR][manageAddInit][saveChatSessionDetail]", err)
		return ms.Error()
	}

	return ms.sendMessage(types.MessageRequest{
		Text: "Memulai Sesi Penambahan Barang\n\nTulis nama barang",
	})
}

func (ms *MessageService) ManageAdd() error {
	if len(ms.chatSessionDetails) > 0 {
		switch ms.chatSessionDetails[0].Topic {
		case types.Topic["manage_add_init"]:
			return ms.manageAddName()
		case types.Topic["manage_add_name"]:
			return ms.manageAddBrand()
		case types.Topic["manage_add_brand"]:
			return ms.manageAddType()
		case types.Topic["manage_add_type"]:
			return ms.manageAddWeight()
		case types.Topic["manage_add_weight"]:
			return ms.manageAddStock()
		case types.Topic["manage_add_stock"]:
			return ms.manageAddInfo()
		case types.Topic["manage_add_info"]:
			return ms.manageAddPhoto()
		case types.Topic["manage_add_photo"]:
			return ms.manageAddConfirm()
		}
	}

	return ms.Unknown()
}

func (ms *MessageService) manageAddName() error {
	gen := helper.NewSessionDataGenerator()
	generatedSessionData := gen.ManageAddName(ms.messageText)

	if err := ms.saveChatSessionDetail(types.Topic["manage_add_name"], generatedSessionData); err != nil {
		log.Println("[ERR][manageAddName][saveChatSessionDetail]", err)
		return ms.Error()
	}

	return ms.sendMessage(types.MessageRequest{
		Text: "Tuliskan merk/brand",
	})
}

func (ms *MessageService) manageAddBrand() error {
	gen := helper.NewSessionDataGenerator()
	generatedSessionData := gen.ManageAddBrand(ms.messageText)

	if err := ms.saveChatSessionDetail(types.Topic["manage_add_brand"], generatedSessionData); err != nil {
		log.Println("[ERR][manageAddBrand][saveChatSessionDetail]", err)
		return ms.Error()
	}

	return ms.sendMessage(types.MessageRequest{
		Text: "Tuliskan tipe barang/alat",
	})
}

func (ms *MessageService) manageAddType() error {
	gen := helper.NewSessionDataGenerator()
	generatedSessionData := gen.ManageAddProductType(ms.messageText)

	if err := ms.saveChatSessionDetail(types.Topic["manage_add_type"], generatedSessionData); err != nil {
		log.Println("[ERR][manageAddType][saveChatSessionDetail]", err)
		return ms.Error()
	}

	return ms.sendMessage(types.MessageRequest{
		Text: "Berapa berat alat tersebut? (dalam gram)",
	})
}

func (ms *MessageService) manageAddWeight() error {
	i, err := strconv.ParseFloat(ms.messageText, 10)
	if err != nil || i < 0 {
		log.Println("[ERR][manageAddWeight][ParseFloat]", err)
		return ms.sendMessage(types.MessageRequest{
			Text: "Mohon sebutkan berat dalam angka.",
		})
	}

	weight := float32(i)

	gen := helper.NewSessionDataGenerator()
	generatedSessionData := gen.ManageAddWeight(weight)

	if err := ms.saveChatSessionDetail(types.Topic["manage_add_weight"], generatedSessionData); err != nil {
		log.Println("[ERR][manageAddWeight][saveChatSessionDetail]", err)
		return ms.Error()
	}

	return ms.sendMessage(types.MessageRequest{
		Text: "Berapa banyak stok yang tersedia untuk dipinjamkan?",
	})
}

func (ms *MessageService) manageAddStock() error {
	i, err := strconv.ParseInt(ms.messageText, 10, 64)
	if err != nil || i < 0 {
		log.Println("[ERR][manageAddStock][ParseInt]", err)
		return ms.sendMessage(types.MessageRequest{
			Text: "Mohon sebutkan jumlah stok dalam angka.",
		})
	}

	gen := helper.NewSessionDataGenerator()
	generatedSessionData := gen.ManageAddStock(i)

	if err := ms.saveChatSessionDetail(types.Topic["manage_add_stock"], generatedSessionData); err != nil {
		log.Println("[ERR][manageAddStock][saveChatSessionDetail]", err)
		return ms.Error()
	}

	return ms.sendMessage(types.MessageRequest{
		Text: "Tuliskan deskripsi lengkap mengenai alat ini",
	})
}

func (ms *MessageService) manageAddInfo() error {
	gen := helper.NewSessionDataGenerator()
	generatedSessionData := gen.ManageAddInfo(ms.messageText)

	if err := ms.saveChatSessionDetail(types.Topic["manage_add_info"], generatedSessionData); err != nil {
		log.Println("[ERR][manageAddInfo][saveChatSessionDetail]", err)
		return ms.Error()
	}

	return ms.sendMessage(types.MessageRequest{
		Text: "Silahkan upload foto barang. (minimal 1, maksimal 10)",
	})
}

func (ms *MessageService) manageAddPhoto() error {
	if len(ms.message.Photo) == 0 {
		return ms.sendMessage(types.MessageRequest{
			Text: "Mohon upload foto barang.",
		})
	}

	pickedPhoto := helper.PickPhoto(ms.message.Photo)

	gen := helper.NewSessionDataGenerator()
	generatedSessionData := gen.ManageAddPhoto(ms.message.MediaGroupID, pickedPhoto.FileID, pickedPhoto.FileUniqueID)

	if err := ms.saveChatSessionDetail(types.Topic["manage_add_photo"], generatedSessionData); err != nil {
		log.Println("[ERR][manageAddPhoto][saveChatSessionDetail]", err)
		return ms.Error()
	}

	tool := helper.GetToolFromChatSessionDetail(types.ManageTypeAdd, ms.chatSessionDetails)

	message := fmt.Sprintf(`Nama : %s
		Brand/Merk : %s
		Tipe Produk : %s
		Berat : %.2f gram
		Stok : %d
		Deskripsi :
		%s
	
		Pastikan data sudah benar kemudian tekan "Lanjutkan".`, tool.Name, tool.Brand, tool.ProductType, tool.Weight, tool.Stock, tool.AdditionalInformation)
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

func (ms *MessageService) manageAddConfirm() error {
	if len(ms.message.MediaGroupID) != 0 {
		pickedPhoto := helper.PickPhoto(ms.message.Photo)

		gen := helper.NewSessionDataGenerator()
		generatedSessionData := gen.ManageAddPhoto(ms.message.MediaGroupID, pickedPhoto.FileID, pickedPhoto.FileUniqueID)

		return ms.saveChatSessionDetail(types.Topic["manage_add_photo"], generatedSessionData)
	}

	var userResponse bool
	if ms.messageText == "yes" {
		userResponse = true
	} else {
		userResponse = false
	}

	gen := helper.NewSessionDataGenerator()
	generatedSessionData := gen.ManageAddConfirm(userResponse)

	if err := ms.saveChatSessionDetail(types.Topic["manage_add_confirm"], generatedSessionData); err != nil {
		log.Println("[ERR][manageAddConfirm][saveChatSessionDetail]", err)
		return ms.Error()
	}

	chatSessionID := ms.chatSessionDetails[0].ChatSessionID
	if err := ms.chatSessionService.UpdateChatSessionStatus(chatSessionID, types.ChatSessionStatus["complete"]); err != nil {
		log.Println("[ERR][manageAddConfirm][UpdateChatSessionStatus]", err)
		return ms.Error()
	}

	if !userResponse {
		return ms.sendMessage(types.MessageRequest{
			Text: "Penambahan barang dibatalkan.",
		})
	}

	tool := helper.GetToolFromChatSessionDetail(types.ManageTypeAdd, ms.chatSessionDetails)
	photos := helper.GetToolPhotosFromChatSessionDetails(ms.chatSessionDetails)

	toolID, err := ms.toolService.SaveTool(tool)
	if err != nil {
		log.Println("[ERR][manageAddConfirm][SaveTool]", err)
		return ms.Error()
	}

	if err = ms.toolService.SaveToolPhotos(toolID, photos); err != nil {
		log.Println("[ERR][manageAddConfirm][SaveToolPhotos]", err)
		return ms.Error()
	}

	return ms.sendMessage(types.MessageRequest{
		Text: fmt.Sprintf("Barang berhasil ditambah dengan ID %d", toolID),
		ReplyMarkup: types.InlineKeyboardMarkup{
			InlineKeyboard: [][]types.InlineKeyboardButton{
				{
					{
						Text:         "Cek Barang",
						CallbackData: fmt.Sprintf("/%s %d", types.CommandCheck, toolID),
					},
				},
			},
		},
	})
}

func (ms *MessageService) manageEditInit(toolID int64) error {
	_, err := ms.toolService.FindByID(toolID)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ERR][manageEditInit][FindByID]", err)
		return ms.Error()
	}

	if err == sql.ErrNoRows {
		return ms.sendMessage(types.MessageRequest{
			Text: "ID barang tidak ditemukan.",
		})
	}

	gen := helper.NewSessionDataGenerator()
	generatedSessionData := gen.ManageEditInit(toolID)

	if err := ms.saveChatSessionDetail(types.Topic["manage_edit_init"], generatedSessionData); err != nil {
		log.Println("[ERR][manageEditInit][saveChatSessionDetail]", err)
		return ms.Error()
	}

	reqBody := types.MessageRequest{
		Text: "Memulai Sesi Pengubahan Barang\n\nSilahkan pilih kolom data yang ingin diubah",
		ReplyMarkup: types.InlineKeyboardMarkup{
			InlineKeyboard: [][]types.InlineKeyboardButton{
				{
					{
						Text:         "Nama",
						CallbackData: "nama",
					},
					{
						Text:         "Brand",
						CallbackData: "brand",
					},
					{
						Text:         "Tipe",
						CallbackData: "tipe",
					},
				},
				{
					{
						Text:         "Berat",
						CallbackData: "berat",
					},
					{
						Text:         "Stok",
						CallbackData: "stok",
					},
					{
						Text:         "Foto",
						CallbackData: "foto",
					},
				},
				{
					{
						Text:         "Keterangan",
						CallbackData: "keterangan",
					},
				},
			},
		},
	}

	return ms.sendMessage(reqBody)
}

func (ms *MessageService) ManageEdit() error {
	if len(ms.chatSessionDetails) > 0 {
		switch ms.chatSessionDetails[0].Topic {
		case types.Topic["manage_edit_init"]:
			return ms.manageEditField()
		case types.Topic["manage_edit_field"]:
			return ms.manageEditComplete()
		}
	}

	return ms.Unknown()
}

func (ms *MessageService) manageEditField() error {
	if ok := helper.IsToolFieldExists(ms.messageText); !ok {
		return ms.sendMessage(types.MessageRequest{
			Text: "Kolom data tidak tersedia. Silahkan pilih kolom data yang akan diubah melalui pilihan menu.",
		})
	}

	sessionTool := helper.GetToolFromChatSessionDetail(types.ManageTypeEdit, ms.chatSessionDetails)
	tool, err := ms.toolService.FindByID(sessionTool.ID)
	if err != nil {
		log.Println("[ERR][manageEditField][FindByID]", err)
		return ms.Error()
	}

	if types.ToolField(ms.messageText) == types.ToolFieldPhoto {
		return ms.managePhotoInit(tool.ID)
	}

	gen := helper.NewSessionDataGenerator()
	generatedSessionData := gen.ManageEditField(ms.messageText)

	if err := ms.saveChatSessionDetail(types.Topic["manage_edit_field"], generatedSessionData); err != nil {
		log.Println("[ERR][manageEditField][saveChatSessionDetail]", err)
		return ms.Error()
	}

	oldValue := helper.GetToolValueByField(tool, ms.messageText)
	messsage := fmt.Sprintf("Data sebelumnya:\n%s\n\nSilahkan tulis data baru", oldValue)

	return ms.sendMessage(types.MessageRequest{
		Text: messsage,
	})
}

func (ms *MessageService) manageEditComplete() error {
	sessionTool := helper.GetToolFromChatSessionDetail(types.ManageTypeEdit, ms.chatSessionDetails)
	tool, err := ms.toolService.FindByID(sessionTool.ID)
	if err != nil {
		log.Println("[ERR][manageEditField][FindByID]", err)
		return ms.Error()
	}

	var choosenField string
	manageEditSession, found := helper.GetChatSessionDetailByTopic(ms.chatSessionDetails, types.Topic["manage_edit_field"])
	if found {
		dataParsed, err := gabs.ParseJSON([]byte(manageEditSession.Data))
		if err != nil {
			return err
		}
		value, ok := dataParsed.Path("field").Data().(string)
		if ok {
			choosenField = value
		}
	}

	updatedTool, err := helper.ChangeToolValueByField(tool, choosenField, ms.messageText)
	if err != nil {
		return ms.sendMessage(types.MessageRequest{
			Text: err.Error(),
		})
	}

	if err := ms.toolService.UpdateTool(updatedTool); err != nil {
		log.Println("[ERR][manageEditComplete][UpdateTool]", err)
		return ms.sendMessage(types.MessageRequest{
			Text: fmt.Sprintf("Terjadi kesalahan. Barang dengan ID %d gagal diubah.", tool.ID),
		})
	}

	chatSessionID := ms.chatSessionDetails[0].ChatSessionID
	if err := ms.chatSessionService.UpdateChatSessionStatus(chatSessionID, types.ChatSessionStatus["complete"]); err != nil {
		log.Println("[ERR][manageEditComplete][UpdateChatSessionStatus]", err)
		return ms.Error()
	}

	return ms.sendMessage(types.MessageRequest{
		Text: fmt.Sprintf("Barang dengan ID %d berhasil diubah", tool.ID),
		ReplyMarkup: types.InlineKeyboardMarkup{
			InlineKeyboard: [][]types.InlineKeyboardButton{
				{
					{
						Text:         "Cek Barang",
						CallbackData: fmt.Sprintf("/%s %d", types.CommandCheck, tool.ID),
					},
				},
			},
		},
	})
}

func (ms *MessageService) manageDeleteInit(toolID int64) error {
	tool, err := ms.toolService.FindByID(toolID)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ERR][managedDeleteInit][FindByID]", err)
		return ms.Error()
	}

	if err == sql.ErrNoRows {
		return ms.sendMessage(types.MessageRequest{
			Text: "ID barang tidak ditemukan.",
		})
	}

	gen := helper.NewSessionDataGenerator()
	generatedSessionData := gen.ManageDeleteInit(toolID)

	if err := ms.saveChatSessionDetail(types.Topic["manage_delete_init"], generatedSessionData); err != nil {
		log.Println("[ERR][manageDeleteInit][saveChatSessionDetail]", err)
		return ms.Error()
	}

	reqBody := types.MessageRequest{
		Text: fmt.Sprintf("Apakah Anda yakin ingin menghapus %s?", tool.Name),
		ReplyMarkup: types.InlineKeyboardMarkup{
			InlineKeyboard: [][]types.InlineKeyboardButton{
				{
					{
						Text:         "Yakin",
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

func (ms *MessageService) ManageDelete() error {
	if len(ms.chatSessionDetails) > 0 {
		switch ms.chatSessionDetails[0].Topic {
		case types.Topic["manage_delete_init"]:
			return ms.manageDeleteComplete()
		}
	}

	return ms.Unknown()
}

func (ms *MessageService) manageDeleteComplete() error {
	var userResponse bool
	if ms.messageText == "yes" {
		userResponse = true
	} else {
		userResponse = false
	}

	gen := helper.NewSessionDataGenerator()
	generatedSessionData := gen.ManageDeleteComplete(userResponse)

	if err := ms.saveChatSessionDetail(types.Topic["manage_delete_complete"], generatedSessionData); err != nil {
		log.Println("[ERR][manageDeleteComplete][saveChatSessionDetail]", err)
		return ms.Error()
	}

	chatSessionID := ms.chatSessionDetails[0].ChatSessionID
	if err := ms.chatSessionService.UpdateChatSessionStatus(chatSessionID, types.ChatSessionStatus["complete"]); err != nil {
		log.Println("[ERR][manageDeleteComplete][UpdateChatSessionStatus]", err)
		return ms.Error()
	}

	if !userResponse {
		return ms.sendMessage(types.MessageRequest{
			Text: "Penghapusan barang dibatalkan.",
		})
	}

	sessionTool := helper.GetToolFromChatSessionDetail(types.ManageTypeDelete, ms.chatSessionDetails)

	if err := ms.toolService.DeleteTool(sessionTool.ID); err != nil {
		log.Println("[ERR][manageDeleteComplete][DeleteTool]", err)
		return ms.Error()
	}

	return ms.sendMessage(types.MessageRequest{
		Text: fmt.Sprintf("Barang dengan ID %d berhasil dihapus.", sessionTool.ID),
	})
}

func (ms *MessageService) managePhotoInit(toolID int64) error {
	_, err := ms.toolService.FindByID(toolID)
	if err != nil && err != sql.ErrNoRows {
		log.Println("[ERR][managePhotoInit][FindByID]", err)
		return ms.Error()
	}

	if err == sql.ErrNoRows {
		return ms.sendMessage(types.MessageRequest{
			Text: "ID barang tidak ditemukan.",
		})
	}

	gen := helper.NewSessionDataGenerator()
	generatedSessionData := gen.ManagePhotoInit(toolID)

	if err := ms.saveChatSessionDetail(types.Topic["manage_photo_init"], generatedSessionData); err != nil {
		log.Println("[ERR][managePhotoInit][saveChatSessionDetail]", err)
		return ms.Error()
	}

	return ms.sendMessage(types.MessageRequest{
		Text: "Silahkan upload foto barang. (minimal 1, maksimal 10)",
	})
}

func (ms *MessageService) ManagePhoto() error {
	if len(ms.chatSessionDetails) > 0 {
		switch ms.chatSessionDetails[0].Topic {
		case types.Topic["manage_photo_init"]:
			return ms.managePhotoUpload()
		case types.Topic["manage_photo_upload"]:
			return ms.managePhotoConfirm()
		}
	}

	return ms.Unknown()
}

func (ms *MessageService) managePhotoUpload() error {
	if len(ms.message.Photo) == 0 {
		return ms.sendMessage(types.MessageRequest{
			Text: "Mohon upload foto barang.",
		})
	}

	pickedPhoto := helper.PickPhoto(ms.message.Photo)

	gen := helper.NewSessionDataGenerator()
	generatedSessionData := gen.ManagePhotoUpload(ms.message.MediaGroupID, pickedPhoto.FileID, pickedPhoto.FileUniqueID)

	if err := ms.saveChatSessionDetail(types.Topic["manage_photo_upload"], generatedSessionData); err != nil {
		log.Println("[ERR][managePhotoUpload][saveChatSessionDetail]", err)
		return ms.Error()
	}

	reqBody := types.MessageRequest{
		Text: "Foto berhasil diunggah.",
		ReplyMarkup: types.InlineKeyboardMarkup{
			InlineKeyboard: [][]types.InlineKeyboardButton{
				{{
					Text:         "Simpan Perubahan",
					CallbackData: "yes",
				}},
				{{
					Text:         "Batalkan",
					CallbackData: "no",
				}},
			},
		},
	}

	return ms.sendMessage(reqBody)
}

func (ms *MessageService) managePhotoConfirm() error {
	if len(ms.message.MediaGroupID) != 0 {
		pickedPhoto := helper.PickPhoto(ms.message.Photo)

		gen := helper.NewSessionDataGenerator()
		generatedSessionData := gen.ManagePhotoUpload(ms.message.MediaGroupID, pickedPhoto.FileID, pickedPhoto.FileUniqueID)

		return ms.saveChatSessionDetail(types.Topic["manage_photo_upload"], generatedSessionData)
	}

	var userResponse bool
	if ms.messageText == "yes" {
		userResponse = true
	} else {
		userResponse = false
	}

	gen := helper.NewSessionDataGenerator()
	generatedSessionData := gen.ManagePhotoConfirm(userResponse)

	if err := ms.saveChatSessionDetail(types.Topic["manage_photo_confirm"], generatedSessionData); err != nil {
		log.Println("[ERR][managePhotoConfirm][saveChatSessionDetail]", err)
		return ms.Error()
	}

	chatSessionID := ms.chatSessionDetails[0].ChatSessionID
	if err := ms.chatSessionService.UpdateChatSessionStatus(chatSessionID, types.ChatSessionStatus["complete"]); err != nil {
		log.Println("[ERR][managePhotoConfirm][UpdateChatSessionStatus]", err)
		return ms.Error()
	}

	if !userResponse {
		return ms.sendMessage(types.MessageRequest{
			Text: "Perubahan foto barang dibatalkan.",
		})
	}

	tool := helper.GetToolFromChatSessionDetail(types.ManageTypePhoto, ms.chatSessionDetails)
	photos := helper.GetToolPhotosFromChatSessionDetails(ms.chatSessionDetails)

	if err := ms.toolService.UpdatePhotos(tool.ID, photos); err != nil {
		log.Println("[ERR][managePhotoConfirm][UpdatePhotos]", err)
		return ms.Error()
	}

	return ms.sendMessage(types.MessageRequest{
		Text: "Foto barang berhasil diubah.",
		ReplyMarkup: types.InlineKeyboardMarkup{
			InlineKeyboard: [][]types.InlineKeyboardButton{
				{
					{
						Text:         "Lihat Foto",
						CallbackData: fmt.Sprintf("/%s %d %s", types.CommandCheck, tool.ID, types.CheckTypePhoto),
					},
				},
			},
		},
	})
}

func (ms *MessageService) Report() error {
	if !ms.isEligibleAdmin() {
		log.Println("[INFO] Not eligible user accessing admin command", ms.messageText)
		return ms.Unknown()
	}

	reportCommands, ok := helper.GetReportCommandOrder(ms.messageText)
	if !ok {
		return ms.reportMenu()
	}

	if reportCommands.Type == types.ReportTypeBorrow {
		return ms.reportBorrow(reportCommands)
	} else if reportCommands.Type == types.ReportTypeToolReturning {
		return ms.reportToolReturning(reportCommands)
	}

	return ms.Unknown()
}

func (ms *MessageService) reportMenu() error {
	return ms.sendMessage(types.MessageRequest{
		Text: "Silahkan pilih menu laporan.",
		ReplyMarkup: types.InlineKeyboardMarkup{
			InlineKeyboard: [][]types.InlineKeyboardButton{
				{{
					Text:         "Peminjaman",
					CallbackData: fmt.Sprintf("/%s %s", types.CommandReport, types.ReportTypeBorrow),
				}},
				{{
					Text:         "Pengembalian",
					CallbackData: fmt.Sprintf("/%s %s", types.CommandReport, types.ReportTypeToolReturning),
				}},
			},
		},
	})
}

func (ms *MessageService) reportBorrow(commands types.ReportCommandOrder) error {
	if len(commands.Text) == 0 {
		message := fmt.Sprintf(`
			Laporan Peminjaman Bulanan dapat dilihat dengan perintah
			"/%s %s [tahun]-[bulan]"
		
			Contoh, laporan peminjaman pada bulan Agustus tahun 2021
			"/%s %s 2021-8"		
		`, types.CommandReport, types.ReportTypeBorrow, types.CommandReport, types.ReportTypeBorrow)
		message = helper.RemoveTab(message)

		currentTime := time.Now()
		currentYear := currentTime.Year()
		currentMonth := int(currentTime.Month())

		return ms.sendMessage(types.MessageRequest{
			Text: message,
			ReplyMarkup: types.InlineKeyboardMarkup{
				InlineKeyboard: [][]types.InlineKeyboardButton{
					{{
						Text:         "Laporan Bulan Ini",
						CallbackData: fmt.Sprintf("/%s %s %d-%d", types.CommandReport, types.ReportTypeBorrow, currentYear, currentMonth),
					}},
					{{
						Text:         "Laporan Bulan Kemarin",
						CallbackData: fmt.Sprintf("/%s %s %d-%d", types.CommandReport, types.ReportTypeBorrow, currentYear, currentMonth-1),
					}},
				},
			},
		})
	}

	year, month, ok := helper.GetReportTimeFromCommand(commands.Text)
	if !ok {
		return ms.sendMessage(types.MessageRequest{
			Text: helper.RemoveTab(fmt.Sprintf(`
				Mohon isi tahun dan bulan dengan format dan nilai yang sesuai.
				Contoh: "/%s %s 2021-8"`,
				types.CommandReport, types.ReportTypeBorrow)),
		})
	}

	borrows, err := ms.borrowService.GetBorrowReport(year, month)
	if err != nil {
		log.Println("[ERR][reportBorrow][GetBorrowReport]", err)
		return ms.Error()
	}

	message := "Tidak ada data peminjaman pada waktu yang dimaksud."
	if len(borrows) > 0 {
		message = fmt.Sprintf("Laporan Peminjaman Bulan %s Tahun %d\n\n", helper.MonthNameSwitcher(month), year)
		message += helper.BuildBorrowReportMessage(borrows)
	}

	return ms.sendMessage(types.MessageRequest{
		Text: message,
	})
}

func (ms *MessageService) reportToolReturning(commands types.ReportCommandOrder) error {
	if len(commands.Text) == 0 {
		message := fmt.Sprintf(`
			Laporan Pengembalian Bulanan dapat dilihat dengan perintah
			"/%s %s [tahun]-[bulan]"

			Contoh, laporan pengembalian pada bulan Agustus tahun 2021
			"/%s %s 2021-8"
		`, types.CommandReport, types.ReportTypeToolReturning, types.CommandReport, types.ReportTypeToolReturning)
		message = helper.RemoveTab(message)

		currentTime := time.Now()
		currentYear := currentTime.Year()
		currentMonth := int(currentTime.Month())

		return ms.sendMessage(types.MessageRequest{
			Text: message,
			ReplyMarkup: types.InlineKeyboardMarkup{
				InlineKeyboard: [][]types.InlineKeyboardButton{
					{{
						Text:         "Laporan Bulan Ini",
						CallbackData: fmt.Sprintf("/%s %s %d-%d", types.CommandReport, types.ReportTypeToolReturning, currentYear, currentMonth),
					}},
					{{
						Text:         "Laporan Bulan Kemarin",
						CallbackData: fmt.Sprintf("/%s %s %d-%d", types.CommandReport, types.ReportTypeToolReturning, currentYear, currentMonth-1),
					}},
				},
			},
		})
	}

	year, month, ok := helper.GetReportTimeFromCommand(commands.Text)
	if !ok {
		return ms.sendMessage(types.MessageRequest{
			Text: helper.RemoveTab(fmt.Sprintf(`
				Mohon isi tahun dan bulan dengan format dan nilai yang sesuai.
				Contoh: "/%s %s 2021-8"`,
				types.CommandReport, types.ReportTypeBorrow)),
		})
	}

	toolReturnings, err := ms.toolReturningService.GetToolReturningReport(year, month)
	if err != nil {
		log.Println("[ERR][reportToolReturning][GetToolReturningReport]", err)
		return ms.Error()
	}

	message := "Tidak ada data pengembalian pada waktu yang dimaksud."
	if len(toolReturnings) > 0 {
		message = fmt.Sprintf("Laporan Pengembalian Bulan %s Tahun %d\n\n", helper.MonthNameSwitcher(month), year)
		message += helper.BuildToolReturningReportMessage(toolReturnings)
	}

	return ms.sendMessage(types.MessageRequest{
		Text: message,
	})
}
