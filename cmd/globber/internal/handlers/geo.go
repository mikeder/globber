package handlers

import (
	"encoding/json"
	"log"
	"net/http"
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

	b, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
		return
	}
	_, err = w.Write(b)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
