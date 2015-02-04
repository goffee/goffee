package controllers

import (
	"net/http"

	"github.com/goffee/goffee/Godeps/_workspace/src/github.com/google/go-github/github"
	"github.com/goffee/goffee/Godeps/_workspace/src/github.com/gorilla/sessions"
	"github.com/goffee/goffee/Godeps/_workspace/src/github.com/zenazn/goji/web"
	"github.com/goffee/goffee/Godeps/_workspace/src/golang.org/x/oauth2"
	"github.com/goffee/goffee/data"
	"github.com/goffee/goffee/web/helpers"
)

var OAuthConf *oauth2.Config

// OAuthAuthorize makes the user login
func OAuthAuthorize(w http.ResponseWriter, req *http.Request) {
	url := OAuthConf.AuthCodeURL("state", oauth2.AccessTypeOnline)
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
	token, err := OAuthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		renderError(c, w, req, "Authentication failed", http.StatusUnauthorized)
		return
	}

	session := c.Env["Session"].(*sessions.Session)

	httpClient := OAuthConf.Client(oauth2.NoContext, token)
	githubClient := github.NewClient(httpClient)
	user, _, err := githubClient.Users.Get("")

	if err != nil {
		renderError(c, w, req, "Authentication failed", http.StatusUnauthorized)
		return
	}

	u := &data.User{
		Name:        *user.Name,
		GitHubId:    int64(*user.ID),
		GitHubLogin: *user.Login,
		Email:       *user.Email,
		OAuthToken:  token.AccessToken,
	}
	u.UpdateOrCreate()

	if u.Email == "" {
		emails, _, err := githubClient.Users.ListEmails(nil)
		if err == nil && len(emails) > 0 {
			email := emails[0]
			u.Email = *email.Email
			u.UpdateOrCreate()
		}
	}

	session.Values["UserId"] = u.Id
	session.Save(req, w)

	http.Redirect(w, req, "/checks", http.StatusFound)
}
