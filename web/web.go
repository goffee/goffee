package web

import (
	"github.com/christopherobin/authy/martini"
	"github.com/go-martini/martini"
	"github.com/goffee/goffee/Godeps/_workspace/src/github.com/zenazn/goji/graceful"
	"github.com/goffee/goffee/web/controllers"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
)

var SessionStore *sessions.CookieStore
var AuthyConfig authy.Config

// SessionMiddleware adds session support to Goffee
// func SessionMiddleware(c *web.C, h http.Handler) http.Handler {
// 	fn := func(w http.ResponseWriter, r *http.Request) {
// 		// Get a session. We're ignoring the error resulted from decoding an
// 		// existing session: Get() always returns a session, even if empty.
// 		session, _ := SessionStore.Get(r, "goffee-session")
//
// 		// Save it.
// 		session.Save(r, w)
//
// 		c.Env["Session"] = session
//
// 		h.ServeHTTP(w, r)
// 	}
// 	return http.HandlerFunc(fn)
// }

// StartServer starts the web server
func StartServer(bind string) {
	// secureMiddleware := secure.New(secure.Options{
	// 	AllowedHosts:       []string{"example.com", "ssl.example.com"},
	// 	FrameDeny:          true,
	// 	ContentTypeNosniff: true,
	// 	BrowserXssFilter:   true,
	// 	IsDevelopment:      true,
	// })
	//
	// m := web.New()
	//
	// m.Use(middleware.RealIP)
	// m.Use(gojistatic.Static("web/public", gojistatic.StaticOptions{SkipLogging: true}))
	// m.Use(middleware.EnvInit)
	// m.Use(secureMiddleware.Handler)
	// m.Use(SessionMiddleware)
	// m.Use(nosurf.NewPure)
	//
	// m.Get("/oauth/authorize", controllers.OAuthAuthorize)
	// m.Get("/oauth/callback", controllers.OAuthCallback)
	// m.Get("/sign_out", controllers.SignOut)
	//
	// m.Get("/checks", controllers.ChecksIndex)
	// m.Get("/checks/new", controllers.NewCheck)
	// m.Get("/checks/:id", controllers.ShowCheck)
	// m.Post("/checks/:id/delete", controllers.DeleteCheck)
	// m.Post("/checks", controllers.CreateCheck)
	//
	// m.Get("/checks/:check_id/results", controllers.ResultsIndex)
	//
	// m.NotFound(controllers.NotFound)

	m := martini.Classic()

	m.Use(martini.Static("web/public"))

	store := sessions.NewCookieStore([]byte("secret123"))
	m.Use(sessions.Sessions("authy", store))

	m.Use(render.Renderer(render.Options{
		Directory: "web/views",
		Delims:    render.Delims{"{{{", "}}}"},
		Layout:    "layout",
	}))

	m.Use(authy.Authy(AuthyConfig))

	m.Get("/callback", authy.LoginRequired(), controllers.Callback)

	m.Get("/", controllers.Home)
	m.Get("/about", controllers.About)
	m.Get("/ip", controllers.IP)

	go m.RunOnAddr(bind)
	// go graceful.ListenAndServe(bind, m)
}

// Wait ...
func Wait() {
	graceful.Wait()
}
