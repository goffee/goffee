package controllers

import (
	"net/http"

	"github.com/gophergala/goffee/web/render"
)

// ChecksIndex render the checks index for the current user
func ChecksIndex(w http.ResponseWriter, req *http.Request) {
	templates := render.GetBaseTemplates()
	templates = append(templates, "web/views/checks.html")
	err := render.Template(w, templates, "layout", map[string]string{"Title": "Checks"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
