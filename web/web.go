package web

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"html/template"

	"github.com/go-martini/martini"
	"github.com/goffee/goffee/data"
	"github.com/goffee/goffee/web/controllers"
	"github.com/goffee/goffee/web/helpers"
	"github.com/goffee/goffee/web/middleware"
	"github.com/martini-contrib/csrf"
	"github.com/martini-contrib/method"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/secure"
	"github.com/martini-contrib/sessions"
	"github.com/tylerb/graceful"
)

var srv *graceful.Server

// SessionSecret ...
var SessionSecret string

func formatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

var templateFuncs = template.FuncMap{
	"formatTime": func(t time.Time) string {
		return t.Format(time.RFC3339)
	},
	// define an empty stub first, otherwise html/template will complain with "missing function"
	"currentUser": func() (data.User, error) {
		return data.User{}, errors.New("User not found")
	},
	"userSignedIn": func() bool {
		return false
	},
}

// middleware to inject the route
func helperFuncs() martini.Handler {
	return func(ren render.Render, s sessions.Session) {
		ren.Template().Funcs(injectHelperFuncs(s))
	}
}

// create the real template helpers
var injectHelperFuncs = func(s sessions.Session) template.FuncMap {
	templateFuncs["currentUser"] = func() (data.User, error) {
		return helpers.CurrentUser(s)
	}
	templateFuncs["userSignedIn"] = func() bool {
		return helpers.UserSignedIn(s)
	}
	return templateFuncs
}

// StartServer starts the web server
func StartServer(bind string) {
	// Equivalent to:
	// m := martini.Classic()
	r := martini.NewRouter()
	ma := martini.New()
	ma.Use(martini.Logger())
	ma.Use(middleware.Recovery())
	ma.Use(martini.Static("public"))
	ma.MapTo(r, (*martini.Routes)(nil))
	ma.Action(r.Handle)
	m := &martini.ClassicMartini{ma, r}

	m.Use(martini.Static("web/public"))

	store := sessions.NewCookieStore([]byte(SessionSecret))
	m.Use(sessions.Sessions("goffee-session", store))

	m.Use(render.Renderer(render.Options{
		Directory: "web/views",
		Delims:    render.Delims{"{{{", "}}}"},
		Layout:    "layout",
		Funcs:     []template.FuncMap{templateFuncs},
	}))
	m.Use(helperFuncs())

	m.Use(method.Override())

	m.Use(secure.Secure(secure.Options{
		AllowedHosts:            []string{"www.goffee.io", "goffee.io", "localhost:8000", "localhost"},
		SSLRedirect:             false,
		STSSeconds:              0, // STSSeconds is the max-age of the Strict-Transport-Security header. Default is 0, which would NOT include the header.
		STSIncludeSubdomains:    false,
		FrameDeny:               true,
		CustomFrameOptionsValue: "SAMEORIGIN",
		ContentTypeNosniff:      true,
		BrowserXssFilter:        true,
	}))

	m.Use(csrf.Generate(&csrf.Options{
		Secret:     SessionSecret,
		SessionKey: "UserId",
		ErrorFunc: func(w http.ResponseWriter) {
			body, err := ioutil.ReadFile("web/public/422.html")
			if err != nil {
				panic(err)
			}

			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(422)
			fmt.Fprint(w, string(body))
		},
	}))

	m.Get("/", controllers.Home)
	m.Get("/about", controllers.About)
	m.Get("/ip", middleware.RealIP(), controllers.IP)

	m.Get("/oauth/authorize", controllers.OAuthAuthorize)
	m.Get("/oauth/callback", controllers.OAuthCallback)
	m.Get("/sign_out", controllers.SignOut)

	m.Get("/checks", controllers.ChecksIndex)
	m.Get("/checks/new", controllers.NewCheck)
	m.Get("/checks/:id", controllers.ShowCheck)
	m.Post("/checks/:id/delete", csrf.Validate, controllers.DeleteCheck)
	m.Post("/checks", csrf.Validate, controllers.CreateCheck)

	m.Get("/checks/:check_id/results", controllers.ResultsIndex)

	m.NotFound(controllers.NotFound)

	srv = &graceful.Server{
		Server:           &http.Server{Addr: bind, Handler: m},
		NoSignalHandling: true,
	}
	go srv.ListenAndServe()
}

func Stop() {
	srv.Stop(10 * time.Second)
}

// Wait ...
func Wait() {
	<-srv.StopChan()
}
