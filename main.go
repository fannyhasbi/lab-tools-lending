package main

import (
	"fmt"
	"net/http"

	"github.com/fannyhasbi/lab-tools-lending/config"
	"github.com/fannyhasbi/lab-tools-lending/handler"
)

func main() {
	fmt.Printf("Server running on port %s\n", config.GetPort())
	http.ListenAndServe(":"+config.GetPort(), http.HandlerFunc(handler.WebhookHandler))
}
