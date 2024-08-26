package main

import (
	"net/http"
	"strconv"
	"strings"

	"harry2an.com/throttler/cmd/proto/users"
)

const (
	Authorization = "Authorization"
	Bearer        = "Bearer"
)

func (app *application) ratelimitHandler(w http.ResponseWriter, r *http.Request) {
	var key string
	var err error
	var user *users.GetUserResponse
	var activated bool

	if auth := r.Header.Get(Authorization); auth != "" {
		segs := strings.Split(auth, " ")
		if len(segs) == 2 && segs[0] == Bearer {
			token := segs[1]

			if id, err := app.models.Users.GetUserIDByToken(token); err == nil && id != "" {
				key = id
			} else {
				result, err := app.cb.Execute(func() (interface{}, error) {
					return app.getUserForToken(token)
				})
				if err != nil {
					if isCircuitBreakerError(err) {
						// fallback in case of circuit breaker is opened
						key, err = getClientIP(r)
						if err != nil {
							app.badRequest(w, r, err)
							return
						}
					} else {
						app.handleGRPCError(w, r, err)
						return
					}
				} else {
					user = result.(*users.GetUserResponse)

					userID := strconv.FormatInt(user.Id, 10)

					err = app.models.Users.InsertTokenWithID(token, userID)
					if err != nil {
						app.serverError(w, r, err)
						return
					}

					key = userID
				}
			}

			if user != nil {
				activated = user.Activated
			}

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

	if ok, err := app.models.Buckets.Allow(key, activated); err != nil {
		app.serverError(w, r, err)
		return
	} else if !ok {
		app.tooManyRequests(w, r, getRateLimitMessage(user == nil || activated))
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"status": true}, nil)
	if err != nil {
		app.serverError(w, r, err)
	}
}
