package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/google/go-github/github"
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

// OAuthCallback handles the callback from GitHub
func OAuthCallback(w http.ResponseWriter, req *http.Request) {
	// Get the code from the response
	code := req.FormValue("code")

	// Exchange the received code for a token
	token, err := conf.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatal(err)
	}

	httpClient := conf.Client(oauth2.NoContext, token)
	githubClient := github.NewClient(httpClient)
	user, _, err := githubClient.Users.Get("")

	fmt.Fprintf(w, "Token: %s!", user.String())
}
