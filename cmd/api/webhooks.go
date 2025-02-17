package main

import (
	"fmt"
	"net/http"
)

func (app *application) webhookHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(w, "create a new movie")
}
