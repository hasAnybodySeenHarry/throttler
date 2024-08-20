package main

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"
)

type envelope map[string]interface{}

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	for k, v := range headers {
		w.Header()[k] = v
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(js)
	if err != nil {
		return err
	}

	return nil
}

func getClientIP(r *http.Request) (string, error) {
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		return strings.TrimSpace(strings.Split(forwarded, ",")[0]), nil
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}

	return host, nil
}

func getRateLimitMessage(activated bool) string {
	if !activated {
		return "Consider activating your account for more quotas."
	}
	return ""
}
