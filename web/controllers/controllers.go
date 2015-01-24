package controllers

import (
	"net/http"

	"github.com/gophergala/goffee/web/render"
)

// Home serves the home page
func Home(w http.ResponseWriter, req *http.Request) {
	templates := render.GetBaseTemplates()
	templates = append(templates, "web/views/home.html")
	err := render.Template(w, templates, "layout", map[string]string{"Title": "Home"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
