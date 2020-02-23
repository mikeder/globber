package handlers

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/favicon.ico")
}

func (s *site) loadTemplates() {
	var allFiles []string
	var templateDir = "./templates/html/"
	files, err := ioutil.ReadDir(templateDir)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		filename := file.Name()
		if strings.HasSuffix(filename, ".html") {
			allFiles = append(allFiles, templateDir+filename)
		}
	}

	templates, err := template.ParseFiles(allFiles...)
	if err != nil {
		log.Fatal(err)
	}

	templates.Funcs(template.FuncMap{
		"html": func(value string) template.HTML {
			return template.HTML(value)
		},
	})
	s.templates = templates
}
