package controllers

import (
	"net/http"

	"github.com/goffee/goffee/Godeps/_workspace/src/github.com/zenazn/goji/web"
	"github.com/goffee/goffee/web/render"
)

func renderError(c web.C, w http.ResponseWriter, req *http.Request, message string, status int) {
	templates := render.GetBaseTemplates()
	templates = append(templates, "web/views/error.html")
	err := render.Template(c, w, req, templates, "layout", map[string]interface{}{"Status": status, "Message": message})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Home serves the home page
func Home(c web.C, w http.ResponseWriter, req *http.Request) {
	templates := render.GetBaseTemplates()
	templates = append(templates, "web/views/home.html")
	err := render.Template(c, w, req, templates, "layout", map[string]interface{}{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// About serves the about page
func About(c web.C, w http.ResponseWriter, req *http.Request) {
	templates := render.GetBaseTemplates()
	templates = append(templates, "web/views/about.html")
	err := render.Template(c, w, req, templates, "layout", map[string]interface{}{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// NotFound serves the 404 page
func NotFound(c web.C, w http.ResponseWriter, req *http.Request) {
	templates := render.GetBaseTemplates()
	templates = append(templates, "web/views/404.html")
	err := render.Template(c, w, req, templates, "layout", map[string]interface{}{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
