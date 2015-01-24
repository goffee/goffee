package web

import (
	"net/http"

	"github.com/gophergala/goffee/data"
	"github.com/gophergala/goffee/web/controllers"
	"github.com/gorilla/sessions"
	"github.com/hypebeast/gojistatic"
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

	goji.Use(gojistatic.Static("web/public", gojistatic.StaticOptions{SkipLogging: true}))
	goji.Use(SessionMiddleware)

	goji.Get("/", controllers.Home)

	goji.Get("/ip", controllers.IP)

	goji.Get("/oauth/authorize", controllers.OAuthAuthorize)
	goji.Get("/oauth/callback", controllers.OAuthCallback)

	goji.Get("/checks", controllers.ChecksIndex)

	goji.Serve()
}
