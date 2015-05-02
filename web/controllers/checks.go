package controllers

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-martini/martini"
	"github.com/goffee/goffee/data"
	"github.com/goffee/goffee/web/helpers"
	"github.com/martini-contrib/csrf"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
)

// ChecksIndex displays the user's checks
func ChecksIndex(s sessions.Session, req *http.Request, r render.Render) {
	user, err := helpers.CurrentUser(s)

	if err != nil {
		r.Redirect("/", http.StatusFound)
		return
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
		r.Redirect("/", http.StatusFound)
		return
	}

	checksCount, err := user.ChecksCount()

	csrf := x.GetToken()

	r.HTML(200, "new_check", map[string]interface{}{"Title": "New Check", "CSRFToken": csrf, "ChecksCount": checksCount})
}

// CreateCheck saves a new check to the DB
func CreateCheck(s sessions.Session, req *http.Request, r render.Render) {
	user, err := helpers.CurrentUser(s)

	if err != nil {
		r.Redirect("/", http.StatusFound)
		return
	}

	checksCount, err := user.ChecksCount()
	if err != nil {
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

// ShowCheck renders a single check
func ShowCheck(s sessions.Session, req *http.Request, r render.Render, params martini.Params, x csrf.CSRF) {
	user, err := helpers.CurrentUser(s)

	if err != nil {
		r.Redirect("/", http.StatusFound)
		return
	}

	checkID, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		panic(err)
	}

	check, err := user.Check(checkID)

	if err != nil {
		panic(err)
	}

	csrf := x.GetToken()

	r.HTML(200, "check", map[string]interface{}{"Title": "Check: " + check.URL, "Check": check, "CSRFToken": csrf})
}

// DeleteCheck deletes a single check
func DeleteCheck(s sessions.Session, req *http.Request, r render.Render, params martini.Params) {
	user, err := helpers.CurrentUser(s)

	if err != nil {
		r.Redirect("/", http.StatusFound)
		return
	}

	checkID, err := strconv.ParseInt(params["id"], 10, 64)
	if err != nil {
		panic(err)
	}

	check, err := user.Check(checkID)

	if err != nil {
		panic(err)
	}

	check.Delete()
	s.AddFlash("Check deleted")

	r.Redirect("/checks", http.StatusSeeOther)
}
