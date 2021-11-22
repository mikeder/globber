package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
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

	ac, rc := newCookies(tokens)

	http.SetCookie(w, &ac)
	http.SetCookie(w, &rc)

	if err := web.Respond(w, resp, http.StatusOK); err != nil {
		log.Println(err)
		return
	}
}

func (a *authAPI) Logout(w http.ResponseWriter, r *http.Request) {
	ac, rc := newCookies(nil)
	http.SetCookie(w, &ac)
	http.SetCookie(w, &rc)
	if err := web.Respond(w, nil, http.StatusOK); err != nil {
		log.Println(err)
		return
	}
}

func (a *authAPI) Refresh(w http.ResponseWriter, r *http.Request) {
	refresh, err := r.Cookie("jwt_refresh")
	if err != nil {
		log.Print(err)
		return
	}

	tokens, err := a.manager.Refresh(r.Context(),
		&auth.Tokens{
			Refresh: &jwt.Token{
				Raw: refresh.Value,
			},
		})

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

	ac, rc := newCookies(tokens)

	http.SetCookie(w, &ac)
	http.SetCookie(w, &rc)

	if err := web.Respond(w, resp, http.StatusOK); err != nil {
		log.Println(err)
		return
	}
}

func (a *authAPI) Tokens(w http.ResponseWriter, r *http.Request) {
	if err := web.Respond(w, a.manager.ListTokens(r.Context()), http.StatusOK); err != nil {
		log.Println(err)
		return
	}
}

func newCookies(t *auth.Tokens) (access, refresh http.Cookie) {
	if t == nil {
		return http.Cookie{
				Name:    "jwt",
				Path:    "/",
				Expires: time.Time{},
			},
			http.Cookie{
				Name:    "jwt_refresh",
				Path:    "/",
				Expires: time.Time{},
			}
	}
	access = http.Cookie{
		Name:     "jwt",
		Path:     "/",
		HttpOnly: true,
		Expires:  t.AccessTTL,
		Value:    t.Access.Raw,
		SameSite: http.SameSiteDefaultMode,
	}
	refresh = http.Cookie{
		Name:     "jwt_refresh",
		Path:     "/",
		HttpOnly: true,
		Expires:  t.RefreshTTL,
		Value:    t.Refresh.Raw,
		SameSite: http.SameSiteDefaultMode,
	}
	return
}
