package main

import (
	"net/http"
)

func (app *application) ratelimitHandler(w http.ResponseWriter, r *http.Request) {
	ip, err := getClientIP(r)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	ok, err := app.models.Buckets.Allow(ip)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if !ok {
		app.tooManyRequests(w, r)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"status": true}, nil)
	if err != nil {
		app.serverError(w, r, err)
	}
}
