package helpers

import (
	"bytes"
	"html/template"
)

// Parse parses and renders a template
func Parse(t *template.Template, name string, data interface{}) string {
	var doc bytes.Buffer
	t.ExecuteTemplate(&doc, name, data)
	return doc.String()
}
