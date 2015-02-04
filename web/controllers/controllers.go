package controllers

import (
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
)

// func renderError(c web.C, w http.ResponseWriter, req *http.Request, message string, status int) {
// 	templates := render.GetBaseTemplates()
// 	templates = append(templates, "web/views/error.html")
// 	err := render.Template(c, w, req, templates, "layout", map[string]interface{}{"Status": status, "Message": message})
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }

// Home serves the home page
// func Home(c web.C, w http.ResponseWriter, req *http.Request) {
// 	templates := render.GetBaseTemplates()
// 	templates = append(templates, "web/views/home.html")
// 	err := render.Template(c, w, req, templates, "layout", map[string]interface{}{})
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }

// Home serves the home page
func Home(s sessions.Session, r render.Render) {
	s.Set("hello", "world")
	r.HTML(200, "home", map[string]interface{}{})
}

// About serves the about page
func About(r render.Render) {
	r.HTML(200, "about", map[string]interface{}{})
}

// About serves the about page
// func About(c web.C, w http.ResponseWriter, req *http.Request) {
// 	templates := render.GetBaseTemplates()
// 	templates = append(templates, "web/views/about.html")
// 	err := render.Template(c, w, req, templates, "layout", map[string]interface{}{})
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }
//
// // NotFound serves the 404 page
// func NotFound(c web.C, w http.ResponseWriter, req *http.Request) {
// 	templates := render.GetBaseTemplates()
// 	templates = append(templates, "web/views/404.html")
// 	err := render.Template(c, w, req, templates, "layout", map[string]interface{}{})
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}
// }
