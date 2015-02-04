package controllers

// import (
// 	"fmt"
// 	"net/http"
// 	"net/url"
// 	"strconv"
//
// 	"github.com/goffee/goffee/Godeps/_workspace/src/github.com/justinas/nosurf"
// 	"github.com/goffee/goffee/Godeps/_workspace/src/github.com/zenazn/goji/web" // ChecksIndex render the checks index for the current user
// 	"github.com/goffee/goffee/data"
// 	"github.com/goffee/goffee/web/helpers"
// 	"github.com/goffee/goffee/web/render"
// )
//
// func ChecksIndex(c web.C, w http.ResponseWriter, req *http.Request) {
// 	user, err := helpers.CurrentUser(c)
//
// 	if err != nil {
// 		renderError(c, w, req, "You need to re-authenticate", http.StatusUnauthorized)
// 		return
// 	}
//
// 	checks, err := user.Checks()
// 	if err != nil {
// 		renderError(c, w, req, "Something went wrong", http.StatusInternalServerError)
// 		return
// 	}
//
// 	templates := render.GetBaseTemplates()
// 	templates = append(templates, "web/views/checks.html")
// 	err = render.Template(c, w, req, templates, "layout", map[string]interface{}{"Title": "Checks", "Checks": checks})
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }
//
// // NewCheck renders the new check form
// func NewCheck(c web.C, w http.ResponseWriter, req *http.Request) {
// 	user, err := helpers.CurrentUser(c)
//
// 	if err != nil {
// 		renderError(c, w, req, "You need to re-authenticate", http.StatusUnauthorized)
// 		return
// 	}
//
// 	checksCount, err := user.ChecksCount()
//
// 	templates := render.GetBaseTemplates()
// 	templates = append(templates, "web/views/new_check.html")
// 	csrf := nosurf.Token(req)
// 	err = render.Template(c, w, req, templates, "layout", map[string]interface{}{"Title": "New Check", "CSRFToken": csrf, "ChecksCount": checksCount})
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }
//
// // CreateCheck saves a new check to the DB
// func CreateCheck(c web.C, w http.ResponseWriter, req *http.Request) {
// 	user, err := helpers.CurrentUser(c)
//
// 	if err != nil {
// 		renderError(c, w, req, "You need to re-authenticate", http.StatusUnauthorized)
// 		return
// 	}
//
// 	checksCount, err := user.ChecksCount()
// 	if err != nil {
// 		renderError(c, w, req, "Something went wrong", http.StatusInternalServerError)
// 		return
// 	}
//
// 	if checksCount >= 5 {
// 		session := helpers.CurrentSession(c)
// 		session.AddFlash("You have too many checks in use, delete some to add more.")
// 		session.Save(req, w)
// 		http.Redirect(w, req, "/checks", http.StatusSeeOther)
// 		return
// 	}
//
// 	u, err := url.Parse(req.FormValue("url"))
// 	if err != nil || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") {
// 		http.Redirect(w, req, "/checks/new", http.StatusSeeOther)
// 		return
// 	}
//
// 	check := &data.Check{URL: u.String(), UserId: user.Id}
// 	check.Create()
//
// 	path := fmt.Sprintf("/checks/%d", check.Id)
//
// 	http.Redirect(w, req, path, http.StatusSeeOther)
// }
//
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
