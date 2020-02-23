package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mikeder/globber/internal/auth"
	"github.com/mikeder/globber/internal/web"
)

func (a *authAPI) Login(w http.ResponseWriter, r *http.Request) {
	creds := new(auth.Credentials)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&creds)

	token, err := a.manager.PasswordLogin(r.Context(), creds)
	if err != nil {
		log.Println(err)
		resp := struct {
			Error  string `json:"error"`
			Reason string `json:"reason,omitempty"`
		}{
			Error: "Forbidden",
		}

		if _, ok := err.(auth.ErrUserMissingField); ok {
			resp.Reason = err.Error()
		}

		if err := web.Respond(w, resp, http.StatusForbidden); err != nil {
			log.Println(err)
		}
		return
	}

	resp := struct {
		Token string `json:"token"`
	}{
		Token: token,
	}

	if err := web.Respond(w, resp, http.StatusOK); err != nil {
		log.Println(err)
		return
	}
}

func (a *authAPI) Logout(w http.ResponseWriter, r *http.Request) {
	return
}
