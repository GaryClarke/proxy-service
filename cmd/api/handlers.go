package main

import (
	"context"
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
	app.logHeaders(r)

	// Decode the JSON request body into a Webhook struct.
	var wh webhook.Webhook
	err := app.readJSON(w, r, &wh)
	if err != nil {
		app.errorResponse(w, r, http.StatusUnprocessableEntity, err.Error())
		return
	}

	err = app.handlerDelegator.Delegate(context.Background(), wh)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	err = app.writeJSON(w, http.StatusNoContent, nil, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
