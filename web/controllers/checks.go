package controllers

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gophergala/goffee/data"
	"github.com/gophergala/goffee/web/helpers"
	"github.com/gophergala/goffee/web/render"
	"github.com/justinas/nosurf"
	"github.com/zenazn/goji/web"
)

// ChecksIndex render the checks index for the current user
func ChecksIndex(c web.C, w http.ResponseWriter, req *http.Request) {
	if !helpers.UserSignedIn(c) {
		http.Error(w, "You need to re-authenticate", http.StatusUnauthorized)
		return
	}

	templates := render.GetBaseTemplates()
	templates = append(templates, "web/views/checks.html")
	err := render.Template(c, w, templates, "layout", map[string]interface{}{"Title": "Checks"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// NewCheck renders the new check form
func NewCheck(c web.C, w http.ResponseWriter, req *http.Request) {
	if !helpers.UserSignedIn(c) {
		http.Error(w, "You need to re-authenticate", http.StatusUnauthorized)
		return
	}

	templates := render.GetBaseTemplates()
	templates = append(templates, "web/views/new_check.html")
	csrf := nosurf.Token(req)
	err := render.Template(c, w, templates, "layout", map[string]interface{}{"Title": "New Check", "CSRFToken": csrf})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// CreateCheck saves a new check to the DB
func CreateCheck(c web.C, w http.ResponseWriter, req *http.Request) {
	user, err := helpers.CurrentUser(c)

	if err != nil {
		http.Error(w, "You need to re-authenticate", http.StatusUnauthorized)
		return
	}

	u, err := url.Parse(req.FormValue("url"))
	if err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") {
		fmt.Println(err)
		http.Redirect(w, req, "/checks/new", http.StatusSeeOther)
		return
	}

	check := &data.Check{URL: u.String(), UserId: user.Id}
	check.Create()

	path := fmt.Sprintf("/checks/%d", check.Id)

	http.Redirect(w, req, path, http.StatusSeeOther)
}

// ShowCheck renders a single check
func ShowCheck(c web.C, w http.ResponseWriter, req *http.Request) {
	user, err := helpers.CurrentUser(c)

	if err != nil {
		http.Error(w, "You need to re-authenticate", http.StatusUnauthorized)
		return
	}

	checkId, err := strconv.ParseInt(c.URLParams["id"], 10, 64)
	if err != nil {
		http.Error(w, "Check not found", http.StatusNotFound)
		return
	}

	check, err := user.Check(checkId)

	templates := render.GetBaseTemplates()
	templates = append(templates, "web/views/check.html")
	err = render.Template(c, w, templates, "layout", map[string]interface{}{"Title": "Check: " + check.URL, "Check": check})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
