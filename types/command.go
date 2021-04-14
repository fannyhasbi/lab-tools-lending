package types

type command struct {
	Register string
	Check    string
	Borrow   string
	Return   string
	Help     string
}

func Command() command {
	return command{
		Register: "daftar",
		Check:    "cek",
		Borrow:   "pinjam",
		Return:   "pengembalian",
		Help:     "help",
	}
}
