package handlers

import (
	"log"
	"net/http"

	"github.com/mikeder/globber/internal/web"
)

// Response is the structure of a minecraft status response.
type Response struct {
	Status string      `json:"status"`
	Reason string      `json:"reason,omitempty"`
	Result interface{} `json:"result"`
}

func (s *site) minecraftPing(w http.ResponseWriter, r *http.Request) {
	resp := Response{
		Status: "ok",
	}
	err := s.mc.PingList()
	if err != nil {
		log.Printf("error: %v\n", err)
		resp.Status = "error"
		resp.Reason = err.Error()
	}
	resp.Result = s.mc.Latency
	web.Respond(w, resp, http.StatusOK)
}

func (s *site) minecraftStatus(w http.ResponseWriter, r *http.Request) {
	resp := Response{
		Status: "ok",
	}
	err := s.mc.PingList()
	if err != nil {
		log.Printf("error: %v\n", err)
		resp.Status = "error"
		resp.Reason = err.Error()
	}
	resp.Result = *s.mc
	web.Respond(w, resp, http.StatusOK)
}
