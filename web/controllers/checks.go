package controllers

import (
	"net/http"

	"github.com/gophergala/goffee/web/render"
	"github.com/zenazn/goji/web"
)

// ChecksIndex render the checks index for the current user
func ChecksIndex(c web.C, w http.ResponseWriter, req *http.Request) {
	if !userSignedIn(c) {
		http.Error(w, "You need to re-authenticate", http.StatusUnauthorized)
		return
	}

	templates := render.GetBaseTemplates()
	templates = append(templates, "web/views/checks.html")
	err := render.Template(c, w, templates, "layout", map[string]string{"Title": "Checks"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
