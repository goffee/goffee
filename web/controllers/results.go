package controllers

import (
	"net/http"
	"strconv"

	"github.com/go-martini/martini"
	"github.com/goffee/goffee/web/helpers"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
)

// ResultsIndex returns the results for a check as JSON
func ResultsIndex(s sessions.Session, req *http.Request, r render.Render, params martini.Params) {
	user, err := helpers.CurrentUser(s)

	if err != nil {
		r.Redirect("/", http.StatusFound)
		return
	}

	checkID, err := strconv.ParseInt(params["check_id"], 10, 64)
	if err != nil {
		panic(err)
	}

	check, err := user.Check(checkID)

	if err != nil {
		panic(err)
	}

	results, err := check.Results()

	if err != nil {
		panic(err)
	}

	r.JSON(200, results)
}
