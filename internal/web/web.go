package web

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
)

// Render executes a given template with the provided data and writes it back to the client.
func Render(w http.ResponseWriter, tmpl *template.Template, data interface{}) {
	if tmpl == nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Page not found."))
		return
	}
	if err := tmpl.Execute(w, data); err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to render requested page."))
		return
	}
}

// Respond converts a Go value to JSON and sends it to the client.
func Respond(w http.ResponseWriter, data interface{}, statusCode int) error {

	// Convert the response value to JSON.
	res, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Respond with the provided JSON.
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	if _, err := w.Write(res); err != nil {
		return err
	}

	return nil
}
