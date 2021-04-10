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
		Register: "register",
		Check:    "check",
		Borrow:   "borrow",
		Return:   "return",
		Help:     "help",
	}
}
