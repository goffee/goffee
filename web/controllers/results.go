package controllers

import (
	"net/http"
	"strconv"

	"github.com/gophergala/goffee/web/helpers"
	"github.com/gophergala/goffee/web/render"
	"github.com/zenazn/goji/web"
)

// ResultsIndex renders results JSON
func ResultsIndex(c web.C, w http.ResponseWriter, req *http.Request) {
	user, err := helpers.CurrentUser(c)

	if err != nil {
		http.Error(w, "You need to re-authenticate", http.StatusUnauthorized)
		return
	}

	checkId, err := strconv.ParseInt(c.URLParams["check_id"], 10, 64)
	if err != nil {
		http.Error(w, "Check not found", http.StatusNotFound)
		return
	}

	check, err := user.Check(checkId)

	if err != nil {
		http.Error(w, "Check not found", http.StatusNotFound)
		return
	}

	results, err := check.Results()

	if err != nil {
		http.Error(w, "Results not found", http.StatusNotFound)
		return
	}

	render.JSON(w, http.StatusOK, results)
}
