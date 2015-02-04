package controllers

import (
	"net/http"
	"strconv"

	"github.com/goffee/goffee/Godeps/_workspace/src/github.com/zenazn/goji/web" // ResultsIndex renders results JSON
	"github.com/goffee/goffee/web/helpers"
	"github.com/goffee/goffee/web/render"
)

func ResultsIndex(c web.C, w http.ResponseWriter, req *http.Request) {
	user, err := helpers.CurrentUser(c)

	if err != nil {
		renderError(c, w, req,"You need to re-authenticate", http.StatusUnauthorized)
		return
	}

	checkId, err := strconv.ParseInt(c.URLParams["check_id"], 10, 64)
	if err != nil {
		renderError(c, w, req,"Check not found", http.StatusNotFound)
		return
	}

	check, err := user.Check(checkId)

	if err != nil {
		renderError(c, w, req,"Check not found", http.StatusNotFound)
		return
	}

	results, err := check.Results()

	if err != nil {
		renderError(c, w, req,"No results found", http.StatusNotFound)
		return
	}

	render.JSON(w, http.StatusOK, results)
}
