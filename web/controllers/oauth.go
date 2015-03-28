package controllers

import (
	"net/http"

	"github.com/goffee/goffee/data"
	"github.com/google/go-github/github"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"golang.org/x/oauth2"
)

// OAuthConf is used to configure OAuth
var OAuthConf *oauth2.Config

// OAuthAuthorize makes the user login
func OAuthAuthorize(s sessions.Session, req *http.Request, r render.Render) {
	url := OAuthConf.AuthCodeURL("state", oauth2.AccessTypeOnline)
	r.Redirect(url, http.StatusFound)
}

// SignOut makes the user sign out
func SignOut(s sessions.Session, req *http.Request, r render.Render) {
	s.Clear()
	r.Redirect("/", http.StatusFound)
}

// OAuthCallback handles the callback from GitHub
func OAuthCallback(s sessions.Session, req *http.Request, r render.Render) {
	// Get the code from the response
	code := req.URL.Query().Get("code")

	// Exchange the received code for a token
	token, err := OAuthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		panic(err)
	}

	httpClient := OAuthConf.Client(oauth2.NoContext, token)
	githubClient := github.NewClient(httpClient)
	user, _, err := githubClient.Users.Get("")

	if err != nil {
		panic(err)
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

	s.Set("UserId", u.Id)

	r.Redirect("/checks", http.StatusFound)
}
