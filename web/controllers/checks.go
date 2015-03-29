package controllers

import (
	"fmt"
	"net/http"
	"net/url"

	// "github.com/goffee/goffee/Godeps/_workspace/src/github.com/justinas/nosurf"
	// "github.com/goffee/goffee/Godeps/_workspace/src/github.com/zenazn/goji/web" // ChecksIndex render the checks index for the current user

	"github.com/goffee/goffee/data"
	"github.com/goffee/goffee/web/helpers"
	"github.com/martini-contrib/csrf"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	// "github.com/goffee/goffee/web/render"
)

// ChecksIndex displays the user's checks
func ChecksIndex(s sessions.Session, req *http.Request, r render.Render) {
	user, err := helpers.CurrentUser(s)

	if err != nil {
		panic(err)
	}

	checks, err := user.Checks()
	if err != nil {
		panic(err)
	}

	r.HTML(200, "checks", map[string]interface{}{"Title": "Checks", "Checks": checks})
}

// NewCheck renders the new check form
func NewCheck(s sessions.Session, req *http.Request, r render.Render, x csrf.CSRF) {
	user, err := helpers.CurrentUser(s)

	if err != nil {
		panic(err)
	}

	checksCount, err := user.ChecksCount()

	csrf := x.GetToken()

	r.HTML(200, "new_check", map[string]interface{}{"Title": "New Check", "CSRFToken": csrf, "ChecksCount": checksCount})
}

// CreateCheck saves a new check to the DB
func CreateCheck(s sessions.Session, req *http.Request, r render.Render) {
	user, err := helpers.CurrentUser(s)

	if err != nil {
		// renderError(c, w, req, "You need to re-authenticate", http.StatusUnauthorized)
		// return
		panic(err)
	}

	checksCount, err := user.ChecksCount()
	if err != nil {
		// renderError(c, w, req, "Something went wrong", http.StatusInternalServerError)
		// return
		panic(err)
	}

	if checksCount >= 5 {
		s.AddFlash("You have too many checks in use, delete some to add more.")
		r.Redirect("/checks", http.StatusFound)
		return
	}

	u, err := url.Parse(req.FormValue("url"))
	if err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") {
		r.Redirect("/checks", http.StatusFound)
		return
	}

	check := &data.Check{URL: u.String(), UserId: user.Id}
	check.Create()

	path := fmt.Sprintf("/checks/%d", check.Id)

	r.Redirect(path, http.StatusSeeOther)
}

// // ShowCheck renders a single check
// func ShowCheck(c web.C, w http.ResponseWriter, req *http.Request) {
// 	user, err := helpers.CurrentUser(c)
//
// 	if err != nil {
// 		renderError(c, w, req, "You need to re-authenticate", http.StatusUnauthorized)
// 		return
// 	}
//
// 	checkId, err := strconv.ParseInt(c.URLParams["id"], 10, 64)
// 	if err != nil {
// 		renderError(c, w, req, "Check not found", http.StatusNotFound)
// 		return
// 	}
//
// 	check, err := user.Check(checkId)
//
// 	if err != nil {
// 		renderError(c, w, req, "Check not found", http.StatusNotFound)
// 		return
// 	}
//
// 	csrf := nosurf.Token(req)
//
// 	templates := render.GetBaseTemplates()
// 	templates = append(templates, "web/views/check.html")
// 	err = render.Template(c, w, req, templates, "layout", map[string]interface{}{"Title": "Check: " + check.URL, "Check": check, "CSRFToken": csrf})
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }
//
// // DeleteCheck deletes a single check
// func DeleteCheck(c web.C, w http.ResponseWriter, req *http.Request) {
// 	user, err := helpers.CurrentUser(c)
//
// 	if err != nil {
// 		renderError(c, w, req, "You need to re-authenticate", http.StatusUnauthorized)
// 		return
// 	}
//
// 	checkId, err := strconv.ParseInt(c.URLParams["id"], 10, 64)
// 	if err != nil {
// 		renderError(c, w, req, "Check not found", http.StatusNotFound)
// 		return
// 	}
//
// 	check, err := user.Check(checkId)
//
// 	if err != nil {
// 		renderError(c, w, req, "Check not found", http.StatusNotFound)
// 		return
// 	}
//
// 	check.Delete()
// 	session := helpers.CurrentSession(c)
// 	session.AddFlash("Check deleted")
// 	session.Save(req, w)
//
// 	http.Redirect(w, req, "/checks", http.StatusSeeOther)
// }
