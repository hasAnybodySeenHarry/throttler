package main

import (
	"fmt"
	"net/http"
)

func (app *application) log(r *http.Request, err error) {
	app.logger.Println(err, "at", r.URL.String())
}

func (app *application) error(w http.ResponseWriter, r *http.Request, status int, msg any) {
	data := envelope{
		"error": msg,
	}

	err := app.writeJSON(w, status, data, nil)
	if err != nil {
		app.log(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *application) methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	app.error(w, r, http.StatusMethodNotAllowed, fmt.Sprintf("The %s method is not allowed", r.Method))
}

func (app *application) notFound(w http.ResponseWriter, r *http.Request) {
	app.error(w, r, http.StatusNotFound, http.StatusText(http.StatusNotFound))
}

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	app.log(r, err)
	app.error(w, r, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
}

func (app *application) badRequest(w http.ResponseWriter, r *http.Request, err error) {
	app.error(w, r, http.StatusBadRequest, err.Error())
}

func (app *application) tooManyRequests(w http.ResponseWriter, r *http.Request, msg string) {
	app.error(w, r, http.StatusTooManyRequests, msg)
}

func (app *application) invalidAuthToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")
	app.error(w, r, http.StatusUnauthorized, "invalid or missing authentication header")
}

func (app *application) invalidCredentials(w http.ResponseWriter, r *http.Request) {
	app.error(w, r, http.StatusUnauthorized, "invalid credentials")
}

func (app *application) badGateway(w http.ResponseWriter, r *http.Request) {
	app.error(w, r, http.StatusBadGateway, http.StatusText(http.StatusBadGateway))
}

func (app *application) gatewayTimeOut(w http.ResponseWriter, r *http.Request) {
	app.error(w, r, http.StatusGatewayTimeout, http.StatusText(http.StatusGatewayTimeout))
}
