package main

import (
	"encoding/json"
	"github.com/garyclarke/proxy-service/internal/webhook"
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
		app.serverErrorResponse(w, r, err)
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

	// Decode the JSON request body into a Webhook struct.
	var wh webhook.Webhook

	err := json.NewDecoder(r.Body).Decode(&wh)
	if err != nil {
		app.errorResponse(w, r, http.StatusUnprocessableEntity, err.Error())
		return
	}

	err = app.writeJSON(w, http.StatusNoContent, nil, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
