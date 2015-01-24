package controllers

import (
	"log"
	"net/http"

	"github.com/google/go-github/github"
	"github.com/gophergala/goffee/data"
	"github.com/gophergala/goffee/web/helpers"
	"github.com/gorilla/sessions"
	"github.com/zenazn/goji/web"
	"golang.org/x/oauth2"
	oauth2github "golang.org/x/oauth2/github"
)

var conf = &oauth2.Config{
	ClientID:     "508322171059309cedad",
	ClientSecret: "8cb47d06cc58c8bc2b4c0d01b870d117cf086c40",
	Scopes:       []string{},
	Endpoint:     oauth2github.Endpoint,
}

// OAuthAuthorize makes the user login
func OAuthAuthorize(w http.ResponseWriter, req *http.Request) {
	url := conf.AuthCodeURL("state", oauth2.AccessTypeOnline)
	http.Redirect(w, req, url, http.StatusFound)
}

// SignOut makes the user sign out
func SignOut(c web.C, w http.ResponseWriter, req *http.Request) {
	session := helpers.CurrentSession(c)
	session.Values = map[interface{}]interface{}{}
	session.Save(req, w)
	http.Redirect(w, req, "/", http.StatusFound)
}

// OAuthCallback handles the callback from GitHub
func OAuthCallback(c web.C, w http.ResponseWriter, req *http.Request) {
	// Get the code from the response
	code := req.FormValue("code")

	// Exchange the received code for a token
	token, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatal(err)
	}

	session := c.Env["Session"].(*sessions.Session)

	httpClient := conf.Client(oauth2.NoContext, token)
	githubClient := github.NewClient(httpClient)
	user, _, err := githubClient.Users.Get("")

	if err != nil {
		session.AddFlash("Authentication failed")
		http.Redirect(w, req, "/", http.StatusSeeOther)
	}

	u := &data.User{
		Name:        *user.Name,
		GitHubId:    int64(*user.ID),
		GitHubLogin: *user.Login,
		Email:       *user.Email,
		OAuthToken:  token.AccessToken,
	}
	u.UpdateOrCreate()

	session.Values["UserId"] = u.Id
	session.Save(req, w)

	http.Redirect(w, req, "/checks", http.StatusFound)
}
