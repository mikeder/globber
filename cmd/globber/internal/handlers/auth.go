package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mikeder/globber/internal/auth"
	"github.com/mikeder/globber/internal/web"
)

func (a *authAPI) Login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	creds := &auth.Credentials{
		Email:    r.Form.Get("email"),
		Password: r.Form.Get("password"),
	}

	if creds.Email == "" || creds.Password == "" {
		decoder := json.NewDecoder(r.Body)
		decoder.Decode(&creds)
	}

	tokens, err := a.manager.PasswordLogin(r.Context(), creds)
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
		Access  string `json:"access_token"`
		Refresh string `json:"refresh_token"`
	}{
		Access:  tokens.Access.Raw,
		Refresh: tokens.Refresh.Raw,
	}

	c := http.Cookie{
		Name:     "jwt",
		Path:     "/",
		Value:    tokens.Access.Raw,
		SameSite: http.SameSiteDefaultMode,
	}

	http.SetCookie(w, &c)

	if err := web.Respond(w, resp, http.StatusOK); err != nil {
		log.Println(err)
		return
	}
}

func (a *authAPI) Logout(w http.ResponseWriter, r *http.Request) {
	return
}
