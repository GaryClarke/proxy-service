package main

import (
	"fmt"
	"net/http"
)

func (app *application) webhookHandler(w http.ResponseWriter, r *http.Request) {
	app.logger.Info("webhook received", "method", r.Method, "url", r.URL.Path)

	// Example: Log headers (optional, depending on your use case)
	for name, values := range r.Header {
		for _, value := range values {
			app.logger.Info("header", "name", name, "value", value)
		}
	}

	fmt.Fprintln(w, "webhook received")
}
