package config

import (
	"os"
)

const port = "3000"

func GetPort() string {
	p, ok := os.LookupEnv("PORT")
	if !ok {
		p = port
	}

	return p
}
