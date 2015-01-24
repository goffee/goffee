package controllers

import (
	"net/http"

	"github.com/gophergala/goffee/web/render"
	"github.com/zenazn/goji/web"
)

// Home serves the home page
func Home(c web.C, w http.ResponseWriter, req *http.Request) {
	templates := render.GetBaseTemplates()
	templates = append(templates, "web/views/home.html")
	err := render.Template(c, w, templates, "layout", map[string]interface{}{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
