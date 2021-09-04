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

	RespondCommandOrder struct {
		Type RespondType
		ID   int64
		Text string
	}

	ManageCommandOrder struct {
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
