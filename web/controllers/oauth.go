package controllers

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var conf = &oauth2.Config{
	ClientID:     "508322171059309cedad",
	ClientSecret: "8cb47d06cc58c8bc2b4c0d01b870d117cf086c40",
	Scopes:       []string{},
	Endpoint:     github.Endpoint,
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

	// client := conf.Client(oauth2.NoContext, token)

	fmt.Fprintf(w, "Token: %s!", token)
}
