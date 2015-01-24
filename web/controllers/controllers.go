package controllers

import (
	"net/http"

	"github.com/gophergala/goffee/data"
	"github.com/gophergala/goffee/web/render"
	"github.com/gorilla/sessions"
	"github.com/zenazn/goji/web"
)

func currentSession(c web.C) *sessions.Session {
	session := c.Env["Session"].(*sessions.Session)
	return session
}

func currentUser(c web.C) (data.User, error) {
	session := currentSession(c)
	userID := session.Values["UserId"].(int64)
	return data.FindUser(userID)
}

func userSignedIn(c web.C) bool {
	_, err := currentUser(c)
	return err != nil
}

// Home serves the home page
func Home(w http.ResponseWriter, req *http.Request) {
	templates := render.GetBaseTemplates()
	templates = append(templates, "web/views/home.html")
	err := render.Template(w, templates, "layout", map[string]string{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
