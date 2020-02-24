package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/mikeder/globber/internal/auth"
	"github.com/mikeder/globber/internal/web"
)

func (a *adminAPI) AddUser(w http.ResponseWriter, r *http.Request) {
	user := new(auth.User)
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&user)

	if err := a.manager.AddUser(r.Context(), user); err != nil {
		resp := struct {
			Error  string `json:"error"`
			Reason string `json:"reason,omitempty"`
		}{
			Error:  "forbidden",
			Reason: err.Error(),
		}
		web.Respond(w, resp, http.StatusForbidden)
		return
	}
	resp := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}
	web.Respond(w, resp, http.StatusCreated)
	return
}
