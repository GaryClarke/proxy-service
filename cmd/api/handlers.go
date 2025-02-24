package main

import (
	"fmt"
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":      "available",
		"environment": app.config.env,
		"version":     version,
	}

	err := app.writeJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}
}

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
