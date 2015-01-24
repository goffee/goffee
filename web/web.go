package web

import (
	"fmt"
	"net/http"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

func root(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Goffee")
}

// StartServer starts the web server
func StartServer() {
	goji.Get("/", root)
	goji.Serve()
}
