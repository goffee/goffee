package helpers

import (
	"errors"

	"github.com/gophergala/goffee/data"
	"github.com/gorilla/sessions"
	"github.com/zenazn/goji/web"
)

// CurrentSession returns the current session
func CurrentSession(c web.C) *sessions.Session {
	session := c.Env["Session"].(*sessions.Session)
	return session
}

// CurrentUser returns the current user
func CurrentUser(c web.C) (data.User, error) {
	session := CurrentSession(c)
	userID := session.Values["UserId"]

	switch userID := userID.(type) {
	case int64:
		return data.FindUser(userID)
	default:
		return data.User{}, errors.New("User not found")
	}
}

// UserSignedIn returns true if there is an authenticated user
func UserSignedIn(c web.C) bool {
	_, err := CurrentUser(c)
	return err == nil
}
