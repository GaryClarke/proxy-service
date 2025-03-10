package main

import "net/http"

func (app *application) routes() http.Handler {
	// create a router (mux)
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthcheck", app.healthcheckHandler)
	mux.HandleFunc("POST /webhook", app.webhookHandler)

	// Wrap the mux with the middleware
	return app.recoverPanic(mux)
}
