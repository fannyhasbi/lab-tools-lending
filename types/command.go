package types

const (
	CommandRegister = "registrasi"
	CommandCheck    = "cek"
	CommandBorrow   = "pinjam"
	CommandReturn   = "pengembalian"
	CommandHelp     = "bantuan"

	// admin stuffs
	CommandRespond = "tanggapi"
)

type (
	RespondType string

	RespondCommands struct {
		Type RespondType
		ID   int64
		Text string
	}
)

var (
	RespondTypeBorrow        RespondType = "pinjam"
	RespondTypeToolReturning RespondType = "kembali"
)
