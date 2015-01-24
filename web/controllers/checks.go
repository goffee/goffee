package controllers

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gophergala/goffee/data"
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

func NewCheck(c web.C, w http.ResponseWriter, req *http.Request) {
	if !userSignedIn(c) {
		http.Error(w, "You need to re-authenticate", http.StatusUnauthorized)
		return
	}

	templates := render.GetBaseTemplates()
	templates = append(templates, "web/views/new_check.html")
	err := render.Template(c, w, templates, "layout", map[string]string{"Title": "New Check"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func CreateCheck(c web.C, w http.ResponseWriter, req *http.Request) {
	user, err := currentUser(c)

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
}

func ShowCheck(c web.C, w http.ResponseWriter, req *http.Request) {
	user, err := currentUser(c)

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
