package main

import (
	"net/http"

	"github.com/Andhika-GIT/go-message-broker-monorepo/app"
)

func main() {
	r := app.InitApp()

	http.ListenAndServe(":3005", r)
}
