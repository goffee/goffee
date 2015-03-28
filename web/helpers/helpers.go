package helpers

import (
	"errors"

	"github.com/goffee/goffee/data"
	"github.com/martini-contrib/sessions"
)

// CurrentUser returns the current user
func CurrentUser(session sessions.Session) (data.User, error) {
	// session := CurrentSession(c)
	userID := session.Get("UserId")

	switch userID := userID.(type) {
	case int64:
		return data.FindUser(userID)
	default:
		return data.User{}, errors.New("User not found")
	}
}

// UserSignedIn returns true if there is an authenticated user
func UserSignedIn(session sessions.Session) bool {
	_, err := CurrentUser(session)
	return err == nil
}
