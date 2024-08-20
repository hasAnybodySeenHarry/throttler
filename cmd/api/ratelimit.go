package main

import (
	"net/http"
	"strings"

	"harry2an.com/throttler/cmd/proto/users"
)

const (
	AuthHeader = "Authorization"
	Bearer     = "Bearer"
)

func (app *application) ratelimitHandler(w http.ResponseWriter, r *http.Request) {
	var key string
	var err error
	var user *users.GetUserResponse

	if auth := r.Header.Get(AuthHeader); auth != "" {
		segs := strings.Split(auth, " ")
		if len(segs) == 2 && segs[0] == Bearer {
			user, err = app.getUserForToken(segs[1])
			if err != nil {
				app.serverError(w, r, err)
				return
			}

			key = segs[1]
		} else {
			app.invalidAuthToken(w, r)
			return
		}
	} else {
		key, err = getClientIP(r)
		if err != nil {
			app.badRequest(w, r, err)
			return
		}
	}

	activated := user != nil && user.Activated

	ok, err := app.models.Buckets.Allow(key, activated)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if !ok {
		app.tooManyRequests(w, r, getRateLimitMessage(activated))
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"status": true}, nil)
	if err != nil {
		app.serverError(w, r, err)
	}
}
