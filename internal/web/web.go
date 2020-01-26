package web

import (
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
