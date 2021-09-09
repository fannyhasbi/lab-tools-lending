package types

const (
	CommandStart    = "start" // automatically sent when user add the bot for the first time
	CommandRegister = "registrasi"
	CommandCheck    = "cek"
	CommandBorrow   = "pinjam"
	CommandReturn   = "pengembalian"
	CommandHelp     = "bantuan"

	// admin stuffs
	CommandRespond = "tanggapi"
	CommandManage  = "kelola"
	CommandReport  = "laporan"
)

type (
	RespondType string
	ManageType  string
	ReportType  string

	CheckCommandOrder struct {
		ID   int64
		Text string
	}

	RespondCommandOrder struct {
		Type RespondType
		ID   int64
		Text string
	}

	ManageCommandOrder struct {
		Type ManageType
		ID   int64
	}

	ReportCommandOrder struct {
		Type ReportType
		Text string
	}
)

var (
	CheckTypePhoto string = "foto"

	RespondTypeBorrow        RespondType = "pinjam"
	RespondTypeToolReturning RespondType = "kembali"

	ManageTypeAdd   ManageType = "tambah"
	ManageTypeEdit  ManageType = "edit"
	ManageTypePhoto ManageType = "foto"

	ReportTypeBorrow        ReportType = "pinjam"
	ReportTypeToolReturning ReportType = "kembali"
)
