package web

import (
	"net/http"

	"github.com/gophergala/goffee/data"
	"github.com/gophergala/goffee/web/controllers"
	"github.com/gorilla/sessions"
	"github.com/hypebeast/gojistatic"
	"github.com/justinas/nosurf"
	"github.com/unrolled/secure"
	"github.com/zenazn/goji/graceful"
	"github.com/zenazn/goji/web"
)

var store = sessions.NewCookieStore([]byte("something-very-secret"))

// SessionMiddleware adds session support to Goffee
func SessionMiddleware(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// Get a session. We're ignoring the error resulted from decoding an
		// existing session: Get() always returns a session, even if empty.
		session, _ := store.Get(r, "goffee-session")

		// Save it.
		session.Save(r, w)

		c.Env = make(map[string]interface{})
		c.Env["Session"] = session

		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// StartServer starts the web server
func StartServer(bind string) {
	data.InitDatabase()

	secureMiddleware := secure.New(secure.Options{
		AllowedHosts:       []string{"example.com", "ssl.example.com"},
		FrameDeny:          true,
		ContentTypeNosniff: true,
		BrowserXssFilter:   true,
		IsDevelopment:      true,
	})

	m := web.New()

	m.Use(gojistatic.Static("web/public", gojistatic.StaticOptions{SkipLogging: true}))
	m.Use(secureMiddleware.Handler)
	m.Use(SessionMiddleware)
	m.Use(nosurf.NewPure)

	m.Get("/", controllers.Home)

	m.Get("/ip", controllers.IP)

	m.Get("/oauth/authorize", controllers.OAuthAuthorize)
	m.Get("/oauth/callback", controllers.OAuthCallback)
	m.Get("/sign_out", controllers.SignOut)

	m.Get("/checks", controllers.ChecksIndex)
	m.Get("/checks/new", controllers.NewCheck)
	m.Get("/checks/:id", controllers.ShowCheck)
	m.Post("/checks", controllers.CreateCheck)

	go graceful.ListenAndServe(bind, m)
}

func Wait() {
	graceful.Wait()
}
