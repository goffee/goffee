package web

import (
	"net/http"

	"github.com/gophergala/goffee/data"
	"github.com/gophergala/goffee/web/controllers"
	"github.com/gorilla/sessions"
	"github.com/hypebeast/gojistatic"
	"github.com/justinas/nosurf"
	"github.com/unrolled/secure"
	"github.com/zenazn/goji"
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

		c.Env["Session"] = session

		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// StartServer starts the web server
func StartServer() {
	data.InitDatabase()

	secureMiddleware := secure.New(secure.Options{
		AllowedHosts:       []string{"example.com", "ssl.example.com"},
		FrameDeny:          true,
		ContentTypeNosniff: true,
		BrowserXssFilter:   true,
		IsDevelopment:      true,
	})

	goji.Use(gojistatic.Static("web/public", gojistatic.StaticOptions{SkipLogging: true}))
	goji.Use(secureMiddleware.Handler)
	goji.Use(SessionMiddleware)
	goji.Use(nosurf.NewPure)

	goji.Get("/", controllers.Home)

	goji.Get("/ip", controllers.IP)

	goji.Get("/oauth/authorize", controllers.OAuthAuthorize)
	goji.Get("/oauth/callback", controllers.OAuthCallback)
	goji.Get("/sign_out", controllers.SignOut)

	goji.Get("/checks", controllers.ChecksIndex)
	goji.Get("/checks/new", controllers.NewCheck)
	goji.Get("/checks/:id", controllers.ShowCheck)
	goji.Post("/checks", controllers.CreateCheck)

	goji.Serve()
}
