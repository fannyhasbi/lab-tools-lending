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
		Register: "registrasi",
		Check:    "cek",
		Borrow:   "pinjam",
		Return:   "pengembalian",
		Help:     "help",
	}
}
