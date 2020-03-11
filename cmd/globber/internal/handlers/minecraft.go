package handlers

import (
	"net/http"
	"time"

	"github.com/mikeder/globber/internal/minecraft"
	"github.com/mikeder/globber/internal/web"
)

func (s *site) minecraftPing(w http.ResponseWriter, r *http.Request) {
	ping, err := s.mc.Ping()
	if err != nil {
		web.Respond(w, struct{ error }{err}, http.StatusInternalServerError)
	}
	web.Respond(w, struct{ time.Duration }{ping}, http.StatusOK)
}

func (s *site) minecraftStatus(w http.ResponseWriter, r *http.Request) {
	err := s.mc.Status()
	if err != nil {
		web.Respond(w, struct{ error }{err}, http.StatusInternalServerError)
	}
	web.Respond(w, struct{ minecraft.Server }{*s.mc}, http.StatusOK)
}
