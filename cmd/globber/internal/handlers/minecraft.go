package handlers

import (
	"net/http"
	"time"

	"github.com/mikeder/globber/internal/web"
)

func (s *site) minecraftPing(w http.ResponseWriter, r *http.Request) {
	err := s.mc.Ping()
	if err != nil {
		web.Respond(w, struct{ error }{err}, http.StatusInternalServerError)
	}
	web.Respond(w, struct {
		Latency time.Duration `json:"latency"`
	}{s.mc.Latency}, http.StatusOK)
}

func (s *site) minecraftStatus(w http.ResponseWriter, r *http.Request) {
	err := s.mc.Status()
	if err != nil {
		web.Respond(w, struct{ error }{err}, http.StatusInternalServerError)
	}
	web.Respond(w, *s.mc, http.StatusOK)
}
