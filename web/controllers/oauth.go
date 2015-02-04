package controllers

import (
	"net/http"

	"github.com/gophergala/authy"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
)

// Callback saves the user to the session and redirects the user
func Callback(s sessions.Session, t authy.Token, r render.Render) {
	// TODO: save the user

	// Token: http://godoc.org/github.com/christopherobin/authy#Token

	// 	// Get the code from the response
	// 	code := r.URL.Query().Get("code")
	//
	// 	// Exchange the received code for a token
	// 	token, err := OAuthConf.Exchange(oauth2.NoContext, code)
	// 	if err != nil {
	// 		log.fatal(err.Error())
	// 		return
	// 	}
	//
	// 	session := c.Env["Session"].(*sessions.Session)
	//
	// 	httpClient := OAuthConf.Client(oauth2.NoContext, token)
	// 	githubClient := github.NewClient(httpClient)
	// 	user, _, err := githubClient.Users.Get("")
	//
	// 	if err != nil {
	// 		// http.Redirect(w, r, PathError, codeRedirect)
	// 		log.fatal(err.Error())
	// 		return
	// 	}
	//
	// 	u := &data.User{
	// 		Name:        *user.Name,
	// 		GitHubId:    int64(*user.ID),
	// 		GitHubLogin: *user.Login,
	// 		Email:       *user.Email,
	// 		OAuthToken:  token.AccessToken,
	// 	}
	// 	u.UpdateOrCreate()
	//
	// 	if u.Email == "" {
	// 		emails, _, err := githubClient.Users.ListEmails(nil)
	// 		if err == nil && len(emails) > 0 {
	// 			email := emails[0]
	// 			u.Email = *email.Email
	// 			u.UpdateOrCreate()
	// 		}
	// 	}
	//
	// 	s.Set("UserId", u.Id)

	r.Redirect("/checks", http.StatusFound)
}
