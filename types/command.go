package types

const (
	CommandRegister = "registrasi"
	CommandCheck    = "cek"
	CommandBorrow   = "pinjam"
	CommandReturn   = "pengembalian"
	CommandHelp     = "bantuan"

	// admin stuffs
	CommandRespond = "tanggapi"
	CommandManage  = "kelola"
)

type (
	RespondType string
	ManageType  string

	RespondCommands struct {
		Type RespondType
		ID   int64
		Text string
	}

	ManageCommands struct {
		Type ManageType
		ID   int64
	}
)

var (
	RespondTypeBorrow        RespondType = "pinjam"
	RespondTypeToolReturning RespondType = "kembali"

	ManageTypeAdd   ManageType = "tambah"
	ManageTypeEdit  ManageType = "edit"
	ManageTypePhoto ManageType = "foto"
)
