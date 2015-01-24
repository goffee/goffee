package render

// https://github.com/hypebeast/goji-boilerplate/blob/master/render/render.go

import (
	"encoding/json"
	"html/template"
	"net/http"
	"time"

	"github.com/gophergala/goffee/web/helpers"
	"github.com/zenazn/goji/web"
)

func formatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

// Template renders HTML templates
func Template(c web.C, w http.ResponseWriter, templates []string, name string, data map[string]interface{}) error {
	funcMap := template.FuncMap{
		"formatTime": formatTime,
	}

	t, err := template.New("").Funcs(funcMap).ParseFiles(templates...)
	if err != nil {
		return err
	}

	var loggedIn bool
	user, err := helpers.CurrentUser(c)

	if err != nil {
		loggedIn = false
	} else {
		loggedIn = true
	}

	data["CurrentUser"] = user
	data["UserSignedIn"] = loggedIn

	err = t.ExecuteTemplate(w, name, data)
	if err != nil {
		return err
	}

	return nil
}

// JSON renders JSON
func JSON(w http.ResponseWriter, status int, v interface{}) {
	var result []byte
	var err error

	result, err = json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// json rendered fine, write out the result
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(result)
}

// GetBaseTemplates does things
func GetBaseTemplates() []string {
	templates := []string{"web/views/layout.html"}
	return templates
}
