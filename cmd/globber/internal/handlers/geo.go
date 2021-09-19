package handlers

import (
	"log"
	"net/http"

	"github.com/mikeder/globber/internal/web"
)

func (s *site) geoLookup(w http.ResponseWriter, r *http.Request) {
	id := s.geo.Lookup(r.RemoteAddr)
	log.Printf("from: %s, in: %s\n", id.ParsedIP, id.CountryAlpha2)

	resp := struct {
		IP      string `json:"ip"`
		Country string `json:"country"`
	}{
		IP:      id.ParsedIP.String(),
		Country: id.CountryAlpha2,
	}

	web.Respond(w, resp, http.StatusOK)
}
