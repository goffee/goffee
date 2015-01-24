package render

// https://github.com/hypebeast/goji-boilerplate/blob/master/render/render.go

import (
	"encoding/json"
	"html/template"
	"net/http"
)

// Template renders HTML templates
func Template(w http.ResponseWriter, templates []string, name string, data interface{}) error {
	t, err := template.ParseFiles(templates...)
	if err != nil {
		return err
	}

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
